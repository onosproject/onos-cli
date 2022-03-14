// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package uenib

import (
	"context"
	"github.com/onosproject/onos-api/go/onos/uenib"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"time"
)

func getDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete {ue}",
		Short: "Delete UE information",
	}
	cmd.AddCommand(getDeleteUECommand())
	return cmd
}

func getDeleteUECommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ue ue-id [args]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Delete UE information",
		RunE:  runDeleteUECommand,
	}
	cmd.Flags().StringSliceP("aspect", "a", []string{}, "UE aspects to delete")
	return cmd
}

func runDeleteUECommand(cmd *cobra.Command, args []string) error {
	aspectTypes, _ := cmd.Flags().GetStringSlice("aspect")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := uenib.CreateUEServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err = client.DeleteUE(ctx, &uenib.DeleteUERequest{ID: uenib.ID(args[0]), AspectTypes: aspectTypes})
	if err != nil {
		cli.Output("Unable to delete UE aspects: %s", err)
		return err
	}
	return nil
}
