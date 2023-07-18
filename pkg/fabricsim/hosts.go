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
	"sort"
)

func createHostCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host <id> [field options]",
		Short: "Create a new simulated host",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreateHostCommand,
	}
	// TODO: Add appropriate options

	return cmd
}

func getHostsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hosts",
		Short: "Get all simulated hosts",
		RunE:  runGetHostsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().Bool("no-nics", false, "disables listing of NICs")
	return cmd
}

func deleteHostCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host <id>",
		Short: "Delete a simulated host",
		Args:  cobra.ExactArgs(1),
		RunE:  runDeleteHostCommand,
	}

	return cmd
}

func getHostCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get a simulated host",
		RunE:  runGetHostCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().Bool("no-nics", false, "disables listing of NICs")
	return cmd
}

func emitARPsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "arp <id> <mac> <ip1> <ip2>...",
		Args:  cobra.MinimumNArgs(3),
		Short: "Emit ARP request(s) via specified host NIC",
		RunE:  runEmitARPsCommand,
	}
	return cmd

}

func getHostClient(cmd *cobra.Command) (simapi.HostServiceClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}
	return simapi.NewHostServiceClient(conn), conn, nil
}

func runCreateHostCommand(_ *cobra.Command, _ []string) error {
	//client, conn, err := getHostClient(cmd)
	//if err != nil {
	//	return err
	//}
	//defer conn.Close()
	//
	//id := simapi.HostID(args[0])
	return nil
}

func runGetHostsCommand(cmd *cobra.Command, _ []string) error {
	client, conn, err := getHostClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noNICs, _ := cmd.Flags().GetBool("no-nics")

	printHostHeaders(noHeaders)

	resp, err := client.GetHosts(context.Background(), &simapi.GetHostsRequest{})
	if err != nil {
		return err
	}

	sort.SliceStable(resp.Hosts, func(i, j int) bool {
		return resp.Hosts[i].ID < resp.Hosts[j].ID
	})
	for _, h := range resp.Hosts {
		printHost(h, noHeaders, noNICs)
	}
	return nil
}

func runGetHostCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getHostClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	id := simapi.HostID(args[0])
	resp, err := client.GetHost(context.Background(), &simapi.GetHostRequest{ID: id})
	if err != nil {
		return err
	}

	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	noNICs, _ := cmd.Flags().GetBool("no-nics")

	printHostHeaders(noHeaders)
	printHost(resp.Host, noHeaders, noNICs)
	return nil
}

func runDeleteHostCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getHostClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.RemoveHost(context.Background(), &simapi.RemoveHostRequest{
		ID: simapi.HostID(args[0]),
	})
	if err != nil {
		cli.Output("Unable to remove host: %+v", err)
	}
	return err
}

func runEmitARPsCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getHostClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.EmitARPs(context.Background(), &simapi.EmitARPsRequest{
		ID:          simapi.HostID(args[0]),
		MacAddress:  args[1],
		IpAddresses: args[2:],
	})
	if err != nil {
		cli.Output("Unable to emit ARP requests: %+v", err)
	}
	return err
}

func printHostHeaders(noHeaders bool) {
	if !noHeaders {
		cli.Output("%-16s %10s\n", "ID", "# of NICs")
	}
}

func printHostNICHeaders(noHeaders bool) {
	if !noHeaders {
		cli.Output("\t%-16s %-18s %-15s %-24s\n", "Port ID", "MAC Address", "IPv4 Address", "IPv6 Address")
	}
}

func printHost(h *simapi.Host, noHeaders bool, noNICs bool) {
	cli.Output("%-16s %10d\n", h.ID, len(h.Interfaces))
	if !noNICs {
		printHostNICHeaders(noHeaders)
		for _, n := range h.Interfaces {
			cli.Output("\t%-16s %-18s %-15s %-24s\n", n.ID, n.MacAddress, n.IpAddress, n.Ipv6Address)
		}
	}
}
