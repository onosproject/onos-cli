// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/onosproject/onos-api/go/onos/config/admin"
	"github.com/onosproject/onos-cli/pkg/format"
	"github.com/onosproject/onos-lib-go/pkg/cli"
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

func runListPluginsCommand(cmd *cobra.Command, _ []string) error {
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
			return tableFormat.Execute(cli.GetOutput(), !noHeaders, 0, allPlugins)
		}
		if err != nil {
			return err
		}
		allPlugins = append(allPlugins, in)
	}

}
