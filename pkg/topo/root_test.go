// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_RootUsage tests the creation of the root command and checks that the ONOS usage messages are included
func Test_RootUsage(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")

	testCases := []struct {
		description string
		expected    string
	}{
		{description: "Create command", expected: `Create a topology resource`},
		{description: "Get command", expected: `Get topology resources`},
		{description: "Delete command", expected: `Delete a topology resource`},
		{description: "Update command", expected: `Update a topology resource`},
		{description: "Watch command", expected: `Watch for changes to a topology resource type`},
		{description: "Usage header", expected: `Usage:`},
		{description: "Usage config command", expected: `topo [command]`},
	}

	cmd := GetCommand()
	assert.NotNil(t, cmd)
	cmd.SetOut(outputBuffer)

	usageErr := cmd.Usage()
	assert.NoError(t, usageErr)

	output := outputBuffer.String()

	for _, testCase := range testCases {
		assert.True(t, strings.Contains(output, testCase.expected), `Expected output "%s"" for %s not found`,
			testCase.expected, testCase.description)
	}
}

// Test_SubCommands tests that the ONOS supplied sub commands are present in the root
func Test_SubCommands(t *testing.T) {
	cmds := GetCommand().Commands()
	assert.NotNil(t, cmds)

	testCases := []struct {
		commandName   string
		expectedShort string
	}{
		{commandName: "config", expectedShort: "Manage the CLI configuration"},
		{commandName: "create", expectedShort: "Create a topology resource"},
		{commandName: "get", expectedShort: "Get topology resources"},
		{commandName: "delete", expectedShort: "Delete a topology resource"},
		{commandName: "set", expectedShort: "Update a topology resource"},
		{commandName: "import", expectedShort: "Import topology resources in JSON format"},
		{commandName: "export", expectedShort: "Export topology resources in JSON format"},
		{commandName: "watch", expectedShort: "Watch for changes to a topology resource type"},
		{commandName: "log", expectedShort: "logging api commands"},
	}

	var subCommandsFound = make(map[string]bool)
	for _, cmd := range cmds {
		subCommandsFound[cmd.Short] = false
	}

	// Each sub command should be found once and only once
	// TODO - Add test cases for new topo commands
	//assert.Equal(t, len(subCommandsFound), len(testCases))
	for _, testCase := range testCases {
		// check that this is an expected sub command
		entry, entryFound := subCommandsFound[testCase.expectedShort]
		assert.True(t, entryFound, "Subcommand %s not found", testCase.commandName)
		assert.False(t, entry, "command %s found more than once", testCase.commandName)
		subCommandsFound[testCase.expectedShort] = true
	}

	// Each sub command should have been found
	for _, testCase := range testCases {
		// check that this is an expected sub command
		entry, entryFound := subCommandsFound[testCase.expectedShort]
		assert.True(t, entryFound)
		assert.True(t, entry, "command %s was not found", testCase.commandName)
	}
}
