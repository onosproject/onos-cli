module github.com/onosproject/onos-cli

go 1.14

require (
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-playground/overalls v0.0.0-20191218162659-7df9f728c018 // indirect
	github.com/mattn/goveralls v0.0.7 // indirect
	github.com/onosproject/onos-api/go v0.7.0
	github.com/onosproject/onos-lib-go v0.7.0
	github.com/onosproject/onos-ric-sdk-go v0.7.1
	github.com/openconfig/gnmi v0.0.0-20200617225440-d2b4e6a45802
	github.com/spf13/cobra v1.1.1
	github.com/yookoala/realpath v1.0.0 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	google.golang.org/grpc v1.33.2
	gotest.tools v2.2.0+incompatible
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
