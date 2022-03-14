// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"github.com/spf13/cobra"
)

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {entity|relation|kind} [args]",
		Short: "Get topology resources",
	}
	cmd.AddCommand(getGetEntityCommand())
	cmd.AddCommand(getGetRelationCommand())
	cmd.AddCommand(getGetKindCommand())
	cmd.AddCommand(getGetObjectsCommand())
	return cmd
}

func getAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create {entity|relation|kind} [args]",
		Short: "Create a topology resource",
	}
	cmd.AddCommand(getCreateEntityCommand())
	cmd.AddCommand(getCreateRelationCommand())
	cmd.AddCommand(getCreateKindCommand())
	return cmd
}

func getUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {entity|relation|kind} [args]",
		Short: "Update a topology resource",
	}
	cmd.AddCommand(getUpdateEntityCommand())
	cmd.AddCommand(getUpdateRelationCommand())
	cmd.AddCommand(getUpdateKindCommand())
	return cmd
}

func getDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete {entity|relation|kind} [args]",
		Short: "Delete a topology resource",
	}
	cmd.AddCommand(getDeleteRelationCommand())
	cmd.AddCommand(getDeleteEntityCommand())
	cmd.AddCommand(getDeleteKindCommand())
	return cmd
}

func getWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch {entity|relation|kind|all} [args]",
		Short: "Watch for changes to a topology resource type",
	}
	cmd.AddCommand(getWatchEntityCommand())
	cmd.AddCommand(getWatchRelationCommand())
	cmd.AddCommand(getWatchKindCommand())
	cmd.AddCommand(getWatchAllCommand())
	return cmd
}
