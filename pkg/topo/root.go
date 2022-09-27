// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	clilib "github.com/onosproject/onos-lib-go/pkg/cli"
	loglib "github.com/onosproject/onos-lib-go/pkg/logging/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "topo"
	defaultAddress = "onos-topo:5150"
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
		Use:   "topo {create,get,set,delete,watch,import,export} [args]",
		Short: "ONOS topology resource commands",
	}

	clilib.AddConfigFlags(cmd, defaultAddress)

	cmd.AddCommand(clilib.GetConfigCommand())
	cmd.AddCommand(getGetCommand())
	cmd.AddCommand(getAddCommand())
	cmd.AddCommand(getUpdateCommand())
	cmd.AddCommand(getDeleteCommand())
	cmd.AddCommand(getWipeoutCommand())
	cmd.AddCommand(getWatchCommand())
	cmd.AddCommand(getImportCommand())
	cmd.AddCommand(getExportCommand())
	cmd.AddCommand(loglib.GetCommand())
	return cmd
}
