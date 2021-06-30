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

package pci

import (
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "pci"
	defaultAddress = "onos-pci:5150"
)

func init() {
	cli.InitConfig(configName)
}

// Init is a hook called after cobra initialization
func Init() {
	// noop for now
}

// GetCommand returns the root command for the RAN service
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pci {get} [args]",
		Short: "ONOS PCI subsystem commands",
	}

	cli.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(getGetCommand())
	return cmd
}
