export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_CLI_VERSION := latest
ONOS_BUILD_VERSION := v0.6.3

build: # @HELP build the Go binaries and run all validations (default)
build: 
	go build -o build/_output/onos ./cmd/onos
	go build -o build/_output/onos-cli-docs-gen ./cmd/onos-cli-docs-gen

build-sdran: # @HELP build the Go binaries and run all validations (default)
build-sdran:
	GOPRIVATE="github.com/onosproject/onos-ric" go build -o build/_output/sdran ./cmd/sdran

test: # @HELP run the unit tests and source code validation
test: build deps license_check linters
	go test github.com/onosproject/onos-cli/pkg/...
	go test github.com/onosproject/onos-cli/cmd/...

coverage: # @HELP generate unit test coverage data
coverage: build build-sdran deps linters license_check
	GOPRIVATE="github.com/onosproject/*" ./../build-tools/build/coveralls/coveralls-coverage onos-cli

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}

update-deps: # @HELP pull updated CLI dependencies
	go get github.com/onosproject/onos-topo
	go get github.com/onosproject/onos-config
	go get github.com/onosproject/onos-ztp

update-sdran-deps: # @HELP pull updated SDRAN CLI dependencies
	GOPRIVATE="github.com/onosproject/onos-ric" go get github.com/onosproject/onos-ric
	GOPRIVATE="github.com/onosproject/onos-ric" go get github.com/onosproject/ran-simulator

onos-cli-docker: # @HELP build onos CLI Docker image
onos-cli-docker: update-deps
	@go mod vendor
	docker build . -f build/onos/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		-t onosproject/onos-cli:${ONOS_CLI_VERSION}
	@rm -rf vendor

onos-sdran-cli-docker: # @HELP build onos CLI Docker image
onos-sdran-cli-docker: update-deps update-sdran-deps
	@go mod vendor
	docker build . -f build/sdran/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		-t onosproject/onos-sdran-cli:${ONOS_CLI_VERSION}
	@rm -rf vendor

images: # @HELP build all Docker images
images: build onos-cli-docker onos-sdran-cli-docker

kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-cli:${ONOS_CLI_VERSION}

all: build images

publish: # @HELP publish version on github and dockerhub
	./../build-tools/publish-version ${VERSION} onosproject/onos-cli onosproject/onos-sdran-cli

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
