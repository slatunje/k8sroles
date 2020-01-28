module k8sroles

go 1.13

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190620084959-7cf5895f2711 // curl -s https://proxy.golang.org/k8s.io/api/@v/kubernetes-1.15.0.info | jq -r .Version
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719 // curl -s https://proxy.golang.org/k8s.io/apimachinery/@v/kubernetes-1.15.0.info | jq -r .Version
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab // curl -s https://proxy.golang.org/k8s.io/client-go/@v/kubernetes-1.15.0.info | jq -r .Version
)

require (
	github.com/gobuffalo/buffalo v0.15.4
	github.com/gobuffalo/envy v1.8.1
	github.com/gobuffalo/mw-contenttype v0.0.0-20190224202710-36c73cc938f3
	github.com/gobuffalo/mw-forcessl v0.0.0-20190224202501-6d1ef7ffb276
	github.com/gobuffalo/mw-paramlogger v0.0.0-20190224201358-0d45762ab655
	github.com/gobuffalo/packr/v2 v2.7.1
	github.com/gobuffalo/suite v2.8.2+incompatible
	github.com/gobuffalo/x v0.1.0
	github.com/googleapis/gnostic v0.2.0
	github.com/imdario/mergo v0.3.8
	github.com/rs/cors v1.7.0
	github.com/unrolled/secure v1.0.7
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
	golang.org/x/net v0.0.0-20191126235420-ef20fe5d7933
	golang.org/x/sys v0.0.0-20191128015809-6d18c012aee9
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15
	gopkg.in/yaml.v2 v2.2.7
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/utils v0.0.0-20200124190032-861946025e34
)
