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

package ransim

import (
	"bytes"
	"strings"
	"testing"

	"gotest.tools/assert"
)

// Test_RootUsage tests the creation of the root command and checks that the ONOS usage messages are included
func Test_RootUsage(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")

	testCases := []struct {
		description string
		expected    string
	}{
		{description: "Change Num UEs", expected: `Change the number of UEs in the RAN simulation`},
		{description: "Reset Metrics", expected: `Reset the metrics counters`},
		{description: "Usage header", expected: `Usage:`},
		{description: "Usage config command", expected: `ransim [command]`},
	}

	cmd := GetCommand()
	assert.Assert(t, cmd != nil)
	cmd.SetOut(outputBuffer)

	usageErr := cmd.Usage()
	assert.NilError(t, usageErr)

	output := outputBuffer.String()

	for _, testCase := range testCases {
		assert.Assert(t, strings.Contains(output, testCase.expected), `Expected output "%s"" for %s not found`,
			testCase.expected, testCase.description)
	}
}

// Test_SubCommands tests that the ONOS supplied sub commands are present in the root
func Test_SubCommands(t *testing.T) {
	cmds := GetCommand().Commands()
	assert.Assert(t, cmds != nil)

	testCases := []struct {
		commandName   string
		expectedShort string
	}{
		{commandName: "config", expectedShort: "Manage the CLI configuration"},
		{commandName: "setnumues", expectedShort: "Change the number of UEs in the RAN simulation"},
		{commandName: "resetmetrics", expectedShort: "Reset the metrics counters"},
		{commandName: "log", expectedShort: "logging api commands"},
	}

	var subCommandsFound = make(map[string]bool)
	for _, cmd := range cmds {
		subCommandsFound[cmd.Short] = false
	}

	// Each sub command should be found once and only once
	assert.Equal(t, len(subCommandsFound), len(testCases))
	for _, testCase := range testCases {
		// check that this is an expected sub command
		entry, entryFound := subCommandsFound[testCase.expectedShort]
		assert.Assert(t, entryFound, "Subcommand %s not found", testCase.commandName)
		assert.Assert(t, entry == false, "command %s found more than once", testCase.commandName)
		subCommandsFound[testCase.expectedShort] = true
	}

	// Each sub command should have been found
	for _, testCase := range testCases {
		// check that this is an expected sub command
		entry, entryFound := subCommandsFound[testCase.expectedShort]
		assert.Assert(t, entryFound)
		assert.Assert(t, entry, "command %s was not found", testCase.commandName)
	}
}
