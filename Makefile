export CGO_ENABLED:=0

GITTAG=$(shell git rev-parse HEAD)
VERSION=v0.1

GOPATH_BIN:=$(shell echo ${GOPATH} | awk 'BEGIN { FS = ":" }; { print $1 }')/bin
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${GITTAG}"

.PHONY: all
all: build

.PHONY: build
build:
	@go build ${LDFLAGS} -o bin/kube-tags2iaas -v github.com/alekssaul/kube-tags2iaas

.PHONY: vendor
vendor:
	@go mod vendor
	
.PHONY: docker-build
docker-build:
	@docker build . -t alekssaul/kube-tags2iaas:dev
	@docker push alekssaul/kube-tags2iaas:dev

.PHONY: test
test:
	@go test ./... 