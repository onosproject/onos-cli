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

package config

import (
	"github.com/spf13/cobra"
	"text/template"
)

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {device-changes,network-changes,plugins,opstate,snapshots} [args]",
		Short: "Get config resources",
	}
	cmd.AddCommand(getListNetworkChangesCommand())
	cmd.AddCommand(getListDeviceChangesCommand())
	cmd.AddCommand(getGetOpstateCommand())
	cmd.AddCommand(getListSnapshotsCommand())
	return cmd
}

func getWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch resource [args]",
		Short: "Watch for updates to a config resource type",
	}
	cmd.AddCommand(getWatchDeviceChangesCommand())
	cmd.AddCommand(getWatchNetworkChangesCommand())
	cmd.AddCommand(getWatchOpstateCommand())
	cmd.AddCommand(getWatchSnapshotsCommand())
	return cmd
}

var funcMapChanges = template.FuncMap{
	"wrappath":      wrapPath,
	"valuetostring": valueToSstring,
}
