# SPDX-License-Identifier: Apache-2.0
# Copyright 2019 Open Networking Foundation
# Copyright 2024 Intel Corporation

export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build license

ONOS_CLI_VERSION ?= latest

GOLANG_CI_VERSION := v1.52.2

all: build docker-build

build: # @HELP build the Go binaries and run all validations (default)
build:
	go build -o build/_output/onos ./cmd/onos
	go build -o build/_output/onos-cli-docs-gen ./cmd/onos-cli-docs-gen
	go build -o build/_output/gnmi_cli ./cmd/gnmi_cli

test: # @HELP run the unit tests and source code validation
test: build lint license
	go test github.com/onosproject/onos-cli/pkg/...
	go test github.com/onosproject/onos-cli/cmd/...

docs: # @HELP generate CLI docs
docs:
	go run cmd/onos-cli-docs-gen/main.go

docker-build-onos-cli: # @HELP build onos CLI Docker image
	@go mod vendor
	docker build . -f build/onos/Dockerfile \
		-t onosproject/onos-cli:${ONOS_CLI_VERSION}
	@rm -rf vendor

docker-build: # @HELP build all Docker images
docker-build: build docker-build-onos-cli

docker-push-onos-cli: # @HELP push onos-cli Docker image
	docker push onosproject/onos-cli:${ONOS_CLI_VERSION}

docker-push: # @HELP push docker images
docker-push: docker-push-onos-cli

lint: # @HELP examines Go source code and reports coding problems
	golangci-lint --version | grep $(GOLANG_CI_VERSION) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b `go env GOPATH`/bin $(GOLANG_CI_VERSION)
	golangci-lint run --timeout 15m

license: # @HELP run license checks
	rm -rf venv
	python3 -m venv venv
	. ./venv/bin/activate;\
	python3 -m pip install --upgrade pip;\
	python3 -m pip install reuse;\
	reuse lint

check-version: # @HELP check version is duplicated
	./build/bin/version_check.sh all

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onos/onos ./cmd/dummy/dummy

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
