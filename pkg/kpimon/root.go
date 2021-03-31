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

package kpimon

import (
	"github.com/onosproject/onos-lib-go/pkg/cli"
	loglib "github.com/onosproject/onos-lib-go/pkg/logging/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "kpimon"
	defaultAddress = "onos-kpimon:5150"
	v1Address      = "onos-kpimon-v1:5150"
	v2Address      = "onos-kpimon-v2:5150"
)

// init initializes the command line
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
		Use:   "kpimon {get/set} [args]",
		Short: "ONOS KPIMON subsystem commands",
	}

	cli.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(getListCommand())
	cmd.AddCommand(getSetCommand())
	cmd.AddCommand(loglib.GetCommand())
	return cmd
}

// GetCommandV1 returns the root command for the RAN service
func GetCommandV1() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kpimonv1 {get/set} [args]",
		Short: "ONOS KPIMON subsystem commands - for KPM V1",
	}

	cli.AddConfigFlags(cmd, v1Address)
	cmd.AddCommand(getListCommand())
	cmd.AddCommand(getSetCommand())
	cmd.AddCommand(loglib.GetCommand())
	return cmd
}

// GetCommandV2 returns the root command for the RAN service
func GetCommandV2() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kpimonv2 {get/set} [args]",
		Short: "ONOS KPIMON subsystem commands - for KPM V2",
	}

	cli.AddConfigFlags(cmd, v2Address)
	cmd.AddCommand(getListCommand())
	cmd.AddCommand(getSetCommand())
	cmd.AddCommand(loglib.GetCommand())
	return cmd
}