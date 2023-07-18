// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

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
	cmd.Flags().String("related-to", "", "use relation filter")
	cmd.Flags().String("related-to-tgt", "", "use relation filter")
	cmd.Flags().String("related-via", "", "use relation filter, must also specify related-to or related-to-tgt")
	cmd.Flags().String("tgt-kind", "", "optional target kind for relation filter")
	cmd.Flags().StringSlice("with-aspect", nil, "aspect entity must have")
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
	cmd.Flags().StringSlice("with-aspect", nil, "aspect relation must have")
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
	cmd.Flags().String("sort-", "unordered", "sort order: ascending|descending|unordered(default)")
	cmd.Flags().StringSlice("with-aspect", nil, "aspect relation must have")
	return cmd
}

func getGetObjectsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "object id",
		Aliases: []string{"objects", "objs"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get Objects",
		RunE:    runGetObjectsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().String("kind", "", "kind query")
	cmd.Flags().String("label", "", "label query")
	cmd.Flags().String("related-to", "", "use relation filter")
	cmd.Flags().String("related-to-tgt", "", "use relation filter")
	cmd.Flags().String("related-via", "", "use relation filter, must also specify related-to or related-to-tgt")
	cmd.Flags().String("tgt-kind", "", "optional target kind for relation filter")
	cmd.Flags().String("scope", "target_only", "target_only|all|source_and_target")
	cmd.Flags().StringSlice("with-aspect", nil, "aspect object must have")
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

func runGetEntityRelationCommand(cmd *cobra.Command, _ []string, to string, toTgt, via string, tgt string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	scopeString, _ := cmd.Flags().GetString("scope")
	scope := topoapi.RelationFilterScope_TARGETS_ONLY
	aspects, _ := cmd.Flags().GetStringSlice("with-aspect")

	if scopeString == "source_and_target" {
		scope = topoapi.RelationFilterScope_SOURCE_AND_TARGETS
	}

	if len(to) > 0 || len(toTgt) > 0 {
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

		err := listObjects(cmd, &topoapi.Filters{RelationFilter: &filter, WithAspects: aspects},
			func(object *topoapi.Object) {
				printObject(writer, *object, verbose, false, false)
			})
		_ = writer.Flush()
		return err
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
	toTgt, _ := cmd.Flags().GetString("related-to-tgt")
	via, _ := cmd.Flags().GetString("related-via")
	tgt, _ := cmd.Flags().GetString("tgt-kind")

	if len(to) != 0 || len(via) != 0 || len(tgt) != 0 {
		return listObjectsRelations(cmd, to, toTgt, via, tgt)
	}
	return listAllObjectTypes(cmd, args)
}

func runGetCommand(cmd *cobra.Command, args []string, objectType topoapi.Object_Type) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	verbose, _ := cmd.Flags().GetBool("verbose")

	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', 0)
	var err error
	if len(args) == 0 {
		filters := compileFilters(cmd, objectType)

		if !noHeaders && !verbose {
			printHeader(writer, objectType, verbose, false)
		}

		err = listObjects(cmd, filters,
			func(object *topoapi.Object) {
				printObject(writer, *object, verbose, false, false)
			})
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
			printObject(writer, *object, verbose, false, true)
		}
	}

	_ = writer.Flush()
	return err
}

func listObjectsRelations(cmd *cobra.Command, to string, toTgt string, via string, tgt string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	scopeString, _ := cmd.Flags().GetString("scope")
	scope := topoapi.RelationFilterScope_TARGETS_ONLY
	aspects, _ := cmd.Flags().GetStringSlice("with-aspect")

	if scopeString == "all" {
		scope = topoapi.RelationFilterScope_ALL
	} else if scopeString == "source_and_target" {
		scope = topoapi.RelationFilterScope_SOURCE_AND_TARGETS
	} else if scopeString == "relations" {
		scope = topoapi.RelationFilterScope_RELATIONS_ONLY
	} else if scopeString == "relations_and_target" {
		scope = topoapi.RelationFilterScope_RELATIONS_AND_TARGETS
	}

	if len(to) > 0 || len(toTgt) > 0 {
		outputWriter := cli.GetOutput()
		writer := new(tabwriter.Writer)
		writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)
		if !noHeaders && !verbose {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "Object Type", "Object ID", "Kind ID", "Source ID", "Target ID", "Labels", "Aspects")
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
		err := listObjects(cmd, &topoapi.Filters{RelationFilter: &filter, WithAspects: aspects},
			func(object *topoapi.Object) {
				printObject(writer, *object, verbose, true, true)
			})
		_ = writer.Flush()
		return err
	}
	return errors.NewInvalid("missing related-to and/or related-via flags")
}

func listAllObjectTypes(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	verbose, _ := cmd.Flags().GetBool("verbose")

	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', 0)

	if !noHeaders && !verbose {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "Object Type", "Object ID", "Kind ID", "Source ID", "Target ID", "Labels", "Aspects")
	}
	if len(args) == 0 {
		filters := compileFilters(cmd, topoapi.Object_ENTITY)
		filters.ObjectTypes = []topoapi.Object_Type{topoapi.Object_ENTITY, topoapi.Object_RELATION, topoapi.Object_KIND}

		if err := listObjects(cmd, filters,
			func(object *topoapi.Object) {
				printObject(writer, *object, verbose, true, true)
			}); err != nil {
			return err
		}
	} else {
		object, err := getObject(cmd, topoapi.ID(args[0]))
		if err != nil {
			return err
		}
		printObject(writer, *object, verbose, true, true)
	}
	_ = writer.Flush()
	return nil
}

func listObjects(cmd *cobra.Command, filters *topoapi.Filters, processObject func(object *topoapi.Object)) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	stream, err := client.Query(context.Background(), &topoapi.QueryRequest{Filters: filters})
	if err != nil {
		return err
	}
	for {
		resp, err1 := stream.Recv()
		if err1 != nil {
			if err1 == io.EOF {
				return nil
			}
			return err1
		}
		processObject(resp.Object)
	}
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

func printObject(writer io.Writer, object topoapi.Object, verbose bool, printType bool, printSrcAndTarget bool) {
	labels := utils.None(labelsAsCSV(object))
	sourceID := utils.None("")
	targetID := utils.None("")

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
			if printSrcAndTarget {
				_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s", object.ID, kindID, sourceID, targetID, labels)
			} else {
				_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", object.ID, kindID, labels)
			}
		} else {
			_, _ = fmt.Fprintf(writer, "ID: %s\nKind ID: %s\nLabels: %s\n", object.ID, kindID, labels)
			if e := object.GetEntity(); e != nil {
				_, _ = fmt.Fprintf(writer, "Source of: ")
				for i, id := range e.SrcRelationIDs {
					if i == 0 {
						_, _ = fmt.Fprintf(writer, "%s", id)
					} else {
						_, _ = fmt.Fprintf(writer, ", %s", id)
					}
				}
				_, _ = fmt.Fprintf(writer, "\n")
				_, _ = fmt.Fprintf(writer, "Target of: ")
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
			if printSrcAndTarget {
				_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s", object.ID, k.GetName(), sourceID, targetID, labels)
			} else {
				_, _ = fmt.Fprintf(writer, "%s\t%s\t%s", object.ID, k.GetName(), labels)
			}
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
