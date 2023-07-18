// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package pci

import (
	"context"
	"fmt"
	"strconv"
	"text/tabwriter"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"

	pciapi "github.com/onosproject/onos-api/go/onos/pci"
)

func getGetConflicts() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "conflicts",
		Aliases: []string{"conflict"},
		Short:   "Get the conflicting cells for a specific cell or all cells if not specified",
		RunE:    runGetConflicts,
		Args:    cobra.MaximumNArgs(1),
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getGetResolvedConflicts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolved",
		Short: "Get the number of resolutions and most recent resolution for all cells",
		RunE:  runGetResolvedConflicts,
		Args:  cobra.NoArgs,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getGetCell() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cell <id>",
		Short: "Get a single cell's info",
		RunE:  runGetCell,
		Args:  cobra.ExactArgs(1),
	}
	return cmd
}

func getGetCells() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cells",
		Short: "Get all cells",
		RunE:  runGetCells,
		Args:  cobra.NoArgs,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runGetConflicts(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	request := pciapi.GetConflictsRequest{}
	if len(args) != 0 {
		id, err := strconv.ParseUint(args[0], 16, 64)
		if err != nil {
			return err
		}
		request.CellId = id
	}
	client := pciapi.NewPciClient(conn)
	response, err := client.GetConflicts(context.Background(), &request)
	if err != nil {
		return err
	}

	printTableHeader(noHeaders, writer)
	for _, cell := range response.GetCells() {
		printTableCell(cell, writer)
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func runGetResolvedConflicts(cmd *cobra.Command, _ []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	request := pciapi.GetResolvedConflictsRequest{}
	client := pciapi.NewPciClient(conn)
	response, err := client.GetResolvedConflicts(context.Background(), &request)
	if err != nil {
		return err
	}

	printResolvedHeader(noHeaders, writer)
	for _, cell := range response.GetCells() {
		printResolvedCell(cell, writer)
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func runGetCell(cmd *cobra.Command, args []string) error {

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	request := pciapi.GetCellRequest{}
	id, err := strconv.ParseUint(args[0], 16, 64)
	if err != nil {
		return err
	}
	request.CellId = id

	client := pciapi.NewPciClient(conn)
	response, err := client.GetCell(context.Background(), &request)
	if err != nil {
		return err
	}

	printSingleCell(response.Cell, writer)
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func runGetCells(cmd *cobra.Command, _ []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	request := pciapi.GetCellsRequest{}
	client := pciapi.NewPciClient(conn)
	response, err := client.GetCells(context.Background(), &request)
	if err != nil {
		return err
	}

	printTableHeader(noHeaders, writer)
	for _, cell := range response.GetCells() {
		printTableCell(cell, writer)
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func printTableHeader(noHeaders bool, writer *tabwriter.Writer) {
	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "ID\tNode ID\tARFCN\tCell Type\tPCI\tPCI Pool\n")
	}
}

func printTableCell(cell *pciapi.PciCell, writer *tabwriter.Writer) {
	_, _ = fmt.Fprintf(writer, "%x\t%s\t%d\t%s\t%d\t", cell.Id, cell.NodeId, cell.Arfcn, cell.CellType.String(), cell.Pci)

	// print pci pools
	_, _ = fmt.Fprint(writer, "[")
	for i, minMax := range cell.PciPool {
		_, _ = fmt.Fprintf(writer, "%d:%d", minMax.Min, minMax.Max)
		if i != len(cell.PciPool)-1 {
			_, _ = fmt.Fprint(writer, ",")
		}
	}
	_, _ = fmt.Fprint(writer, "]\n")
}

func printResolvedHeader(noHeaders bool, writer *tabwriter.Writer) {
	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "ID\tTotal Resolved Conflicts\tMost Recent Resolution\n")
	}
}

func printResolvedCell(cell *pciapi.CellResolution, writer *tabwriter.Writer) {
	_, _ = fmt.Fprintf(writer, "%x\t%d\t%d=>%d\n", cell.Id, cell.ResolvedConflicts, cell.OriginalPci, cell.ResolvedPci)
}

func printSingleCell(cell *pciapi.PciCell, writer *tabwriter.Writer) {
	_, _ = fmt.Fprintf(writer, "ID: %x\nNode ID: %s\nARFCN: %d\nCell Type: %s\nPCI: %d\n", cell.Id, cell.NodeId, cell.Arfcn, cell.CellType.String(), cell.Pci)

	// print neighbors
	_, _ = fmt.Fprint(writer, "Neighbors: [")
	for i, neighbor := range cell.NeighborIds {
		_, _ = fmt.Fprintf(writer, "%x", neighbor)
		if i != len(cell.NeighborIds)-1 {
			_, _ = fmt.Fprint(writer, ",")
		}
	}
	_, _ = fmt.Fprint(writer, "]\n")

	// print pci pools
	_, _ = fmt.Fprint(writer, "PCI Pool: [")
	for i, minMax := range cell.PciPool {
		_, _ = fmt.Fprintf(writer, "%d:%d", minMax.Min, minMax.Max)
		if i != len(cell.PciPool)-1 {
			_, _ = fmt.Fprint(writer, ",")
		}
	}
	_, _ = fmt.Fprint(writer, "]\n")
}
