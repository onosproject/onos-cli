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

package e2sub

import "github.com/spf13/cobra"

func getListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list {endpoints | subscriptions} [args]",
		Short: "List E2Sub resources",
	}
	cmd.AddCommand(getListEndPointsCommand())
	cmd.AddCommand(getListSubscriptionsCommand())
	return cmd
}

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {endpoint | subscription} [args]",
		Short: "Get E2Sub resources",
	}
	cmd.AddCommand(getGetEndPointCommand())
	cmd.AddCommand(getGetSubscriptionCommand())
	return cmd
}

func getAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add {endpoint | subscription} [args]",
		Short: "Add E2Sub resources",
	}
	cmd.AddCommand(getAddEndPointCommand())
	cmd.AddCommand(getAddSubscriptionCommand())
	return cmd
}

func getRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove {endpoint | subscription} [args]",
		Short: "Remove E2Sub resources",
	}
	cmd.AddCommand(getRemoveEndPointCommand())
	cmd.AddCommand(getRemoveSubscriptionCommand())
	return cmd
}
