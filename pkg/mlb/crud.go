// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mlb

import "github.com/spf13/cobra"

func getListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list {parameters/ocns}",
		Short: "List MLB resources",
	}
	cmd.AddCommand(getListParameters())
	cmd.AddCommand(getListOcns())
	return cmd
}

func getSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {parameters}",
		Short: "Set MLB resources",
	}
	cmd.AddCommand(getSetParameters())
	return cmd
}
