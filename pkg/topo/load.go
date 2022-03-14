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
	"time"
)

func getLoadTopoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load [jsonFilePath|-]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Load topology resources in JSON format",
		RunE:  runLoadTopoCommand,
	}
	cmd.Flags().StringP("data", "d", "{}", "JSON data")
	return cmd
}

func runLoadTopoCommand(cmd *cobra.Command, args []string) error {
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

func loadFromBytes(cmd *cobra.Command, jsonData []byte) error {
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Iterate over the top-level objects and create either kind, entity or relation accordingly
	jsonObjects := data.(map[string]interface{})
	for k, v := range jsonObjects {
		object, err := parseObject(topoapi.ID(k), v)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(os.Stdout, "Creating %s...\n", object.ID)
		_, err = client.Create(ctx, &topoapi.CreateRequest{Object: object})
		if err != nil {
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
		if objectType == "kind" {
			return createKind(id, jsonObject), nil
		} else if objectType == "entity" {
			return createEntity(id, jsonObject), nil
		} else if objectType == "relation" {
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
