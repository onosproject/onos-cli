// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2t

import "github.com/spf13/cobra"

func getListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list {connections} [args]",
		Short: "List E2T resources",
	}
	cmd.AddCommand(getListConnectionsCommand())
	return cmd
}

func getWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch {indications} [args]",
		Short: "Watch E2T resources",
	}
	cmd.AddCommand(getWatchIndicationsCommand())
	return cmd
}
