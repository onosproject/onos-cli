// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package perf

import (
	clilib "github.com/onosproject/onos-lib-go/pkg/cli"
	loglib "github.com/onosproject/onos-lib-go/pkg/logging/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "perf"
	defaultAddress = "onos-perf:5150"
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
		Use:   "perf {ping,stream} [args]",
		Short: "simple gRPC performance measurement client",
	}

	clilib.AddConfigFlags(cmd, defaultAddress)

	cmd.AddCommand(clilib.GetConfigCommand())
	cmd.AddCommand(getPingCommand())
	cmd.AddCommand(getStreamCommand())
	cmd.AddCommand(loglib.GetCommand())
	return cmd
}
