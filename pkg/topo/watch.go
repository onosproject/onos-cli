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
	"os"
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
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay past topo updates")
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
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay past topo updates")
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
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay past topo updates")
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
		// TODO: Filtering for ID and object type is still client-side; labels and kinds are server-side now
		if id == topoapi.NullID || id == event.Object.ID {
			if event.Object.Type == topoapi.Object_UNSPECIFIED || objectType == event.Object.Type {
				printUpdateType(writer, event.Type, event.Object.Type)
				printObject(writer, event.Object, true)
			}
		}
	}

	return nil
}

func printUpdateType(writer io.Writer, eventType topoapi.EventType, objectType topoapi.Object_Type) {
	if eventType == topoapi.EventType_NONE {
		_, _ = fmt.Fprintf(writer, "%-12s\t%-10s\t", "REPLAY", objectType)
	} else {
		_, _ = fmt.Fprintf(writer, "%-12s\t%-10s\t", eventType, objectType)
	}
}
