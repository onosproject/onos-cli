// Copyright 2021-present Open Networking Foundation.
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
	"github.com/spf13/cobra"
)

func getLoadYamlEntitiesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load {file}",
		Args:  cobra.ExactArgs(1),
		Short: "Load topo data from a YAML file",
		RunE:  runLoadYamlEntitiesCommand,
	}
	cmd.Flags().StringArray("attr", []string{""}, "Extra attributes to add to each device in k=v format")
	return cmd
}

func runLoadYamlEntitiesCommand(cmd *cobra.Command, args []string) error {
	/*
		var filename string
		if len(args) > 0 {
			filename = args[0]
		}

		extraAttrs, err := cmd.Flags().GetStringArray("attr")
		if err != nil {
			return err
		}
		for _, x := range extraAttrs {
			cli.Output("runLoadYamlEntitiesCommand %v", x)
			split := strings.Split(x, "=")
			if len(split) != 2 {
				return fmt.Errorf("expect extra args to be in the format a=b. Rejected: %s", x)
			}
		}

		topoConfig, err := load.GetTopoConfig(filename)
		if err != nil {
			return err
		}

		conn, err := cli.GetConnection(cmd)
		if err != nil {
			return err
		}
		defer conn.Close()
		client := topoapi.CreateTopoClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		for _, kind := range topoConfig.TopoKinds {
			if kind.Attributes == nil {
				a := make(map[string]string)
				kind.Attributes = &a
			}
			for _, x := range extraAttrs {
				split := strings.Split(x, "=")
				(*kind.Attributes)[split[0]] = split[1]
			}

			kind := kind // pin

			object := load.TopoKindToTopoObject(&kind)
			_, err = client.Create(ctx, &topoapi.CreateRequest{Object: object})
			if err != nil {
				return err
			}
		}

		for _, entity := range topoConfig.TopoEntities {
			if entity.Attributes == nil {
				a := make(map[string]string)
				entity.Attributes = &a
			}
			for _, x := range extraAttrs {
				split := strings.Split(x, "=")
				(*entity.Attributes)[split[0]] = split[1]
			}

			entity := entity // pin
			object := load.TopoEntityToTopoObject(&entity)
			_, err = client.Create(ctx, &topoapi.CreateRequest{Object: object})
			if err != nil {
				return err
			}
		}

		for _, relation := range topoConfig.TopoRelations {
			if relation.Attributes == nil {
				a := make(map[string]string)
				relation.Attributes = &a
			}
			for _, x := range extraAttrs {
				split := strings.Split(x, "=")
				(*relation.Attributes)[split[0]] = split[1]
			}

			relation := relation // pin
			object := load.TopoRelationToTopoObject(&relation)
			_, err = client.Create(ctx, &topoapi.CreateRequest{Object: object})
			if err != nil {
				return err
			}
		}

		for _, relation := range topoConfig.TopoRelations {
			if relation.Attributes == nil {
				a := make(map[string]string)
				relation.Attributes = &a
			}
			for _, x := range extraAttrs {
				split := strings.Split(x, "=")
				(*relation.Attributes)[split[0]] = split[1]
			}

			relation := relation // pin
			object := load.TopoRelationToTopoObject(&relation)
			_, err = client.Create(ctx, &topoapi.CreateRequest{Object: object})
			if err != nil {
				return err
			}
		}

		fmt.Printf("Loaded %d topo devices from %s\n", len(topoConfig.TopoEntities), filename)
	*/

	return nil
}
