package actions

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"gopkg.in/yaml.v2"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	yamlFormat      = "yaml"
	yamlContentType = "application/yaml"
)

const (
	roleKind         = "RoleBinding"
	svcAccountKind   = "ServiceAccount"
	clusterKind      = "ClusterRoleBinding"
	clusterWideScope = "cluster-wide"
)

func HomeHandler(c buffalo.Context) error {

	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusMethodNotAllowed, r.JSON(map[string]string{"message": "method not support. please use a post a request"}))
	}

	// ====================================================================

	var (
		ep string
		kc string
	)

	type config struct {
		Subjects  string `json:"subject"`
		Format    string `json:"format"`
		Context   string `json:"context"`
		Kind      string `json:"kind"`
		Namespace string `json:"namespace"`
	}

	sc := &config{}
	if err := c.Bind(sc); err != nil {
		return err
	}

	var subjects = strings.Split(sc.Subjects, ",")

	// ===============================================================================

	cr := clientcmd.NewDefaultClientConfigLoadingRules()
	cr.ExplicitPath = ep
	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(cr, &clientcmd.ConfigOverrides{CurrentContext: kc}, )
	nc, err := cc.ClientConfig()
	if err != nil {
		log.Println("error reading client config")
		return c.Render(http.StatusInternalServerError, r.JSON(map[string]string{"message": "server error"}))
	}

	cs, err := kubernetes.NewForConfig(nc)
	if err != nil {
		log.Printf("error generating kubernetes clientset from client config: %v\n", err)
		return c.Render(http.StatusInternalServerError, r.JSON(map[string]string{"message": "server error"}))
	}

	// TODO: topic.Kind, subjectName, scope, sRole.Kind, sRole.Name, sRole.Source.Kind, sRole.Source.Name

	type response struct {
		Subject string
		Scope   string
		Role    []string
	}

	var roles = make([]string, 0)
	var results = make(map[string]map[string]response, 0)

	// ====================================================================
	//
	//	TODO: Call each search in a routine to speed up the process
	//
	// ====================================================================

	sort.Strings(subjects)

	for _, filter := range subjects {

		if len(subjects) > 1 && filter == "" {
			continue
		}

		var result = make(map[string]response, 0)

		results[filter] = result

		sr := search{
			search: filter,
			kind:   strings.ToLower(sc.Kind),
			k8si:   cs,
			scope:  make(map[string]topic),
		}

		if err := sr.roles(); err != nil {
			log.Printf("unable to load bindings: %v\n", err)
			return c.Render(http.StatusInternalServerError, r.JSON(map[string]string{"message": "server error"}))
		}

		if len(sr.scope) < 1 {
			return c.Render(200, r.JSON(map[string]string{"subjects": sc.Subjects, "roles": "none"}))
		}

		names := make([]string, 0, len(sr.scope))
		for n := range sr.scope {
			names = append(names, n)
		}

		sort.Strings(names)
		for _, sn := range names {
			rsj := sr.scope[sn]
			for scope, sRoles := range rsj.RolesByScope {
				for _, sr := range sRoles {
					roles = append(roles, fmt.Sprintf("%s/%s", sr.Source.Kind, sr.Source.Name))
					if len(sRoles)-1 == 0 {
						sort.Sort(byLen(roles))
					}
					results[filter][sn] = response{Subject: sn, Scope: scope, Role: roles}
				}
			}
		}

	}

	// ====================================================================
	//
	//	TODO: Wait here to collect the results
	//
	// ====================================================================

	if sc.Format == yamlFormat {
		data, err := yaml.Marshal(results)
		if err != nil {
			log.Printf("unable to load bindings: %v\n", err)
			return c.Render(http.StatusInternalServerError, r.JSON(map[string]string{"message": "server error"}))
		}
		return c.Render(http.StatusOK, r.Func(yamlContentType, func(w io.Writer, d render.Data) error {
			_, err := w.Write(data)
			return err
		}))
	}

	return c.Render(http.StatusOK, r.JSON(results))
}

// ==========================
//
//			Search
//
// ==========================

type search struct {
	k8si   kubernetes.Interface
	kind   string
	search string
	scope  map[string]topic
}

