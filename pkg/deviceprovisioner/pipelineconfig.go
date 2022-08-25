// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package deviceprovisioner

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-api/go/onos/device-provisioner/admin"
	p4rtapi "github.com/onosproject/onos-api/go/onos/p4rt/v1"
	"github.com/onosproject/onos-cli/pkg/format"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"io"
	"time"
)

const pipelineListTemplate = "table{{.ID}}\t{{.TargetID}}\t{{.Status.State}}\t{{.Action}}"

var pipelineListTemplateVerbose = fmt.Sprintf("%s\t{{.PipelineConfigSpec}}", pipelineListTemplate)

const pipelineEventTemplate = "table{{.PipelineConfig.ID}}\t{{.PipelineConfig.TargetID}}\t{{.PipelineConfig.Status.State}}\t{{.PipelineConfig.Action}}"

var pipelineEventTemplateVerbose = fmt.Sprintf("%s\t{{.PipelineConfig.PipelineConfigSpec}}", pipelineListTemplate)

type pipelineEventWidths struct {
	PipelineConfig struct {
		ID       int
		TargetID int
		Status   struct {
			State int
		}
		Action int
	}
}

var pipelineWidths = pipelineEventWidths{

	PipelineConfig: struct {
		ID       int
		TargetID int
		Status   struct{ State int }
		Action   int
	}{
		ID:       13,
		TargetID: 13,
		Status:   struct{ State int }{State: 40},
		Action:   15,
	},
}

func getListPipelinesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pipelines [pipelineConfigID]",
		Short:   "List target pipeline configs",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"pipeline"},
		RunE:    runListPipelinesCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the change with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchPipelinesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pipelines [pipelineConfigID]",
		Short:   "Watch pipelines ",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"pipelines"},
		RunE:    runWatchPipelinesCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the change with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("no-replay", "r", false, "do not replay existing configurations")
	return cmd
}

func runListPipelinesCommand(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewPipelineConfigServiceClient(conn)
	ctx, cancel := context.WithTimeout(cli.NewContextWithAuthHeaderFromFlag(cmd.Context(), cmd.Flag(cli.AuthHeaderFlag)), 15*time.Second)
	defer cancel()

	if len(args) > 0 {
		return getPipelines(ctx, client, p4rtapi.PipelineConfigID(args[0]), noHeaders, verbose)
	}
	return listPipelines(ctx, client, noHeaders, verbose)
}

func listPipelines(ctx context.Context, client admin.PipelineConfigServiceClient, noHeaders bool, verbose bool) error {
	stream, err := client.ListPipelines(ctx, &admin.ListPipelinesRequest{})
	if err != nil {
		cli.Output("Unable to list piplines %s", err)
		return err
	}

	var tableFormat format.Format
	if verbose {
		tableFormat = format.Format(pipelineListTemplateVerbose)
	} else {
		tableFormat = format.Format(pipelineListTemplate)
	}

	var allPipelineConfigs []*p4rtapi.PipelineConfig

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			if e := tableFormat.Execute(cli.GetOutput(), !noHeaders, 0, allPipelineConfigs); e != nil {
				return e
			}
			return nil
		} else if err != nil {
			cli.Output("Unable to read pipeline config: %s", err)
			return err
		}
		allPipelineConfigs = append(allPipelineConfigs, resp.Pipelineconfig)
	}

}

func getPipelines(ctx context.Context, client admin.PipelineConfigServiceClient, id p4rtapi.PipelineConfigID, noHeaders bool, verbose bool) error {
	resp, err := client.GetPipeline(ctx, &admin.GetPipelineRequest{PipelineConfigID: id})
	if err != nil {
		cli.Output("Unable to get pipeline: %s", err)
		return err
	}
	var tableFormat format.Format
	if verbose {
		tableFormat = format.Format(pipelineListTemplateVerbose)
	} else {
		tableFormat = format.Format(pipelineListTemplate)
	}

	if e := tableFormat.Execute(cli.GetOutput(), !noHeaders, 0, resp.Pipelineconfig); e != nil {
		return e
	}
	return nil

}

func runWatchPipelinesCommand(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noReplay, _ := cmd.Flags().GetBool("no-replay")

	id := p4rtapi.PipelineConfigID("")
	if len(args) > 0 {
		id = p4rtapi.PipelineConfigID(args[0])
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewPipelineConfigServiceClient(conn)
	request := &admin.WatchPipelinesRequest{Noreplay: noReplay, PipelineConfigID: id}
	stream, err := client.WatchPipelines(cli.NewContextWithAuthHeaderFromFlag(cmd.Context(), cmd.Flag(cli.AuthHeaderFlag)), request)
	if err != nil {
		return err
	}

	f := format.Format(pipelineEventTemplate)
	if verbose {
		f = format.Format(pipelineEventTemplateVerbose)
	}
	if !noHeaders {
		output, err := f.ExecuteFixedWidth(pipelineWidths, true, nil)
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

		event := res.PipelineConfig
		if len(id) == 0 || id == event.ID {
			output, err := f.ExecuteFixedWidth(pipelineWidths, false, res)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", output)
		}
	}

	return nil
}
