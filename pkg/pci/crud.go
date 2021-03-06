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
