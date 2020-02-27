module github.com/onosproject/onos-cli

go 1.13

require (
	github.com/onosproject/onos-config v0.0.0-20200227222707-4f6d09fd6c96
	github.com/onosproject/onos-ric v0.0.0-20200225182040-dcf370614b8e
	github.com/onosproject/onos-topo v0.0.0-20200227183114-fc86211b9a1d
	github.com/onosproject/onos-ztp v0.0.0-20200218172126-3375e0509f99
	github.com/spf13/cobra v0.0.6
)

replace github.com/onosproject/onos-topo => ../onos-topo

replace github.com/onosproject/onos-lib-go => ../onos-lib-go
