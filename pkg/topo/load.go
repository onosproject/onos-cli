// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/types"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

func getImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "import [jsonFilePath|-]",
		Aliases: []string{"load"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Import topology resources in JSON format",
		RunE:    runImportCommand,
	}
	cmd.Flags().StringP("data", "d", "{}", "JSON data")
	cmd.Flags().BoolP("ignore-errors", "i", false, "ignore errors and continue")
	return cmd
}

func runImportCommand(cmd *cobra.Command, args []string) error {
	var err error
	var data []byte

	if len(args) > 0 {
		if args[0] == "-" {
			data, err = ioutil.ReadAll(os.Stdin)
		} else {
			data, err = ioutil.ReadFile(args[0])
		}
	} else {
		sdata, _ := cmd.Flags().GetString("data")
		data = []byte(sdata)
	}

	if err != nil {
		return err
	}
	return loadFromBytes(cmd, data)
}

func getExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export [jsonFilePath|-]",
		Aliases: []string{"unload"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Export topology resources in JSON format",
		RunE:    runExportCommand,
	}
	cmd.Flags().StringP("data", "d", "{}", "JSON data")
	return cmd
}

func runExportCommand(cmd *cobra.Command, _ []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	resp, err := client.List(context.Background(), &topoapi.ListRequest{})
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	for _, o := range resp.Objects {
		d := make(map[string]interface{})
		switch o.Type {
		case topoapi.Object_ENTITY:
			e := o.GetEntity()
			d["type"] = "entity"
			d["kind"] = e.KindID
		case topoapi.Object_RELATION:
			r := o.GetRelation()
			d["type"] = "relation"
			d["source"] = r.SrcEntityID
			d["target"] = r.TgtEntityID
			d["kind"] = r.KindID
		case topoapi.Object_KIND:
			k := o.GetKind()
			d["name"] = k.Name
		}
		if err = exportAspects(d, o); err != nil {
			return err
		}
		data[string(o.ID)] = d
	}

	b, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("%s", string(b))
	return nil
}

func exportAspects(d map[string]interface{}, o topoapi.Object) error {
	for k, a := range o.Aspects {
		var ad interface{}
		err := json.Unmarshal(a.Value, &ad)
		if err != nil {
			return err
		}
		ajo := ad.(map[string]interface{})
		d[k] = ajo
	}
	return nil
}

func loadFromBytes(cmd *cobra.Command, jsonData []byte) error {
	ignoreErrors, _ := cmd.Flags().GetBool("ignore-errors")

	// Load the JSON data
	var data interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}

	// Get the topo client
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)
	ctx := context.Background()

	// Iterate over the top-level objects and create either kind, entity or relation accordingly
	jsonObjects := data.(map[string]interface{})
	relations := make([]*topoapi.Object, 0)
	for k, v := range jsonObjects {
		object, err := parseObject(topoapi.ID(k), v)
		if err != nil {
			return err
		}
		if object.Type == topoapi.Object_RELATION {
			// Stow away relations for last
			relations = append(relations, object)
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "Creating %s...\n", object.ID)
			_, err = client.Create(ctx, &topoapi.CreateRequest{Object: object})
			if !ignoreErrors && err != nil {
				return err
			}
		}
	}

	// Create all relations that have been recorded
	for _, object := range relations {
		_, _ = fmt.Fprintf(os.Stdout, "Creating %s...\n", object.ID)
		_, err = client.Create(ctx, &topoapi.CreateRequest{Object: object})
		if !ignoreErrors && err != nil {
			return err
		}
	}

	return nil
}

func parseObject(id topoapi.ID, v interface{}) (*topoapi.Object, error) {
	switch v.(type) {
	case map[string]interface{}:
		jsonObject := v.(map[string]interface{})
		objectType := jsonObject["type"].(string)
		switch objectType {
		case "kind":
			return createKind(id, jsonObject), nil
		case "entity":
			return createEntity(id, jsonObject), nil
		case "relation":
			return createRelation(id, jsonObject), nil
		}
	}
	return nil, errors.New("invalid json")
}

func getString(k string, jsonObject map[string]interface{}) string {
	v := jsonObject[k]
	if v == nil {
		return ""
	}
	return v.(string)
}

func createKind(id topoapi.ID, jsonObject map[string]interface{}) *topoapi.Object {
	name := getString("name", jsonObject)
	return &topoapi.Object{
		ID:   id,
		Type: topoapi.Object_KIND,
		Obj: &topoapi.Object_Kind{
			Kind: &topoapi.Kind{Name: name},
		},
		Aspects: getAspects(jsonObject),
		Labels:  getLabels(jsonObject),
	}
}

func createEntity(id topoapi.ID, jsonObject map[string]interface{}) *topoapi.Object {
	kindID := getString("kind", jsonObject)
	return &topoapi.Object{
		ID:   id,
		Type: topoapi.Object_ENTITY,
		Obj: &topoapi.Object_Entity{
			Entity: &topoapi.Entity{
				KindID: topoapi.ID(kindID),
			},
		},
		Aspects: getAspects(jsonObject),
		Labels:  getLabels(jsonObject),
	}
}

func createRelation(id topoapi.ID, jsonObject map[string]interface{}) *topoapi.Object {
	kindID := getString("kind", jsonObject)
	srcID := getString("source", jsonObject)
	tgtID := getString("target", jsonObject)
	return &topoapi.Object{
		ID:   id,
		Type: topoapi.Object_RELATION,
		Obj: &topoapi.Object_Relation{
			Relation: &topoapi.Relation{
				KindID:      topoapi.ID(kindID),
				SrcEntityID: topoapi.ID(srcID),
				TgtEntityID: topoapi.ID(tgtID),
			},
		},
		Aspects: getAspects(jsonObject),
		Labels:  getLabels(jsonObject),
	}
}

func getAspects(jsonObject map[string]interface{}) map[string]*types.Any {
	aspects := map[string]*types.Any{}
	for k, v := range jsonObject {
		// Any key with a "." in it is considered an aspect type
		if strings.Contains(k, ".") {
			aspectValue, err := json.Marshal(v)
			if err == nil {
				aspects[k] = &types.Any{
					TypeUrl: k,
					Value:   aspectValue,
				}
			}
		}
	}
	return aspects
}

func getLabels(jsonObject map[string]interface{}) map[string]string {
	labels := map[string]string{}
	jsonLabels := jsonObject["labels"]
	if jsonLabels != nil {
		labelsObject := jsonLabels.(map[string]interface{})
		for k, v := range labelsObject {
			labels[k] = v.(string)
		}
	}
	return labels
}
