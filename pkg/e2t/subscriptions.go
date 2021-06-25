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
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/onosproject/onos-api/go/onos/e2t/e2/v1beta1"
	subapi "github.com/onosproject/onos-api/go/onos/e2t/e2/v1beta1"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	subscriptionHeaders = "ID\tRevision\tApp ID\tService Model ID\tE2 NodeID\tEncodStatus"
	subscriptionFormat  = "%s\t%d\t%s\t%s:%s\t%s\t%s\n"
)

func displaySubscriptionHeaders(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, subscriptionHeaders)
}

func displaySubscription(writer io.Writer, sub *subapi.Subscription) {
	_, _ = fmt.Fprintf(writer, subscriptionFormat,
		sub.ID, sub.Revision, sub.SubscriptionMeta.ServiceModel.Name, sub.SubscriptionMeta.E2NodeID, sub.SubscriptionMeta.Encoding,
		string(sub.Status.Phase), string(sub.Status.State))
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
