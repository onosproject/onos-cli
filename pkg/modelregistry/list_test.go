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
	"bytes"
	"fmt"
	"github.com/onosproject/onos-api/go/onos/configmodel"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"gotest.tools/assert"
	"strings"
	"testing"
	"time"
)

var testModels []*configmodel.ConfigModel

func generateModelData(count int) {
	models := make([]*configmodel.ConfigModel, count)
	now := time.Now()
	revision := now.Format("2006-01-02")
	for modelIndex := range models {
		model := new(configmodel.ConfigModel)
		model.Name = fmt.Sprintf("Model %d", modelIndex)
		model.Version = "1.0.0"
		model.Modules = make([]*configmodel.ConfigModule, count)
		for moduleIndex := range model.Modules {
			module := new(configmodel.ConfigModule)
			module.Name = fmt.Sprintf("Module %d", modelIndex)
			module.Revision = revision
			module.File = fmt.Sprintf("test-file-%d-%d.yang", modelIndex, moduleIndex)
			model.Modules[moduleIndex] = module
		}
		models[modelIndex] = model
	}

	testModels = models
}

func Test_ListPlugins(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")
	cli.CaptureOutput(outputBuffer)
	generateModelData(4)

	setupMockClients()
	listCmd := getListCommand()
	err := listCmd.RunE(listCmd, []string{})
	assert.NilError(t, err)
	output := outputBuffer.String()
	assert.Equal(t, strings.Count(output, "YANGS"), len(testModels))
	assert.Equal(t, strings.Count(output, "Module"), len(testModels)*len(testModels))

	t.Logf("Output\n%v", output)
}
