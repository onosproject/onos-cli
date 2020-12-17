// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2sub

import (
	"context"
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	epapi "github.com/onosproject/onos-api/go/onos/e2sub/endpoint"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	endPointHeaders = "ID\tIP\tPort\n"
	endPointFormat  = "%s\t%s\t%d\n"
)

func displayEndPointHeaders(writer io.Writer) {
	_, _ = fmt.Fprint(writer, endPointHeaders)
}

func displayEndPoint(writer io.Writer, ep epapi.TerminationEndpoint) {
	_, _ = fmt.Fprintf(writer, endPointFormat, ep.ID, ep.IP, ep.Port)
}

func getListEndPointsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoints",
		Short: "Get endpoints",
		RunE:  runListEndpointsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getAddEndPointCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoint",
		Short: "Add endpoint",
		RunE:  runAddEndpointCommand,
	}
	cmd.Flags().String("IP", "", "IP address")
	cmd.Flags().Int32("port", 0, "Port number")
	cmd.Flags().String("ID", "", "Identifier")

	return cmd
}

func getRemoveEndPointCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoint",
		Short: "Remove endpoint",
		Args:  cobra.ExactArgs(1),
		RunE:  runRemoveEndpointCommand,
	}

	return cmd
}

func getGetEndPointCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoint",
		Short: "Get endpoint",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetEndpointCommand,
	}

	return cmd
}

func runListEndpointsCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		displayEndPointHeaders(writer)
	}

	request := epapi.ListTerminationsRequest{}

	client := epapi.NewE2RegistryServiceClient(conn)

	response, err := client.ListTerminations(context.Background(), &request)
	if err != nil {
		return err
	}

	for _, ep := range response.Endpoints {
		displayEndPoint(writer, ep)
	}

	_ = writer.Flush()

	return nil
}

func runAddEndpointCommand(cmd *cobra.Command, args []string) error {
	IP, _ := cmd.Flags().GetString("IP")
	if IP == "" {
		return errors.New("IP address must be specified with --IP")
	}
	ID, _ := cmd.Flags().GetString("ID")
	if ID == "" {
		return errors.New("identifier must be specified with --ID")
	}
	port, _ := cmd.Flags().GetInt32("port")
	if port == 0 {
		return errors.New("port must be specified with --port")
	}
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	ep := epapi.TerminationEndpoint{
		ID:   epapi.ID(ID),
		IP:   epapi.IP(IP),
		Port: epapi.Port(port),
	}
	request := epapi.AddTerminationRequest{Endpoint: &ep}

	client := epapi.NewE2RegistryServiceClient(conn)

	_, err = client.AddTermination(context.Background(), &request)

	return err
}

func runRemoveEndpointCommand(cmd *cobra.Command, args []string) error {
	ID := args[0]
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := epapi.RemoveTerminationRequest{ID: epapi.ID(ID)}

	client := epapi.NewE2RegistryServiceClient(conn)

	_, err = client.RemoveTermination(context.Background(), &request)

	return err
}

func runGetEndpointCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	ID := epapi.ID(args[0])
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		displayEndPointHeaders(writer)
	}

	request := epapi.GetTerminationRequest{ID: ID}

	client := epapi.NewE2RegistryServiceClient(conn)

	response, err := client.GetTermination(context.Background(), &request)
	if err != nil {
		return err
	}

	if response.Endpoint == nil {
		return errors.New("endpoint not found")
	}

	displayEndPoint(writer, *response.Endpoint)
	_ = writer.Flush()

	return nil
}
