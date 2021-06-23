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
	"io"
	"text/tabwriter"
	"time"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-lib-go/pkg/errors"
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
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().String("kind", "", "kind query")
	cmd.Flags().String("label", "", "label query")
	cmd.Flags().String("related-to", "", "use relation filter, must also specify related-via")
	cmd.Flags().String("related-via", "", "use relation filter, must also specify related-to")
	cmd.Flags().String("tgt-kind", "", "optional target kind for relation filter")
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
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().String("kind", "", "kind query")
	cmd.Flags().String("label", "", "label query")
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
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().String("label", "", "label query")
	return cmd
}

func runGetEntityCommand(cmd *cobra.Command, args []string) error {
	// if any flag relating to the entity-relation filter is set, call the corresponding function (which checks if all necessary flags are set)
	to, _ := cmd.Flags().GetString("related-to")
	via, _ := cmd.Flags().GetString("related-via")
	tgt, _ := cmd.Flags().GetString("tgt-kind")
	if len(to) != 0 || len(via) != 0 || len(tgt) != 0 {
		return runGetEntityRelationCommand(cmd, args, to, via, tgt)
	}
	return runGetCommand(cmd, args, topoapi.Object_ENTITY)
}

func runGetEntityRelationCommand(cmd *cobra.Command, args []string, to string, via string, tgt string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")

	if len(to) > 0 && len(via) > 0 {
		outputWriter := cli.GetOutput()
		writer := new(tabwriter.Writer)
		writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)
		if len(tgt) == 0 {
			tgt = ""
		}
		objects, err := listObjects(cmd, &topoapi.Filters{
			RelationFilter: &topoapi.RelationFilter{
				SrcId:        to,
				RelationKind: via,
				TargetKind:   tgt,
			},
		})
		if err == nil {
			for _, object := range objects {
				printObject(writer, object, verbose)
			}
		}
		_ = writer.Flush()
		return nil
	}
	return errors.NewInvalid("missing related-to and/or related-via flags")
}

func runGetRelationCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topoapi.Object_RELATION)
}

func runGetKindCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topoapi.Object_KIND)
}

func runGetCommand(cmd *cobra.Command, args []string, objectType topoapi.Object_Type) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	verbose, _ := cmd.Flags().GetBool("verbose")

	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)
	if len(args) == 0 {
		filters := compileFilters(cmd, objectType)

		if !noHeaders {
			printHeader(writer, objectType, verbose, false)
		}

		objects, err := listObjects(cmd, filters)
		if err == nil {
			for _, object := range objects {
				printObject(writer, object, verbose)
			}
		}
	} else {
		id := args[0]
		object, err := getObject(cmd, topoapi.ID(id))
		if err != nil {
			return err
		}
		if object != nil {
			printObject(writer, *object, verbose)
		}
	}

	_ = writer.Flush()
	return nil
}

func listObjects(cmd *cobra.Command, filters *topoapi.Filters) ([]topoapi.Object, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	resp, err := client.List(context.Background(), &topoapi.ListRequest{Filters: filters})
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

func printHeader(writer io.Writer, objectType topoapi.Object_Type, verbose bool, printUpdateType bool) {
	if printUpdateType {
		_, _ = fmt.Fprintf(writer, "%s\t%s", "Update Type", "Object Type")
	}

	if objectType == topoapi.Object_ENTITY {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", "Entity ID", "Kind ID", "Labels")
	} else if objectType == topoapi.Object_RELATION {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s", "Relation ID", "Kind ID", "Source ID", "Target ID", "Labels")
	} else if objectType == topoapi.Object_KIND {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", "Kind ID", "Name", "Labels")
	} else {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", "ID", "Kind ID/Name", "Labels")
	}

	if !verbose {
		_, _ = fmt.Fprintf(writer, "\tAspects\n")
	} else {
		_, _ = fmt.Fprintf(writer, "\n")
	}
}

func printObject(writer io.Writer, object topoapi.Object, verbose bool) {
	labels := labelsAsCSV(object)
	switch object.Type {
	case topoapi.Object_ENTITY:
		var kindID topoapi.ID
		if e := object.GetEntity(); e != nil {
			kindID = e.KindID
		}
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", object.ID, kindID, labels)
		printAspects(writer, object, verbose)

	case topoapi.Object_RELATION:
		r := object.GetRelation()
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s", object.ID, r.KindID, r.SrcEntityID, r.TgtEntityID, labels)
		printAspects(writer, object, verbose)

	case topoapi.Object_KIND:
		k := object.GetKind()
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", object.ID, k.GetName(), labels)
		printAspects(writer, object, verbose)

	default:
		_, _ = fmt.Fprintf(writer, "\n")
	}
}

func labelsAsCSV(object topoapi.Object) string {
	var buffer bytes.Buffer
	first := true
	for k, v := range object.Labels {
		if first {
			buffer.WriteString(",")
		}
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(v)
		first = false
	}
	return buffer.String()
}

func printAspects(writer io.Writer, object topoapi.Object, verbose bool) {
	first := true
	if object.Aspects != nil {
		for aspectType, aspect := range object.Aspects {
			if verbose {
				if first {
					_, _ = fmt.Fprintf(writer, "\n")
				}
				_, _ = fmt.Fprintf(writer, "\t%s=%s\n", aspectType, bytes.NewBuffer(aspect.Value).String())
			} else {
				if !first {
					_, _ = fmt.Fprintf(writer, ",")
				} else {
					_, _ = fmt.Fprintf(writer, "\t")
				}
				_, _ = fmt.Fprintf(writer, "%s", aspectType)
			}
			first = false
		}
	}

	if object.Aspects == nil || !verbose {
		_, _ = fmt.Fprintf(writer, "\n")
	}
}
