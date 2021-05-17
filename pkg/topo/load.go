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
	"encoding/json"
	"errors"
	"github.com/gogo/protobuf/types"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

func getLoadTopoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load [jsonFilePath|-]",
		Args:  cobra.ExactArgs(0),
		Short: "Load topology resources in JSON format",
		RunE:  runLoadTopoCommand,
	}
	cmd.Flags().StringP("data", "d", "{}", "JSON data")
	return cmd
}

func runLoadTopoCommand(cmd *cobra.Command, args []string) error {
	sdata, err := cmd.Flags().GetString("data")
	if err != nil {
		return err
	}

	// Load the JSON data
	var data interface{}
	err = json.Unmarshal([]byte(sdata), &data)
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

func createKind(id topoapi.ID, jsonObject map[string]interface{}) *topoapi.Object {
	name := jsonObject["name"].(string)
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
	kindID := jsonObject["kind"].(string)
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
	kindID := jsonObject["kind"].(string)
	srcID := jsonObject["source"].(string)
	tgtID := jsonObject["target"].(string)
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

func getLabels(jsonObject map[string]interface{}) []string {
	labels := make([]string, 0, 0)
	jsonLabels := jsonObject["labels"]
	if jsonLabels != nil {
		for _, l := range jsonLabels.([]interface{}) {
			labels = append(labels, l.(string))
		}
	}
	return labels
}
