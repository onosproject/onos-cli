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

package kpimon

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	kpimonapi "github.com/onosproject/onos-api/go/onos/kpimon"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getListMetricsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Get metrics",
		RunE:  runListMetricsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runListMetricsCommand(cmd *cobra.Command, args []string) error {
	results := make(map[string]map[string]string)
	var types []string

	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	request := kpimonapi.GetRequest{
		Id: "kpimon",
	}
	client := kpimonapi.NewKpimonClient(conn)
	respGetHeader, err := client.GetMetricTypes(context.Background(), &request)
	if err != nil {
		return err
	}
	respGetMetrics, err := client.GetMetrics(context.Background(), &request)
	if err != nil {
		return err
	}

	for k, v := range respGetMetrics.GetObject().GetAttributes() {
		ids := strings.Split(k, ":")
		cellId := fmt.Sprintf("%s:%s", ids[0], ids[1])
		tmpMetricType := ids[2]
		if _, ok := results[cellId]; !ok {
			results[cellId] = make(map[string]string)
		}
		results[cellId][tmpMetricType] = v
	}

	for key := range respGetHeader.GetObject().Attributes {
		types = append(types, key)
	}
	sort.Strings(types)

	header := "Cell ID"

	for _, key := range types {
		tmpHeader := header
		header = fmt.Sprintf("%s\t%s", tmpHeader, key)
	}

	if !noHeaders {
		_, _ = fmt.Fprintln(writer, header)
	}

	for k1, v1 := range results {
		resultLine := k1
		for _, v2 := range types {
			tmpResultLine := resultLine
			var tmpValue string
			if _, ok := v1[v2]; !ok {
				tmpValue = "N/A"
			} else {
				tmpValue = v1[v2]
			}
			resultLine = fmt.Sprintf("%s\t%s", tmpResultLine, tmpValue)
		}
		_, _ = fmt.Fprintln(writer, resultLine)
	}

	_ = writer.Flush()

	return nil
}
