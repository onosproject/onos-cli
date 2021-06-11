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
	"text/tabwriter"
	"time"

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
	var types []string
	results := make(map[string]map[uint64]map[string]string)

	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	request := kpimonapi.GetRequest{}
	client := kpimonapi.NewKpimonClient(conn)

	respGetMeasurement, err := client.GetMeasurement(context.Background(), &request)
	if err != nil {
		return err
	}

	attr := make(map[string]string)
	for cellID, measItems := range respGetMeasurement.GetMeasurements() {
		for _, measItem := range measItems.MeasurementItems {
			for _, measRecord := range measItem.MeasurementRecords {
				timeStamp := measRecord.Timestamp
				measName := measRecord.MeasurementName
				measValue := measRecord.MeasurementValue

				if _, ok := attr[measName]; !ok {
					attr[measName] = measName
				}

				if _, ok1 := results[cellID]; !ok1 {
					results[cellID] = make(map[uint64]map[string]string)
				}
				if _, ok2 := results[cellID][timeStamp]; !ok2 {
					results[cellID][timeStamp] = make(map[string]string)
				}
				results[cellID][timeStamp][measName] = fmt.Sprintf("%v", measValue)
			}
		}
	}

	for key := range attr {
		types = append(types, key)
	}
	sort.Strings(types)

	header := "Cell ID\tTime"

	for _, key := range types {
		tmpHeader := header
		header = fmt.Sprintf("%s\t%s", tmpHeader, key)
	}

	if !noHeaders {
		_, _ = fmt.Fprintln(writer, header)
	}

	for keyID, metrics := range results {
		// sort 2nd map with timestamp
		timeKeySlice := make([]uint64, 0, len(metrics))
		for timeStampKey := range metrics {
			timeKeySlice = append(timeKeySlice, timeStampKey)
		}

		sort.Slice(timeKeySlice, func(i, j int) bool { return timeKeySlice[i] < timeKeySlice[j] })

		for _, timeStamp := range timeKeySlice {
			timeObj := time.Unix(0, int64(timeStamp))
			tsFormat := fmt.Sprintf("%02d:%02d:%02d.%d", timeObj.Hour(), timeObj.Minute(), timeObj.Second(), timeObj.Nanosecond()/1000000)

			resultLine := fmt.Sprintf("%s\t%s", keyID, tsFormat)
			for _, typeValue := range types {
				tmpResultLine := resultLine
				var tmpValue string
				if _, ok := metrics[timeStamp][typeValue]; !ok {
					tmpValue = "N/A"
				} else {
					tmpValue = metrics[timeStamp][typeValue]
				}
				resultLine = fmt.Sprintf("%s\t%s", tmpResultLine, tmpValue)
			}
			_, _ = fmt.Fprintln(writer, resultLine)
		}
	}

	_ = writer.Flush()

	return nil
}
