export GO111MODULE = on
export COMPOSE_HTTP_TIMEOUT ?= 480



.PHONY: lint

GOPATH := $(shell go env GOPATH)
LINTERS := \
	github.com/golang/lint/golint \
	github.com/kisielk/errcheck \
	honnef.co/go/tools/cmd/staticcheck \
	honnef.co/go/tools/cmd/unused

PACKAGES = $(shell go list ./... | grep -v /vendor/)

lint:
	env GO111MODULE=off go fmt ./...
	env GO111MODULE=on go vet -mod=vendor ./...
	gofmt -w ./

install-build-deps:
	go install -v $(LINTERS)

test: lint
	mkdir -p builds
	env GO111MODULE=on go test -mod=vendor -race -coverprofile=${UNIT_COVERAGE_OUTPUT} ./...

modtidy:
	go mod tidy

modvendor:
	go mod vendor

modall: modtidy modvendor