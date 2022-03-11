// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package uenib

import (
	"context"
	"github.com/gogo/protobuf/types"
	"github.com/onosproject/onos-api/go/onos/uenib"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"time"
)

func getUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update {ue}",
		Short: "Update UE information",
	}
	cmd.AddCommand(getUpdateUECommand())
	return cmd
}

func getUpdateUECommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ue ue-id [args]",
		Short: "Update UE information",
		RunE:  runUpdateUECommand,
	}
	cmd.Flags().StringToStringP("aspect", "a", map[string]string{}, "UE aspect to update")
	_ = cmd.MarkFlagRequired("aspect")
	return cmd
}

func runUpdateUECommand(cmd *cobra.Command, args []string) error {
	aspects, _ := cmd.Flags().GetStringToString("aspect")

	ue := &uenib.UE{ID: uenib.ID(args[0]), Aspects: map[string]*types.Any{}}
	for aspectType, aspectValue := range aspects {
		ue.Aspects[aspectType] = &types.Any{TypeUrl: aspectType, Value: []byte(aspectValue)}
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := uenib.CreateUEServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err = client.UpdateUE(ctx, &uenib.UpdateUERequest{UE: *ue})
	if err != nil {
		cli.Output("Unable to update UE aspects: %s", err)
		return err
	}
	return nil
}
