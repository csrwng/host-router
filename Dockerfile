FROM openshift/origin-release:golang-1.13 as builder
RUN mkdir -p /go/src/github.com/csrwng/host-router
WORKDIR /go/src/github.com/csrwng/host-router
COPY . .
RUN go build -mod=vendor -o ./bin/host-router ./cmd/host-router/main.go

FROM registry.access.redhat.com/ubi7/ubi
COPY --from=builder /go/src/github.com/csrwng/host-router/bin/host-router /usr/bin/host-router
ENTRYPOINT /usr/bin/host-router
