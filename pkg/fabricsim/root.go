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
		Use:     "fabric-sim {create,delete,get,start,stop,enable,disable,emit} [args]",
		Short:   "ONOS fabric simulator commands",
		Aliases: []string{"fabricsim", "fabsim", "fsim"},
	}

	cli.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(cli.GetConfigCommand())

	cmd.AddCommand(getCreateCommand())
	cmd.AddCommand(getDeleteCommand())
	cmd.AddCommand(getGetCommand())

	cmd.AddCommand(getStartCommand())
	cmd.AddCommand(getStopCommand())

	cmd.AddCommand(getEnableCommand())
	cmd.AddCommand(getDisableCommand())
	cmd.AddCommand(getEmitCommand())

	cmd.AddCommand(loglib.GetCommand())
	return cmd
}

func getCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create {device,link} [args]",
		Short: "Commands for creating simulated entities",
	}

	cmd.AddCommand(createDeviceCommand())
	cmd.AddCommand(createLinkCommand())
	cmd.AddCommand(createHostCommand())
	return cmd
}

func getDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete {device,link} [args]",
		Short: "Commands for deleting simulated entities",
	}

	cmd.AddCommand(deleteDeviceCommand())
	cmd.AddCommand(deleteLinkCommand())
	cmd.AddCommand(deleteHostCommand())
	return cmd
}

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {device(s),link(s),host(s), stats} [args]",
		Short: "Commands for retrieving simulated entities related information",
	}

	cmd.AddCommand(getDevicesCommand())
	cmd.AddCommand(getDeviceCommand())

	cmd.AddCommand(getLinksCommand())
	cmd.AddCommand(getLinkCommand())

	cmd.AddCommand(getHostsCommand())
	cmd.AddCommand(getHostCommand())

	cmd.AddCommand(getStatsCommand())
	return cmd
}

func getStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start {device} [args]",
		Short: "Commands for starting simulated entities",
	}

	cmd.AddCommand(startDeviceCommand())
	return cmd
}

func getStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop {device} [args]",
		Short: "Commands for stopping simulated entities",
	}

	cmd.AddCommand(stopDeviceCommand())
	return cmd
}

func getEnableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable {port} [args]",
		Short: "Commands for enabling simulated entities",
	}

	cmd.AddCommand(enablePortCommand())
	return cmd
}

func getDisableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable {port} [args]",
		Short: "Commands for disabling simulated entities",
	}

	cmd.AddCommand(disablePortCommand())
	return cmd
}

func getEmitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "emit {arp} [args]",
		Short: "Emit ARP, DHCP requests, etc",
	}

	cmd.AddCommand(emitARPsCommand())
	return cmd
}
