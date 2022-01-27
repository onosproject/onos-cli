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
	"io"
	"os"
	"time"

	"github.com/onosproject/onos-api/go/onos/config/admin"
	v2 "github.com/onosproject/onos-api/go/onos/config/v2"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getListConfigurationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configurations [targetID]",
		Short: "List target configurations",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runListConfigurationsCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the change with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchConfigurationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configurations [targetID]",
		Short: "Watch target configurations",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runWatchConfigurationsCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the change with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runListConfigurationsCommand(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	writer := os.Stdout
	if !noHeaders {
		printConfigurationHeader(writer, verbose, false)
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := admin.NewConfigurationServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	response, err := client.ListConfigurations(ctx, &admin.ListConfigurationsRequest{})
	if err != nil {
		cli.Output("Unable to list configurations: %s", err)
		return err
	}

	for {
		resp, err := response.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			cli.Output("Unable to read configuration: %s", err)
			return err
		} else {
			printConfiguration(writer, resp.Configuration, verbose)
		}
	}

	return nil
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
	stream, err := client.WatchConfigurations(context.Background(), &admin.WatchConfigurationsRequest{Noreplay: noReplay})
	if err != nil {
		return err
	}

	writer := os.Stdout
	if !noHeaders {
		printConfigurationHeader(writer, verbose, true)
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
		if len(id) == 0 || id == event.Configuration.ID {
			printConfigurationUpdateType(writer, event.Type)
			printConfiguration(writer, &event.Configuration, false)
		}
	}

	return nil
}

func printConfiguration(writer io.Writer, c *v2.Configuration, verbose bool) {
	if verbose {
		_, _ = fmt.Fprintf(writer, "%-12s\t%-12s\t%-8s\t%-10s\t%-10s\t%-8d\t%v\t\n",
			c.ID, c.TargetID, c.TargetVersion, c.TargetType, c.Status.State, c.Revision, c.Values)
	} else {
		_, _ = fmt.Fprintf(writer, "%-12s\t%-12s\t%-8s\t%-10s\t%-10s\t%-8d\t\n",
			c.ID, c.TargetID, c.TargetVersion, c.TargetType, c.Status.State, c.Revision)
	}
}

func printConfigurationUpdateType(writer io.Writer, eventType v2.ConfigurationEventType) {
	if eventType == v2.ConfigurationEventType_CONFIGURATION_REPLAYED {
		_, _ = fmt.Fprintf(writer, "%-12s\t", "REPLAY")
	} else {
		_, _ = fmt.Fprintf(writer, "%-12s\t", eventType)
	}
}

func printConfigurationHeader(writer *os.File, verbose bool, event bool) {
	if event {
		_, _ = fmt.Fprintf(writer, "%-12s\t", "Event Type")
	}
	if verbose {
		_, _ = fmt.Fprintf(writer, "%-12s\t%-12s\t%-8s\t%-10s\t%-10s\t%-8s\t%-8s\t\n",
			"ID", "Target ID", "Version", "Type", "Status", "Revision", "Values")
	} else {
		_, _ = fmt.Fprintf(writer, "%-12s\t%-12s\t%-8s\t%-10s\t%-10s\t%-8s\t\n",
			"ID", "Target ID", "Version", "Type", "Status", "Revision")
	}
}
