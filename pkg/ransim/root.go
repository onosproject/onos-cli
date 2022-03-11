// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ransim

import (
	"github.com/onosproject/onos-lib-go/pkg/cli"
	loglib "github.com/onosproject/onos-lib-go/pkg/logging/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "ransim"
	defaultAddress = "ran-simulator:5150"
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
		Use:   "ransim {get,set,create,delete,starts,stop,load,clear} [args]",
		Short: "ONOS RAN simulator commands",
	}

	cli.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(cli.GetConfigCommand())

	cmd.AddCommand(getCreateCommand())
	cmd.AddCommand(getDeleteCommand())
	cmd.AddCommand(getGetCommand())
	cmd.AddCommand(getSetCommand())

	cmd.AddCommand(startNodeCommand())
	cmd.AddCommand(stopNodeCommand())

	cmd.AddCommand(loadCommand())
	cmd.AddCommand(clearCommand())

	cmd.AddCommand(loglib.GetCommand())
	return cmd
}

func getCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create {node,cell,route} [args]",
		Short: "Commands for creating simulated entities",
	}

	cmd.AddCommand(createNodeCommand())
	cmd.AddCommand(createCellCommand())
	cmd.AddCommand(createRouteCommand())
	return cmd
}

func getDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete {node,cell,route} [args]",
		Short: "Commands for deleting simulated entities",
	}
	cmd.AddCommand(deleteNodeCommand())
	cmd.AddCommand(deleteCellCommand())
	cmd.AddCommand(deleteRouteCommand())
	cmd.AddCommand(deleteMetricCommand())
	cmd.AddCommand(deleteMetricsCommand())
	return cmd
}

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {plmnid,layout,node(s),cell(s),ue(s),ueCount,route(s)} [args]",
		Short: "Commands for retrieving RAN simulator model and other information",
	}

	cmd.AddCommand(getPlmnIDCommand())
	cmd.AddCommand(getLayoutCommand())

	cmd.AddCommand(getNodesCommand())
	cmd.AddCommand(getNodeCommand())

	cmd.AddCommand(getCellsCommand())
	cmd.AddCommand(getCellCommand())

	cmd.AddCommand(getUEsCommand())
	cmd.AddCommand(getUECommand())
	cmd.AddCommand(getUECountCommand())

	cmd.AddCommand(getRouteCommand())
	cmd.AddCommand(getRoutesCommand())

	cmd.AddCommand(getMetricCommand())
	cmd.AddCommand(getMetricsCommand())
	return cmd
}

func getSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {node,cell,ue,ueCount,metric} [args]",
		Short: "Commands for setting RAN simulator model metrics and other information",
	}

	cmd.AddCommand(updateNodeCommand())
	cmd.AddCommand(updateCellCommand())
	cmd.AddCommand(updateUECommand())
	cmd.AddCommand(setMetricCommand())
	cmd.AddCommand(setUECountCommand())
	return cmd
}
