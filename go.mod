module github.com/onosproject/onos-cli

go 1.14

require (
	github.com/onosproject/onos-cli/pkg/sdrancli v0.0.0
	github.com/onosproject/onos-config v0.6.11
	github.com/onosproject/onos-lib-go v0.6.15
	github.com/onosproject/onos-ric v0.6.14 // indirect
	github.com/onosproject/onos-topo v0.6.15
	github.com/onosproject/onos-ztp v0.6.0
	github.com/onosproject/ran-simulator v0.6.6 // indirect
	github.com/spf13/cobra v0.0.6
)

replace github.com/onosproject/onos-cli/pkg/sdrancli => ./pkg/sdrancli
