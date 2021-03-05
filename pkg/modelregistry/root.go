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

// Package modelregistry holds ONOS command-line command implementations for onos-config
package modelregistry

import (
	clilib "github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "modelregistry"
	defaultAddress = "onos-config:5151"
)

// init initializes the command line
func init() {
	clilib.InitConfig(configName)
}

// Init is a hook called after cobra initialization
func Init() {
	// noop for now
}

// GetCommand returns the root command for the config service.
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modelregistry {list|get} [args]",
		Short: "ONOS Config Model Registry subsystem commands",
	}

	clilib.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(getListCommand())
	cmd.AddCommand(getGetCommand())
	return cmd
}
