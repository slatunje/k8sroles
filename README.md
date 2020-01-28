# K8s Roles Search


    kubectl create -f https://k8s.io/examples/admin/namespace-dev.json
    kubectl create -f https://k8s.io/examples/admin/namespace-prod.json
    kubectl apply -f delploy/k8s/resource.yaml
    kubectl create -f https://raw.githubusercontent.com/kubernetes/kops/master/addons/kubernetes-dashboard/v1.6.3.yaml


 ## Example

    curl -X "POST" "http://127.0.0.1:3000/" \
         -H 'Content-Type: application/json' \
         -d $'{
      "subject": "system:kube-.*,manager",
      "format": "json"
    }'

    ## Request
    curl -X "POST" "http://127.0.0.1:3000/" \
         -H 'Content-Type: application/json' \
         -d $'{
      "subject": "system:kube-.*,manager,jane",
      "format": "json"
    }'

## Deploy

    kubectl apply -f deploy/k8s/deploy.yaml

    kubectl apply -f deploy/k8s/service.yaml


