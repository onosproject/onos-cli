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
	"bytes"
	"context"
	"fmt"
	"github.com/gogo/protobuf/types"
	"io"
	"text/tabwriter"
	"time"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getGetEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "entity <id>",
		Aliases: []string{"entities"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get Entity",
		RunE:    runGetEntityCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getGetRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "relation <id>",
		Aliases: []string{"relations"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get Relation",
		RunE:    runGetRelationCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getGetKindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kind <id>",
		Aliases: []string{"kinds"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get Kind",
		RunE:    runGetKindCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getAddEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity <id> [args]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Add Entity",
		RunE:  runAddEntityCommand,
	}
	cmd.Flags().StringP("kind", "k", "", "Kind ID")
	//_ = cmd.MarkFlagRequired("kind")
	cmd.Flags().StringToStringP("attributes", "a", map[string]string{}, "an user defined mapping of entity attributes")
	return cmd
}

func getAddRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <id> <src-entity-id> <tgt-entity-id> [args]",
		Args:  cobra.MinimumNArgs(3),
		Short: "Add Relation",
		RunE:  runAddRelationCommand,
	}
	cmd.Flags().StringP("kind", "k", "", "Kind ID")
	//_ = cmd.MarkFlagRequired("kind")
	return cmd
}

func getAddKindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kind <id> <name> [args]",
		Args:  cobra.MinimumNArgs(2),
		Short: "Add Kind",
		RunE:  runAddKindCommand,
	}
	return cmd
}

func getRemoveObjectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "object <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove an object",
		RunE:  runRemoveObjectCommand,
	}
}

func getWatchEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity [id] [args]",
		Short: "Watch Entities",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchEntityCommand,
	}
	cmd.Flags().BoolP("noreplay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <id> [args]",
		Short: "Watch Relations",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchRelationCommand,
	}
	cmd.Flags().BoolP("noreplay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchKindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kind [id] [args]",
		Short: "Watch Kinds",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runWatchKindCommand,
	}
	cmd.Flags().BoolP("noreplay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all [args]",
		Short: "Watch Entities and Relations",
		RunE:  runWatchAllCommand,
	}
	cmd.Flags().BoolP("noreplay", "r", false, "do not replay past topo updates")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runGetEntityCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topoapi.Object_ENTITY)
}

func runGetRelationCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topoapi.Object_RELATION)
}

func runGetKindCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topoapi.Object_KIND)
}

func runAddEntityCommand(cmd *cobra.Command, args []string) error {
	return writeObject(cmd, args, topoapi.Object_ENTITY)
}

func runAddRelationCommand(cmd *cobra.Command, args []string) error {
	return writeObject(cmd, args, topoapi.Object_RELATION)
}

func runAddKindCommand(cmd *cobra.Command, args []string) error {
	return writeObject(cmd, args, topoapi.Object_KIND)
}

func runRemoveObjectCommand(cmd *cobra.Command, args []string) error {
	id := args[0]

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Delete(ctx, &topoapi.DeleteRequest{ID: topoapi.ID(id)})
	if err != nil {
		return err
	}
	cli.Output("Removed object %s", id)
	return nil
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

func runGetCommand(cmd *cobra.Command, args []string, objectType topoapi.Object_Type) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	if !noHeaders {
		printHeader(false)
	}

	if len(args) == 0 {
		objects, err := listObjects(cmd)
		if err == nil {
			for _, object := range objects {
				if objectType == topoapi.Object_UNSPECIFIED || objectType == object.Type {
					printRow(object, false, noHeaders)
				}
			}
		}
	} else {
		id := args[0]
		object, err := getObject(cmd, topoapi.ID(id))
		if err != nil {
			return err
		}
		if object != nil {
			if objectType == topoapi.Object_UNSPECIFIED || objectType == object.Type {
				printRow(*object, false, noHeaders)
			}
		}
	}

	return nil
}

