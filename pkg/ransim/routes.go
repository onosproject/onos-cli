// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ransim

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"strconv"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func getRoutesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routes",
		Short: "Get all UE routes",
		RunE:  runGetRoutesCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("watch", "w", false, "watch route changes")
	return cmd
}

func createRouteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route <imsi> [field options]",
		Args:  cobra.ExactArgs(1),
		Short: "Create an E2 node",
		RunE:  runCreateRouteCommand,
	}
	cmd.Flags().String("color", "gray", "route color")
	cmd.Flags().Float64("speed-avg", 80.0, "average speed in km/h")
	cmd.Flags().Float64("speed-stddev", 0.0, "speed std. deviation in km/h")
	cmd.Flags().Float64Slice("lat", []float64{}, "waypoint latitude")
	cmd.Flags().Float64Slice("lng", []float64{}, "waypoint longitude")
	return cmd
}

func getRouteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route <imsi>",
		Args:  cobra.ExactArgs(1),
		Short: "Get a UE route",
		RunE:  runGetRouteCommand,
	}
	return cmd
}

func deleteRouteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route <imsi>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a UE route",
		RunE:  runDeleteRouteCommand,
	}
	return cmd
}

func getRouteClient(cmd *cobra.Command) (modelapi.RouteModelClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}
	return modelapi.NewRouteModelClient(conn), conn, nil
}

func waypointsToString(points []*types.Point) string {
	s := ""
	for _, p := range points {
		s = s + fmt.Sprintf("; (%8.4f,%8.4f)", p.Lat, p.Lng)
	}
	if len(s) > 2 {
		return s[2:]
	}
	return s
}

func runGetRoutesCommand(cmd *cobra.Command, _ []string) error {
	client, conn, err := getRouteClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	if noHeaders, _ := cmd.Flags().GetBool("no-headers"); !noHeaders {
		cli.Output("%-16s %-8s %-5s -%-5s %s\n", "IMSI", "Color", "µkm/h", "∂km/h", "Waypoints")
	}

	if watch, _ := cmd.Flags().GetBool("watch"); watch {
		stream, err := client.WatchRoutes(context.Background(), &modelapi.WatchRoutesRequest{NoReplay: false})
		if err != nil {
			return err
		}
		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			route := r.Route
			cli.Output("%-16d %-8s %5.1f %5.1f %s\n", route.RouteID, route.Color,
				float64(route.SpeedAvg)/1000, float64(route.SpeedStdev)/1000, waypointsToString(route.Waypoints))
		}

	} else {

		stream, err := client.ListRoutes(context.Background(), &modelapi.ListRoutesRequest{})
		if err != nil {
			return err
		}

		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			route := r.Route
			cli.Output("%-16d %-8s %5.1f %5.1f %s\n", route.RouteID, route.Color,
				float64(route.SpeedAvg)/1000, float64(route.SpeedStdev)/1000, waypointsToString(route.Waypoints))
		}
	}

	return nil
}

func runCreateRouteCommand(cmd *cobra.Command, args []string) error {
	imsi, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}
	color, _ := cmd.Flags().GetString("color")
	speedAvg, _ := cmd.Flags().GetFloat64("speed-avg")
	speedStddev, _ := cmd.Flags().GetFloat64("speed-stddev")
	waypoints, err := waypointsFromOptions(cmd)
	if err != nil {
		return err
	}

	client, conn, err := getRouteClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	route := &types.Route{
		RouteID:    types.IMSI(imsi),
		Color:      color,
		SpeedAvg:   uint32(speedAvg * 1000),
		SpeedStdev: uint32(speedStddev * 1000),
		Waypoints:  waypoints,
	}

	_, err = client.CreateRoute(context.Background(), &modelapi.CreateRouteRequest{Route: route})
	if err != nil {
		return err
	}
	cli.Output("Route %d created\n", imsi)
	return nil
}

func waypointsFromOptions(cmd *cobra.Command) ([]*types.Point, error) {
	points := make([]*types.Point, 0)

	lats, _ := cmd.Flags().GetFloat64Slice("lat")
	lngs, _ := cmd.Flags().GetFloat64Slice("lng")
	if len(lats) != len(lngs) {
		return nil, errors.NewInvalid("lat/lng mismatch")
	}

	for i := range lats {
		points = append(points, &types.Point{Lat: lats[i], Lng: lngs[i]})
	}
	return points, nil
}

func outputRoute(route *types.Route) {
	cli.Output("IMSI: %-16d\nColor: %s\nAvg km/h: %5.1f\nStd. Dev km/h: %5.1f\nNextPoint: %d\nReverse: %t\nWaypoints: %s\n",
		route.RouteID, route.Color, float64(route.SpeedAvg)/1000, float64(route.SpeedStdev)/1000, route.NextPoint, route.Reverse,
		waypointsToString(route.Waypoints))
}

func runGetRouteCommand(cmd *cobra.Command, args []string) error {
	imsi, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getRouteClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	res, err := client.GetRoute(context.Background(), &modelapi.GetRouteRequest{IMSI: types.IMSI(imsi)})
	if err != nil {
		return err
	}

	outputRoute(res.Route)
	return nil
}

func runDeleteRouteCommand(cmd *cobra.Command, args []string) error {
	imsi, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	client, conn, err := getRouteClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.DeleteRoute(context.Background(), &modelapi.DeleteRouteRequest{IMSI: types.IMSI(imsi)})
	if err != nil {
		return err
	}

	cli.Output("Route %d deleted\n", imsi)
	return nil
}
