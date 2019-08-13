module github.com/onosproject/onos-cli

go 1.12

require (
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/onosproject/onos-config v0.0.0-20190813184213-3f8105f6a4a2
	github.com/onosproject/onos-topo v0.0.0-20190809180402-2931d9c31bf5
	github.com/onosproject/onos-ztp v0.0.0-20190813150048-e3acd4e2b902
	github.com/spf13/cobra v0.0.5
)

replace github.com/onosproject/onos-config => ../onos-config
