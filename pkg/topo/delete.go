// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"context"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"time"
)

func getDeleteEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete Entity",
		RunE:  runDeleteEntityCommand,
	}
	cmd.Flags().Uint64P("revision", "r", 0, "revision to delete")
	return cmd
}

func getDeleteRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete Relation",
		RunE:  runDeleteRelationCommand,
	}
	cmd.Flags().Uint64P("revision", "r", 0, "revision to delete")
	return cmd
}

func getDeleteKindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kind <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete kind",
		RunE:  runDeleteKindCommand,
	}
	cmd.Flags().Uint64P("revision", "r", 0, "revision to delete")
	return cmd
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
	revision, _ := cmd.Flags().GetUint64("revision")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Delete(ctx, &topoapi.DeleteRequest{ID: topoapi.ID(id), Revision: topoapi.Revision(revision)})
	if err != nil {
		return err
	}
	cli.Output("Deleted %s %s\n", typeName, id)
	return nil
}
