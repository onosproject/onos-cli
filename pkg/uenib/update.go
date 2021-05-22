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
	"github.com/gogo/protobuf/types"
	"github.com/onosproject/onos-api/go/onos/uenib"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"time"
)

func getUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create {ue}",
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
