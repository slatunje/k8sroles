
dep:
	kubectl create -f https://k8s.io/examples/admin/namespace-dev.json
	kubectl create -f https://k8s.io/examples/admin/namespace-prod.json
	kubectl apply -f ./deploy/k8s/deploy.yaml
	kubectl apply -f ./deploy/k8s/service.yaml

build: build.binary build.docker

build.binary:
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	buffalo build  \
	-o bin/k8sroles \
	&& chmod -R 0777 bin

build.docker:
	docker build -t k8sroles .

start:
	GO_ENV=production ./bin/k8sroles

