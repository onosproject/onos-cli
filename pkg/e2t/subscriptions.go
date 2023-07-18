// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package e2t

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-cli/pkg/utils"
	"io"
	"strings"
	"text/tabwriter"

	subapi "github.com/onosproject/onos-api/go/onos/e2t/e2/v1beta1"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	subscriptionHeaders = "Subscription ID\tRevision\tService Model ID\tE2 NodeID\tEncoding\tPhase\tState"
	subscriptionFormat  = "%s\t%d\t%s:%s\t%s\t%s\t%s\t%s\n"
)

func displaySubscriptionHeaders(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, subscriptionHeaders)
}

func displaySubscription(writer io.Writer, sub *subapi.Subscription) {
	_, _ = fmt.Fprintf(writer, subscriptionFormat,
		sub.ID, sub.Revision, sub.SubscriptionMeta.ServiceModel.Name, sub.SubscriptionMeta.ServiceModel.Version, utils.None(string(sub.SubscriptionMeta.E2NodeID)),
		utils.None(sub.SubscriptionMeta.Encoding.String()), utils.None(sub.Status.Phase.String()), utils.None(sub.Status.State.String()))
}

func getGetSubscriptionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "Get SB subscriptions",
		RunE:  runGetSubscriptionsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runGetSubscriptionsCommand(cmd *cobra.Command, _ []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := subapi.NewSubscriptionAdminServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.ListSubscriptions(ctx, &subapi.ListSubscriptionsRequest{})
	if err != nil {
		return err
	}

	if !noHeaders {
		displaySubscriptionHeaders(writer)
	}

	for _, sub := range response.Subscriptions {
		pin := sub
		displaySubscription(writer, &pin)
	}

	_ = writer.Flush()
	return nil
}

func getGetSubscriptionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription",
		Short: "Get SB subscription",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetSubscriptionCommand,
	}
	return cmd
}

func runGetSubscriptionCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := subapi.NewSubscriptionAdminServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.GetSubscription(ctx, &subapi.GetSubscriptionRequest{
		SubscriptionID: subapi.SubscriptionID(args[0]),
	})
	if err != nil {
		return err
	}

	sub := response.Subscription
	_, _ = fmt.Fprintf(writer, "Subscription ID:\t%s\nRevision:\t%d\nService Model:\t%s\nService Model Version:\t%s\n",
		sub.ID, sub.Revision, sub.SubscriptionMeta.ServiceModel.Name, sub.SubscriptionMeta.ServiceModel.Version)
	_, _ = fmt.Fprintf(writer, "E2 Node ID:\t%s\nEncoding:\t%s\n",
		utils.None(string(sub.SubscriptionMeta.E2NodeID)), utils.None(sub.SubscriptionMeta.Encoding.String()))
	_, _ = fmt.Fprintf(writer, "Phase:\t%s\nStatus:\t%s\n", utils.None(sub.Status.Phase.String()), utils.None(sub.Status.State.String()))
	_, _ = fmt.Fprintf(writer, "Actions:\t%v\nTrigger:\t%v\n", sub.Spec.Actions, sub.Spec.EventTrigger)
	_ = writer.Flush()
	return nil
}

func getWatchSubscriptionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "Watch SB subscriptions",
		RunE:  runWatchSubscriptionsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().Bool("no-replay", false, "disables replay of existing state")
	return cmd
}

func runWatchSubscriptionsCommand(cmd *cobra.Command, _ []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noReplay, _ := cmd.Flags().GetBool("no-replay")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := subapi.NewSubscriptionAdminServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.WatchSubscriptions(ctx, &subapi.WatchSubscriptionsRequest{
		NoReplay: noReplay,
	})
	if err != nil {
		return err
	}

	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "Event Type\t")
		displaySubscriptionHeaders(writer)
		_ = writer.Flush()
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			cli.Output("Error receiving notification : %v", err)
			return err
		}

		event := res.Event
		_, _ = fmt.Fprintf(writer, "%s\t", strings.Replace(event.Type.String(), "SUBSCRIPTION_", "", 1))
		displaySubscription(writer, &event.Subscription)
		_ = writer.Flush()
	}

	return nil
}
