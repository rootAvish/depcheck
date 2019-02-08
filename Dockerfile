FROM openshift/origin-cli:latest

FROM registry.svc.ci.openshift.org/openshift/release:golang-1.11 AS builder
WORKDIR /go/src/github.com/mfojtik/depcheck
COPY . .
ENV GO_PACKAGE github.com/mfojtik/depcheck
RUN go build .

FROM docker.io/openshift/origin-cli:latest
COPY --from=builder /go/src/github.com/mfojtik/depcheck/depcheck /usr/bin/
