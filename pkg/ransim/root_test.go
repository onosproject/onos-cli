// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ransim

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// Test_RootUsage tests the creation of the root command and checks that the ONOS usage messages are included
func Test_RootUsage(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")

	testCases := []struct {
		description string
		expected    string
	}{
		{description: "Usage header", expected: `Usage:`},
		{description: "Usage config command", expected: `ransim [command]`},
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
		{commandName: "log", expectedShort: "logging api commands"},
		{commandName: "get", expectedShort: "Commands for retrieving RAN simulator model and other information"},
		{commandName: "set", expectedShort: "Commands for setting RAN simulator model metrics and other information"},
		{commandName: "create", expectedShort: "Commands for creating simulated entities"},
		{commandName: "delete", expectedShort: "Commands for deleting simulated entities"},
		{commandName: "start", expectedShort: "Start E2 node agent"},
		{commandName: "stop", expectedShort: "Stop E2 node agent"},
		{commandName: "load", expectedShort: "Load model and/or metric data"},
		{commandName: "clear", expectedShort: "Clear the simulated nodes, cells and metrics"},
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
