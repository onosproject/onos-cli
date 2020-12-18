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
	"text/tabwriter"
)

func getListNumActiveUEsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "numues",
		Short: "Get the number of active UEs",
		RunE:  runListNumActiveUEsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runListNumActiveUEsCommand(cmd *cobra.Command, args []string) error {
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
		_, _ = fmt.Fprintln(writer, "Key[PLMNID, nodeID]\tnum(Active UEs)")
	}

	request := kpimonapi.GetRequest{
		Id: "kpimon",
	}

	client := kpimonapi.NewKpimonClient(conn)

	response, err := client.GetNumActiveUEs(context.Background(), &request)

	if err != nil {
		return err
	}

	for k, v := range response.GetObject().GetAttributes() {
		_, _ = fmt.Fprintf(writer, "%s\t%v\n", k, v)
	}

	_ = writer.Flush()

	return nil
}
