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

package e2t

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	adminapi "github.com/onosproject/onos-api/go/onos/e2t/admin"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getListConnectionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connections",
		Short: "Get SB connections",
		RunE:  runListConnectionsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runListConnectionsCommand(cmd *cobra.Command, args []string) error {
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
		_, _ = fmt.Fprintln(writer, "Global ID\tPLNM ID\tIP Addr\tPort\tConn Type")
	}

	request := adminapi.ListE2NodeConnectionsRequest{}

	client := adminapi.NewE2TAdminServiceClient(conn)

	stream, err := client.ListE2NodeConnections(context.Background(), &request)
	if err != nil {
		return err
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			cli.Output("recv error")
			return err
		}

		_, _ = fmt.Fprintf(
			writer,
			"%s\t%s\t%s\t%d\t%s\n",
			response.Id, response.PlmnId, strings.Join(response.RemoteIp, ","), response.RemotePort,
			response.ConnectionType.String())
	}

	_ = writer.Flush()

	return nil
}
