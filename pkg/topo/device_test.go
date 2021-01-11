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

// Unit tests for device CLI
package topo

import (
	"bytes"
	"fmt"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"gotest.tools/assert"
	"strings"
	"testing"
)

const version = "1.0.0"
const deviceType = "TestDevice"

func generateDeviceData(count int) []*topoapi.Object {
	devices := make([]*topoapi.Object, count)

	protocol := new(topoapi.ProtocolState)
	protocol.Protocol = topoapi.Protocol_GNMI
	protocol.ConnectivityState = topoapi.ConnectivityState_REACHABLE
	protocol.ChannelState = topoapi.ChannelState_CONNECTED
	protocol.ServiceState = topoapi.ServiceState_AVAILABLE

	for devIdx := range devices {
		o := &topoapi.Object{
			ID:       topoapi.ID(fmt.Sprintf("test-device-%d", devIdx)),
			Revision: topoapi.Revision(devIdx),
			Type:     topoapi.Object_ENTITY,
			Obj: &topoapi.Object_Entity{
				Entity: &topoapi.Entity{
					KindID:    topoapi.ID(deviceType),
					Protocols: []*topoapi.ProtocolState{protocol},
				},
			},
		}
		o.GetEntity().KindID = topoapi.ID("legacy")
		o.GetEntity().GetProtocols()

		o.Attributes = make(map[string]string)
		setAttribute(o, topoapi.Type, deviceType)
		setAttribute(o, topoapi.Role, "leaf")
		setAttribute(o, topoapi.Address, fmt.Sprintf("192.168.0.%d", devIdx))
		setAttribute(o, topoapi.Target, "")
		setAttribute(o, topoapi.Version, topoapi.Version)
		setAttribute(o, topoapi.Timeout, "1")

		devices[devIdx] = o
	}
	return devices
}

func Test_GetDevice(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")
	cli.CaptureOutput(outputBuffer)

	setUpMockClients()
	getDevices := getGetDeviceCommand()
	args := make([]string, 1)
	args[0] = "test-device-1"
	getDevices.SetArgs(args)
	err := getDevices.Execute()
	assert.NilError(t, err)
	output := outputBuffer.String()
	assert.Assert(t, strings.Contains(output, "test-device"))
}

func Test_AddDevice(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")
	cli.CaptureOutput(outputBuffer)

	setUpMockClients()
	addDevice := getAddDeviceCommand()
	args := make([]string, 7)
	args[0] = "test-device-1" // Name
	args[1] = fmt.Sprintf("--type=%s", deviceType)
	args[2] = fmt.Sprintf("--version=%s", version)
	args[3] = "--address=192.168.0.1"
	args[4] = "--timeout=1s"
	args[5] = "--user=test"
	args[6] = "--role=leaf"
	addDevice.SetArgs(args)
	err := addDevice.Execute()
	assert.NilError(t, err)
	output := outputBuffer.String()
	assert.Assert(t, strings.Contains(output, "Added device test-device-1"))
}

func Test_UpdateDevice(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")
	cli.CaptureOutput(outputBuffer)

	setUpMockClients()
	updateDevice := getUpdateDeviceCommand()
	args := make([]string, 7)
	args[0] = "test-device-1" // Name
	args[1] = fmt.Sprintf("--type=%s", deviceType)
	args[2] = fmt.Sprintf("--version=%s", version)
	args[3] = "--address=192.168.0.1"
	args[4] = "--timeout=1s"
	args[5] = "--user=test"
	args[6] = "--role=leaf"
	updateDevice.SetArgs(args)
	err := updateDevice.Execute()
	assert.NilError(t, err)
	output := outputBuffer.String()
	assert.Assert(t, strings.Contains(output, "Updated device test-device-1"))
}

func Test_RemoveDevice(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")
	cli.CaptureOutput(outputBuffer)

	setUpMockClients()
	removeDevice := getRemoveDeviceCommand()
	args := make([]string, 1)
	args[0] = "test-device-1" // Name
	removeDevice.SetArgs(args)
	err := removeDevice.Execute()
	assert.NilError(t, err)
	output := outputBuffer.String()
	assert.Assert(t, strings.Contains(output, "Removed device test-device-1"))
}
