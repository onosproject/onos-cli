// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package discovery

import (
	"github.com/onosproject/onos-lib-go/pkg/cli"
	loglib "github.com/onosproject/onos-lib-go/pkg/logging/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "discovery"
	defaultAddress = "topo-discovery:5150"
)

// init initializes the command line
func init() {
	cli.InitConfig(configName)
}

// Init is a hook called after cobra initialization
func Init() {
	// noop for now
}

// GetCommand returns the root command for the device provisioner service
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "discovery {add} [args]",
		Aliases: []string{"topo-discovery"},
		Short:   "Topology discovery subsystem commands",
	}

	cli.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(cli.GetConfigCommand())

	cmd.AddCommand(getAddCommand())
	cmd.AddCommand(loglib.GetCommand())
	return cmd
}
