// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ransim

import (
	"context"
	"fmt"
	"strconv"

	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func getPlmnIDCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plmnid",
		Short: "Get the PLMNID",
		RunE:  runGetPlmnIDCommand,
	}
	cmd.Flags().BoolP("hex", "x", false, "show PLMNID in hex")
	return cmd
}

func getNodesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Get all E2 nodes",
		RunE:  runGetNodesCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().BoolP("watch", "w", false, "watch node changes")
	return cmd
}

func createNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node <enbid> [field options]",
		Args:  cobra.ExactArgs(1),
		Short: "Create an E2 node",
		RunE:  runCreateNodeCommand,
	}
	cmd.Flags().UintSlice("cells", []uint{}, "cell NCGIs")
	cmd.Flags().StringSlice("service-models", []string{}, "supported service models")
	cmd.Flags().StringSlice("controllers", []string{}, "E2T controller")
	return cmd
}

func getNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Get an E2 node",
		RunE:  runGetNodeCommand,
	}
	return cmd
}

func updateNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node <enbid> [field options]",
		Args:  cobra.ExactArgs(1),
		Short: "Update an E2 node",
		RunE:  runUpdateNodeCommand,
	}
	cmd.Flags().UintSlice("cells", []uint{}, "cell NCGIs")
	cmd.Flags().StringSlice("service-models", []string{}, "supported service models")
	cmd.Flags().StringSlice("controllers", []string{}, "E2T controller")
	return cmd
}

func deleteNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete an E2 node",
		RunE:  runDeleteNodeCommand,
	}
	return cmd
}

func startNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Start E2 node agent",
		RunE:  runStartNodeCommand,
	}
	return cmd
}

func stopNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop <enbid>",
		Args:  cobra.ExactArgs(1),
		Short: "Stop E2 node agent",
		RunE:  runStopNodeCommand,
	}
	return cmd
}

func getNodeClient(cmd *cobra.Command) (modelapi.NodeModelClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}
	return modelapi.NewNodeModelClient(conn), conn, nil
}

func runGetPlmnIDCommand(cmd *cobra.Command, _ []string) error {
	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	resp, err := client.GetPlmnID(context.Background(), &modelapi.PlmnIDRequest{})
	if err != nil {
		return err
	}
	if hex, _ := cmd.Flags().GetBool("hex"); hex {
		cli.Output("%x\n", resp.PlmnID)
	} else {
		cli.Output("%d\n", resp.PlmnID)
	}
	return nil
}

func runGetNodesCommand(cmd *cobra.Command, _ []string) error {
	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	if noHeaders, _ := cmd.Flags().GetBool("no-headers"); !noHeaders {
		cli.Output("%-16s %-8s %-16s %-20s %s\n", "GnbID", "Status", "Service Models", "E2T Controllers", "Cell NCGIs")
	}

	if watch, _ := cmd.Flags().GetBool("watch"); watch {
		stream, err := client.WatchNodes(context.Background(), &modelapi.WatchNodesRequest{NoReplay: false})
		if err != nil {
			return err
		}
		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			node := r.Node
			cli.Output("%-16x %-8s %-16s %-20s %s\n", node.GnbID, node.Status,
				catStrings(node.ServiceModels), catStrings(node.Controllers), catNCGIs(node.CellNCGIs))
		}

	} else {

		stream, err := client.ListNodes(context.Background(), &modelapi.ListNodesRequest{})
		if err != nil {
			return err
		}

		for {
			r, err := stream.Recv()
			if err != nil {
				break
			}
			node := r.Node
			cli.Output("%-16x %-8s %-16s %-20s %s\n", node.GnbID, node.Status,
				catStrings(node.ServiceModels), catStrings(node.Controllers), catNCGIs(node.CellNCGIs))
		}
	}

	return nil
}

func optionsToNode(cmd *cobra.Command, node *types.Node, update bool) (*types.Node, error) {
	cells, _ := cmd.Flags().GetUintSlice("cells")
	if !update || cmd.Flags().Changed("cells") {
		node.CellNCGIs = toNCGIs(cells)
	}

	models, _ := cmd.Flags().GetStringSlice("service-models")
	if !update || cmd.Flags().Changed("service-models") {
		node.ServiceModels = models
	}

	controllers, _ := cmd.Flags().GetStringSlice("controllers")
	if !update || cmd.Flags().Changed("controllers") {
		node.Controllers = controllers
	}
	return node, nil
}

