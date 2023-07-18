// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package ransim

import (
	"context"
	modelapi "github.com/onosproject/onos-api/go/onos/ransim/model"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
)

func loadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load",
		Args:  cobra.MaximumNArgs(1),
		Short: "Load model and/or metric data",
		RunE:  runLoadCommand,
	}
	cmd.Flags().StringSlice("data-name", []string{}, "data set names")
	cmd.Flags().StringSlice("data", []string{}, "data set file paths")
	return cmd
}

func clearCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Args:  cobra.ExactArgs(0),
		Short: "Clear the simulated nodes, cells and metrics",
		RunE:  runClearCommand,
	}
	return cmd
}

func getModelClient(cmd *cobra.Command) (modelapi.ModelServiceClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}
	return modelapi.NewModelServiceClient(conn), conn, nil
}

func runLoadCommand(cmd *cobra.Command, args []string) error {
	names, _ := cmd.Flags().GetStringSlice("data-name")
	paths, _ := cmd.Flags().GetStringSlice("data")

	if len(names) != len(paths) {
		return errors.NewInvalid("Number of 'data-name' and 'data' options must be equal")
	}

	// Handle the default case
	if len(names) == 0 {
		if len(args) == 1 {
			names = append(names, "model")
			paths = append(paths, args[0])
		} else {
			return errors.NewInvalid("At least the path of the model YAML file should be given")
		}
	}

	// Slurp all data files into memory as bytes; each as a separate data set
	dataSet := make([]*modelapi.DataSet, 0, len(names))
	for i, name := range names {
		var err error
		var data []byte
		if paths[i] == "-" {
			data, err = ioutil.ReadAll(os.Stdin)
		} else {
			data, err = ioutil.ReadFile(paths[i])
		}
		if err != nil {
			return err
		}
		dataSet = append(dataSet, &modelapi.DataSet{Type: name, Data: data})
	}

	client, conn, err := getModelClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.Load(context.Background(), &modelapi.LoadRequest{DataSet: dataSet, Resume: true})
	if err != nil {
		return err
	}
	return nil
}

func runClearCommand(cmd *cobra.Command, _ []string) error {
	client, conn, err := getModelClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.Clear(context.Background(), &modelapi.ClearRequest{Resume: true})
	if err != nil {
		return err
	}
	return nil
}
