// Copyright 2022-present Open Networking Foundation.
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

package config

import (
	"github.com/onosproject/onos-api/go/onos/config/admin"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/opencord/cordctl/pkg/format" // NOTE cordctl is not really maintained anymore, consider importing this code
	"github.com/spf13/cobra"
	"io"
)

const pluginListTemplate = "table{{.Id}}\t{{.Status}}\t{{.Endpoint}}\t{{.Info.Name}}\t{{.Info.Version}}\t{{.Error}}"
const pluginListTemplateVerbose = "table{{.Id}}\t{{.Status}}\t{{.Endpoint}}\t{{.Info.Name}}\t{{.Info.Version}}\t{{.Error}}\t{{.Info.ModelData}}"

func getListPluginsCommand() *cobra.Command {
	// TODO support model filtering
	cmd := &cobra.Command{
		Use:   "plugins",
		Short: "plugins",
		RunE:  runListPluginsCommand,
	}

	cmd.Flags().BoolP("verbose", "v", false, "prints all the models in a plugin")
	cmd.Flags().Bool("no-headers", false, "disables output headers")

	return cmd
}

func runListPluginsCommand(cmd *cobra.Command, args []string) error {
	connection, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	client := admin.CreateConfigAdminServiceClient(connection)

	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	req := admin.ListModelsRequest{
		Verbose: verbose,
	}

	ctx := cli.NewContextWithAuthHeaderFromFlag(cmd.Context(), cmd.Flag(cli.AuthHeaderFlag))
	stream, err := client.ListRegisteredModels(ctx, &req)
	if err != nil {
		return err
	}

	var tableFormat format.Format
	if verbose {
		tableFormat = pluginListTemplateVerbose
	} else {
		tableFormat = pluginListTemplate
	}

	allPlugins := []*admin.ModelPlugin{}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			if e := tableFormat.Execute(cli.GetOutput(), !noHeaders, allPlugins); e != nil {
				return e
			}
			return nil
		}
		if err != nil {
			return err
		}
		allPlugins = append(allPlugins, in)
	}

}
