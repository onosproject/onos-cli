// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package modelregistry

import (
	"github.com/onosproject/onos-config-model/api/onos/configmodel"
	"google.golang.org/grpc"
)

// ConfigModelRegistryServiceClientFactory : Default ConfigAdminClient creation.
// TODO move this to onos-api
var ConfigModelRegistryServiceClientFactory = func(cc *grpc.ClientConn) configmodel.ConfigModelRegistryServiceClient {
	return configmodel.NewConfigModelRegistryServiceClient(cc)
}

// CreateConfigModelRegistryServiceClient creates and returns a new config admin client
func CreateConfigModelRegistryServiceClient(cc *grpc.ClientConn) configmodel.ConfigModelRegistryServiceClient {
	return ConfigModelRegistryServiceClientFactory(cc)
}
