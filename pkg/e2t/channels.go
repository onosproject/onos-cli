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

package e2t

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-cli/pkg/utils"
	"io"
	"strings"
	"text/tabwriter"

	subapi "github.com/onosproject/onos-api/go/onos/e2t/e2/v1beta1"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	channelHeaders = "Channel ID\tRevision\tService Model ID\tE2 NodeID\tEncoding\tPhase\tState"
	channelFormat  = "%s\t%d\t%s:%s\t%s\t%s\t%s\t%s\n"
)

func displayChannelHeaders(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, channelHeaders)
}

func displayChannel(writer io.Writer, sub *subapi.Channel) {
	_, _ = fmt.Fprintf(writer, channelFormat,
		sub.ID, sub.Revision, sub.ChannelMeta.ServiceModel.Name, sub.ChannelMeta.ServiceModel.Version, utils.None(string(sub.ChannelMeta.E2NodeID)),
		utils.None(sub.ChannelMeta.Encoding.String()), utils.None(sub.Status.Phase.String()), utils.None(sub.Status.State.String()))
}

func getGetChannelsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channels",
		Short: "Get NB channels",
		RunE:  runGetChannelsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runGetChannelsCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := subapi.NewSubscriptionAdminServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.ListChannels(ctx, &subapi.ListChannelsRequest{})
	if err != nil {
		return err
	}

	if !noHeaders {
		displayChannelHeaders(writer)
	}

	for _, sub := range response.Channels {
		pin := sub
		displayChannel(writer, &pin)
	}

	_ = writer.Flush()
	return nil
}

func getGetChannelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channel",
		Short: "Get NB channel",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetChannelCommand,
	}
	return cmd
}

func runGetChannelCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := subapi.NewSubscriptionAdminServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.GetChannel(ctx, &subapi.GetChannelRequest{
		ChannelID: subapi.ChannelID(args[0]),
	})
	if err != nil {
		return err
	}

	chn := response.Channel
	_, _ = fmt.Fprintf(writer, "Channel ID:\t%s\nRevision:\t%d\nService Model:\t%s\nService Model Version:\t%s\n",
		chn.ID, chn.Revision, chn.ChannelMeta.ServiceModel.Name, chn.ChannelMeta.ServiceModel.Version)
	_, _ = fmt.Fprintf(writer, "E2 Node ID:\t%s\nEncoding:\t%s",
		utils.None(string(chn.ChannelMeta.E2NodeID)), utils.None(chn.ChannelMeta.Encoding.String()))
	_, _ = fmt.Fprintf(writer, "Phase:\t%s\nStatus:\t%s\n", utils.None(chn.Status.Phase.String()), utils.None(chn.Status.State.String()))
	_, _ = fmt.Fprintf(writer, "Actions:\t%v\nTrigger:\t%v\n", chn.Spec.Actions, chn.Spec.EventTrigger)
	_ = writer.Flush()
	return nil
}

func getWatchChannelsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channels",
		Short: "Watch NB channels",
		RunE:  runWatchChannelsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().Bool("no-replay", false, "disables replay of existing state")
	return cmd
}

func runWatchChannelsCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noReplay, _ := cmd.Flags().GetBool("no-replay")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := subapi.NewSubscriptionAdminServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.WatchChannels(ctx, &subapi.WatchChannelsRequest{
		NoReplay: noReplay,
	})
	if err != nil {
		return err
	}

	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "Event Type\t")
		displayChannelHeaders(writer)
		_ = writer.Flush()
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

		event := res.Event
		_, _ = fmt.Fprintf(writer, "%s\t", strings.Replace(event.Type.String(), "SUBSCRIPTION_", "", 1))
		displayChannel(writer, &event.Channel)
		_ = writer.Flush()
	}

	return nil
}
