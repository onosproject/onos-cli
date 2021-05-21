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

package uenib

import (
	clilib "github.com/onosproject/onos-lib-go/pkg/cli"
	loglib "github.com/onosproject/onos-lib-go/pkg/logging/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "uenib"
	defaultAddress = "onos-uenib:5150"
)

// init initializes the command line
func init() {
	clilib.InitConfig(configName)
}

// Init is a hook called after cobra initialization
func Init() {
	// noop for now
}

// GetCommand returns the root command for the topo service
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uenib {create,get,set,delete,watch} [args]",
		Short: "ONOS UE-NIB subsystem commands",
	}

	clilib.AddConfigFlags(cmd, defaultAddress)

	cmd.AddCommand(clilib.GetConfigCommand())
	cmd.AddCommand(getGetCommand())
	cmd.AddCommand(getCreateCommand())
	cmd.AddCommand(getUpdateCommand())
	cmd.AddCommand(getDeleteCommand())
	cmd.AddCommand(getWatchCommand())
	cmd.AddCommand(loglib.GetCommand())
	return cmd
}
