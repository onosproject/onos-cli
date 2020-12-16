module github.com/onosproject/onos-cli/pkg/sdrancli

go 1.14

require (
    github.com/onosproject/onos-cli v0.6.6
    github.com/onosproject/onos-config v0.7.0
    github.com/onosproject/onos-e2t v0.7.0
    github.com/onosproject/onos-e2sub v0.7.0
    github.com/onosproject/onos-ric-sdk-go v0.7.0
    github.com/onosproject/onos-lib-go v0.7.0
    github.com/onosproject/onos-kpimon v0.1.3
    github.com/onosproject/onos-topo v0.7.0
    github.com/onosproject/onos-ztp v0.6.0
    github.com/onosproject/ran-simulator v0.7.0
    github.com/spf13/cobra v0.0.6
)

replace github.com/onosproject/onos-cli/ => ../cli
