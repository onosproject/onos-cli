// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package kpimon

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	prototypes "github.com/gogo/protobuf/types"
	"github.com/prometheus/common/log"

	kpimonapi "github.com/onosproject/onos-api/go/onos/kpimon"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	nodeIDHeader       = "Node ID"
	cellObjIDHeader    = "Cell Object ID"
	cellGlobalIDHeader = "Cell Global ID"
	timeHeader         = "Time"
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

func getWatchMetricsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Watch metrics",
		RunE:  runWatchMetricsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runListMetricsCommand(cmd *cobra.Command, _ []string) error {
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

	respGetMeasurement, err := client.ListMeasurements(context.Background(), &request)
	if err != nil {
		return err
	}

	attr := make(map[string]string)
	for key, measItems := range respGetMeasurement.GetMeasurements() {
		for _, measItem := range measItems.MeasurementItems {
			for _, measRecord := range measItem.MeasurementRecords {
				timeStamp := measRecord.Timestamp
				measName := measRecord.MeasurementName
				measValue := measRecord.MeasurementValue

				if _, ok := attr[measName]; !ok {
					attr[measName] = measName
				}

				if _, ok1 := results[key]; !ok1 {
					results[key] = make(map[uint64]map[string]string)
				}
				if _, ok2 := results[key][timeStamp]; !ok2 {
					results[key][timeStamp] = make(map[string]string)
				}

				var value interface{}

				switch {
				case prototypes.Is(measValue, &kpimonapi.IntegerValue{}):
					v := kpimonapi.IntegerValue{}
					err := prototypes.UnmarshalAny(measValue, &v)
					if err != nil {
						log.Warn(err)
					}
					value = v.GetValue()

				case prototypes.Is(measValue, &kpimonapi.RealValue{}):
					v := kpimonapi.RealValue{}
					err := prototypes.UnmarshalAny(measValue, &v)
					if err != nil {
						log.Warn(err)
					}
					value = v.GetValue()

				case prototypes.Is(measValue, &kpimonapi.NoValue{}):
					v := kpimonapi.NoValue{}
					err := prototypes.UnmarshalAny(measValue, &v)
					if err != nil {
						log.Warn(err)
					}
					value = v.GetValue()

				}

				results[key][timeStamp][measName] = fmt.Sprintf("%v", value)
			}
		}
	}

	for key := range attr {
		types = append(types, key)
	}
	sort.Strings(types)

	header := fmt.Sprintf("%-10s %20s %20s %15s", nodeIDHeader, cellObjIDHeader, cellGlobalIDHeader, timeHeader)
	//header := fmt.Sprintf("%-10s %20s %20s", "Node ID", "Cell Object ID", "Time")

	for _, key := range types {
		tmpHeader := header
		header = fmt.Sprintf(fmt.Sprintf("%%s %%%ds", len(key)+3), tmpHeader, key)
		//header = fmt.Sprintf("%s %25s", tmpHeader, key)
	}

	if !noHeaders {
		_, _ = fmt.Fprintln(writer, header)
	}

	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, keyID := range keys {
		metrics := results[keyID]
		// sort 2nd map with timestamp
		timeKeySlice := make([]uint64, 0, len(metrics))
		for timeStampKey := range metrics {
			timeKeySlice = append(timeKeySlice, timeStampKey)
		}

		sort.Slice(timeKeySlice, func(i, j int) bool { return timeKeySlice[i] < timeKeySlice[j] })

		for _, timeStamp := range timeKeySlice {
			timeObj := time.Unix(0, int64(timeStamp))
			tsFormat := fmt.Sprintf("%02d:%02d:%02d.%d", timeObj.Hour(), timeObj.Minute(), timeObj.Second(), timeObj.Nanosecond()/1000000)

			ids := strings.Split(keyID, ":")
			e2ID, nodeID, cellID, cellGlobalID := ids[0], ids[1], ids[2], ids[3]
			resultLine := fmt.Sprintf("%-10s %20s %20s %15s", fmt.Sprintf("%s:%s", e2ID, nodeID), cellID, cellGlobalID, tsFormat)
			//resultLine := fmt.Sprintf("%-10s %20s %20s", nodeID, fmt.Sprintf("%x", cellNum), tsFormat)
			for _, typeValue := range types {
				tmpResultLine := resultLine
				var tmpValue string
				if _, ok := metrics[timeStamp][typeValue]; !ok {
					tmpValue = "N/A"
				} else {
					tmpValue = metrics[timeStamp][typeValue]
				}
				resultLine = fmt.Sprintf(fmt.Sprintf("%%s %%%ds", len(typeValue)+3), tmpResultLine, tmpValue)
			}
			_, _ = fmt.Fprintln(writer, resultLine)
		}
		_ = writer.Flush()
	}
	return nil
}

