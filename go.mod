module github.com/onosproject/onos-cli

go 1.12

require (
	github.com/onosproject/onos-topo v0.0.0-20190806004156-537a9862c203
	github.com/spf13/cobra v0.0.5
)

replace github.com/onosproject/onos-topo => ../onos-topo
