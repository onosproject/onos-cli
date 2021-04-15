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
	kpimonapi "github.com/onosproject/onos-api/go/onos/kpimon"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
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
	results := make(map[string]map[string]map[string]string)
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
		tmpCid := ids[0]
		tmpPlmnID := ids[1]
		tmpEgnbID := ids[2]
		tmpMetricType := ids[3]
		tmpTimestamp := ids[4]
		tmpKey := fmt.Sprintf("%s:%s:%s", tmpCid, tmpPlmnID, tmpEgnbID)

		if _, ok1 := results[tmpKey]; !ok1 {
			results[tmpKey] = make(map[string]map[string]string)
		}
		if _, ok2 := results[tmpKey][tmpTimestamp]; !ok2 {
			results[tmpKey][tmpTimestamp] = make(map[string]string)
		}
		results[tmpKey][tmpTimestamp][tmpMetricType] = v
	}

	for key := range respGetHeader.GetObject().Attributes {
		types = append(types, key)
	}
	sort.Strings(types)

	header := "PlmnID\tegNB ID\tCell ID\tTime"

	for _, key := range types {
		tmpHeader := header
		header = fmt.Sprintf("%s\t%s", tmpHeader, key)
	}

	if !noHeaders {
		_, _ = fmt.Fprintln(writer, header)
	}

	for keyID, metrics := range results {
		// sort 2nd map with timestamp
		timeKeySlice := make([]string, 0, len(metrics))
		for timeStampKey := range metrics {
			timeKeySlice = append(timeKeySlice, timeStampKey)
		}
		sort.Strings(timeKeySlice)

		for _, timeKey := range timeKeySlice {
			timeStamp, err := strconv.ParseUint(timeKey, 10, 64)
			if err != nil {
				return err
			}
			timeObj := time.Unix(0, int64(timeStamp))
			tsFormat := fmt.Sprintf("%02d:%02d:%02d.%d", timeObj.Hour(), timeObj.Minute(), timeObj.Second(), timeObj.Nanosecond()/1000000)
			ids := strings.Split(keyID, ":")

			resultLine := fmt.Sprintf("%s\t%s\t%s\t%s", ids[1], ids[2], ids[0], tsFormat)
			for _, typeValue := range types {
				tmpResultLine := resultLine
				var tmpValue string
				if _, ok := metrics[timeKey][typeValue]; !ok {
					tmpValue = "N/A"
				} else {
					tmpValue = metrics[timeKey][typeValue]
				}
				resultLine = fmt.Sprintf("%s\t%s", tmpResultLine, tmpValue)
			}
			_, _ = fmt.Fprintln(writer, resultLine)
		}
	}

	_ = writer.Flush()

	return nil
}
