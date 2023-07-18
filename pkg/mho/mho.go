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

func runGetUesCommand(cmd *cobra.Command, _ []string) error {
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
		cli.Output("%-20s %-16s %-8s\n", "AMF-UE-NGAP-ID", "CellGlobalID", "HOState")
	}

	for _, ue := range ueList.Ues {
		ueID, _ := strconv.Atoi(ue.UeId)
		hoState := ue.HoState
		if len(hoState) >= 10 {
			hoState = hoState[10:]
		}
		cli.Output("%-20x %-16s %-8s\n", ueID, ue.Cgi, hoState)
	}

	return nil
}

func runGetCellsCommand(cmd *cobra.Command, _ []string) error {
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
		cli.Output("%-16s %-16s\n", "CGI", "Num UEs")
	}

	for _, cell := range cellList.Cells {
		cli.Output("%-16s %-16d\n", cell.Cgi, cell.NumUes)
	}

	return nil
}
