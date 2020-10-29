module github.com/onosproject/onos-cli

go 1.14

require (
	github.com/onosproject/onos-cli/pkg/sdrancli v0.0.0
	github.com/onosproject/onos-config v0.6.16
	github.com/onosproject/onos-lib-go v0.6.21
	github.com/onosproject/onos-ric v0.6.20 // indirect
	github.com/onosproject/onos-topo v0.6.19
	github.com/onosproject/onos-ztp v0.6.0
	github.com/onosproject/ran-simulator v0.6.6 // indirect
	github.com/spf13/cobra v0.0.6
	github.com/onosproject/onos-e2t v0.6.2
)

replace github.com/onosproject/onos-cli/pkg/sdrancli => ./pkg/sdrancli
