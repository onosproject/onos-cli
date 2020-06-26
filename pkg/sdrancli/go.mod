module github.com/onosproject/onos-cli/pkg/sdrancli

go 1.14

require (
	github.com/onosproject/onos-cli v0.6.6
	github.com/onosproject/onos-config v0.6.5
	github.com/onosproject/onos-lib-go v0.6.6
	github.com/onosproject/onos-ric v0.6.8
	github.com/onosproject/onos-topo v0.6.9
	github.com/onosproject/onos-ztp v0.6.0
	github.com/onosproject/ran-simulator v0.6.5
	github.com/spf13/cobra v0.0.6
)

replace github.com/onosproject/onos-cli/ => ../cli
