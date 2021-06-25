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

package e2t

import (
	"context"
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/onosproject/onos-api/go/onos/e2t/e2/v1beta1"
	subapi "github.com/onosproject/onos-api/go/onos/e2t/e2/v1beta1"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	subscriptionHeaders = "ID\tRevision\tService Model ID\tE2 NodeID\tEncoding\tPhase\tState"
	subscriptionFormat  = "%s\t%d\t%s:%s\t%s\t%s\t%s\t%s\n"
)

func displaySubscriptionHeaders(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, subscriptionHeaders)
}

func displaySubscription(writer io.Writer, sub *subapi.Subscription) {
	_, _ = fmt.Fprintf(writer, subscriptionFormat,
		sub.ID, sub.Revision, sub.SubscriptionMeta.ServiceModel.Name, sub.SubscriptionMeta.ServiceModel.Version, none(string(sub.SubscriptionMeta.E2NodeID)),
		none(sub.SubscriptionMeta.Encoding.String()), none(sub.Status.Phase.String()), none(sub.Status.State.String()))
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
	cmd.Flags().String("appInstanceID", "", "Application Identifier")
	cmd.Flags().String("e2NodeID", "", "Identifier of the E2 node")
	cmd.Flags().String("smID", "", "Identifier of the service model")
	cmd.Flags().String("smVer", "", "Version of the service model")
	return cmd
}

func getRemoveSubscriptionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription",
		Short: "Remove subscription",
		RunE:  runRemoveSubscriptionCommand,
	}
	cmd.Flags().String("transactionID", "", "Identifier")
	cmd.Flags().String("appID", "", "Application Identifier")
	cmd.Flags().String("appInstanceID", "", "Application Identifier")
	cmd.Flags().String("e2NodeID", "", "Identifier of the E2 node")
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

	client := v1beta1.NewSubscriptionAdminServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.ListSubscriptions(ctx, &v1beta1.ListSubscriptionsRequest{})
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

func runAddSubscriptionCommand(cmd *cobra.Command, args []string) error {
	return errors.New("Unimplemented")
	// ID, _ := cmd.Flags().GetString("ID")
	// if ID == "" {
	// 	return errors.New("identifier must be specified with --ID")
	// }
	// appID, _ := cmd.Flags().GetString("appID")
	// if appID == "" {
	// 	return errors.New("appID must be specified with --appID")
	// }
	// appInstanceID, _ := cmd.Flags().GetString("appInstanceID")
	// if appInstanceID == "" {
	// 	return errors.New("appInstanceID must be specified with --appInstanceID")
	// }
	// e2NodeID, _ := cmd.Flags().GetString("e2NodeID")
	// if e2NodeID == "" {
	// 	return errors.New("e2NodeID must be specified with --e2NodeID")
	// }
	// smID, _ := cmd.Flags().GetString("smID")
	// if smID == "" {
	// 	return errors.New("service model ID must be specified with --smID")
	// }
	// smVer, _ := cmd.Flags().GetString("smVer")
	// if smVer == "" {
	// 	return errors.New("service model version must be specified with --smVer")
	// }

	// conn, err := cli.GetConnection(cmd)
	// if err != nil {
	// 	return err
	// }
	// defer conn.Close()
	// outputWriter := cli.GetOutput()
	// writer := new(tabwriter.Writer)
	// writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	// client := v1beta1.NewSubscriptionServiceClient(conn)

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// request := v1beta1.SubscribeRequest{
	// 	Headers: subapi.RequestHeaders{
	// 		AppID:         subapi.AppID(appID),
	// 		AppInstanceID: subapi.AppInstanceID(appInstanceID),
	// 		E2NodeID:      subapi.E2NodeID(e2NodeID),
	// 		ServiceModel: subapi.ServiceModel{
	// 			Name:    subapi.ServiceModelName(smID),
	// 			Version: subapi.ServiceModelVersion(smVer),
	// 		},
	// 		Encoding: 0,
	// 	},
	// 	TransactionID: subapi.TransactionID(ID),
	// 	Subscription:  subapi.SubscriptionSpec{},
	// }
	// _, err = client.Subscribe(ctx, &request)

	// if err != nil {
	// 	return err
	// }

	// _ = writer.Flush()
	// return nil
}

func runRemoveSubscriptionCommand(cmd *cobra.Command, args []string) error {
	return errors.New("Unimplemented")
	// transactionID, _ := cmd.Flags().GetString("transactionID")
	// if transactionID == "" {
	// 	return errors.New("identifier must be specified with --transactionID")
	// }
	// appID, _ := cmd.Flags().GetString("appID")
	// if appID == "" {
	// 	return errors.New("appID must be specified with --appID")
	// }
	// appInstanceID, _ := cmd.Flags().GetString("appInstanceID")
	// if appInstanceID == "" {
	// 	return errors.New("appInstanceID must be specified with --appInstanceID")
	// }
	// e2NodeID, _ := cmd.Flags().GetString("e2NodeID")
	// if e2NodeID == "" {
	// 	return errors.New("e2NodeID must be specified with --e2NodeID")
	// }
	// conn, err := cli.GetConnection(cmd)
	// if err != nil {
	// 	return err
	// }
	// defer conn.Close()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// client := v1beta1.NewSubscriptionServiceClient(conn)
	// _, err = client.Unsubscribe(ctx, &subapi.UnsubscribeRequest{
	// 	Headers: subapi.RequestHeaders{
	// 		AppID:         subapi.AppID(appID),
	// 		AppInstanceID: subapi.AppInstanceID(appInstanceID),
	// 		E2NodeID:      subapi.E2NodeID(e2NodeID),
	// 		ServiceModel:  subapi.ServiceModel{},
	// 		Encoding:      0,
	// 	},
	// 	TransactionID: subapi.TransactionID(transactionID),
	// })

	// return err
}

func runGetSubscriptionCommand(cmd *cobra.Command, args []string) error {
	return errors.New("Unimplemented")
	// noHeaders, _ := cmd.Flags().GetBool("no-headers")
	// ID := v1beta1.ID(args[0])
	// conn, err := cli.GetConnection(cmd)
	// if err != nil {
	// 	return err
	// }
	// defer conn.Close()
	// outputWriter := cli.GetOutput()
	// writer := new(tabwriter.Writer)
	// writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	// if !noHeaders {
	// 	displaySubscriptionHeaders(writer)
	// }

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// client := subscription.NewClient(conn)

	// sub, err := client.Get(ctx, ID)
	// if err != nil {
	// 	return err
	// }

	// displaySubscription(writer, sub)
	// _ = writer.Flush()

	// return nil
}

func none(s string) string {
	if s == "" {
		return "<None>"
	}
	return s
}
