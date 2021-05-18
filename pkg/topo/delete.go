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

func getDeleteEntityCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "entity <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete Entity",
		RunE:  runDeleteEntityCommand,
	}
}

func getDeleteRelationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "object <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete Relation",
		RunE:  runDeleteRelationCommand,
	}
}

func getDeleteKindCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "kind <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete kind",
		RunE:  runDeleteKindCommand,
	}
}

func runDeleteEntityCommand(cmd *cobra.Command, args []string) error {
	return runDeleteObjectCommand(cmd, args, "entity")
}

func runDeleteRelationCommand(cmd *cobra.Command, args []string) error {
	return runDeleteObjectCommand(cmd, args, "relation")
}

func runDeleteKindCommand(cmd *cobra.Command, args []string) error {
	return runDeleteObjectCommand(cmd, args, "kind")
}

func runDeleteObjectCommand(cmd *cobra.Command, args []string, typeName string) error {
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
	cli.Output("Deleted %s %s", typeName, id)
	return nil
}
