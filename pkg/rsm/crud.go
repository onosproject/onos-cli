// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package rsm

import "github.com/spf13/cobra"

func getSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {association}",
		Short: "Set RSM resources",
	}
	cmd.AddCommand(getSetAssociation())
	return cmd
}

func getCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create {slice}",
		Short: "Create RSM resources",
	}
	cmd.AddCommand(getCreateSlice())
	return cmd
}

func getUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update {slice}",
		Short: "Update RSM resources",
	}
	cmd.AddCommand(getUpdateSlice())
	return cmd
}

func getDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete {slice}",
		Short: "Delete RSM resources",
	}
	cmd.AddCommand(getDeleteSlice())
	return cmd
}