func writeObject(cmd *cobra.Command, args []string, objectType topoapi.Object_Type) error {
	var object *topoapi.Object
	id := args[0]
	//attributes, _ := cmd.Flags().GetStringToString("attributes")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	if objectType == topoapi.Object_ENTITY {
		kindID, _ := cmd.Flags().GetString("kind")
		entity := &topoapi.Object_Entity{
			Entity: &topoapi.Entity{
				KindID: topoapi.ID(kindID),
			},
		}

		object = &topoapi.Object{
			ID:   topoapi.ID(id),
			Type: objectType,
			Obj:  entity,
			// Deal with aspects
		}
	} else if objectType == topoapi.Object_RELATION {
		kindID, _ := cmd.Flags().GetString("kind")
		relation := &topoapi.Object_Relation{
			Relation: &topoapi.Relation{
				KindID:      topoapi.ID(kindID),
				SrcEntityID: topoapi.ID(args[1]),
				TgtEntityID: topoapi.ID(args[2]),
			},
		}

		object = &topoapi.Object{
			ID:   topoapi.ID(id),
			Type: objectType,
			Obj:  relation,
		}
	} else if objectType == topoapi.Object_KIND {
		kind := &topoapi.Object_Kind{
			Kind: &topoapi.Kind{
				Name: args[1],
			},
		}

		object = &topoapi.Object{
			ID:   topoapi.ID(id),
			Type: objectType,
			Obj:  kind,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Create(ctx, &topoapi.CreateRequest{Object: object})
	if err != nil {
		return err
	}
	return nil
}

func listObjects(cmd *cobra.Command) ([]topoapi.Object, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	resp, err := client.List(context.Background(), &topoapi.ListRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Objects, nil
}

func getObject(cmd *cobra.Command, id topoapi.ID) (*topoapi.Object, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	response, err := client.Get(ctx, &topoapi.GetRequest{ID: id})
	if err != nil {
		cli.Output("get error")
		return nil, err
	}
	return response.Object, nil
}

func watch(cmd *cobra.Command, args []string, objectType topoapi.Object_Type) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noreplay, _ := cmd.Flags().GetBool("noreplay")

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

	if !noHeaders {
		printHeader(true)
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
				printUpdateType(event.Type)
				printRow(event.Object, true, noHeaders)
			}
		}
	}

	return nil
}

func printHeader(printUpdateType bool) {
	var width = 16
	var prec = width - 1
	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)

	if printUpdateType {
		_, _ = fmt.Fprintf(writer, "%-*.*s", width, prec, "Update Type")
	}
	_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s%-*.*s\n", width, prec, "Object Type", width, prec, "Object ID", width, prec, "Kind ID", width, prec, "Attributes")
}

func printUpdateType(eventType topoapi.EventType) {
	var width = 16
	var prec = width - 1
	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)
	if eventType == topoapi.EventType_NONE {
		_, _ = fmt.Fprintf(writer, "%-*.*s", width, prec, "REPLAY")
	} else {
		_, _ = fmt.Fprintf(writer, "%-*.*s", width, prec, eventType)
	}
	_ = writer.Flush()
}

func printRow(object topoapi.Object, watch bool, noHeaders bool) {
	var width = 16
	var prec = width - 1
	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)

	switch object.Type {
	case topoapi.Object_ENTITY:
		var kindID topoapi.ID
		if e := object.GetEntity(); e != nil {
			kindID = e.KindID
		}
		// printUpdateType()
		_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s%s\n", width, prec, object.Type, width, prec, object.ID, width, prec, kindID, attrsToString(object.Aspects))
	case topoapi.Object_RELATION:
		r := object.GetRelation()
		// printUpdateType()
		_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s", width, prec, object.Type, width, prec, object.ID, width, prec, r.KindID)
		_, _ = fmt.Fprintf(writer, "src=%s, tgt=%s, %s\n", r.SrcEntityID, r.TgtEntityID, attrsToString(object.Aspects))
	case topoapi.Object_KIND:
		k := object.GetKind()
		// printUpdateType()
		_, _ = fmt.Fprintf(writer, "%-*.*s%-*.*s%-*.*s\n", width, prec, object.Type, width, prec, object.ID, width, prec, k.GetName())
	default:
		_, _ = fmt.Fprintf(writer, "\n")
	}
	_ = writer.Flush()
}

func attrsToString(attrs map[string]*types.Any) string {
	attributesBuf := bytes.Buffer{}
	first := true
	for key, attribute := range attrs {
		if !first {
			attributesBuf.WriteString(", ")
		} else {
			first = false
		}
		attributesBuf.WriteString(key)
		attributesBuf.WriteString(":")
		attributesBuf.WriteString(attribute.String()) // FIXME
	}
	return attributesBuf.String()
}
