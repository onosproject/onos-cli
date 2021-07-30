// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
