// SPDX-FileCopyrightText: 2020-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package o1t

import "github.com/spf13/cobra"

func getListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list {sessions} [args]",
		Short: "List o1t sessions",
	}
	cmd.AddCommand(getListSessionsCommand())
	return cmd
}

func getWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch {sessions} [args]",
		Short: "Watch o1t sessions",
	}
	cmd.AddCommand(getWatchSessionsCommand())
	return cmd
}
