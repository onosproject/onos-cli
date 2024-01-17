// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ransim

import (
	"context"
	"strconv"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func getCellsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cells",
		Short: "Get all cells",
		RunE:  runGetCellsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("watch", "w", false, "watch cell changes")

	return cmd
}

func createCellCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cell <enbid> [field options]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a cell",
		RunE:  runCreateCellCommand,
	}
	cmd.Flags().Uint32("max-ues", 10000, "maximum number of UEs connected")
	cmd.Flags().Float64("tx-power", 11.0, "transmit power (dB)")
	cmd.Flags().Float64("lat", 11.0, "geo location latitude")
	cmd.Flags().Float64("lng", 11.0, "geo location longitude")
	cmd.Flags().Int32("azimuth", 0, "azimuth of the coverage arc")
	cmd.Flags().Int32("arc", 120, "angle width of the coverage arc")
	cmd.Flags().UintSlice("neighbors", []uint{}, "neighbor cell NCGIs")
	cmd.Flags().String("color", "blue", "color label")
	cmd.Flags().Int32("a3-offset", int32(0), "A3 offset")
	cmd.Flags().Int32("a3-ttt", int32(0), "Time-To-Trigger")
	cmd.Flags().Int32("a3-hyst", int32(0), "A3 hysteresis")
	cmd.Flags().Int32("a3-celloffset", int32(0), "A3 cell Offset")
	cmd.Flags().Int32("a3-freqoffset", int32(0), "A3 frequency offset")
	return cmd
}

func getCellCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cell <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Get a cell",
		RunE:  runGetCellCommand,
	}
	return cmd
}

func updateCellCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cell <enbid> [field options]",
		Args:  cobra.ExactArgs(1),
		Short: "Update a cell",
		RunE:  runUpdateCellCommand,
	}
	cmd.Flags().Uint32("max-ues", 10000, "maximum number of UEs connected")
	cmd.Flags().Float64("tx-power", 11.0, "transmit power (dB)")
	cmd.Flags().Float64("lat", 11.0, "geo location latitude")
	cmd.Flags().Float64("lng", 11.0, "geo location longitude")
	cmd.Flags().Int32("azimuth", 0, "azimuth of the coverage arc")
	cmd.Flags().Int32("arc", 120, "angle width of the coverage arc")
	cmd.Flags().UintSlice("neighbors", []uint{}, "neighbor cell NCGIs")
	cmd.Flags().String("color", "blue", "color label")
	cmd.Flags().Int32("a3-offset", int32(0), "A3 offset")
	cmd.Flags().Int32("a3-ttt", int32(0), "Time-To-Trigger")
	cmd.Flags().Int32("a3-hyst", int32(0), "A3 hysteresis")
	cmd.Flags().Int32("a3-celloffset", int32(0), "A3 cell Offset")
	cmd.Flags().Int32("a3-freqoffset", int32(0), "A3 frequency offset")
	cmd.Flags().Uint32("pci", uint32(0), "new PCI value")
	return cmd
}

func deleteCellCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cell <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a cell",
		RunE:  runDeleteCellCommand,
	}
	return cmd
}

func getCellClient(cmd *cobra.Command) (modelapi.CellModelClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}

	return modelapi.NewCellModelClient(conn), conn, nil
}

