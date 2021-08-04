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
	"time"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/spf13/cobra"
)

func getWipeoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wipeout please",
		Args:  cobra.ExactArgs(1),
		Short: "Delete All Relations and Entities",
		RunE:  runWipeoutCommand,
	}
	cmd.Flags().Bool("include-kinds", false, "delete kinds as well as entities and relations")
	return cmd
}

func runWipeoutCommand(cmd *cobra.Command, args []string) error {
	includeKinds, _ := cmd.Flags().GetBool("include-kinds")
	if args[0] != "please" {
		return errors.NewInvalid("Wipeout requires the string 'please'")
	}

	filters := topoapi.Filters{
		ObjectTypes: []topoapi.Object_Type{topoapi.Object_ENTITY, topoapi.Object_RELATION},
	}
	if includeKinds {
		filters.ObjectTypes = append(filters.ObjectTypes, topoapi.Object_KIND)
	}

	objects, err := listObjects(cmd, &filters, topoapi.SortOrder_UNORDERED)
	if err != nil {
		return err
	}
	for _, object := range objects {
		err = deleteObject(cmd, object)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteObject(cmd *cobra.Command, object topoapi.Object) error {
	id := object.ID

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Delete(ctx, &topoapi.DeleteRequest{ID: id})
	if err != nil {
		return err
	}
	cli.Output("Deleted %s %s\n", object.Type, id)
	return nil
}
