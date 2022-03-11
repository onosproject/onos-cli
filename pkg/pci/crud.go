// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package pci

import "github.com/spf13/cobra"

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {conflicts/resolved/cell/cells}",
		Short: "Get PCI resources",
	}
	cmd.AddCommand(getGetConflicts())
	cmd.AddCommand(getGetResolvedConflicts())
	cmd.AddCommand(getGetCell())
	cmd.AddCommand(getGetCells())
	return cmd
}
