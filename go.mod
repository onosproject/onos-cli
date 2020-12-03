module github.com/onosproject/onos-cli

go 1.14

require (
	github.com/onosproject/onos-cli/pkg/sdrancli v0.0.0
	github.com/onosproject/onos-config v0.6.18
	github.com/onosproject/onos-lib-go v0.6.25
	github.com/onosproject/onos-ric v0.6.20 // indirect
	github.com/onosproject/onos-topo v0.6.21
	github.com/onosproject/onos-ztp v0.6.0
	github.com/spf13/cobra v0.0.6
)

replace github.com/onosproject/onos-cli/pkg/sdrancli => ./pkg/sdrancli

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
