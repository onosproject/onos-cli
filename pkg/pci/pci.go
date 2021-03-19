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
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"strings"
	"text/tabwriter"

	pciapi "github.com/onosproject/onos-api/go/onos/pci"
)

func getListNumConflicts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "numconflicts",
		Short: "Get the number of conflicts for a specific cell",
		RunE:  runListNumConflicts,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getListNumConflictsAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "numconflicts",
		Short: "Get the number of conflicts for all cells",
		RunE:  runListNumConflictsAll,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getListNeighbors() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neighbors",
		Short: "Get neighbors for a specific cell",
		RunE:  runListNeighbors,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getListNeighborsAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neighbors",
		Short: "Get neighbors for all cells",
		RunE:  runListNeighborsAll,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getListMetric() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metric",
		Short: "Get the metric for a specific cell",
		RunE:  runListMetric,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getListMetricAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metric",
		Short: "Get the metrics for all cells",
		RunE:  runListMetricAll,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getListPci() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pci",
		Short: "Get the PCI for a specific cell",
		RunE:  runListPci,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getListPciAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pci",
		Short: "Get the PCIs for all cells",
		RunE:  runListPciAll,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runListNumConflicts(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "ID\tnum(conflicts)")
	}

	request := pciapi.GetRequest{
		Id: args[0],
	}

	client := pciapi.NewPciClient(conn)

	response, err := client.GetNumConflicts(context.Background(), &request)

	if err != nil {
		return err
	}

	for k, v := range response.GetObject().GetAttributes() {
		_, _ = fmt.Fprintf(writer, "%s\t%v\n", k, v)
	}

	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}

func runListNumConflictsAll(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "ID\tnum(conflicts)")
	}

	request := pciapi.GetRequest{
		Id: "pci",
	}

	client := pciapi.NewPciClient(conn)

	response, err := client.GetNumConflictsAll(context.Background(), &request)

	if err != nil {
		return err
	}

	for k, v := range response.GetObject().GetAttributes() {
		_, _ = fmt.Fprintf(writer, "%s\t%v\n", k, v)
	}

	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}

func runListNeighbors(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "ID\tNeighbors")
	}

	request := pciapi.GetRequest{
		Id: args[0],
	}

	client := pciapi.NewPciClient(conn)

	response, err := client.GetNeighbors(context.Background(), &request)

	if err != nil {
		return err
	}

	nMap := make(map[string]string)

	for k, v := range response.GetObject().GetAttributes() {
		if _, ok := nMap[strings.Split(k, ":")[0]]; !ok {
			nMap[strings.Split(k, ":")[0]] = v
		} else {
			nMap[strings.Split(k, ":")[0]] = fmt.Sprintf("%s,%s", nMap[strings.Split(k, ":")[0]], v)
		}
	}

	for k, v := range nMap {
		_, _ = fmt.Fprintf(writer, "%s\t%s\n", k, v)
	}

	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}

func runListNeighborsAll(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "ID\tNeighbors")
	}

	request := pciapi.GetRequest{
		Id: "pci",
	}

	client := pciapi.NewPciClient(conn)

	response, err := client.GetNeighborsAll(context.Background(), &request)

	if err != nil {
		return err
	}

	nMap := make(map[string]string)

	for k, v := range response.GetObject().GetAttributes() {
		if _, ok := nMap[strings.Split(k, ":")[0]]; !ok {
			nMap[strings.Split(k, ":")[0]] = v
		} else {
			nMap[strings.Split(k, ":")[0]] = fmt.Sprintf("%s,%s", nMap[strings.Split(k, ":")[0]], v)
		}
	}

	for k, v := range nMap {
		_, _ = fmt.Fprintf(writer, "%s\t%s\n", k, v)
	}

	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}

func runListMetric(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "ID\tPCI\tDlEarfcn\tCellSize")
	}

	request := pciapi.GetRequest{
		Id: args[0],
	}

	client := pciapi.NewPciClient(conn)

	response, err := client.GetMetric(context.Background(), &request)

	if err != nil {
		return err
	}

	keys := make(map[string]bool)
	pcis := make(map[string]string)
	dlearfcns := make(map[string]string)
	cellsizes := make(map[string]string)

	for k, v := range response.GetObject().GetAttributes() {
		keys[strings.Split(k, ":")[0]] = true
		if strings.Split(k, ":")[1] == "PCI" {
			pcis[strings.Split(k, ":")[0]] = v
		} else if strings.Split(k, ":")[1] == "DlEarfcn" {
			dlearfcns[strings.Split(k, ":")[0]] = v
		} else if strings.Split(k, ":")[1] == "CellSize" {
			cellsizes[strings.Split(k, ":")[0]] = v
		}
	}

	for k := range keys {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", k, pcis[k], dlearfcns[k], cellsizes[k])
	}

	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}

func runListMetricAll(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "ID\tPCI\tDlEarfcn\tCellSize")
	}

	request := pciapi.GetRequest{
		Id: "pci",
	}

	client := pciapi.NewPciClient(conn)

	response, err := client.GetMetricAll(context.Background(), &request)

	if err != nil {
		return err
	}

	keys := make(map[string]bool)
	pcis := make(map[string]string)
	dlearfcns := make(map[string]string)
	cellsizes := make(map[string]string)

	for k, v := range response.GetObject().GetAttributes() {
		keys[strings.Split(k, ":")[0]] = true
		if strings.Split(k, ":")[1] == "PCI" {
			pcis[strings.Split(k, ":")[0]] = v
		} else if strings.Split(k, ":")[1] == "DlEarfcn" {
			dlearfcns[strings.Split(k, ":")[0]] = v
		} else if strings.Split(k, ":")[1] == "CellSize" {
			cellsizes[strings.Split(k, ":")[0]] = v
		}
	}

	for k := range keys {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", k, pcis[k], dlearfcns[k], cellsizes[k])
	}

	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}

func runListPci(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "ID\tPCI")
	}

	request := pciapi.GetRequest{
		Id: args[0],
	}

	client := pciapi.NewPciClient(conn)

	response, err := client.GetPci(context.Background(), &request)

	if err != nil {
		return err
	}

	for k, v := range response.GetObject().GetAttributes() {
		_, _ = fmt.Fprintf(writer, "%s\t%s\n", k, v)
	}

	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}

func runListPciAll(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		_, _ = fmt.Fprintf(writer, "ID\tPCI")
	}

	request := pciapi.GetRequest{
		Id: "pci",
	}

	client := pciapi.NewPciClient(conn)

	response, err := client.GetPciAll(context.Background(), &request)

	if err != nil {
		return err
	}

	for k, v := range response.GetObject().GetAttributes() {
		_, _ = fmt.Fprintf(writer, "%s\t%s\n", k, v)
	}

	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}
