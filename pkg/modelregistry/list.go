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

package modelregistry

import (
	"github.com/onosproject/onos-api/go/onos/configmodel"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"text/template"
)

const modellistTemplate = "{{.Name}}: {{.Version}}     {{len .Modules}} YANGS\n" +
	"Name                      File                           Revision   Organization\n" +
	"{{range .Modules}}" +
	"{{printf \"%-25s %-30s %-12s %-25s\" .Name .File .Revision .Organization}}\n" +
	"{{end}}\n"

func getListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all models in config model registry",
		Args:  cobra.NoArgs,
		RunE:  runListCommand,
	}
	return cmd
}

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <name> <version>",
		Short: "Get a model in config model registry by name and version",
		Args:  cobra.ExactArgs(2),
		RunE:  runGetCommand,
	}
	return cmd
}

func runListCommand(cmd *cobra.Command, args []string) error {
	tmplModelList, _ := template.New("model").Parse(modellistTemplate)
	clientConnection, clientConnectionError := cli.GetConnection(cmd)

	if clientConnectionError != nil {
		return clientConnectionError
	}
	client := configmodel.CreateConfigModelRegistryServiceClient(clientConnection)
	request := &configmodel.ListModelsRequest{}
	ctx := cli.NewContextWithAuthHeaderFromFlag(cmd.Context(), cmd.Flag(cli.AuthHeaderFlag))
	models, err := client.ListModels(ctx, request)
	if err != nil {
		return err
	}
	for _, model := range models.GetModels() {
		_ = tmplModelList.Execute(cli.GetOutput(), model)
	}
	return nil
}

func runGetCommand(cmd *cobra.Command, args []string) error {
	name := args[0]    // Argument is mandatory
	version := args[1] // Argument is mandatory

	tmplModelList, _ := template.New("model").Parse(modellistTemplate)
	clientConnection, clientConnectionError := cli.GetConnection(cmd)

	if clientConnectionError != nil {
		return clientConnectionError
	}
	client := configmodel.NewConfigModelRegistryServiceClient(clientConnection)
	request := &configmodel.GetModelRequest{
		Name:    name,
		Version: version,
	}
	ctx := cli.NewContextWithAuthHeaderFromFlag(cmd.Context(), cmd.Flag(cli.AuthHeaderFlag))
	model, err := client.GetModel(ctx, request)
	if err != nil {
		return err
	}
	_ = tmplModelList.Execute(cli.GetOutput(), model.GetModel())
	return nil
}
