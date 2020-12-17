// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

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
