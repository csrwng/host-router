SRC_DIRS = cmd pkg

.PHONY: default
default: all

.PHONY: all
all: verify build

.PHONY: build
build:  
	go build -mod=vendor -o bin/host-router github.com/csrwng/host-router/cmd/host-router

.PHONY: verify-gofmt
verify-gofmt:
	@echo Verifying gofmt
	@gofmt -l -s $(SRC_DIRS)>.out 2>&1 || true
	@[ ! -s .out ] || \
	  (echo && echo "*** Please run 'make fmt' in order to fix the following:" && \
	  cat .out && echo && rm .out && false)
	@rm .out

.PHONY: fmt
fmt:
	gofmt -w -s $(SRC_DIRS)

.PHONY: verify
verify: verify-gofmt
