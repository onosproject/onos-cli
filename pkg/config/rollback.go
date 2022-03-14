// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

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
