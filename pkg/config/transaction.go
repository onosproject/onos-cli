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
	"context"
	"fmt"
	"github.com/onosproject/onos-api/go/onos/config/admin"
	v2 "github.com/onosproject/onos-api/go/onos/config/v2"
	"github.com/onosproject/onos-cli/pkg/format"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"io"
	"time"
)

const transactionListTemplate = "table{{.ID}}\t{{.Index}}\t{{.Revision}}\t{{.Status.State}}\t{{.Created}}\t{{.Updated}}\t{{.Deleted}}\t{{.Username}}\t{{.Atomic}}"
var transactionListTemplateVerbose = fmt.Sprintf("%s\t{{.Transaction}}", transactionListTemplate)
const transactionEventTemplate = "table{{.Type}}\t{{.Transaction.Index}}\t{{.Transaction.Revision}}\t{{.Transaction.Status.State}}\t{{.Transaction.Created}}\t{{.Transaction.Updated}}\t{{.Transaction.Deleted}}\t{{.Transaction.Username}}\t{{.Transaction.Atomic}}"
var transactionEventTemplateVerbose = fmt.Sprintf("%s\t{{.Transaction}}", transactionListTemplate)

type transactionEventWidths struct {
	Type          int
	Transaction struct {
		ID       int
		Created  int
		Updated  int
		Deleted  int
		Username int
		Atomic   int
		Status   struct {
			State int
		}
		Revision int
		Index    int
		Transaction   int
	}
}

var transactionWidths = transactionEventWidths{
	Type: 30,
	Transaction: struct {
		ID       int
		Created  int
		Updated  int
		Deleted  int
		Username int
		Atomic   int
		Status   struct {
			State int
		}
		Revision    int
		Index       int
		Transaction int
	}{
		ID: 13,
		Created: 13,
		Updated: 13,
		Deleted: 13,
		Username: 13,
		Atomic: 6,
		Status: struct{ State int }{State: 40},
		Revision: 5,
		Index: 5,
		Transaction: 50,
	},
}

func getListTransactionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transactions [transactionID]",
		Short: "Get list of configuration transactions",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runListTransactionsCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the change with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchTransactionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transactions [transactionID wildcard]",
		Short: "Watch configuration transaction changes",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runWatchTransactionsCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the change with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay existing UE state")
	return cmd
}

func runListTransactionsCommand(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	f := format.Format(transactionListTemplate)
	if verbose {
		f = format.Format(transactionListTemplateVerbose)
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewTransactionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	response, err := client.ListTransactions(ctx, &admin.ListTransactionsRequest{})
	if err != nil {
		cli.Output("Unable to list transactions: %s", err)
		return err
	}

	allTx := []*v2.Transaction{}

	for {
		resp, err := response.Recv()
		if err == io.EOF {
			if e := f.Execute(cli.GetOutput(), !noHeaders, 0, allTx); e != nil {
				return e
			}
			return nil
		} else if err != nil {
			cli.Output("Unable to read transaction: %s", err)
			return err
		}
		allTx = append(allTx, resp.Transaction)
	}
}

func runWatchTransactionsCommand(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noReplay, _ := cmd.Flags().GetBool("no-replay")

	id := v2.TransactionID("")
	if len(args) > 0 {
		id = v2.TransactionID(args[0])
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewTransactionServiceClient(conn)
	stream, err := client.WatchTransactions(context.Background(), &admin.WatchTransactionsRequest{Noreplay: noReplay})
	if err != nil {
		return err
	}

	f := format.Format(transactionEventTemplate)
	if verbose {
		f = format.Format(transactionEventTemplateVerbose)
	}

	if !noHeaders {
		output, err := f.ExecuteFixedWidth(transactionWidths, true, nil)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", output)
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

		event := res.TransactionEvent
		if len(id) == 0 || id == event.Transaction.ID {
			output, err := f.ExecuteFixedWidth(transactionWidths, false, res)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", output)
		}
	}

	return nil
}
