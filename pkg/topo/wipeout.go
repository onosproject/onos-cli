// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

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

	// delete relations first, to avoid an error where relations are already deleted

	err := listObjects(cmd, &topoapi.Filters{ObjectTypes: []topoapi.Object_Type{topoapi.Object_RELATION}},
		func(object *topoapi.Object) {
			_ = deleteObject(cmd, *object)
		})
	if err != nil {
		return err
	}

	filters := topoapi.Filters{
		ObjectTypes: []topoapi.Object_Type{topoapi.Object_ENTITY},
	}
	if includeKinds {
		filters.ObjectTypes = append(filters.ObjectTypes, topoapi.Object_KIND)
	}

	err = listObjects(cmd, &filters,
		func(object *topoapi.Object) {
			_ = deleteObject(cmd, *object)
		})
	return err
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
