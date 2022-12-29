// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getUpdateEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity <id> [args]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Update Entity",
		RunE:  runUpdateEntityCommand,
	}
	cmd.Flags().StringToStringP("aspect", "a", map[string]string{}, "aspect of this entity")
	cmd.Flags().StringToStringP("label", "l", map[string]string{}, "classification label")
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
	cmd.Flags().StringToStringP("label", "l", map[string]string{}, "classification label")
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
	cmd.Flags().StringToStringP("label", "l", map[string]string{}, "classification label")
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

const deleteKeyword = "--delete"

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

	// Apply label changes
	labels, err := cmd.Flags().GetStringToString("label")
	if err != nil {
		return err
	}

	if object.Labels == nil && len(labels) > 0 {
		object.Labels = make(map[string]string)
	}
	for labelKey, labelValue := range labels {
		if len(labelValue) > 0 && labelValue != deleteKeyword {
			object.Labels[labelKey] = labelValue
		} else {
			delete(object.Labels, labelKey)
		}
	}

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
			if len(aspectValue) > 0 && aspectValue != deleteKeyword {
				if object.Aspects == nil {
					object.Aspects = make(map[string]*types.Any)
				}
				object.Aspects[aspectType] = &types.Any{
					TypeUrl: aspectType,
					Value:   []byte(aspectValue),
				}
			} else {
				delete(object.Aspects, aspectType)
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
