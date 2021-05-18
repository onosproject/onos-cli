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
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"time"
)

func getCreateEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity <id> [args]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Create Entity",
		RunE:  runCreateEntityCommand,
	}
	cmd.Flags().StringP("kind", "k", "", "Kind ID")
	//_ = cmd.MarkFlagRequired("kind")
	cmd.Flags().StringToStringP("aspect", "a", map[string]string{}, "aspect of this entity")
	return cmd
}

func getCreateRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <id> <src-entity-id> <tgt-entity-id> [args]",
		Args:  cobra.MinimumNArgs(3),
		Short: "Create Relation",
		RunE:  runCreateRelationCommand,
	}
	cmd.Flags().StringP("kind", "k", "", "Kind ID")
	//_ = cmd.MarkFlagRequired("kind")
	cmd.Flags().StringToStringP("aspect", "a", map[string]string{}, "aspect of this relation")
	return cmd
}

func getCreateKindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kind <id> <name> [args]",
		Args:  cobra.MinimumNArgs(2),
		Short: "Create Kind",
		RunE:  runCreateKindCommand,
	}
	cmd.Flags().StringToStringP("aspect", "a", map[string]string{}, "default aspect for entities of this kind")
	return cmd
}

func runCreateEntityCommand(cmd *cobra.Command, args []string) error {
	kindID, _ := cmd.Flags().GetString("kind")
	return createObject(&topoapi.Object{
		ID:   topoapi.ID(args[0]),
		Type: topoapi.Object_ENTITY,
		Obj: &topoapi.Object_Entity{
			Entity: &topoapi.Entity{
				KindID: topoapi.ID(kindID),
			},
		},
	}, cmd)
}

func runCreateRelationCommand(cmd *cobra.Command, args []string) error {
	kindID, _ := cmd.Flags().GetString("kind")
	return createObject(&topoapi.Object{
		ID:   topoapi.ID(args[0]),
		Type: topoapi.Object_RELATION,
		Obj: &topoapi.Object_Relation{
			Relation: &topoapi.Relation{
				KindID:      topoapi.ID(kindID),
				SrcEntityID: topoapi.ID(args[1]),
				TgtEntityID: topoapi.ID(args[2]),
			},
		},
	}, cmd)
}

func runCreateKindCommand(cmd *cobra.Command, args []string) error {
	return createObject(&topoapi.Object{
		ID:   topoapi.ID(args[0]),
		Type: topoapi.Object_KIND,
		Obj: &topoapi.Object_Kind{
			Kind: &topoapi.Kind{
				Name: args[1],
			},
		},
	}, cmd)
}

func createObject(object *topoapi.Object, cmd *cobra.Command) error {
	aspects, err := cmd.Flags().GetStringToString("aspect")
	if err != nil {
		return err
	}

	// Apply all aspect values
	for aspectType, aspectValue := range aspects {
		err := object.SetAspectBytes(aspectType, []byte(aspectValue))
		if err != nil {
			return err
		}
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Create(ctx, &topoapi.CreateRequest{Object: object})
	if err != nil {
		return err
	}
	return nil
}
