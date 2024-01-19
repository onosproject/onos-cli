// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package uenib

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/onosproject/onos-api/go/onos/uenib"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {ue|ues}",
		Short: "Get UE information",
	}
	cmd.AddCommand(getGetUECommand())
	cmd.AddCommand(getGetUEsCommand())
	return cmd
}

func getGetUECommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ue ue-id [args]",
		Args:  cobra.ExactArgs(1),
		Short: "Get UE information",
		RunE:  runGetUECommand,
	}
	cmd.Flags().StringSliceP("aspect", "a", []string{}, "UE aspects to get")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the change with verbose output")
	return cmd
}

func getGetUEsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ues [args]",
		Args:  cobra.ExactArgs(0),
		Short: "Get list of UE information",
		RunE:  runGetUEsCommand,
	}
	cmd.Flags().StringSliceP("aspect", "a", []string{}, "UE aspects to get")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the change with verbose output")
	return cmd
}

func runGetUECommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	aspectTypes, _ := cmd.Flags().GetStringSlice("aspect")
	verbose, _ := cmd.Flags().GetBool("verbose")
	// headers do not make sense when printing flat
	if verbose {
		noHeaders = true
	}

	writer := os.Stdout
	if !noHeaders {
		printHeader(writer, false)
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := uenib.CreateUEServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	response, err := client.GetUE(ctx, &uenib.GetUERequest{ID: uenib.ID(args[0]), AspectTypes: aspectTypes})
	if err != nil {
		cli.Output("Unable to get UE aspects: %s", err)
		return err
	}

	printUE(writer, response.UE, verbose)
	return nil
}

func runGetUEsCommand(cmd *cobra.Command, _ []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	aspectTypes, _ := cmd.Flags().GetStringSlice("aspect")
	verbose, _ := cmd.Flags().GetBool("verbose")
	// headers do not make sense when printing flat
	if verbose {
		noHeaders = true
	}

	writer := os.Stdout
	if !noHeaders {
		printHeader(writer, false)
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := uenib.CreateUEServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	response, err := client.ListUEs(ctx, &uenib.ListUERequest{AspectTypes: aspectTypes})
	if err != nil {
		cli.Output("Unable to list UEs: %s", err)
		return err
	}

	for {
		resp, err := response.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			cli.Output("Unable to read UE: %s", err)
			return err
		}
		printUE(writer, resp.UE, verbose)
	}

	return nil
}

func printHeader(writer *os.File, replay bool) {
	if replay {
		_, _ = fmt.Fprintf(writer, "%-12s\t%-16s\t%-20s\t%s\n", "Event Type", "UE ID", "Aspect Type", "Aspect Value")
	} else {
		_, _ = fmt.Fprintf(writer, "%-16s\t%-20s\n", "UE ID", "Aspect Types")
	}
}

func printUE(writer io.Writer, ue uenib.UE, verbose bool) {
	if !verbose {
		_, _ = fmt.Fprintf(writer, "%-16s\t", ue.ID)
		aspectTypes := make([]string, 0, len(ue.Aspects))
		for k := range ue.Aspects {
			aspectTypes = append(aspectTypes, k)
		}
		_, _ = fmt.Fprintf(writer, "%s\n", strings.Join(aspectTypes[:], ","))
	} else {
		_, _ = fmt.Fprintf(writer, "ID: %s\n", ue.ID)
		_, _ = fmt.Fprintf(writer, "Aspects:\n")
		for aspectType, any := range ue.Aspects {
			_, _ = fmt.Fprintf(writer, "- %s=%s\n", aspectType, any.Value)
		}
	}
}
