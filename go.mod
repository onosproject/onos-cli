module github.com/onosproject/onos-cli

go 1.16

require (
	github.com/aybabtme/uniplot v0.0.0-20151203143629-039c559e5e7e
	github.com/gogo/protobuf v1.3.2
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.2
	github.com/onosproject/onos-api/go v0.9.13
	github.com/onosproject/onos-lib-go v0.8.13
	github.com/onosproject/onos-ric-sdk-go v0.8.9
	github.com/openconfig/gnmi v0.0.0-20200617225440-d2b4e6a45802
	github.com/openconfig/ygot v0.12.4 // indirect
	github.com/prometheus/common v0.4.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20220128200615-198e4374d7ed
	google.golang.org/grpc v1.41.0
)

//replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