func runGetCellsCommand(cmd *cobra.Command, _ []string) error {
	if noHeaders, _ := cmd.Flags().GetBool("no-headers"); !noHeaders {
		cli.Output("%-17s %7s %7s %7s %9s %9s %7s %7s %10s %7s %7s %10s %10s %8s %8s %4s %4s %s\n",
			"NCGI", "#UEs", "Max UEs", "TxDB", "Lat", "Lng", "Azimuth", "Arc",
			"A3Offset", "TTT", "A3Hyst", "PCellOffset", "FreqOffset", "PCI", "Color", "Idle", "Conn", "Neighbors(NCellOffset)")
	}

	client, conn, err := getCellClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	if watch, _ := cmd.Flags().GetBool("watch"); watch {
		stream, err := client.WatchCells(context.Background(), &modelapi.WatchCellsRequest{NoReplay: false})
		if err != nil {
			return err
		}

		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			cell := r.Cell
			cli.Output("%-17x %7d %7d %7.2f %9.3f %9.3f %7d %7d %10d %7d %7d %10d %10d %8d %8s %4d %4d %s\n",
				cell.NCGI, len(cell.CrntiMap), cell.MaxUEs, cell.TxPowerdB,
				cell.Location.Lat, cell.Location.Lng, cell.Sector.Azimuth, cell.Sector.Arc,
				cell.MeasurementParams.EventA3Params.A3Offset, cell.MeasurementParams.TimeToTrigger, cell.MeasurementParams.Hysteresis,
				cell.MeasurementParams.EventA3Params.A3Offset, cell.MeasurementParams.FrequencyOffset, cell.Pci, cell.Color,
				cell.RrcIdleCount, cell.RrcConnectedCount, catNCGIsWithOcn(cell.Neighbors, cell.MeasurementParams.NcellIndividualOffsets))
		}

	} else {

		stream, err := client.ListCells(context.Background(), &modelapi.ListCellsRequest{})
		if err != nil {
			return err
		}

		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			cell := r.Cell
			cli.Output("%-17x %7d %7d %7.2f %9.3f %9.3f %7d %7d %10d %7d %7d %10d %10d %8d %8s %4d, %4d, %s\n",
				cell.NCGI, len(cell.CrntiMap), cell.MaxUEs, cell.TxPowerdB,
				cell.Location.Lat, cell.Location.Lng, cell.Sector.Azimuth, cell.Sector.Arc,
				cell.MeasurementParams.EventA3Params.A3Offset, cell.MeasurementParams.TimeToTrigger, cell.MeasurementParams.Hysteresis,
				cell.MeasurementParams.EventA3Params.A3Offset, cell.MeasurementParams.FrequencyOffset, cell.Pci, cell.Color,
				cell.RrcIdleCount, cell.RrcConnectedCount, catNCGIsWithOcn(cell.Neighbors, cell.MeasurementParams.NcellIndividualOffsets))
		}
	}
	return nil
}

func optionsToCell(cmd *cobra.Command, cell *types.Cell, update bool) (*types.Cell, error) {
	geoToCell(cmd, cell, update)
	a3FreqOffset, _ := cmd.Flags().GetInt32("a3-freqoffset")
	if !update || cmd.Flags().Changed("a3-freqoffset") {
		cell.MeasurementParams.FrequencyOffset = a3FreqOffset
	}

	neighbors, _ := cmd.Flags().GetUintSlice("neighbors")
	if !update || cmd.Flags().Changed("neighbors") {
		cell.Neighbors = toNCGIs(neighbors)
	}
	pci, _ := cmd.Flags().GetUint32("pci")
	if !update || cmd.Flags().Changed("pci") {
		cell.Pci = pci
	}
	txPower, _ := cmd.Flags().GetFloat64("tx-power")
	if !update || cmd.Flags().Changed("tx-power") {
		cell.TxPowerdB = txPower
	}
	return cell, nil
}

func geoToCell(cmd *cobra.Command, cell *types.Cell, update bool) {
	arc, _ := cmd.Flags().GetInt32("arc")
	azimuth, _ := cmd.Flags().GetInt32("azimuth")
	lat, _ := cmd.Flags().GetFloat64("lat")
	lng, _ := cmd.Flags().GetFloat64("lng")

	if cell.Location == nil {
		cell.Location = &types.Point{Lat: lat, Lng: lng}
	} else {
		if !update || cmd.Flags().Changed("lat") {
			cell.Location.Lng = lng
		}
		if !update || cmd.Flags().Changed("lng") {
			cell.Location.Lng = lng
		}
	}

	if cell.Sector == nil {
		cell.Sector = &types.Sector{Centroid: cell.Location, Azimuth: azimuth, Arc: arc}
	} else {
		cell.Sector.Centroid = cell.Location
		if !update || cmd.Flags().Changed("arc") {
			cell.Sector.Arc = arc
		}
		if !update || cmd.Flags().Changed("azimuth") {
			cell.Sector.Azimuth = azimuth
		}
	}

	color, _ := cmd.Flags().GetString("color")
	if !update || cmd.Flags().Changed("color") {
		cell.Color = color
	}
}

