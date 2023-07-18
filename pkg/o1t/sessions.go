// SPDX-FileCopyrightText: 2020-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package o1t

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	o1tapi "github.com/onosproject/onos-api/go/onos/o1t"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getListSessionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sessions",
		Short: "Get sessions",
		RunE:  runListSessionsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchSessionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sessions",
		Short: "Watch sessions",
		RunE:  runWatchSessionsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runListSessionsCommand(cmd *cobra.Command, _ []string) error {
	results := make(map[string]string)

	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	request := o1tapi.GetRequest{}
	client := o1tapi.NewNetconfSessionsClient(conn)

	respSessions, err := client.List(context.Background(), &request)
	if err != nil {
		return err
	}

	for _, sessionItem := range respSessions.GetSessions() {
		for _, sessionOp := range sessionItem.GetOperations() {

			sessionID := sessionItem.GetSessionID()
			timestamp := sessionOp.GetTimestamp()

			operation := sessionOp.GetName()
			namespace := sessionOp.GetNamespace()
			status := sessionOp.GetStatus()

			timeObj := time.Unix(0, int64(timestamp))
			tsFormat := fmt.Sprintf("%02d:%02d:%02d.%d", timeObj.Hour(), timeObj.Minute(), timeObj.Second(), timeObj.Nanosecond()/1000000)

			key := fmt.Sprintf("%s/%s", sessionID, tsFormat)
			value := fmt.Sprintf("%s/%s/%v", namespace, operation, status)
			results[key] = value

		}
	}

	// types := []string{"timestamp", "session", "namespace", "operation", "status"}

	header := fmt.Sprintf("%-15s %25s %25s %20s \t%s", "Time", "Operation", "Namespace", "Status", "Session")

	if !noHeaders {
		_, _ = fmt.Fprintln(writer, header)
	}

	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, keyID := range keys {
		valueID := results[keyID]

		ids := strings.Split(keyID, "/")
		values := strings.Split(valueID, "/")

		sessionID, timestamp := ids[0], ids[1]
		namespace, operation, status := values[0], values[1], values[2]

		resultLine := fmt.Sprintf("%-15s %25s %25s %20s \t%s", timestamp, operation, namespace, status, sessionID)

		_, _ = fmt.Fprintln(writer, resultLine)
		_ = writer.Flush()
	}
	return nil
}

func runWatchSessionsCommand(cmd *cobra.Command, _ []string) error {
	var results map[string]string

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

	request := o1tapi.GetRequest{}
	client := o1tapi.NewNetconfSessionsClient(conn)

	respWatchMeasurement, err := client.Watch(context.Background(), &request)
	if err != nil {
		return err
	}

	for {
		respSessions, err := respWatchMeasurement.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
		results = make(map[string]string)

		for _, sessionItem := range respSessions.GetSessions() {
			for _, sessionOp := range sessionItem.GetOperations() {

				sessionID := sessionItem.GetSessionID()
				timestamp := sessionOp.GetTimestamp()

				operation := sessionOp.GetName()
				namespace := sessionOp.GetNamespace()
				status := sessionOp.GetStatus()

				timeObj := time.Unix(0, int64(timestamp))
				tsFormat := fmt.Sprintf("%02d:%02d:%02d.%d", timeObj.Hour(), timeObj.Minute(), timeObj.Second(), timeObj.Nanosecond()/1000000)

				key := fmt.Sprintf("%s/%s", sessionID, tsFormat)
				value := fmt.Sprintf("%s/%s/%v", namespace, operation, status)
				results[key] = value
			}
		}

		// types := []string{"timestamp", "session", "namespace", "operation", "status"}

		if !headerPrinted {
			header := fmt.Sprintf("\n%-15s %25s %25s %20s \t%s", "Time", "Operation", "Namespace", "Status", "Session")

			if !noHeaders {
				_, _ = fmt.Fprintln(writer, header)
			}
		}
		// headerPrinted = true

		keys := make([]string, 0, len(results))
		for k := range results {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, keyID := range keys {
			valueID := results[keyID]

			ids := strings.Split(keyID, "/")
			values := strings.Split(valueID, "/")

			sessionID, timestamp := ids[0], ids[1]
			namespace, operation, status := values[0], values[1], values[2]

			resultLine := fmt.Sprintf("%-15s %25s %25s %20s \t%s", timestamp, operation, namespace, status, sessionID)

			_, _ = fmt.Fprintln(writer, resultLine)
			_ = writer.Flush()
		}

	}
	return nil
}
