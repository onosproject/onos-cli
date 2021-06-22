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
	"strconv"

	simapi "github.com/onosproject/onos-api/go/onos/ransim/trafficsim"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"google.golang.org/grpc"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
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

func getUECommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ue <imsi>",
		Short: "Get UE",
		RunE:  runGetUECommand,
	}
	return cmd
}

func updateUECommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ue <imsi> [field options]",
		Args:  cobra.ExactArgs(1),
		Short: "Update a UE NCGI assignment and/or geo location",
		RunE:  runUpdateUECommand,
	}
	cmd.Flags().Uint64("ncgi", 0, "serving cell NCGI")
	cmd.Flags().Float64("lat", 0.0, "new coordinate latitude")
	cmd.Flags().Float64("lng", 0.0, "new coordinate longitude")
	cmd.Flags().Uint32("heading", 0, "new heading")
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

func setUECountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ueCount <count>",
		Args:  cobra.ExactArgs(1),
		Short: "Set UE count",
		RunE:  runSetUECountCommand,
	}
	return cmd
}

func getUEClient(cmd *cobra.Command) (modelapi.UEModelClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}
	return modelapi.NewUEModelClient(conn), conn, nil
}

func runGetUEsCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getUEClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	if noHeaders, _ := cmd.Flags().GetBool("no-headers"); !noHeaders {
		cli.Output("%-16s %-16s %-10s %-5s\n", "IMSI", "Serving Cell", "CRNTI", "Admitted")
	}

	if watch, _ := cmd.Flags().GetBool("watch"); watch {
		stream, err := client.WatchUEs(context.Background(), &modelapi.WatchUEsRequest{NoReplay: false})
		if err != nil {
			return err
		}
		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			ue := r.Ue
			cli.Output("%-16d %-16x %-10d %-5t\n", ue.IMSI, ue.ServingTower, ue.CRNTI, ue.Admitted)
		}

	} else {
		stream, err := client.ListUEs(context.Background(), &modelapi.ListUEsRequest{})
		if err != nil {
			return err
		}

		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			ue := r.Ue
			cli.Output("%-16d %-16x %-10d %-10t %-5d\n", ue.IMSI, ue.ServingTower, ue.CRNTI, ue.Admitted, ue.RrcState)
		}
	}

	return nil
}

func outputUE(ue *types.Ue) {
	cli.Output("IMSI:      %-16d\nNCGI:      %-16x\nStrength: %.4f\nCRNTI:     %-16d\nAdmitted:  %t\nLat:       %.4f\nLng:       %.4f\nHeading:   %3d\nRrc:       %d\n",
		ue.IMSI, ue.ServingTower, ue.ServingTowerStrength, ue.CRNTI, ue.Admitted, ue.Position.Lat, ue.Position.Lng, ue.Rotation, ue.RrcState)
	cli.Output("Candidate Cells: %x, %x, %x\n", ue.Tower1, ue.Tower2, ue.Tower3)
}

func runGetUECommand(cmd *cobra.Command, args []string) error {
	imsi, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getUEClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := client.GetUE(context.Background(), &modelapi.GetUERequest{IMSI: types.IMSI(imsi)})
	if err != nil {
		return err
	}

	outputUE(res.Ue)
	return nil
}

func runUpdateUECommand(cmd *cobra.Command, args []string) error {
	imsi, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getUEClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	ncgi, _ := cmd.Flags().GetUint64("ncgi")
	if ncgi != 0 {
		_, err := client.MoveToCell(context.Background(),
			&modelapi.MoveToCellRequest{
				IMSI: types.IMSI(imsi),
				NCGI: types.NCGI(ncgi),
			})
		if err != nil {
			return err
		}
		cli.Output("UE %d cell updated\n", imsi)
	}

	lat, _ := cmd.Flags().GetFloat64("lat")
	lng, _ := cmd.Flags().GetFloat64("lng")
	heading, _ := cmd.Flags().GetUint32("heading")
	if lat != 0 || lng != 0 || heading != 0 {
		_, err := client.MoveToLocation(context.Background(),
			&modelapi.MoveToLocationRequest{
				IMSI:     types.IMSI(imsi),
				Location: &types.Point{Lat: lat, Lng: lng},
				Heading:  heading,
			})
		if err != nil {
			return err
		}
		cli.Output("UE %d location updated\n", imsi)
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

func runSetUECountCommand(cmd *cobra.Command, args []string) error {
	count, err := strconv.ParseUint(args[0], 10, 16)
	if err != nil {
		return err
	}
	client, conn, err := getUEClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.SetUECount(context.Background(), &modelapi.SetUECountRequest{Count: uint32(count)})
	if err != nil {
		return err
	}

	return nil
}
