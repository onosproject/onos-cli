// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

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

const transactionListTemplate = "table{{.ID}}\t{{.Index}}\t{{.Status.State}}\t{{.TransactionType}}\t{{.Created}}\t{{.Updated}}\t{{.Deleted}}\t{{.Username}}\t{{.TransactionStrategy.Isolation}}\t{{.TransactionStrategy.Synchronicity}}"

var transactionListTemplateVerbose = fmt.Sprintf("%s\t{{.Transaction}}", transactionListTemplate)

const transactionEventTemplate = "table{{.EventType}}\t{{.ID}}\t{{.Index}}\t{{.Status.State}}\t{{.TransactionType}}\t{{.Created}}\t{{.Updated}}\t{{.Deleted}}\t{{.Username}}\t{{.TransactionStrategy.Isolation}}\t{{.TransactionStrategy.Synchronicity}}"

type cliTransaction struct {
	v2.Transaction
	TransactionType string
	EventType       v2.TransactionEvent_EventType
}

type transactionEventWidths struct {
	EventType       int
	ID              int
	Index           int
	TransactionType int
	Created         int
	Updated         int
	Deleted         int
	Username        int
	Status          struct {
		State int
	}
	TransactionStrategy struct {
		Synchronicity int
		Isolation     int
	}
}

var transactionWidths = transactionEventWidths{
	EventType:       12,
	ID:              42,
	Index:           10,
	TransactionType: 10,
	Created:         19,
	Updated:         19,
	Deleted:         19,
	Username:        10,
	Status:          struct{ State int }{State: 12},
	TransactionStrategy: struct {
		Synchronicity int
		Isolation     int
	}{
		Synchronicity: 13,
		Isolation:     12,
	},
}

func getListTransactionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transactions [transactionID]",
		Short:   "Get list of configuration transactions",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"transaction"},
		RunE:    runListTransactionsCommand,
	}
	cmd.Flags().Uint64("index", 0, "optional index for transaction lookup; takes precedence over ID")
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the change with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchTransactionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transactions [transactionID]",
		Short:   "Watch configuration transaction changes",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"transaction"},
		RunE:    runWatchTransactionsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay existing transactions")
	return cmd
}

func runListTransactionsCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	verbose, _ := cmd.Flags().GetBool("verbose")
	index, _ := cmd.Flags().GetUint64("index")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewTransactionServiceClient(conn)
	ctx, cancel := context.WithTimeout(cli.NewContextWithAuthHeaderFromFlag(cmd.Context(), cmd.Flag(cli.AuthHeaderFlag)), 15*time.Second)
	defer cancel()

	if index > 0 {
		return getTransaction(ctx, client, &admin.GetTransactionRequest{Index: v2.Index(index)}, noHeaders, verbose)
	}
	if len(args) > 0 {
		return getTransaction(ctx, client, &admin.GetTransactionRequest{ID: v2.TransactionID(args[0])}, noHeaders, verbose)
	}
	return listTransactions(ctx, client, noHeaders, verbose)
}

func getTransaction(ctx context.Context, client admin.TransactionServiceClient,
	req *admin.GetTransactionRequest, noHeaders bool, verbose bool) error {
	resp, err := client.GetTransaction(ctx, req)
	if err != nil {
		cli.Output("Unable to list transactions: %s", err)
		return err
	}

	f := format.Format(transactionListTemplate)
	if verbose {
		f = format.Format(transactionListTemplateVerbose)
	}
	return f.Execute(cli.GetOutput(), !noHeaders, 0, prepareTransactionOutput(resp.Transaction, v2.TransactionEvent_UNKNOWN))
}

func listTransactions(ctx context.Context, client admin.TransactionServiceClient, noHeaders bool, verbose bool) error {
	stream, err := client.ListTransactions(ctx, &admin.ListTransactionsRequest{})
	if err != nil {
		cli.Output("Unable to list transactions: %s", err)
		return err
	}

	f := format.Format(transactionListTemplate)
	if verbose {
		f = format.Format(transactionListTemplateVerbose)
	}

	var allTx []*cliTransaction

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return f.Execute(cli.GetOutput(), !noHeaders, 0, allTx)
		} else if err != nil {
			cli.Output("Unable to read transaction: %s", err)
			return err
		}

		tx := prepareTransactionOutput(resp.Transaction, v2.TransactionEvent_UNKNOWN)

		allTx = append(allTx, tx)
	}
}

func prepareTransactionOutput(tx *v2.Transaction, eventType v2.TransactionEvent_EventType) *cliTransaction {
	var txType string

	switch tx.GetDetails().(type) {
	case *v2.Transaction_Change:
		txType = "Change"
	case *v2.Transaction_Rollback:
		txType = "Rollback"
	}

	return &cliTransaction{
		TransactionType: txType,
		Transaction:     *tx,
		EventType:       eventType,
	}
}

func runWatchTransactionsCommand(cmd *cobra.Command, args []string) error {
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
	request := &admin.WatchTransactionsRequest{Noreplay: noReplay, ID: id}
	stream, err := client.WatchTransactions(cli.NewContextWithAuthHeaderFromFlag(cmd.Context(), cmd.Flag(cli.AuthHeaderFlag)), request)
	if err != nil {
		return err
	}

	f := format.Format(transactionEventTemplate)

	if !noHeaders {
		output, err := f.ExecuteFixedWidth(transactionWidths, true, nil)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", output)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			cli.Output("Error receiving notification : %v", err)
			return err
		}

		event := resp.TransactionEvent
		if len(id) == 0 || id == event.Transaction.ID {
			output, err := f.ExecuteFixedWidth(transactionWidths, false, prepareTransactionOutput(&resp.TransactionEvent.Transaction, event.Type))
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", output)
		}
	}

	return nil
}
