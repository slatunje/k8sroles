# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.15.4 as builder

RUN mkdir -p $GOPATH/src/github.com/slatunje
WORKDIR $GOPATH/src/github.com/slatunje

ADD . .
ENV GO111MODULES=on
RUN go get ./...
RUN buffalo build --static -o /bin/k8sroles

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/k8sroles .

ENV GO_ENV=production

# outside aaccess
ENV ADDR=0.0.0.0

EXPOSE 3000

CMD exec /bin/k8sroles
