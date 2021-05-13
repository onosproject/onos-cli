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

package topo

import (
	"context"
	"fmt"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"io"
	"text/tabwriter"
)

func getWatchEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity [id] [args]",
		Short: "Watch Entities",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchEntityCommand,
	}
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	return cmd
}

func getWatchRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <id> [args]",
		Short: "Watch Relations",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchRelationCommand,
	}
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	return cmd
}

func getWatchKindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kind [id] [args]",
		Short: "Watch Kinds",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchKindCommand,
	}
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	return cmd
}

func getWatchAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all [args]",
		Short: "Watch Entities, Relations and Kinds",
		RunE:  runWatchAllCommand,
	}
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
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

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	req := &topoapi.WatchRequest{
		Noreplay: noreplay,
	}

	stream, err := client.Watch(context.Background(), req)
	if err != nil {
		return err
	}

	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		printHeader(writer, objectType, verbose, true)
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
		// FIXME: for now doing client-side filtering of events
		if id == topoapi.NullID || id == event.Object.ID {
			if event.Object.Type == topoapi.Object_UNSPECIFIED || objectType == event.Object.Type {
				printUpdateType(writer, event.Type)
				printObject(writer, event.Object, true)
				_ = writer.Flush()
			}
		}
	}

	return nil
}

func printUpdateType(writer io.Writer, eventType topoapi.EventType) {
	if eventType == topoapi.EventType_NONE {
		_, _ = fmt.Fprintf(writer, "%-*.*s", width, prec, "REPLAY")
	} else {
		_, _ = fmt.Fprintf(writer, "%-*.*s", width, prec, eventType)
	}
}