func runWatchMetricsCommand(cmd *cobra.Command, _ []string) error {
	var types []string
	var results map[string]map[uint64]map[string]string

	headerPrinted := false

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

	respWatchMeasurement, err := client.WatchMeasurements(context.Background(), &request)
	if err != nil {
		return err
	}

	for {
		respGetMeasurement, err := respWatchMeasurement.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
		results = make(map[string]map[uint64]map[string]string)

		attr := make(map[string]string)
		for key, measItems := range respGetMeasurement.GetMeasurements() {
			for _, measItem := range measItems.MeasurementItems {
				for _, measRecord := range measItem.MeasurementRecords {
					timeStamp := measRecord.Timestamp
					measName := measRecord.MeasurementName
					measValue := measRecord.MeasurementValue

					if _, ok := attr[measName]; !ok {
						attr[measName] = measName
					}

					if _, ok1 := results[key]; !ok1 {
						results[key] = make(map[uint64]map[string]string)
					}
					if _, ok2 := results[key][timeStamp]; !ok2 {
						results[key][timeStamp] = make(map[string]string)
					}

					var value interface{}

					switch {
					case prototypes.Is(measValue, &kpimonapi.IntegerValue{}):
						v := kpimonapi.IntegerValue{}
						err := prototypes.UnmarshalAny(measValue, &v)
						if err != nil {
							log.Warn(err)
						}
						value = v.GetValue()

					case prototypes.Is(measValue, &kpimonapi.RealValue{}):
						v := kpimonapi.RealValue{}
						err := prototypes.UnmarshalAny(measValue, &v)
						if err != nil {
							log.Warn(err)
						}
						value = v.GetValue()

					case prototypes.Is(measValue, &kpimonapi.NoValue{}):
						v := kpimonapi.NoValue{}
						err := prototypes.UnmarshalAny(measValue, &v)
						if err != nil {
							log.Warn(err)
						}
						value = v.GetValue()

					}

					results[key][timeStamp][measName] = fmt.Sprintf("%v", value)
				}
			}
		}

		types = []string{}

		for key := range attr {
			types = append(types, key)
		}
		sort.Strings(types)

		if !headerPrinted {

			header := fmt.Sprintf("%-10s %20s %20s %15s", nodeIDHeader, cellObjIDHeader, cellGlobalIDHeader, timeHeader)

			for _, key := range types {
				tmpHeader := header
				header = fmt.Sprintf(fmt.Sprintf("%%s %%%ds", len(key)+3), tmpHeader, key)
			}

			if !noHeaders {
				_, _ = fmt.Fprintln(writer, header)
				_ = writer.Flush()
			}
		}

		headerPrinted = true

		keys := make([]string, 0, len(results))
		for k := range results {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, keyID := range keys {
			metrics := results[keyID]
			// sort 2nd map with timestamp
			timeKeySlice := make([]uint64, 0, len(metrics))
			for timeStampKey := range metrics {
				timeKeySlice = append(timeKeySlice, timeStampKey)
			}

			sort.Slice(timeKeySlice, func(i, j int) bool { return timeKeySlice[i] < timeKeySlice[j] })

			for _, timeStamp := range timeKeySlice {
				timeObj := time.Unix(0, int64(timeStamp))
				tsFormat := fmt.Sprintf("%02d:%02d:%02d.%d", timeObj.Hour(), timeObj.Minute(), timeObj.Second(), timeObj.Nanosecond()/1000000)

				ids := strings.Split(keyID, ":")
				e2id, nodeID, cellID, cellGlobalID := ids[0], ids[1], ids[2], ids[3]
				resultLine := fmt.Sprintf("%-10s %20s %20s %15s", fmt.Sprintf("%s:%s", e2id, nodeID), cellID, cellGlobalID, tsFormat)

				for _, typeValue := range types {
					tmpResultLine := resultLine

					var tmpValue string
					if _, ok := metrics[timeStamp][typeValue]; !ok {
						tmpValue = "N/A"
					} else {
						tmpValue = metrics[timeStamp][typeValue]
					}
					resultLine = fmt.Sprintf(fmt.Sprintf("%%s %%%ds", len(typeValue)+3), tmpResultLine, tmpValue)
				}
				_, _ = fmt.Fprintln(writer, resultLine)

			}
			_ = writer.Flush()
		}

	}
	return nil
}
