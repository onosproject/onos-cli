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

func createLinkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link <src-port-id> <tgt-port-id> [field options]",
		Short: "Create a new simulated link",
		Args:  cobra.ExactArgs(2),
		RunE:  runCreateLinkCommand,
	}
	cmd.Flags().Bool("bidirectional", true, "create inverse link too")

	return cmd
}

func getLinksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "links",
		Short: "Get all simulated links",
		RunE:  runGetLinksCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func deleteLinkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link {<id>|<src-port-id>} [<tgt-port-id>",
		Short: "Delete a simulated link",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runDeleteLinkCommand,
	}
	cmd.Flags().Bool("bidirectional", true, "create inverse link too")

	return cmd
}

//func getLinkCommand() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "link <id>",
//		Args:  cobra.ExactArgs(1),
//		Short: "Get a simulated link",
//		RunE:  runGetLinkCommand,
//	}
//	return cmd
//}

func getLinkClient(cmd *cobra.Command) (simapi.LinkServiceClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}
	return simapi.NewLinkServiceClient(conn), conn, nil
}

func runCreateLinkCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getLinkClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	srcPortID := simapi.PortID(args[0])
	tgtPortID := simapi.PortID(args[1])
	id := simapi.NewLinkID(srcPortID, tgtPortID)

	if err := createLink(client, id, srcPortID, tgtPortID); err != nil {
		return err
	}

	bilink, _ := cmd.Flags().GetBool("bidirectional")
	if bilink {
		id = simapi.NewLinkID(tgtPortID, srcPortID)
		if err := createLink(client, id, tgtPortID, srcPortID); err != nil {
			return err
		}
	}
	return nil
}

func createLink(client simapi.LinkServiceClient, id simapi.LinkID, srcPortID simapi.PortID, tgtPortID simapi.PortID) error {
	_, err := client.AddLink(context.Background(), &simapi.AddLinkRequest{
		Link: &simapi.Link{
			ID:    id,
			SrcID: srcPortID,
			TgtID: tgtPortID,
		},
	})
	if err != nil {
		cli.Output("Unable to create link: %+v", err)
		return err
	}
	return nil
}

func runGetLinksCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getLinkClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	if !noHeaders {
		cli.Output("%-24s %-16s %-16s %10s\n", "ID", "Src Port ID", "Tgt Port ID", "Status")
	}

	resp, err := client.GetLinks(context.Background(), &simapi.GetLinksRequest{})
	if err != nil {
		return err
	}

	for _, link := range resp.Links {
		cli.Output("%-24s %-16s %-16s %10s\n", link.ID, link.SrcID, link.TgtID, link.Status)
	}
	return nil
}

func runDeleteLinkCommand(cmd *cobra.Command, args []string) error {
	client, conn, err := getLinkClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	id := simapi.LinkID(args[0])
	if len(args) > 1 {
		id = simapi.NewLinkID(simapi.PortID(args[0]), simapi.PortID(args[1]))
	}

	if _, err = client.RemoveLink(context.Background(), &simapi.RemoveLinkRequest{ID: id}); err != nil {
		cli.Output("Unable to remove link: %+v", err)
	}
	return err
}
