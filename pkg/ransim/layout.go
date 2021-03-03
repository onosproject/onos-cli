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

package ransim

import (
	"context"

	simapi "github.com/onosproject/onos-api/go/onos/ransim/trafficsim"
	"github.com/onosproject/onos-lib-go/pkg/cli"

	"github.com/spf13/cobra"
)

func getLayoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "layout",
		Short: "Get Layout",
		RunE:  runGetLayoutCommand,
	}
	return cmd
}

func runGetLayoutCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := simapi.NewTrafficClient(conn)
	ml, err := client.GetMapLayout(context.Background(), &simapi.MapLayoutRequest{})
	if err != nil {
		return err
	}

	cli.Output("Center: %7.3f,%7.3f\nZoom: %5.2f\nFade: %v\nShowRoutes: %v\nShowPower: %v\nLocationsScale: %5.2f\n",
		ml.Center.Lat, ml.Center.Lng, ml.Zoom, ml.Fade, ml.ShowRoutes, ml.ShowPower, ml.LocationsScale)
	return nil
}
