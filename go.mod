module github.com/onosproject/onos-cli

go 1.14

require (
	github.com/onosproject/onos-api/go v0.7.0
	github.com/onosproject/onos-config v0.7.0
	github.com/onosproject/onos-lib-go v0.7.0
	github.com/onosproject/onos-ric-sdk-go v0.7.0
	github.com/onosproject/onos-topo v0.7.0
	github.com/onosproject/onos-ztp v0.6.0
	github.com/openconfig/gnmi v0.0.0-20200617225440-d2b4e6a45802
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
	gotest.tools v2.2.0+incompatible
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
