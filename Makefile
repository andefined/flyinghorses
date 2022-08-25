REGISTRY=reg.cvcio.org
PROJECT=flyinghorses
TAG:=$(shell git rev-parse HEAD)
BRANCH:=$(shell git rev-parse --abbrev-ref HEAD)
BUF_VERSION:=1.7.0

init: keys tools buf-install buf-get-proto

keys:
	openssl req \
		-newkey rsa:2048 \
		-nodes -keyout server-key.key \
		-x509 -out server-cert.pem \
		-addext "subjectAltName = DNS:localhost"

tools:
	go get github.com/oxequa/realize
	go get github.com/golangci/golangci-lint

buf-install:
	curl -sSL \
    	"https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(shell uname -s)-$(shell uname -m)" \
    	-o "$(shell go env GOPATH)/bin/buf" && \
  	chmod +x "$(shell go env GOPATH)/bin/buf"

buf-generate:
	buf generate --template buf.gen.yaml .

buf-update:
	buf mod update

buf-lint:
	buf lint

run:
	realize start

test:
	go test -v ./...

lint:
	golangci-lint run -e vendor
	buf-lint

include $(INCLUDE_MAKEFILE)

.PHONY: release
release: custom 