func runCreateCellCommand(cmd *cobra.Command, args []string) error {
	ncgi, err := strconv.ParseUint(args[0], 16, 64)
	if err != nil {
		return err
	}

	client, conn, err := getCellClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	cell, err := optionsToCell(cmd, &types.Cell{NCGI: types.NCGI(ncgi)}, false)
	if err != nil {
		return err
	}

	_, err = client.CreateCell(context.Background(), &modelapi.CreateCellRequest{Cell: cell})
	if err != nil {
		return err
	}
	cli.Output("Cell %x created\n", ncgi)
	return nil
}

func runUpdateCellCommand(cmd *cobra.Command, args []string) error {
	ncgi, err := strconv.ParseUint(args[0], 16, 64)
	if err != nil {
		return err
	}

	client, conn, err := getCellClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Get the cell first to prime the update cell with existing values and allow sparse update
	gres, err := client.GetCell(context.Background(), &modelapi.GetCellRequest{NCGI: types.NCGI(ncgi)})
	if err != nil {
		return err
	}

	cell, err := optionsToCell(cmd, gres.Cell, true)
	if err != nil {
		return err
	}
	cli.Output("Cell is now %v", cell)

	_, err = client.UpdateCell(context.Background(), &modelapi.UpdateCellRequest{Cell: cell})
	if err != nil {
		return err
	}
	cli.Output("Cell %x updated\n", ncgi)
	return nil
}

func runGetCellCommand(cmd *cobra.Command, args []string) error {
	ncgi, err := strconv.ParseUint(args[0], 16, 64)
	if err != nil {
		return err
	}

	client, conn, err := getCellClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	res, err := client.GetCell(context.Background(), &modelapi.GetCellRequest{NCGI: types.NCGI(ncgi)})
	if err != nil {
		return err
	}

	cell := res.Cell
	cli.Output("NCGI:       %-17x\nUE Count:   %-5d\nMax UEs:    %-5d\nTxPower dB: %.2f\n",
		cell.NCGI, len(cell.CrntiMap), cell.MaxUEs, cell.TxPowerdB)
	cli.Output("Latitude:   %.3f\nLongitude:  %.3f\nAzimuth:    %d\nArc:        %d\nPCI:        %d\nColor:      %s\nNeighbors:  %s\n",
		cell.Location.Lat, cell.Location.Lng, cell.Sector.Azimuth, cell.Sector.Arc, cell.Pci, cell.Color,
		catNCGIs(cell.Neighbors))
	cli.Output("A3offset:          %7d\nA3TimeToTrigger:   %7d\nA3Hystereis:       %7d\nA3CellOffset:      %7d\nA3FrequencyOffset: %7d\n",
		cell.MeasurementParams.EventA3Params.A3Offset, cell.MeasurementParams.TimeToTrigger, cell.MeasurementParams.Hysteresis,
		cell.MeasurementParams.PcellIndividualOffset, cell.MeasurementParams.FrequencyOffset)
	cli.Output("RrcIdleCount: %d\nRrcConnectedCount: %d\n", cell.RrcIdleCount, cell.RrcConnectedCount)
	return nil
}

func runDeleteCellCommand(cmd *cobra.Command, args []string) error {
	ncgi, err := strconv.ParseUint(args[0], 16, 64)
	if err != nil {
		return err
	}

	client, conn, err := getCellClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.DeleteCell(context.Background(), &modelapi.DeleteCellRequest{NCGI: types.NCGI(ncgi)})
	if err != nil {
		return err
	}

	cli.Output("Cell %x deleted\n", ncgi)
	return nil
}
