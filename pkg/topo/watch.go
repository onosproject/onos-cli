// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"context"
	"fmt"
	"io"
	"os"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getWatchEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity [id] [args]",
		Short: "Watch Entities",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchEntityCommand,
	}
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay existing topo state")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().String("kind", "", "kind query")
	cmd.Flags().String("label", "", "label query")
	return cmd
}

func getWatchRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <id> [args]",
		Short: "Watch Relations",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchRelationCommand,
	}
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay exiting topo state")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().String("kind", "", "kind query")
	cmd.Flags().String("label", "", "label query")
	return cmd
}

func getWatchKindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kind [id] [args]",
		Short: "Watch Kinds",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchKindCommand,
	}
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay exiting topo state")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().String("kind", "", "kind query")
	cmd.Flags().String("label", "", "label query")
	return cmd
}

func getWatchAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all [args]",
		Short: "Watch Entities, Relations and Kinds",
		RunE:  runWatchAllCommand,
	}
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay exiting topo state")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().String("kind", "", "kind query")
	cmd.Flags().String("label", "", "label query")
	return cmd
}

func runWatchEntityCommand(cmd *cobra.Command, args []string) error {
	return watch(cmd, args, topoapi.Object_ENTITY)
}

func runWatchRelationCommand(cmd *cobra.Command, args []string) error {
	return watch(cmd, args, topoapi.Object_RELATION)
}

func runWatchKindCommand(cmd *cobra.Command, args []string) error {
	return watch(cmd, args, topoapi.Object_KIND)
}

func runWatchAllCommand(cmd *cobra.Command, args []string) error {
	return watch(cmd, args, topoapi.Object_UNSPECIFIED)
}

func watch(cmd *cobra.Command, args []string, objectType topoapi.Object_Type) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noreplay, _ := cmd.Flags().GetBool("noreplay")
	verbose, _ := cmd.Flags().GetBool("verbose")

	var id topoapi.ID
	if len(args) > 0 {
		id = topoapi.ID(args[0])
	} else {
		id = topoapi.NullID
	}

	filters := compileFilters(cmd, objectType)

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	req := &topoapi.WatchRequest{
		Filters:  filters,
		Noreplay: noreplay,
	}

	stream, err := client.Watch(context.Background(), req)
	if err != nil {
		return err
	}

	writer := os.Stdout
	if !noHeaders {
		printHeader(writer, objectType, true, verbose)
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
		// TODO: Filtering for ID and object type is still client-side; labels and kinds are server-side now
		if id == topoapi.NullID || id == event.Object.ID {
			if event.Object.Type == topoapi.Object_UNSPECIFIED || objectType == event.Object.Type {
				printUpdateType(writer, event.Type, event.Object.Type, verbose)
				printObject(writer, event.Object, verbose, false, false)
			}
		}
	}

	return nil
}

func printUpdateType(writer io.Writer, eventType topoapi.EventType, objectType topoapi.Object_Type, verbose bool) {
	if verbose {
		if eventType == topoapi.EventType_NONE {
			fmt.Fprintf(writer, "Update Type:\t%s\n", "REPLAY")
		} else {
			fmt.Fprintf(writer, "Update Type:\t%s\n", eventType)
		}
		fmt.Fprintf(writer, "Object Type:\t%s\n", objectType)
	} else {
		if eventType == topoapi.EventType_NONE {
			_, _ = fmt.Fprintf(writer, "%-12s\t%-10s\t", "REPLAY", objectType)
		} else {
			_, _ = fmt.Fprintf(writer, "%-12s\t%-10s\t", eventType, objectType)
		}
	}

}
