module github.com/onosproject/onos-cli

go 1.14

require (
	github.com/onosproject/onos-api/go v0.7.0
	github.com/onosproject/onos-config v0.7.0
	github.com/onosproject/onos-lib-go v0.7.0
	github.com/onosproject/onos-ric v0.6.20 // indirect
	github.com/onosproject/onos-ric-sdk-go v0.7.0
	github.com/onosproject/onos-topo v0.7.0
	github.com/onosproject/onos-ztp v0.6.0
	github.com/spf13/cobra v1.1.1
	google.golang.org/grpc v1.33.2
	gotest.tools v2.2.0+incompatible
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
