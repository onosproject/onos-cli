// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

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
