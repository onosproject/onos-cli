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

func getListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list {numconflicts/neighbors/metric/pci}",
		Short: "List PCI resources for a specific cell",
	}
	cmd.AddCommand(getListNumConflicts())
	cmd.AddCommand(getListNeighbors())
	cmd.AddCommand(getListMetric())
	cmd.AddCommand(getListPci())
	return cmd
}

func getListAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listall {numconflicts/neighbors/metric/pci}",
		Short: "List PCI resources for all cells",
	}
	cmd.AddCommand(getListNumConflictsAll())
	cmd.AddCommand(getListNeighborsAll())
	cmd.AddCommand(getListMetricAll())
	cmd.AddCommand(getListPciAll())
	return cmd
}
