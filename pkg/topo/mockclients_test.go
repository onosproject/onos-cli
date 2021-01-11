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

// Client Mocks
package topo

import (
	"context"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"google.golang.org/grpc"
)

type mockTopoService struct {
	test string
}

func (m *mockTopoService) Create(ctx context.Context, request *topoapi.CreateRequest, opts ...grpc.CallOption) (*topoapi.CreateResponse, error) {
	addedDevice := request.Object
	return &topoapi.CreateResponse{Object: addedDevice}, nil // Just reflect it
}

func (m *mockTopoService) Update(ctx context.Context, request *topoapi.UpdateRequest, opts ...grpc.CallOption) (*topoapi.UpdateResponse, error) {
	updatedDevice := request.Object
	return &topoapi.UpdateResponse{Object: updatedDevice}, nil
}

func (m *mockTopoService) Get(ctx context.Context, request *topoapi.GetRequest, opts ...grpc.CallOption) (*topoapi.GetResponse, error) {
	return &topoapi.GetResponse{Object: generateDeviceData(1)[0]}, nil
}

func (m *mockTopoService) Delete(ctx context.Context, request *topoapi.DeleteRequest, opts ...grpc.CallOption) (*topoapi.DeleteResponse, error) {
	return &topoapi.DeleteResponse{}, nil
}

func (m *mockTopoService) List(ctx context.Context, in *topoapi.ListRequest, opts ...grpc.CallOption) (*topoapi.ListResponse, error) {
	return nil, nil
}

func (m *mockTopoService) Watch(ctx context.Context, in *topoapi.WatchRequest, opts ...grpc.CallOption) (topoapi.Topo_WatchClient, error) {
	return nil, nil
}

// setUpMockClients sets up factories to create mocks of top level clients used by the CLI
func setUpMockClients() {
	topoapi.TopoClientFactory = func(cc *grpc.ClientConn) topoapi.TopoClient {
		return &mockTopoService{test: ""}
	}
}