func runCreateNodeCommand(cmd *cobra.Command, args []string) error {
	enbid, err := strconv.ParseUint(args[0], 16, 64)
	if err != nil {
		return err
	}

	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	node, err := optionsToNode(cmd, &types.Node{GnbID: types.GnbID(enbid)}, false)
	if err != nil {
		return err
	}

	_, err = client.CreateNode(context.Background(), &modelapi.CreateNodeRequest{Node: node})
	if err != nil {
		return err
	}
	cli.Output("Node %d created\n", enbid)
	return nil
}

func runUpdateNodeCommand(cmd *cobra.Command, args []string) error {
	enbid, err := strconv.ParseUint(args[0], 16, 64)
	if err != nil {
		return err
	}

	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Get the node first to prime the update node with existing values and allow sparse update
	gres, err := client.GetNode(context.Background(), &modelapi.GetNodeRequest{GnbID: types.GnbID(enbid)})
	if err != nil {
		return err
	}

	node, err := optionsToNode(cmd, gres.Node, true)
	if err != nil {
		return err
	}

	_, err = client.UpdateNode(context.Background(), &modelapi.UpdateNodeRequest{Node: node})
	if err != nil {
		return err
	}
	cli.Output("Node %d updated\n", enbid)
	return nil
}

func outputNode(node *types.Node) {
	cli.Output("GnbID: %-16d\nStatus: %s\nService Models: %s\nControllers: %s\nCell EGGIs: %s\n",
		node.GnbID, node.Status, catStrings(node.ServiceModels), catStrings(node.Controllers), catNCGIs(node.CellNCGIs))
}

func runGetNodeCommand(cmd *cobra.Command, args []string) error {
	enbid, err := strconv.ParseUint(args[0], 16, 64)
	if err != nil {
		return err
	}

	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	res, err := client.GetNode(context.Background(), &modelapi.GetNodeRequest{GnbID: types.GnbID(enbid)})
	if err != nil {
		return err
	}

	outputNode(res.Node)
	return nil
}

func runDeleteNodeCommand(cmd *cobra.Command, args []string) error {
	enbid, err := strconv.ParseUint(args[0], 16, 64)
	if err != nil {
		return err
	}

	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.DeleteNode(context.Background(), &modelapi.DeleteNodeRequest{GnbID: types.GnbID(enbid)})
	if err != nil {
		return err
	}

	cli.Output("Node %d deleted\n", enbid)
	return nil
}

func runControlCommand(command string, cmd *cobra.Command, args []string) error {
	enbid, err := strconv.ParseUint(args[0], 16, 64)
	if err != nil {
		return err
	}

	client, conn, err := getNodeClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := &modelapi.AgentControlRequest{GnbID: types.GnbID(enbid), Command: command}
	res, err := client.AgentControl(context.Background(), request)
	if err != nil {
		return err
	}
	outputNode(res.Node)
	return nil
}

func runStartNodeCommand(cmd *cobra.Command, args []string) error {
	return runControlCommand("start", cmd, args)
}

func runStopNodeCommand(cmd *cobra.Command, args []string) error {
	return runControlCommand("stop", cmd, args)
}

func toNCGIs(ids []uint) []types.NCGI {
	ecgis := make([]types.NCGI, 0, len(ids))
	for _, id := range ids {
		ecgis = append(ecgis, types.NCGI(id))
	}
	return ecgis
}

func catNCGIs(ecgis []types.NCGI) string {
	s := ""
	for _, ncgi := range ecgis {
		s = s + fmt.Sprintf(",%x", ncgi)
	}
	if len(s) > 1 {
		return s[1:]
	}
	return s
}

func catNCGIsWithOcn(ecgis []types.NCGI, ocns map[types.NCGI]int32) string {
	s := ""
	for _, ncgi := range ecgis {
		s = s + fmt.Sprintf(",%x(%d)", ncgi, ocns[ncgi])
	}
	if len(s) > 1 {
		return s[1:]
	}
	return s
}

func catStrings(strings []string) string {
	s := ""
	for _, string := range strings {
		s = s + "," + string
	}
	if len(s) > 1 {
		return s[1:]
	}
	return s
}
