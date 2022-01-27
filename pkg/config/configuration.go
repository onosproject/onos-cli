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

package config

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-api/go/onos/config/admin"
	v2 "github.com/onosproject/onos-api/go/onos/config/v2"
	"github.com/onosproject/onos-cli/pkg/format"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"io"
	"time"
)

const configurationListTemplate = "table{{.ID}}\t{{.TargetID}}\t{{.TargetVersion}}\t{{.TargetType}}\t{{.Status.State}}"

var configurationListTemplateVerbose = fmt.Sprintf("%s\t{{.Values}}", configurationListTemplate)

const configurationEventTemplate = "table{{.Type}}\t{{.Configuration.ID}}\t{{.Configuration.TargetID}}\t{{.Configuration.TargetVersion}}\t{{.Configuration.TargetType}}\t{{.Configuration.Status.State}}"

var configurationEventTemplateVerbose = fmt.Sprintf("%s\t{{.Configuration.Values}}", configurationEventTemplate)

type configurationEventWidths struct {
	Type          int
	Configuration struct {
		ID            int
		TargetID      int
		TargetVersion int
		TargetType    int
		Status        struct {
			State int
		}
		Revision int
		Index    int
		Values   int
	}
}

var configWidths = configurationEventWidths{
	Type: 30,
	Configuration: struct {
		ID            int
		TargetID      int
		TargetVersion int
		TargetType    int
		Status        struct{ State int }
		Revision      int
		Index         int
		Values        int
	}{
		ID:            13,
		TargetID:      13,
		TargetVersion: 15,
		TargetType:    13,
		Status:        struct{ State int }{State: 40},
		Revision:      5,
		Index:         5,
		Values:        50,
	},
}

func getListConfigurationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "configurations [configurationID]",
		Short:   "List target configurations",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"configuration"},
		RunE:    runListConfigurationsCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the change with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchConfigurationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "configurations [configurationID]",
		Short:   "Watch target configurations",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"configuration"},
		RunE:    runWatchConfigurationsCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the change with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay existing configurations")
	return cmd
}

func runListConfigurationsCommand(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewConfigurationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if len(args) > 0 {
		return getConfigurations(ctx, client, v2.ConfigurationID(args[0]), noHeaders, verbose)
	}
	return listConfigurations(ctx, client, noHeaders, verbose)
}

func getConfigurations(ctx context.Context, client admin.ConfigurationServiceClient, id v2.ConfigurationID, noHeaders bool, verbose bool) error {
	resp, err := client.GetConfiguration(ctx, &admin.GetConfigurationRequest{ConfigurationID: id})
	if err != nil {
		cli.Output("Unable to get configuration: %s", err)
		return err
	}
	var tableFormat format.Format
	if verbose {
		tableFormat = format.Format(configurationListTemplateVerbose)
	} else {
		tableFormat = format.Format(configurationListTemplate)
	}

	if e := tableFormat.Execute(cli.GetOutput(), !noHeaders, 0, resp.Configuration); e != nil {
		return e
	}
	return nil

}

func listConfigurations(ctx context.Context, client admin.ConfigurationServiceClient, noHeaders bool, verbose bool) error {
	stream, err := client.ListConfigurations(ctx, &admin.ListConfigurationsRequest{})
	if err != nil {
		cli.Output("Unable to list configurations: %s", err)
		return err
	}

	var tableFormat format.Format
	if verbose {
		tableFormat = format.Format(configurationListTemplateVerbose)
	} else {
		tableFormat = format.Format(configurationListTemplate)
	}

	allConfigurations := []*v2.Configuration{}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			if e := tableFormat.Execute(cli.GetOutput(), !noHeaders, 0, allConfigurations); e != nil {
				return e
			}
			return nil
		} else if err != nil {
			cli.Output("Unable to read configuration: %s", err)
			return err
		}
		allConfigurations = append(allConfigurations, resp.Configuration)
	}

}

func runWatchConfigurationsCommand(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noReplay, _ := cmd.Flags().GetBool("no-replay")

	id := v2.ConfigurationID("")
	if len(args) > 0 {
		id = v2.ConfigurationID(args[0])
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewConfigurationServiceClient(conn)
	request := &admin.WatchConfigurationsRequest{Noreplay: noReplay, ConfigurationID: id}
	stream, err := client.WatchConfigurations(context.Background(), request)
	if err != nil {
		return err
	}

	f := format.Format(configurationEventTemplate)
	if verbose {
		f = format.Format(configurationEventTemplateVerbose)
	}
	if !noHeaders {
		output, err := f.ExecuteFixedWidth(configWidths, true, nil)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", output)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			cli.Output("Error receiving notification : %v", err)
			return err
		}

		event := res.ConfigurationEvent
		if len(id) == 0 || id == event.Configuration.ID {
			output, err := f.ExecuteFixedWidth(configWidths, false, res)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", output)
		}
	}

	return nil
}
