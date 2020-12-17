// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2sub

import (
	"context"
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	subapi "github.com/onosproject/onos-api/go/onos/e2sub/subscription"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/subscription"
	"github.com/spf13/cobra"
)

const (
	subscriptionHeaders = "ID\tRevision\tApp ID\tService Model ID\tE2 NodeID\tStatus"
	subscriptionFormat  = "%s\t%d\t%s\t%s\t%s\t%d\n"
)

func displaySubscriptionHeaders(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, subscriptionHeaders)
}

func displaySubscription(writer io.Writer, sub *subapi.Subscription) {
	_, _ = fmt.Fprintf(writer, subscriptionFormat,
		sub.ID, sub.Revision, sub.AppID, sub.Details.ServiceModel.ID, sub.Details.E2NodeID,
		sub.Lifecycle.Status)
}

func getListSubscriptionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "List subscriptions",
		RunE:  runListSubscriptionsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getAddSubscriptionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription",
		Short: "Get subscription",
		RunE:  runAddSubscriptionCommand,
	}
	cmd.Flags().String("ID", "", "Identifier")
	cmd.Flags().String("appID", "", "Application Identifier")
	cmd.Flags().String("e2NodeID", "", "Identifier of the E2 node")
	cmd.Flags().String("smID", "", "Identifier of the service model")
	cmd.Flags().Int32("revision", 0, "Revision")
	return cmd
}

func getRemoveSubscriptionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription",
		Short: "Remove subscription",
		Args:  cobra.ExactArgs(1),
		RunE:  runRemoveSubscriptionCommand,
	}

	return cmd
}

func getGetSubscriptionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription",
		Short: "Get subscription",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetSubscriptionCommand,
	}

	return cmd
}

func runListSubscriptionsCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := subscription.NewClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.List(ctx)
	if err != nil {
		return err
	}

	if !noHeaders {
		displaySubscriptionHeaders(writer)
	}

	for _, sub := range response {
		displaySubscription(writer, &sub)
	}

	_ = writer.Flush()
	return nil
}

func runAddSubscriptionCommand(cmd *cobra.Command, args []string) error {
	ID, _ := cmd.Flags().GetString("ID")
	if ID == "" {
		return errors.New("identifier must be specified with --ID")
	}
	appID, _ := cmd.Flags().GetString("appID")
	if appID == "" {
		return errors.New("appID must be specified with --appID")
	}
	e2NodeID, _ := cmd.Flags().GetString("e2NodeID")
	if e2NodeID == "" {
		return errors.New("e2NodeID must be specified with --e2NodeID")
	}
	smID, _ := cmd.Flags().GetString("smId")
	revision, _ := cmd.Flags().GetInt32("revision")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := subscription.NewClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub := &subapi.Subscription{
		ID:       subapi.ID(ID),
		Revision: subapi.Revision(revision),
		AppID:    subapi.AppID(appID),
		Details: &subapi.SubscriptionDetails{
			E2NodeID:     subapi.E2NodeID(e2NodeID),
			ServiceModel: subapi.ServiceModel{ID: subapi.ServiceModelID(smID)},
		},
		Lifecycle: subapi.Lifecycle{Status: subapi.Status_ACTIVE},
	}

	err = client.Add(ctx, sub)
	if err != nil {
		return err
	}

	_ = writer.Flush()
	return nil
}

func runRemoveSubscriptionCommand(cmd *cobra.Command, args []string) error {
	ID := args[0]
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := subscription.NewClient(conn)

	sub, err := client.Get(context.Background(), subapi.ID(ID))
	if err != nil {
		return nil
	}

	err = client.Remove(ctx, sub)

	return err
}

func runGetSubscriptionCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	ID := subapi.ID(args[0])
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		displaySubscriptionHeaders(writer)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := subscription.NewClient(conn)

	sub, err := client.Get(ctx, ID)
	if err != nil {
		return err
	}

	displaySubscription(writer, sub)
	_ = writer.Flush()

	return nil
}
