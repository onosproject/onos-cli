export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_CLI_VERSION := latest
ONOS_BUILD_VERSION := stable

build: # @HELP build the Go binaries and run all validations (default)
build:
	go build -o build/_output/onos ./cmd/onos

test: # @HELP run the unit tests and source code validation
test: build deps license_check linters
	go test github.com/onosproject/onos-cli/pkg/...
	go test github.com/onosproject/onos-cli/cmd/...

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	./build/bin/coveralls-coverage

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	./build/licensing/boilerplate.py -v

onos-cli-docker: # @HELP build onos CLI Docker image
	@go mod vendor
	docker build . -f build/onos/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		-t onosproject/onos-cli:${ONOS_CLI_VERSION}
	@rm -rf vendor

images: # @HELP build all Docker images
images: build onos-cli-docker

all: build images

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
