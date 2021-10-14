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

	"github.com/onosproject/onos-cli/pkg/utils"

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
	cmd.Flags().String("related-to-tgt", "", "use relation filter, must also specify related-via")
	cmd.Flags().String("related-via", "", "use relation filter, must also specify related-to or related-to-tgt")
	cmd.Flags().String("tgt-kind", "", "optional target kind for relation filter")
	cmd.Flags().String("sort-order", "unordered", "sort order: ascending|descending|unordered(default)")
	cmd.Flags().String("scope", "target_only", "target_only|source_and_target")
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
	cmd.Flags().String("sort-order", "unordered", "sort order: ascending|descending|unordered(default)")
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
	cmd.Flags().String("sort-order", "unordered", "sort order: ascending|descending|unordered(default)")
	return cmd
}

func getGetObjectsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "objects id",
		Aliases: []string{"objs"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get Objects",
		RunE:    runGetObjectsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().String("kind", "", "kind query")
	cmd.Flags().String("label", "", "label query")
	cmd.Flags().String("related-to", "", "use relation filter, must also specify related-via")
	cmd.Flags().String("related-to-tgt", "", "use relation filter, must also specify related-via")
	cmd.Flags().String("related-via", "", "use relation filter, must also specify related-to or related-to-tgt")
	cmd.Flags().String("tgt-kind", "", "optional target kind for relation filter")
	cmd.Flags().String("sort-order", "unordered", "sort order: ascending|descending|unordered(default)")
	cmd.Flags().String("scope", "target_only", "target_only|all|source_and_target")
	return cmd
}

func runGetEntityCommand(cmd *cobra.Command, args []string) error {
	// if any flag relating to the entity-relation filter is set, call the corresponding function (which checks if all necessary flags are set)
	to, _ := cmd.Flags().GetString("related-to")
	toTgt, _ := cmd.Flags().GetString("related-to-tgt")
	via, _ := cmd.Flags().GetString("related-via")
	tgt, _ := cmd.Flags().GetString("tgt-kind")

	if len(to) != 0 && len(toTgt) != 0 {
		return errors.NewInvalid("only 'related-to' or 'related-to-tgt' flag can be specified; not both")
	}
	if len(to) != 0 || len(toTgt) != 0 || len(via) != 0 || len(tgt) != 0 {
		return runGetEntityRelationCommand(cmd, args, to, toTgt, via, tgt)
	}
	return runGetCommand(cmd, args, topoapi.Object_ENTITY)
}

func runGetEntityRelationCommand(cmd *cobra.Command, args []string, to string, toTgt, via string, tgt string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	scopeString, _ := cmd.Flags().GetString("scope")
	scope := topoapi.RelationFilterScope_TARGET_ONLY

	if scopeString == "source_and_target" {
		scope = topoapi.RelationFilterScope_SOURCE_AND_TARGET
	}

	if (len(to) > 0 || len(toTgt) > 0) && len(via) > 0 {
		outputWriter := cli.GetOutput()
		writer := new(tabwriter.Writer)
		writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)
		if !noHeaders {
			printHeader(writer, topoapi.Object_ENTITY, verbose, false)
		}
		if len(tgt) == 0 {
			tgt = ""
		}

		filter := topoapi.RelationFilter{
			RelationKind: via,
			TargetKind:   tgt,
			Scope:        scope,
		}
		if len(to) > 0 {
			filter.SrcId = to
		} else {
			filter.TargetId = toTgt
		}

		objects, err := listObjects(cmd, &topoapi.Filters{RelationFilter: &filter}, topoapi.SortOrder_UNORDERED)
		if err == nil {
			for _, object := range objects {
				printObject(writer, object, verbose, false)
			}
		}
		_ = writer.Flush()
		return nil
	}
	return errors.NewInvalid("missing 'related-to', 'related-to-tgt' and/or 'related-via' flags")
}

func runGetRelationCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topoapi.Object_RELATION)
}

func runGetKindCommand(cmd *cobra.Command, args []string) error {
	return runGetCommand(cmd, args, topoapi.Object_KIND)
}

