// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

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

const configurationListTemplate = "table{{.ID}}\t{{.TargetID}}\t{{.Status.State}}\t{{.Index}}"

var configurationListTemplateVerbose = fmt.Sprintf("%s\t{{.Values}}", configurationListTemplate)

const configurationEventTemplate = "table{{.Type}}\t{{.Configuration.ID}}\t{{.Configuration.TargetID}}\t{{.Configuration.Status.State}}\t{{.Configuration.Index}}"

var configurationEventTemplateVerbose = fmt.Sprintf("%s\t{{.Configuration.Values}}", configurationEventTemplate)

type configurationEventWidths struct {
	Type          int
	Configuration struct {
		ID       int
		TargetID int
		Status   struct {
			State int
		}
		Index  int
		Values int
	}
}

var configWidths = configurationEventWidths{
	Type: 30,
	Configuration: struct {
		ID       int
		TargetID int
		Status   struct{ State int }
		Index    int
		Values   int
	}{
		ID:       13,
		TargetID: 13,
		Status:   struct{ State int }{State: 40},
		Index:    5,
		Values:   50,
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
	ctx, cancel := context.WithTimeout(cli.NewContextWithAuthHeaderFromFlag(cmd.Context(), cmd.Flag(cli.AuthHeaderFlag)), 15*time.Second)
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

	return tableFormat.Execute(cli.GetOutput(), !noHeaders, 0, resp.Configuration)

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

	var allConfigurations []*v2.Configuration

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return tableFormat.Execute(cli.GetOutput(), !noHeaders, 0, allConfigurations)
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
	stream, err := client.WatchConfigurations(cli.NewContextWithAuthHeaderFromFlag(cmd.Context(), cmd.Flag(cli.AuthHeaderFlag)), request)
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
