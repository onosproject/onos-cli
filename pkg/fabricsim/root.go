// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package fabricsim

import (
	"github.com/onosproject/onos-lib-go/pkg/cli"
	loglib "github.com/onosproject/onos-lib-go/pkg/logging/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "fabricsim"
	defaultAddress = "fabric-sim:5150"
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
		Use:   "fabricsim {get,start,stop,enable,disable} [args]",
		Short: "ONOS fabric simulator commands",
	}

	cli.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(cli.GetConfigCommand())

	cmd.AddCommand(getGetCommand())

	//cmd.AddCommand(startDeviceCommand())
	//cmd.AddCommand(stopDeviceCommand())

	//cmd.AddCommand(enablePortCommand())
	//cmd.AddCommand(disablePortCommand())

	cmd.AddCommand(loglib.GetCommand())
	return cmd
}

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {device(s),link(s),host(s)} [args]",
		Short: "Commands for retrieving simulated entities related information",
	}

	cmd.AddCommand(getDevicesCommand())
	//cmd.AddCommand(getDeviceCommand())
	//
	//cmd.AddCommand(getLinksCommand())
	//cmd.AddCommand(getLinkCommand())
	//
	//cmd.AddCommand(getHostsCommand())
	//cmd.AddCommand(getHostCommand())
	return cmd
}
