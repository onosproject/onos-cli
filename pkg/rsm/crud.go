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

