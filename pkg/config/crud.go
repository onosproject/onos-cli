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
)

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {transactions,configurations,plugins} [args]",
		Short: "Get config resources",
	}
	cmd.AddCommand(getListTransactionsCommand())
	cmd.AddCommand(getListConfigurationsCommand())
	cmd.AddCommand(getListPluginsCommand())
	return cmd
}

func getWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch {transactions,configurations} [args]",
		Short: "Watch for updates to a config resource type",
	}
	cmd.AddCommand(getWatchConfigurationsCommand())
	cmd.AddCommand(getWatchTransactionsCommand())
	return cmd
}
