// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package fabricsim

import (
	"context"
	"fmt"
	simapi "github.com/onosproject/onos-api/go/onos/fabricsim"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func createDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "device <id> [field options]",
		Short: "Create a new simulated device",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreateDeviceCommand,
	}
	cmd.Flags().String("type", "switch", "switch (default) or IPU")
	cmd.Flags().Uint16("agent-port", 20000, "agent gRPC (TCP) port")
	cmd.Flags().Uint16("port-count", 32, "number of ports to create; default 32")
	cmd.Flags().Bool("start-agent", true, "starts agent upon creation")

	return cmd
}

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

func deleteDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "device <id>",
		Short: "Delete a simulated device",
		Args:  cobra.ExactArgs(1),
		RunE:  runDeleteDeviceCommand,
	}

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

func runCreateDeviceCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getDeviceClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	id := simapi.DeviceID(args[0])
	dtype, _ := cmd.Flags().GetString("type")
	deviceType := simapi.DeviceType_SWITCH
	if dtype == "IPU" {
		deviceType = simapi.DeviceType_IPU
	}

	agentPort, _ := cmd.Flags().GetUint16("agent-port")

	// FIXME: This is just a quick hack to allow creating device ports en masse; implement proper creation later
	portCount, _ := cmd.Flags().GetUint16("port-count")
	ports := make([]*simapi.Port, 0, portCount)
	for pn := uint16(1); pn <= portCount; pn++ {
		ports = append(ports, &simapi.Port{
			ID:             simapi.PortID(fmt.Sprintf("%s/%d", id, pn)),
			Name:           fmt.Sprintf("%d", pn),
			Number:         uint32(pn),
			InternalNumber: uint32(1024 + pn),
			Speed:          "100Gbps",
		})
	}

	_, err = client.AddDevice(context.Background(), &simapi.AddDeviceRequest{
		Device: &simapi.Device{
			ID:          id,
			Type:        deviceType,
			Ports:       ports,
			ControlPort: int32(agentPort),
		},
	})
	if err != nil {
		cli.Output("Unable to create device: %+v", err)
		return err
	}

	startAgent, _ := cmd.Flags().GetBool("start-agent")
	if startAgent {
		_, err := client.StartDevice(context.Background(), &simapi.StartDeviceRequest{
			ID: id,
		})
		if err != nil {
			cli.Output("Unable to start device agent: %+v", err)
			return err
		}
	}
	return nil
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
		cli.Output("%-16s %-8s %-16s %10s\n", "ID", "Type", "Agent Port", "# of Ports")
	}

	resp, err := client.GetDevices(context.Background(), &simapi.GetDevicesRequest{})
	if err != nil {
		return err
	}

	for _, d := range resp.Devices {
		cli.Output("%-16s %-8s %8d %10d\n", d.ID, d.Type, d.ControlPort, len(d.Ports))
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

func runDeleteDeviceCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getDeviceClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.RemoveDevice(context.Background(), &simapi.RemoveDeviceRequest{
		ID: simapi.DeviceID(args[0]),
	})
	if err != nil {
		cli.Output("Unable to remove device: %+v", err)
	}
	return err
}
