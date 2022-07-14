// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package fabricsim

import (
	"context"
	simapi "github.com/onosproject/onos-api/go/onos/fabricsim"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func getDevicesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices",
		Short: "Get all simulated devices",
		RunE:  runGetDevicesCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().Bool("no-ports", false, "disables listing of ports")
	return cmd
}

//func getDeviceCommand() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "device <id>",
//		Args:  cobra.ExactArgs(1),
//		Short: "Get a simulated device",
//		RunE:  runGetDeviceCommand,
//	}
//	return cmd
//}

func getDeviceClient(cmd *cobra.Command) (simapi.DeviceServiceClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}
	return simapi.NewDeviceServiceClient(conn), conn, nil
}

func runGetDevicesCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getDeviceClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noPorts, _ := cmd.Flags().GetBool("no-ports")

	if !noHeaders {
		cli.Output("%-16s %-8s %-16s %5d\n", "ID", "Type", "Agent Port", "# of Ports")
	}

	resp, err := client.GetDevices(context.Background(), &simapi.GetDevicesRequest{})
	if err != nil {
		return err
	}

	for _, d := range resp.Devices {
		cli.Output("%-16s %-8s %8d\n", d.ID, d.Type, d.ControlPort, len(d.Ports))
		if !noPorts {
			if !noHeaders {
				cli.Output("\t%-16s %8s %8s %-12s %16s\n", "Port ID", "Port #", "SDN #", "Name", "Speed")
			}
			for _, p := range d.Ports {
				cli.Output("\t%-16s %8d %8d %-12s %-16s\n", p.ID, p.Number, p.InternalNumber, p.Name, p.Speed)
			}
		}
	}

	return nil
}