func (ss *search) roles() error {
	if err := ss.fromRole(); err != nil {
		return err
	}
	if err := ss.fromClusterRole(); err != nil {
		return err
	}
	return nil
}

func (ss *search) fromRole() error {
	rb, err := ss.k8si.RbacV1().RoleBindings("").List(metav1.ListOptions{})
	if err != nil {
		log.Println("error loading role bindings")
		return err
	}
	for _, r := range rb.Items {
		for _, subj := range r.Subjects {
			ok, err := ss.name(subj.Name)
			if err != nil {
				return err
			}
			if ok && ss.kinds(subj.Kind) {
				if rs, e := ss.scope[subj.Name]; e {
					rs.addRole(&r)
				} else {
					s := topic{
						Kind:         subj.Kind,
						RolesByScope: make(map[string][]sRole),
					}
					s.addRole(&r)
					k := subj.Name
					if subj.Kind == svcAccountKind {
						k = fmt.Sprintf("%s:%s", subj.Namespace, subj.Name)
					}
					ss.scope[k] = s
				}
			}
		}
	}
	return nil
}

var regexTokens = []string{
	"^",
	"*",
	".",
	"$",
	"|",
	"\\",
	"?",
	"+",
	"?",
	"{",
	"[",
	"(",
}

func (ss *search) canDoRegex() bool {
	for _, t := range regexTokens {
		if strings.Contains(ss.search, t) {
			return true
		}
	}
	return false
}

func (ss *search) fromClusterRole() error {
	rb, err := ss.k8si.RbacV1().ClusterRoleBindings().List(metav1.ListOptions{})
	if err != nil {
		log.Println("error loading cluster role bindings")
		return err
	}
	for _, r := range rb.Items {
		for _, subj := range r.Subjects {
			ok, err := ss.name(subj.Name)
			if err != nil {
				return err
			}
			if ok && ss.kinds(subj.Kind) {
				if s, exist := ss.scope[subj.Name]; exist {
					s.addClusterRole(&r)
				} else {
					s := topic{Kind: subj.Kind, RolesByScope: make(map[string][]sRole),}
					s.addClusterRole(&r)
					k := subj.Name
					if s.Kind == svcAccountKind {
						k = fmt.Sprintf("%s:%s", subj.Namespace, subj.Name)
					}
					ss.scope[k] = s
				}
			}
		}
	}
	return nil
}

func (ss *search) name(n string) (ok bool, err error) {
	if ss.search == "" {
		return true, nil
	}
	if ss.canDoRegex() {
		log.Print()
		reg, err := regexp.Compile(ss.search)
		if err != nil {
			return false, fmt.Errorf("invalid regex token provided")
		}
		return reg.MatchString(n), nil
	}
	return strings.Contains(n, ss.search), nil
}

func (ss *search) kinds(k string) bool {
	if ss.kind == "" {
		return true
	}
	return ss.kind == strings.ToLower(k)
}

// ==========================
//
//		  Role Topic
//
// ==========================

type topic struct {
	Kind         string
	RolesByScope map[string][]sRole
}

type sRole struct {
	Kind   string
	Name   string
	Source rSource
}

type rSource struct {
	Kind string
	Name string
}

func (t *topic) addRole(rb *rbacv1.RoleBinding) {
	s := sRole{
		Name:   rb.RoleRef.Name,
		Source: rSource{Name: rb.Name, Kind: roleKind},
	}
	s.Kind = rb.RoleRef.Kind
	t.RolesByScope[rb.Namespace] = append(t.RolesByScope[rb.Namespace], s)
}

func (t *topic) addClusterRole(rb *rbacv1.ClusterRoleBinding) {
	s := sRole{
		Name:   rb.RoleRef.Name,
		Source: rSource{Name: rb.Name, Kind: clusterKind},
	}
	s.Kind = rb.RoleRef.Kind
	t.RolesByScope[clusterWideScope] = append(t.RolesByScope[clusterWideScope], s)
}

// ==========================
//
//			Sorter
//
// ==========================

type byLen []string

func (s byLen) Len() int {
	return len(s)
}
func (s byLen) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}
func (s byLen) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
