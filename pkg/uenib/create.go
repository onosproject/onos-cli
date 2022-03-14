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

func getCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create {ue}",
		Short: "Create UE information",
	}
	cmd.AddCommand(getCreateUECommand())
	return cmd
}

func getCreateUECommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ue ue-id [args]",
		Short: "Create UE information",
		RunE:  runCreateUECommand,
	}
	cmd.Flags().StringToStringP("aspect", "a", map[string]string{}, "UE aspect to create")
	_ = cmd.MarkFlagRequired("aspect")
	return cmd
}

func runCreateUECommand(cmd *cobra.Command, args []string) error {
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
	_, err = client.CreateUE(ctx, &uenib.CreateUERequest{UE: *ue})
	if err != nil {
		cli.Output("Unable to create UE aspects: %s", err)
		return err
	}
	return nil
}
