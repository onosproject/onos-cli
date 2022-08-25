// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package deviceprovisioner

import "github.com/spf13/cobra"

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {pipelines} [args]",
		Short: "Get pipelines",
	}
	cmd.AddCommand(getListPipelinesCommand())
	return cmd
}

func getWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch {pipelines} [args]",
		Short: "Watch for updates to a config resource type",
	}
	cmd.AddCommand(getWatchPipelinesCommand())
	return cmd
}