func runGetObjectsCommand(cmd *cobra.Command, args []string) error {
	// if any flag relating to the entity-relation filter is set, call the corresponding function (which checks if all necessary flags are set)
	to, _ := cmd.Flags().GetString("related-to")
	via, _ := cmd.Flags().GetString("related-via")
	tgt, _ := cmd.Flags().GetString("tgt-kind")

	if len(to) != 0 || len(via) != 0 || len(tgt) != 0 {
		return listObjectsRelations(cmd, to, via, tgt)
	}
	return listAllObjectTypes(cmd, args)
}

func runGetCommand(cmd *cobra.Command, args []string, objectType topoapi.Object_Type) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	verbose, _ := cmd.Flags().GetBool("verbose")
	sortString, _ := cmd.Flags().GetString("sort-order")
	// sort order selection
	sortOrder := topoapi.SortOrder_UNORDERED
	if sortString == "ascending" {
		sortOrder = topoapi.SortOrder_ASCENDING
	} else if sortString == "descending" {
		sortOrder = topoapi.SortOrder_DESCENDING
	} else if sortString == "unordered" {
		sortOrder = topoapi.SortOrder_UNORDERED
	}

	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', 0)
	if len(args) == 0 {
		filters := compileFilters(cmd, objectType)

		if !noHeaders && !verbose {
			printHeader(writer, objectType, verbose, false)
		}

		objects, err := listObjects(cmd, filters, sortOrder)
		if err == nil {
			for _, object := range objects {
				printObject(writer, object, verbose, false)
			}
		}
	} else {
		id := args[0]
		object, err := getObject(cmd, topoapi.ID(id))
		if !noHeaders {
			printHeader(writer, objectType, verbose, false)
		}
		if err != nil {
			return err
		}
		if object != nil {
			printObject(writer, *object, verbose, false)
		}
	}

	_ = writer.Flush()
	return nil
}

func listObjectsRelations(cmd *cobra.Command, to string, via string, tgt string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	scopeString, _ := cmd.Flags().GetString("scope")
	scope := topoapi.RelationFilterScope_TARGET_ONLY

	if scopeString == "all" {
		scope = topoapi.RelationFilterScope_ALL
	} else if scopeString == "source_and_target" {
		scope = topoapi.RelationFilterScope_SOURCE_AND_TARGET
	}

	if len(to) > 0 && len(via) > 0 {
		outputWriter := cli.GetOutput()
		writer := new(tabwriter.Writer)
		writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)
		if !noHeaders {
			printHeader(writer, topoapi.Object_RELATION, verbose, false)
		}
		if len(tgt) == 0 {
			tgt = ""
		}

		objects, err := listObjects(cmd, &topoapi.Filters{
			RelationFilter: &topoapi.RelationFilter{
				SrcId:        to,
				RelationKind: via,
				TargetKind:   tgt,
				Scope:        scope,
			},
		}, topoapi.SortOrder_UNORDERED)
		if err == nil {
			for _, object := range objects {
				printObject(writer, object, verbose, false)
			}
		}
		_ = writer.Flush()
		return nil
	}
	return errors.NewInvalid("missing related-to and/or related-via flags")
}

func listAllObjectTypes(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	verbose, _ := cmd.Flags().GetBool("verbose")
	sortString, _ := cmd.Flags().GetString("sort-order")
	sortOrder := topoapi.SortOrder_UNORDERED
	if sortString == "ascending" {
		sortOrder = topoapi.SortOrder_ASCENDING
	} else if sortString == "descending" {
		sortOrder = topoapi.SortOrder_DESCENDING
	} else if sortString == "unordered" {
		sortOrder = topoapi.SortOrder_UNORDERED
	}

	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', 0)
	objects := make([]topoapi.Object, 0)
	if len(args) == 0 {
		filters := compileFilters(cmd, topoapi.Object_ENTITY)
		filters.ObjectTypes = []topoapi.Object_Type{topoapi.Object_ENTITY, topoapi.Object_RELATION, topoapi.Object_KIND}

		listedObjects, err := listObjects(cmd, filters, sortOrder)
		if err != nil {
			return err
		}
		objects = append(objects, listedObjects...)
	} else {
		id := args[0]
		object, err := getObject(cmd, topoapi.ID(id))
		if err != nil {
			return err
		}
		objects = append(objects, *object)
	}

	if !noHeaders && !verbose {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "Object Type", "Object ID", "Kind ID", "Source ID", "Target ID", "Labels", "Aspects")
	}
	for _, object := range objects {
		printObject(writer, object, verbose, true)
	}
	_ = writer.Flush()
	return nil
}

