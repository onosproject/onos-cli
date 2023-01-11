// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package discovery

import (
	"context"
	"github.com/onosproject/onos-api/go/onos/discovery"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	podIDFlag  = "pod"
	rackIDFlag = "rack"

	p4endpointFlag       = "p4rt-endpoint"
	gNMIendpointFlag     = "gnmi-endpoint"
	pipelineConfigIDFlag = "pipeline-config"
	chassisConfigIDFlag  = "chassis-config"
)

func getAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add {pod|rack|switch|ipu} <id> [args]",
		Short: "Add new topology discovery seed entity",
	}

	cmd.AddCommand(getAddPodCommand())
	cmd.AddCommand(getAddRackCommand())
	cmd.AddCommand(getAddSwitchCommand())
	cmd.AddCommand(getAddServerIPUCommand())
	return cmd
}

func getAddPodCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pod <id>",
		Short: "Add a new pod",
		Args:  cobra.ExactArgs(1),
		RunE:  runAddPodCommand,
	}
	return cmd
}

func getAddRackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rack <id>",
		Short: "Add a new rack to a pod",
		Args:  cobra.ExactArgs(1),
		RunE:  runAddRackCommand,
	}
	cmd.Flags().String(podIDFlag, "", "ID of the parent pod")
	return cmd
}

func getAddSwitchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch <id>",
		Short: "Add a new switch to a rack",
		Args:  cobra.ExactArgs(1),
		RunE:  runAddSwitchCommand,
	}
	addStratumFlags(cmd)
	return cmd
}

func getAddServerIPUCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ipu <id>",
		Short:   "Add a new server with an IPU to a rack",
		Aliases: []string{"server", "server-ipu"},
		Args:    cobra.ExactArgs(1),
		RunE:    runAddServerIPUCommand,
	}
	addStratumFlags(cmd)
	return cmd
}

func addStratumFlags(cmd *cobra.Command) {
	cmd.Flags().String(podIDFlag, "", "ID of the parent pod")
	cmd.Flags().String(rackIDFlag, "", "ID of the parent rack")
	cmd.Flags().String(p4endpointFlag, "", "P4Runtime endpoint as host:port")
	cmd.Flags().String(gNMIendpointFlag, "", "gNMI endpoint as host:port")
	cmd.Flags().String(pipelineConfigIDFlag, "", "pipeline configuration ID")
	cmd.Flags().String(chassisConfigIDFlag, "", "chassis configuration ID")
}

func getDiscoveryClient(cmd *cobra.Command) (discovery.DiscoveryServiceClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}
	return discovery.NewDiscoveryServiceClient(conn), conn, nil
}

func runAddPodCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getDiscoveryClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.AddPod(context.Background(), &discovery.AddPodRequest{ID: args[0]})
	return err
}

func runAddRackCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getDiscoveryClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.AddRack(context.Background(), &discovery.AddRackRequest{
		ID:    args[0],
		PodID: getFlag(cmd, podIDFlag),
	})
	return err
}

func runAddSwitchCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getDiscoveryClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.AddSwitch(context.Background(), &discovery.AddSwitchRequest{
		ID:               args[0],
		PodID:            getFlag(cmd, podIDFlag),
		RackID:           getFlag(cmd, rackIDFlag),
		P4Endpoint:       getFlag(cmd, p4endpointFlag),
		GNMIEndpoint:     getFlag(cmd, gNMIendpointFlag),
		PipelineConfigID: getFlag(cmd, pipelineConfigIDFlag),
		ChassisConfigID:  getFlag(cmd, chassisConfigIDFlag),
	})
	return err
}

func runAddServerIPUCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getDiscoveryClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.AddServerIPU(context.Background(), &discovery.AddServerIPURequest{
		ID:               args[0],
		PodID:            getFlag(cmd, podIDFlag),
		RackID:           getFlag(cmd, rackIDFlag),
		P4Endpoint:       getFlag(cmd, p4endpointFlag),
		GNMIEndpoint:     getFlag(cmd, gNMIendpointFlag),
		PipelineConfigID: getFlag(cmd, pipelineConfigIDFlag),
		ChassisConfigID:  getFlag(cmd, chassisConfigIDFlag),
	})
	return err
}

func getFlag(cmd *cobra.Command, flag string) string {
	v, _ := cmd.Flags().GetString(flag)
	return v
}
