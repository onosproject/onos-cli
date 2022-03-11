// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/spf13/cobra"
)

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {transactions,configurations,plugins} [args]",
		Short: "Get config resources",
	}
	cmd.AddCommand(getListTransactionsCommand())
	cmd.AddCommand(getListConfigurationsCommand())
	cmd.AddCommand(getListPluginsCommand())
	return cmd
}

func getWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch {transactions,configurations} [args]",
		Short: "Watch for updates to a config resource type",
	}
	cmd.AddCommand(getWatchConfigurationsCommand())
	cmd.AddCommand(getWatchTransactionsCommand())
	return cmd
}