func listObjects(cmd *cobra.Command, filters *topoapi.Filters, order topoapi.SortOrder) ([]topoapi.Object, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	resp, err := client.List(context.Background(), &topoapi.ListRequest{Filters: filters, SortOrder: order})
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
		_, _ = fmt.Fprintf(writer, "%s\t%s\t", "Update Type", "Object Type")
	}

	if !verbose {
		if objectType == topoapi.Object_ENTITY {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", "Entity ID", "Kind ID", "Labels")
		} else if objectType == topoapi.Object_RELATION {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s", "Relation ID", "Kind ID", "Source ID", "Target ID", "Labels")
		} else if objectType == topoapi.Object_KIND {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", "Kind ID", "Name", "Labels")
		} else {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", "ID", "Kind ID/Name", "Labels")
		}
		_, _ = fmt.Fprintf(writer, "\tAspects")
	}
	_, _ = fmt.Fprintf(writer, "\n")
}

func printObject(writer io.Writer, object topoapi.Object, verbose bool, printType bool) {
	labels := utils.None(labelsAsCSV(object))

	if printType {
		if verbose {
			_, _ = fmt.Fprintf(writer, "Object Type: %s\n", object.Type)
		} else {
			_, _ = fmt.Fprintf(writer, "%s\t", object.Type)
		}
	}

	switch object.Type {
	case topoapi.Object_ENTITY:
		var kindID topoapi.ID
		if e := object.GetEntity(); e != nil {
			kindID = e.KindID
		}
		if !verbose {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", object.ID, kindID, labels)
		} else {
			_, _ = fmt.Fprintf(writer, "ID: %s\nKind ID: %s\nLabels: %s\n", object.ID, kindID, labels)
			if e := object.GetEntity(); e != nil {
				_, _ = fmt.Fprintf(writer, "Source Id's: ")
				for i, id := range e.SrcRelationIDs {
					if i == 0 {
						_, _ = fmt.Fprintf(writer, "%s", id)
					} else {
						_, _ = fmt.Fprintf(writer, ", %s", id)
					}
				}
				_, _ = fmt.Fprintf(writer, "\n")
				_, _ = fmt.Fprintf(writer, "Target Id's: ")
				for i, id := range e.TgtRelationIDs {
					if i == 0 {
						_, _ = fmt.Fprintf(writer, "%s", id)
					} else {
						_, _ = fmt.Fprintf(writer, ", %s", id)
					}
				}
				_, _ = fmt.Fprintf(writer, "\n")
			}
		}

		printAspects(writer, object, verbose)

	case topoapi.Object_RELATION:
		r := object.GetRelation()
		if !verbose {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s", object.ID, r.KindID, r.SrcEntityID, r.TgtEntityID, labels)
		} else {
			_, _ = fmt.Fprintf(writer, "ID:\t%s\nKind ID:\t%s\nSource Entity ID:\t%s\nTarget Entity ID:\t%s\nLabels:\t%s\n", object.ID, r.KindID, r.SrcEntityID, r.TgtEntityID, labels)
		}
		printAspects(writer, object, verbose)

	case topoapi.Object_KIND:
		k := object.GetKind()
		if !verbose {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", object.ID, k.GetName(), labels)
		} else {
			_, _ = fmt.Fprintf(writer, "ID:\t%s\nName:\t%s\nLabels:\t%s\n", object.ID, k.GetName(), labels)
		}
		printAspects(writer, object, verbose)

	default:
		_, _ = fmt.Fprintf(writer, "\n")
	}
}

func labelsAsCSV(object topoapi.Object) string {
	var buffer bytes.Buffer
	first := true
	for k, v := range object.Labels {
		if !first {
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
	if verbose {
		_, _ = fmt.Fprintf(writer, "Aspects:\n")
	}
	if object.Aspects != nil {
		for aspectType, aspect := range object.Aspects {
			if verbose {
				_, _ = fmt.Fprintf(writer, "- %s=%s\n", aspectType, aspect.Value)
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

	if object.Aspects == nil {
		_, _ = fmt.Fprintf(writer, "\t%s", utils.None(""))
	}

	_, _ = fmt.Fprintf(writer, "\n")
}
