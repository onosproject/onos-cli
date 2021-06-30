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
		id, err := strconv.ParseUint(args[0], 10, 64)
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
		// _, _ = fmt.Fprintf(writer, "%x\t%v\n", cell.Id, cell)
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

func runGetCells(cmd *cobra.Command, args []string) error {
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
		_, _ = fmt.Fprintf(writer, "ID\tNode ID\tDlearfcn\tCell Type\tPCI\tPCI Pool\n")
	}
}

func printTableCell(cell *pciapi.PciCell, writer *tabwriter.Writer) {
	_, _ = fmt.Fprintf(writer, "%x\t%s\t%d\t%s\t%d\t", cell.Id, cell.NodeId, cell.Dlearfcn, cell.CellType.String(), cell.Pci)

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

func printSingleCell(cell *pciapi.PciCell, writer *tabwriter.Writer) {
	_, _ = fmt.Fprintf(writer, "ID: %x\nNode ID: %s\nDlearfcn: %d\nCell Type: %s\nPCI: %d\n", cell.Id, cell.NodeId, cell.Dlearfcn, cell.CellType.String(), cell.Pci)

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
