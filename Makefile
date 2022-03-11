# SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build

ONOS_CLI_VERSION := latest

build: # @HELP build the Go binaries and run all validations (default)
build:
	go build -o build/_output/onos ./cmd/onos
	go build -o build/_output/onos-cli-docs-gen ./cmd/onos-cli-docs-gen
	go build -o build/_output/gnmi_cli ./cmd/gnmi_cli

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

mod-update: # @HELP Download the dependencies to the vendor folder
	go mod tidy
	go mod vendor
mod-lint: mod-update # @HELP ensure that the required dependencies are in place
	# dependencies are vendored, but not committed, go.sum is the only thing we need to check
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

test: # @HELP run the unit tests and source code validation
test: mod-lint build license_check_apache linters license
	go test github.com/onosproject/onos-cli/pkg/...
	go test github.com/onosproject/onos-cli/cmd/...

jenkins-test:  # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: mod-lint build license_check_apache linters license
	TEST_PACKAGES=github.com/onosproject/onos-cli/... ./../build-tools/build/jenkins/make-unit

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	./../build-tools/build/coveralls/coveralls-coverage onos-cli

onos-cli-docker: # @HELP build onos CLI Docker image
onos-cli-docker:
	@go mod vendor
	docker build . -f build/onos/Dockerfile \
		-t onosproject/onos-cli:${ONOS_CLI_VERSION}
	@rm -rf vendor

images: # @HELP build all Docker images
images: build onos-cli-docker

kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-cli:${ONOS_CLI_VERSION}

all: build images

publish: # @HELP publish version on github and dockerhub
	./../build-tools/publish-version ${VERSION} onosproject/onos-cli

jenkins-publish: build-tools jenkins-tools # @HELP Jenkins calls this to publish artifacts
	./build/bin/push-images
	../build-tools/release-merge-commit
	../build-tools/build/docs/push-docs

clean:: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor
