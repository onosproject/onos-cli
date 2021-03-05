// Copyright 2021-present Open Networking Foundation.
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
	"context"
	"github.com/onosproject/onos-config-model/api/onos/configmodel"
	"google.golang.org/grpc"
)

type mockConfigModelRegistryServiceClient struct {
}

var LastCreatedClient *mockConfigModelRegistryServiceClient

func (c *mockConfigModelRegistryServiceClient) GetModel(ctx context.Context, in *configmodel.GetModelRequest, opts ...grpc.CallOption) (*configmodel.GetModelResponse, error) {
	return nil, nil
}

func (c *mockConfigModelRegistryServiceClient) ListModels(ctx context.Context, in *configmodel.ListModelsRequest, opts ...grpc.CallOption) (*configmodel.ListModelsResponse, error) {
	return &configmodel.ListModelsResponse{
		Models: testModels,
	}, nil
}

func (c *mockConfigModelRegistryServiceClient) PushModel(ctx context.Context, in *configmodel.PushModelRequest, opts ...grpc.CallOption) (*configmodel.PushModelResponse, error) {
	return nil, nil
}

func (c *mockConfigModelRegistryServiceClient) DeleteModel(ctx context.Context, in *configmodel.DeleteModelRequest, opts ...grpc.CallOption) (*configmodel.DeleteModelResponse, error) {
	return nil, nil
}

func setupMockClients() {
	ConfigModelRegistryServiceClientFactory = func(cc *grpc.ClientConn) configmodel.ConfigModelRegistryServiceClient {
		LastCreatedClient = &mockConfigModelRegistryServiceClient{}
		return LastCreatedClient
	}

}
