// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package provisioner

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"github.com/onosproject/onos-api/go/onos/provisioner"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
	"os"
	"strings"
)

const (
	kindFlag          = "kind"
	artifactsPathFlag = "artifacts"
	noHeadersFlag     = "no-headers"
)

func getAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add config <config-id> [args]",
		Short: "Add new device configuration",
		Args:  cobra.ExactArgs(1),
		RunE:  runAddConfigCommand,
	}
	cmd.Flags().String(kindFlag, provisioner.PipelineConfigKind, "kind of configuration: pipeline or chassis")
	cmd.Flags().String(artifactsPathFlag, "-", "artifacts tar file (- for stdin)")
	return cmd
}

func getDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete config <config-id> [args]",
		Short: "Delete device configuration",
		Args:  cobra.ExactArgs(1),
		RunE:  runDeleteConfigCommand,
	}
	return cmd

}

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get config {config-id} [args]",
		Short:   "Get device configurations",
		Aliases: []string{"get configs"},
		Args:    cobra.MaximumNArgs(1),
		RunE:    runGetConfigCommand,
	}
	cmd.Flags().String(kindFlag, provisioner.PipelineConfigKind, "kind of configuration: pipeline or chassis")
	cmd.Flags().String(artifactsPathFlag, "", "artifacts tar file; - for stdin")
	cmd.Flags().Bool(noHeadersFlag, false, "disables output headers")
	return cmd
}

func getProvisionerClient(cmd *cobra.Command) (provisioner.ProvisionerServiceClient, *grpc.ClientConn, error) {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil, nil, err
	}
	return provisioner.NewProvisionerServiceClient(conn), conn, nil
}

func runAddConfigCommand(cmd *cobra.Command, args []string) error {
	configID := provisioner.ConfigID(args[0])
	kind, _ := cmd.Flags().GetString(kindFlag)
	artifactsPath, _ := cmd.Flags().GetString(artifactsPathFlag)

	artifacts, err := readArtifacts(artifactsPath)
	if err != nil {
		return err
	}

	client, conn, err := getProvisionerClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	config := &provisioner.Config{
		Record:    &provisioner.ConfigRecord{ConfigID: configID, Kind: kind},
		Artifacts: artifacts,
	}
	_, err = client.Add(context.Background(), &provisioner.AddConfigRequest{Config: config})
	return err
}

func runDeleteConfigCommand(cmd *cobra.Command, args []string) error {
	configID := provisioner.ConfigID(args[0])

	client, conn, err := getProvisionerClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.Delete(context.Background(), &provisioner.DeleteConfigRequest{ConfigID: configID})
	return err
}

func runGetConfigCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		return getConfig(cmd, provisioner.ConfigID(args[0]))
	}
	return listConfigs(cmd)
}

func listConfigs(cmd *cobra.Command) error {
	noHeaders, _ := cmd.Flags().GetBool(noHeadersFlag)

	client, conn, err := getProvisionerClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	stream, err := client.List(context.Background(), &provisioner.ListConfigsRequest{})
	if err != nil {
		return err
	}

	printConfigHeaders(noHeaders)
	for {
		resp, err1 := stream.Recv()
		if err1 != nil {
			if err1 == io.EOF {
				return nil
			}
			return err1
		}
		printConfigRecord(resp.Config.Record)
	}
}

func getConfig(cmd *cobra.Command, configID provisioner.ConfigID) error {
	noHeaders, _ := cmd.Flags().GetBool(noHeadersFlag)
	artifactsPath, _ := cmd.Flags().GetString(artifactsPathFlag)
	includeArtifacts := len(artifactsPath) > 0

	client, conn, err := getProvisionerClient(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	req := &provisioner.GetConfigRequest{ConfigID: configID, IncludeArtifacts: includeArtifacts}
	resp, err := client.Get(context.Background(), req)
	if err != nil {
		return err
	}
	if includeArtifacts {
		return writeArtifacts(artifactsPath, resp.Config.Artifacts)
	}

	printConfigHeaders(noHeaders)
	printConfigRecord(resp.Config.Record)
	return nil
}

func printConfigHeaders(noHeaders bool) {
	if !noHeaders {
		cli.Output("%-32s\t%-12s\t%s\n", "Config ID", "Kind", "Artifacts")
	}
}

func printConfigRecord(record *provisioner.ConfigRecord) {
	artifacts := strings.Join(record.Artifacts, ", ")
	cli.Output("%-32s\t%-12s\t%s\n", record.ConfigID, record.Kind, artifacts)
}

// Reads artifacts from the given gzip tar archive (stdin if "-") into an artifact map
func readArtifacts(path string) (map[string][]byte, error) {
	var input io.Reader
	var err error

	input = os.Stdin
	if path != "-" {
		input, err = os.Open(path)
		if err != nil {
			return nil, err
		}
	}

	gzr, err := gzip.NewReader(input)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	artifacts := make(map[string][]byte)
	for {
		header, err1 := tr.Next()
		switch {
		case err1 == io.EOF:
			return artifacts, nil
		case err1 != nil:
			return nil, err1
		case header == nil:
			continue
		}

		data := make([]byte, header.Size)
		for total := 0; total < int(header.Size); {
			l, err1 := tr.Read(data[total:])
			if err1 != nil && err1 != io.EOF {
				return nil, err1
			}
			total = total + l
		}
		artifacts[header.Name] = data
	}
}

// Write the configuration's artifacts into the specified file; stdout if "-"
func writeArtifacts(path string, artifacts map[string][]byte) error {
	var output io.Writer
	var err error

	output = os.Stdout
	if path != "-" {
		output, err = os.Create(path)
		if err != nil {
			return err
		}
	}

	gzw := gzip.NewWriter(output)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	for artifact, data := range artifacts {
		header := &tar.Header{
			Name: artifact,
			Size: int64(len(data)),
		}

		if err = tw.WriteHeader(header); err != nil {
			return err
		}

		if _, err = tw.Write(data); err != nil {
			return err
		}
	}
	_ = tw.Close()
	_ = gzw.Close()

	return nil
}
