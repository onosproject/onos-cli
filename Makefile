export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build

ONOS_CLI_VERSION := latest

build: # @HELP build the Go binaries and run all validations (default)
build:
	go build -o build/_output/onos ./cmd/onos
	go build -o build/_output/onos-cli-docs-gen ./cmd/onos-cli-docs-gen

test: # @HELP run the unit tests and source code validation
test: build deps license_check linters
	go test github.com/onosproject/onos-cli/pkg/...
	go test github.com/onosproject/onos-cli/cmd/...

jenkins-test:  # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: build-tools deps license_check linters
	TEST_PACKAGES=github.com/onosproject/onos-cli/... ./../build-tools/build/jenkins/make-unit

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	./../build-tools/build/coveralls/coveralls-coverage onos-cli

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: golang-ci # @HELP examines Go source code and reports coding problems
	golangci-lint run

build-tools: # @HELP install the ONOS build tools if needed
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi

jenkins-tools: # @HELP installs tooling needed for Jenkins
	cd .. && go get -u github.com/jstemmer/go-junit-report && go get github.com/t-yuki/gocover-cobertura

golang-ci: # @HELP install golang-ci if not present
	golangci-lint --version || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b `go env GOPATH`/bin v1.36.0

license_check: # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}

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

bumponosdeps: # @HELP update "onosproject" go dependencies and push patch to git.
	./../build-tools/bump-onos-deps ${VERSION}

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
