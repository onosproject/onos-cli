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

package config

import (
	"github.com/onosproject/onos-api/go/onos/config/admin"
	v2 "github.com/onosproject/onos-api/go/onos/config/v2"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"strconv"
)

func getRollbackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback <index>",
		Short: "Rolls-back a transaction",
		Args:  cobra.ExactArgs(1),
		RunE:  runRollbackCommand,
	}
	return cmd
}

func runRollbackCommand(cmd *cobra.Command, args []string) error {
	index, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		cli.Output("Index must be a number: %+v", err)
	}

	clientConnection, clientConnectionError := cli.GetConnection(cmd)

	if clientConnectionError != nil {
		return clientConnectionError
	}
	client := admin.CreateConfigAdminServiceClient(clientConnection)

	ctx := cli.NewContextWithAuthHeaderFromFlag(cmd.Context(), cmd.Flag(cli.AuthHeaderFlag))
	resp, err := client.RollbackTransaction(ctx, &admin.RollbackRequest{Index: v2.Index(index)})
	if err != nil {
		cli.Output("Rollback failed: %+v\n", err)
		return err
	}
	cli.Output("Rollback transaction ID %s; Index %d\n", resp.ID, resp.Index)
	return nil
}
