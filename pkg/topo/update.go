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
	"github.com/gogo/protobuf/types"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"time"
)

func getUpdateEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity <id> [args]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Update Entity",
		RunE:  runUpdateEntityCommand,
	}
	cmd.Flags().StringToStringP("aspect", "a", map[string]string{}, "aspect of this entity")
	cmd.Flags().StringToStringP("set", "s", map[string]string{}, "set single attribute of an aspect")
	return cmd
}

func getUpdateRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <id> [args]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Update Relation",
		RunE:  runUpdateRelationCommand,
	}
	cmd.Flags().StringToStringP("aspect", "a", map[string]string{}, "aspect of this entity")
	cmd.Flags().StringToStringP("set", "s", map[string]string{}, "set single attribute of an aspect")
	return cmd
}

func getUpdateKindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kind <id> [args]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Update Kind",
		RunE:  runUpdateKindCommand,
	}
	cmd.Flags().StringP("name", "n", "", "Kind Name")
	cmd.Flags().StringToStringP("aspect", "a", map[string]string{}, "aspect of this entity")
	cmd.Flags().StringToStringP("set", "s", map[string]string{}, "set single attribute of an aspect")
	return cmd
}

func runUpdateEntityCommand(cmd *cobra.Command, args []string) error {
	return updateObject(cmd, args, topoapi.Object_ENTITY)
}

func runUpdateRelationCommand(cmd *cobra.Command, args []string) error {
	return updateObject(cmd, args, topoapi.Object_RELATION)
}

func runUpdateKindCommand(cmd *cobra.Command, args []string) error {
	return updateObject(cmd, args, topoapi.Object_KIND)
}

func updateObject(cmd *cobra.Command, args []string, objectType topoapi.Object_Type) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Fetch the object
	id := topoapi.ID(args[0])
	response, err := client.Get(ctx, &topoapi.GetRequest{ID: id})
	if err != nil {
		return err
	}

	object := response.Object

	// If kind, change name if needed
	if objectType == topoapi.Object_KIND {
		if cmd.Flags().Changed("name") {
			name, err := cmd.Flags().GetString("name")
			if err == nil {
				object.GetKind().Name = name
			}
		}
	}

	// Apply changed aspects
	aspects, err := cmd.Flags().GetStringToString("aspect")
	if err == nil {
		for aspectType, aspectValue := range aspects {
			object.Aspects[aspectType] = &types.Any{
				TypeUrl: aspectType,
				Value:   bytes.NewBufferString(aspectValue).Bytes(),
			}
		}
	}

	// TODO: Apply individual aspect attribute changes

	// Update the object
	_, err = client.Update(ctx, &topoapi.UpdateRequest{Object: object})
	if err != nil {
		return err
	}
	return nil
}
