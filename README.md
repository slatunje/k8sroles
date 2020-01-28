# K8s Roles Search

## Deploy

    kubectl apply -f deploy/k8s/deploy.yaml
    kubectl apply -f deploy/k8s/service.yaml
    kubectl proxy --port=3000

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

## Tested with Deployments

    kubectl create -f https://k8s.io/examples/admin/namespace-dev.json
    kubectl create -f https://k8s.io/examples/admin/namespace-prod.json
    kubectl apply -f delploy/k8s/resource.yaml
    kubectl create -f https://raw.githubusercontent.com/kubernetes/kops/master/addons/kubernetes-dashboard/v1.6.3.yaml

## Build

    GOOS=linux GOARCH=amd64 buffalo build --tags="k8sroles"

    docker build -t slatunje/k8sroles .
    docker push slatunje/k8sroles:latest

    kubectl apply -f deploy/k8s/deploy.yaml







