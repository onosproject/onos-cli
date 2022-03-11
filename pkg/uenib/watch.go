// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package uenib

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/onosproject/onos-api/go/onos/uenib"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch {ue|ues}",
		Short: "Watch for changes to UE information",
	}
	cmd.AddCommand(getWatchUECommand())
	cmd.AddCommand(getWatchUEsCommand())
	return cmd
}

func getWatchUECommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ue ue-id [args]",
		Short: "Watch for changes to a specific UE information",
		Args:  cobra.ExactArgs(1),
		RunE:  runWatchUEsCommand,
	}
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay existing UE state")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().StringSliceP("aspect", "a", []string{}, "UE aspects to watch")
	return cmd
}

func getWatchUEsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ues [args]",
		Args:  cobra.ExactArgs(0),
		Short: "Watch for changes to any UE information",
		RunE:  runWatchUEsCommand,
	}
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay existing UE state")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().StringSliceP("aspect", "a", []string{}, "UE aspects to watch")
	return cmd
}

func runWatchUEsCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noReplay, _ := cmd.Flags().GetBool("no-replay")

	aspectTypes, _ := cmd.Flags().GetStringSlice("aspect")

	var id uenib.ID
	if len(args) > 0 {
		id = uenib.ID(args[0])
	} else {
		id = uenib.NullID
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := uenib.CreateUEServiceClient(conn)

	req := &uenib.WatchUERequest{
		AspectTypes: aspectTypes,
		Noreplay:    noReplay,
	}

	stream, err := client.WatchUEs(context.Background(), req)
	if err != nil {
		return err
	}

	writer := os.Stdout
	if !noHeaders {
		printHeader(writer, true)
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
		// TODO: Filtering for ID is still client-side; we need to fix this
		if id == uenib.NullID || id == event.UE.ID {
			printUpdateType(writer, event.Type)
			printUE(writer, event.UE, false)
		}
	}

	return nil
}

func printUpdateType(writer io.Writer, eventType uenib.EventType) {
	if eventType == uenib.EventType_NONE {
		_, _ = fmt.Fprintf(writer, "%-12s\t", "REPLAY")
	} else {
		_, _ = fmt.Fprintf(writer, "%-12s\t", eventType)
	}
}
