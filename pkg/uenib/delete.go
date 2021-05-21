// Copyright 2021-present Open Networking Foundation.
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
