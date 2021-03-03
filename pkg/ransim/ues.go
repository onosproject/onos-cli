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

package ransim

import (
	"context"

	simapi "github.com/onosproject/onos-api/go/onos/ransim/trafficsim"
	"github.com/onosproject/onos-lib-go/pkg/cli"

	"github.com/spf13/cobra"
)

func getUEsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ues",
		Short: "Get UEs",
		RunE:  runGetUEsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("watch", "w", false, "watch ue changes")
	return cmd
}

func getUECountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ueCount",
		Short: "Get UE count",
		RunE:  runGetUECountCommand,
	}
	return cmd
}

func runGetUEsCommand(cmd *cobra.Command, args []string) error {
	if noHeaders, _ := cmd.Flags().GetBool("no-headers"); !noHeaders {
		cli.Output("%-16s %-16s %-5s\n", "IMSI", "Serving Cell", "Admitted")
	}
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := simapi.NewTrafficClient(conn)

	if watch, _ := cmd.Flags().GetBool("watch"); watch {

	} else {
		stream, err := client.ListUes(context.Background(), &simapi.ListUesRequest{})
		if err != nil {
			return err
		}

		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			ue := r.Ue
			cli.Output("%-16d %-16d %-5t\n", ue.IMSI, ue.ServingTower, ue.Admitted)
		}
	}

	return nil
}

func countUEs(stream simapi.Traffic_ListUesClient) int {
	count := 0
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
		count = count + 1
	}
	return count
}

func runGetUECountCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := simapi.NewTrafficClient(conn)
	stream, err := client.ListUes(context.Background(), &simapi.ListUesRequest{})
	if err != nil {
		return err
	}

	cli.Output("%d\n", countUEs(stream))
	return nil
}
