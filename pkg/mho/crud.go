// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mho

import "github.com/spf13/cobra"

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {ues|cells} [args]",
		Short: "Get UE, Cell info",
	}
	cmd.AddCommand(getGetUesCommand())
	cmd.AddCommand(getGetCellsCommand())
	return cmd
}
