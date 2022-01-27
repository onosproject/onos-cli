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
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"io"
	"os"
	"time"
)

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
		Use:   "transactions",
		Short: "Watch configuration transaction changes",
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

	writer := os.Stdout
	if !noHeaders {
		printTransactionHeader(writer, verbose, false)
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewTransactionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if len(args) > 0 {
		return getTransactions(ctx, writer, client, v2.TransactionID(args[0]), verbose)
	}
	return listTransactions(ctx, writer, client, verbose)
}

func getTransactions(ctx context.Context, writer *os.File, client admin.TransactionServiceClient, id v2.TransactionID, verbose bool) error {
	resp, err := client.GetTransaction(ctx, &admin.GetTransactionRequest{ID: id})
	if err != nil {
		cli.Output("Unable to list transactions: %s", err)
		return err
	}
	printTransaction(writer, resp.Transaction, verbose)
	return nil
}

func listTransactions(ctx context.Context, writer *os.File, client admin.TransactionServiceClient, verbose bool) error {
	stream, err := client.ListTransactions(ctx, &admin.ListTransactionsRequest{})
	if err != nil {
		cli.Output("Unable to list transactions: %s", err)
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			cli.Output("Unable to read transaction: %s", err)
			return err
		} else {
			printTransaction(writer, resp.Transaction, verbose)
		}
	}

	return nil
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

	writer := os.Stdout
	if !noHeaders {
		printTransactionHeader(writer, verbose, true)
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
		if len(id) == 0 || id == event.Transaction.ID {
			printTransactionUpdateType(writer, event.Type)
			printTransaction(writer, &event.Transaction, false)
		}
	}

	return nil
}

func printTransaction(writer io.Writer, t *v2.Transaction, verbose bool) {
	if verbose {
		_, _ = fmt.Fprintf(writer, "%-12s\t%16d\t%8d\t%-8d\t%-12s\t%-12s\t%-12s\t%-10s\t%-8t\n",
			t.ID, t.Index, t.Status.State, t.Revision, t.Created, t.Updated, t.Deleted, t.Username, t.Atomic)
	} else {
		_, _ = fmt.Fprintf(writer, "%-12s\t%16d\t%8d\t%-8d\n",
			t.ID, t.Index, t.Status.State, t.Revision)
	}
}

func printTransactionUpdateType(writer io.Writer, eventType v2.TransactionEventType) {
	if eventType == v2.TransactionEventType_TRANSACTION_REPLAYED {
		_, _ = fmt.Fprintf(writer, "%-12s\t", "REPLAY")
	} else {
		_, _ = fmt.Fprintf(writer, "%-12s\t", eventType)
	}
}

func printTransactionHeader(writer *os.File, verbose bool, event bool) {
	if event {
		_, _ = fmt.Fprintf(writer, "%-12s\t", "Event Type")
	}
	if verbose {
		_, _ = fmt.Fprintf(writer, "%-12s\t%-16s\t%-8s\t%-8s\t%-12s\t%-12s\t%-12s\t%-10s\t%-8s\n",
			"Transaction ID", "Index", "Status", "Revision", "Created", "Updated", "Deleted", "User Name", "Atomic")
	} else {
		_, _ = fmt.Fprintf(writer, "%-12s\t%-16s\t%-8s\t%-8s\n",
			"Transaction ID", "Index", "Status", "Revision")
	}
}
