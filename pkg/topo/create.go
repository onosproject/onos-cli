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

func getCreateEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity <id> [args]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Create Entity",
		RunE:  runCreateEntityCommand,
	}
	cmd.Flags().StringP("kind", "k", "", "Kind ID")
	cmd.Flags().StringToStringP("aspect", "a", map[string]string{}, "aspect of this entity")
	cmd.Flags().StringToStringP("label", "l", map[string]string{}, "classification label")
	return cmd
}

func getCreateRelationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relation <src-entity-id> <tgt-entity-id> [args]",
		Args:  cobra.MinimumNArgs(2),
		Short: "Create Relation",
		RunE:  runCreateRelationCommand,
	}
	cmd.Flags().StringP("kind", "k", "", "Kind ID")
	cmd.Flags().StringToStringP("aspect", "a", map[string]string{}, "aspect of this relation")
	cmd.Flags().StringToStringP("label", "l", map[string]string{}, "classification label")
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
	cmd.Flags().StringToStringP("label", "l", map[string]string{}, "classification label")
	return cmd
}

func runCreateEntityCommand(cmd *cobra.Command, args []string) error {
	kindID, _ := cmd.Flags().GetString("kind")
	return createObject(topoapi.NewEntity(topoapi.ID(args[0]), topoapi.ID(kindID)), cmd)
}

func runCreateRelationCommand(cmd *cobra.Command, args []string) error {
	kindID, _ := cmd.Flags().GetString("kind")
	return createObject(topoapi.NewRelation(topoapi.ID(args[0]), topoapi.ID(args[2]), topoapi.ID(kindID)), cmd)
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
	labels, err := cmd.Flags().GetStringToString("label")
	if err != nil {
		return err
	}
	object.Labels = labels

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
