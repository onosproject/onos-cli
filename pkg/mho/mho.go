// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mho

import (
	"context"
	mhoapi "github.com/onosproject/onos-api/go/onos/mho"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"strconv"
	"text/tabwriter"
)

func getGetUesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ues",
		Short: "Get ues",
		RunE:  runGetUesCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getGetCellsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cells",
		Short: "Get cells",
		RunE:  runGetCellsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runGetUesCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	request := mhoapi.GetRequest{}
	client := mhoapi.NewMhoClient(conn)

	ueList, err := client.GetUes(context.Background(), &request)
	if err != nil {
		return err
	}

	if noHeaders, _ := cmd.Flags().GetBool("no-headers"); !noHeaders {
		cli.Output("%-8s %-16s %-8s\n", "UeID", "CellGlobalID", "RrcState")
	}

	for _, ue := range ueList.Ues {
		ueID, _ := strconv.Atoi(ue.UeId)
		rrcState := ue.RrcState
		if len(rrcState) >= 10 {
			rrcState = rrcState[10:]
		}
		cli.Output("%-8x %-16s %-8s\n", ueID, ue.Cgi, rrcState)
	}

	return nil
}

func runGetCellsCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	request := mhoapi.GetRequest{}
	client := mhoapi.NewMhoClient(conn)

	cellList, err := client.GetCells(context.Background(), &request)
	if err != nil {
		return err
	}

	if noHeaders, _ := cmd.Flags().GetBool("no-headers"); !noHeaders {
		cli.Output("%-16s %-16s %-16s %-16s\n", "CGI", "Num UEs", "Handovers-in", "Handovers-out")
	}

	for _, cell := range cellList.Cells {
		cumulativeHandoversIn := cell.CumulativeHandoversIn
		cumulativeHandoversOut := cell.CumulativeHandoversOut
		cli.Output("%-16s %-16d %-16d %-16d\n", cell.Cgi, cell.NumUes, cumulativeHandoversIn, cumulativeHandoversOut)
	}

	return nil
}
