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

package a1t

import (
	"context"
	"fmt"
	a1 "github.com/onosproject/onos-api/go/onos/a1t/admin"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"io"
	"text/tabwriter"
)

const (
	subscriptionFormat = "%-50s\t%-30s\t%-30s\t%s\n"
)

func displaySubscriptionHeaders(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, subscriptionFormat,
		"xApp ID", "xApp A1 Interface", "A1 Service", "A1 Service Type ID\t")
}

func displaySubscription(writer io.Writer, resp *a1.GetXAppConnectionResponse) {
	_, _ = fmt.Fprintf(writer, subscriptionFormat,
		resp.XappId, resp.XappA1Endpoint, resp.SupportedA1Service, resp.SupportedA1ServiceTypeId)
}

func getGetSubscriptionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "subscription",
		Short: "Get A1 subscription(s)",
		RunE: runGetSubscriptionCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().String("xAppID", "", "xApp ID (optional)")
	return cmd
}

func runGetSubscriptionCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	defer conn.Close()
	if err != nil {
		return err
	}
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		displaySubscriptionHeaders(writer)
		_ = writer.Flush()
	}

	client := a1.NewA1TAdminServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutTimer)
	defer cancel()

	xAppID := ""

	if cmd.Flags().Changed("xAppID") {
		xAppID, err = cmd.Flags().GetString("xAppID")
		if err != nil {
			return err
		}
	}

	req := &a1.GetXAppConnectionsRequest{
		XappId: xAppID,
	}

	stream, err := client.GetXAppConnections(ctx, req)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			cli.Output("Error receiving notification: %v", err)
			return err
		}
		displaySubscription(writer, resp)
		_ = writer.Flush()
	}
	return nil
}
