// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getGetDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id>",
		Aliases: []string{"devices"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get a device",
		RunE:    runGetDeviceCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the device with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func outputDeviceHeader(writer io.Writer, verbose bool) {
	if verbose {
		_, _ = fmt.Fprintln(writer, "ID\tDISPLAYNAME\tADDRESS\tVERSION\tTYPE\tSTATE\tATTRIBUTES")
	} else {
		_, _ = fmt.Fprintln(writer, "ID\tDISPLAYNAME\tADDRESS\tVERSION\tTYPE\tSTATE")
	}
}

func outputDevice(obj topoapi.Object, writer io.Writer, verbose bool) {
	ent := obj.GetEntity()
	if ent != nil {
		state := stateString(ent)
		adhoc := obj.GetAspect(&topoapi.AdHoc{}).(*topoapi.AdHoc)
		asset := obj.GetAspect(&topoapi.Asset{}).(*topoapi.Asset)
		configurable := obj.GetAspect(&topoapi.Configurable{}).(*topoapi.Configurable)

		if verbose {
			attributesBuf := bytes.Buffer{}
			for key, attribute := range adhoc.Properties {
				attributesBuf.WriteString(key)
				attributesBuf.WriteString(": ")
				attributesBuf.WriteString(attribute)
				attributesBuf.WriteString(", ")
			}
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", obj.ID,
				asset.Name, configurable.Address, configurable.Version, configurable.Type, state, attributesBuf.String())
		} else {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n", obj.ID,
				asset.Name, configurable.Address, configurable.Version, configurable.Type, state)
		}

	}
}

func outputDeviceFull(obj topoapi.Object, writer io.Writer, verbose bool) {
	ent := obj.GetEntity()
	if ent != nil {
		state := stateString(ent)
		asset := obj.GetAspect(&topoapi.Asset{}).(*topoapi.Asset)
		configurable := obj.GetAspect(&topoapi.Configurable{}).(*topoapi.Configurable)

		_, _ = fmt.Fprintf(writer, "ID\t%s\n", obj.ID)
		_, _ = fmt.Fprintf(writer, "DisplayName\t%s\n", asset.Name)
		_, _ = fmt.Fprintf(writer, "ADDRESS\t%s\n", configurable.Address)
		_, _ = fmt.Fprintf(writer, "VERSION\t%s\n", configurable.Version)
		_, _ = fmt.Fprintf(writer, "TYPE\t%s\n", configurable.Type)
		_, _ = fmt.Fprintf(writer, "STATE\t%s\n", state)

		if verbose {
			adhoc := obj.GetAspect(&topoapi.AdHoc{}).(*topoapi.AdHoc)
			for key, attribute := range adhoc.Properties {
				_, _ = fmt.Fprintf(writer, "%s\t%s\n", strings.ToUpper(key), attribute)
			}
		}
	}
}

func runGetDeviceCommand(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := topoapi.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if len(args) == 0 {
		resp, err := client.List(ctx, &topoapi.ListRequest{})
		if err != nil {
			cli.Output("list error")
			return err
		}

		if !noHeaders {
			outputDeviceHeader(writer, verbose)
		}

		for _, obj := range resp.Objects {
			outputDevice(obj, writer, verbose)
		}
	} else {
		response, err := client.Get(ctx, &topoapi.GetRequest{ID: topoapi.ID(args[0])})
		if err != nil {
			cli.Output("get error")
			return err
		}
		outputDeviceFull(*response.Object, writer, verbose)
	}
	return writer.Flush()
}

func stateString(ent *topoapi.Entity) string {
	stateBuf := bytes.Buffer{}
	if ent.Protocols != nil {
		for index, protocol := range ent.Protocols {
			stateBuf.WriteString(protocol.Protocol.String())
			stateBuf.WriteString(": {Connectivity: ")
			stateBuf.WriteString(protocol.ConnectivityState.String())
			stateBuf.WriteString(", Channel: ")
			stateBuf.WriteString(protocol.ChannelState.String())
			stateBuf.WriteString(", Service: ")
			stateBuf.WriteString(protocol.ServiceState.String())
			stateBuf.WriteString("}")
			if index != len(ent.Protocols) && len(ent.Protocols) != 1 {
				stateBuf.WriteString("\n")
			}
		}
	}
	return stateBuf.String()
}

func getAddDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.ExactArgs(1),
		Short:   "Add a device",
		RunE:    runAddDeviceCommand,
	}
	cmd.Flags().StringP("type", "t", "", "the type of the device")
	cmd.Flags().StringP("role", "r", "", "the device role")
	cmd.Flags().StringP("target", "g", "", "the device target name")
	cmd.Flags().StringP("address", "a", "", "the address of the device")
	cmd.Flags().StringP("user", "u", "", "the device username")
	cmd.Flags().StringP("password", "p", "", "the device password")
	cmd.Flags().StringP("version", "v", "", "the device software version")
	cmd.Flags().StringP("displayname", "d", "", "A user friendly display name")
	cmd.Flags().String("key", "", "the TLS key")
	cmd.Flags().String("cert", "", "the TLS certificate")
	cmd.Flags().String("ca-cert", "", "the TLS CA certificate")
	cmd.Flags().Bool("plain", false, "whether to connect over a plaintext connection")
	cmd.Flags().Bool("insecure", false, "whether to enable skip verification")
	cmd.Flags().Duration("timeout", 5*time.Second, "the device connection timeout")
	cmd.Flags().StringToString("attributes", map[string]string{}, "an arbitrary mapping of device attributes")

	_ = cmd.MarkFlagRequired("version")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func runAddDeviceCommand(cmd *cobra.Command, args []string) error {
	id := args[0]
	deviceType, _ := cmd.Flags().GetString("type")
	deviceRole, _ := cmd.Flags().GetString("role")
	deviceTarget, _ := cmd.Flags().GetString("target")
	address, _ := cmd.Flags().GetString("address")
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	version, _ := cmd.Flags().GetString("version")
	displayName, _ := cmd.Flags().GetString("displayname")
	key, _ := cmd.Flags().GetString("key")
	cert, _ := cmd.Flags().GetString("cert")
	caCert, _ := cmd.Flags().GetString("ca-cert")
	plain, _ := cmd.Flags().GetBool("plain")
	insecure, _ := cmd.Flags().GetBool("insecure")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	attributes, _ := cmd.Flags().GetStringToString("attributes")

	// Target defaults to the ID
	if deviceTarget == "" {
		deviceTarget = id
	}

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	obj := &topoapi.Object{
		ID:   topoapi.ID(id),
		Type: topoapi.Object_ENTITY,
	}

	_ = obj.SetAspect(&topoapi.Asset{
		Name: displayName,
		Role: deviceRole,
	})

	_ = obj.SetAspect(&topoapi.TLSOptions{
		Plain:    plain,
		Insecure: insecure,
		Key:      key,
		CaCert:   caCert,
		Cert:     cert,
	})

	_ = obj.SetAspect(&topoapi.Configurable{
		Type:    deviceType,
		Address: address,
		Target:  deviceTarget,
		Version: version,
		Timeout: uint64(timeout.Milliseconds()),
	})

	adhoc := &topoapi.AdHoc{Properties: attributes}
	adhoc.Properties["user"] = user
	adhoc.Properties["password"] = password
	_ = obj.SetAspect(adhoc)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Create(ctx, &topoapi.CreateRequest{
		Object: obj,
	})
	if err != nil {
		return err
	}
	cli.Output("Added device %s \n", id)
	return nil
}

func getUpdateDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.ExactArgs(1),
		Short:   "Update a device",
		RunE:    runUpdateDeviceCommand,
	}
	cmd.Flags().StringP("type", "t", "", "the type of the device")
	cmd.Flags().StringP("role", "r", "", "the device role")
	cmd.Flags().StringP("target", "g", "", "the device target name")
	cmd.Flags().StringP("address", "a", "", "the address of the device")
	cmd.Flags().StringP("user", "u", "", "the device username")
	cmd.Flags().StringP("password", "p", "", "the device password")
	cmd.Flags().StringP("version", "v", "", "the device software version")
	cmd.Flags().StringP("displayname", "d", "", "A user friendly display name")
	cmd.Flags().String("key", "", "the TLS key")
	cmd.Flags().String("cert", "", "the TLS certificate")
	cmd.Flags().String("ca-cert", "", "the TLS CA certificate")
	cmd.Flags().Bool("plain", false, "whether to connect over a plaintext connection")
	cmd.Flags().Bool("insecure", false, "whether to enable skip verification")
	cmd.Flags().Duration("timeout", 30*time.Second, "the device connection timeout")
	cmd.Flags().StringToString("attributes", map[string]string{}, "an arbitrary mapping of device attributes")
	return cmd
}

func runUpdateDeviceCommand(cmd *cobra.Command, args []string) error {
	id := args[0]

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return nil
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	response, err := client.Get(ctx, &topoapi.GetRequest{
		ID: topoapi.ID(id),
	})
	if err != nil {
		return err
	}

	cancel()
	obj := response.Object

	// Ad-hoc properties
	updateAdHocAspect(cmd, obj)

	// Asset aspect properties
	updateAssetAspect(cmd, obj)

	// Configurable aspect properties
	updateConfigurableAspect(cmd, obj)

	// TLSInfo aspect properties
	updateTLSOptions(cmd, obj)

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Update(ctx, &topoapi.UpdateRequest{
		Object: obj,
	})
	if err != nil {
		return err
	}
	cli.Output("Updated device %s", id)
	return nil
}

func updateAdHocAspect(cmd *cobra.Command, obj *topoapi.Object) {
	defaultAdhoc := &topoapi.AdHoc{}
	adhoc := obj.GetAspect(defaultAdhoc).(*topoapi.AdHoc)
	if adhoc == nil {
		adhoc = defaultAdhoc
	}

	if cmd.Flags().Changed("attributes") {
		attributes, err := cmd.Flags().GetStringToString("attributes")
		if err == nil {
			adhoc.Properties = attributes
			_ = obj.SetAspect(adhoc)
		}
	}

	if adhoc.Properties == nil {
		adhoc.Properties = make(map[string]string)
	}

	if cmd.Flags().Changed("user") {
		user, err := cmd.Flags().GetString("user")
		if err == nil {
			adhoc.Properties["user"] = user
			_ = obj.SetAspect(adhoc)
		}
	}
	if cmd.Flags().Changed("password") {
		password, err := cmd.Flags().GetString("password")
		if err == nil {
			adhoc.Properties["password"] = password
			_ = obj.SetAspect(adhoc)
		}
	}
}

func updateAssetAspect(cmd *cobra.Command, obj *topoapi.Object) {
	defaultAsset := &topoapi.Asset{}
	asset := obj.GetAspect(defaultAsset).(*topoapi.Asset)
	if asset == nil {
		asset = defaultAsset
	}

	if cmd.Flags().Changed("role") {
		deviceRole, err := cmd.Flags().GetString("role")
		if err == nil {
			asset.Role = deviceRole
			_ = obj.SetAspect(asset)
		}
	}
	if cmd.Flags().Changed("displayname") {
		displayName, err := cmd.Flags().GetString("displayname")
		if err == nil {
			asset.Name = displayName
			_ = obj.SetAspect(asset)
		}
	}
}

func updateConfigurableAspect(cmd *cobra.Command, obj *topoapi.Object) {
	defaultCfg := &topoapi.Configurable{}
	cfg := obj.GetAspect(defaultCfg).(*topoapi.Configurable)
	if cfg == nil {
		cfg = defaultCfg
	}

	if cmd.Flags().Changed("type") {
		deviceType, err := cmd.Flags().GetString("type")
		if err == nil {
			cfg.Type = deviceType
			_ = obj.SetAspect(cfg)
		}
	}
	if cmd.Flags().Changed("target") {
		deviceTarget, err := cmd.Flags().GetString("target")
		if err == nil {
			cfg.Target = deviceTarget
			_ = obj.SetAspect(cfg)
		}
	}
	if cmd.Flags().Changed("address") {
		address, err := cmd.Flags().GetString("address")
		if err == nil {
			cfg.Address = address
			_ = obj.SetAspect(cfg)
		}
	}
	if cmd.Flags().Changed("version") {
		version, err := cmd.Flags().GetString("version")
		if err == nil {
			cfg.Version = version
			_ = obj.SetAspect(cfg)
		}
	}
	if cmd.Flags().Changed("timeout") {
		timeout, err := cmd.Flags().GetDuration("timeout")
		if err == nil {
			cfg.Timeout = uint64(timeout.Milliseconds())
			_ = obj.SetAspect(cfg)
		}
	}
}

func updateTLSOptions(cmd *cobra.Command, obj *topoapi.Object) {
	defaultTLS := &topoapi.TLSOptions{}
	tls := obj.GetAspect(defaultTLS).(*topoapi.TLSOptions)
	if tls == nil {
		tls = defaultTLS
	}

	if cmd.Flags().Changed("key") {
		key, err := cmd.Flags().GetString("key")
		if err == nil {
			tls.Key = key
			_ = obj.SetAspect(tls)
		}
	}
	if cmd.Flags().Changed("cert") {
		cert, err := cmd.Flags().GetString("cert")
		if err == nil {
			tls.Cert = cert
			_ = obj.SetAspect(tls)
		}
	}
	if cmd.Flags().Changed("ca-cert") {
		caCert, err := cmd.Flags().GetString("ca-cert")
		if err == nil {
			tls.CaCert = caCert
			_ = obj.SetAspect(tls)
		}
	}
	if cmd.Flags().Changed("plain") {
		plain, err := cmd.Flags().GetBool("plain")
		if err == nil {
			tls.Plain = plain
			_ = obj.SetAspect(tls)
		}
	}
	if cmd.Flags().Changed("insecure") {
		insecure, err := cmd.Flags().GetBool("insecure")
		if err == nil {
			tls.Insecure = insecure
			_ = obj.SetAspect(tls)
		}
	}
}

func getRemoveDeviceCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.ExactArgs(1),
		Short:   "Remove a device",
		RunE:    runRemoveDeviceCommand,
	}
}

func runRemoveDeviceCommand(cmd *cobra.Command, args []string) error {
	id := args[0]

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = client.Delete(ctx, &topoapi.DeleteRequest{ID: topoapi.ID(id)})
	if err != nil {
		return err
	}
	cli.Output("Removed device %s", id)
	return nil
}

func getWatchDeviceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device <id> [args]",
		Aliases: []string{"devices"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Watch for device changes",
		RunE:    runWatchDeviceCommand,
	}
	cmd.Flags().BoolP("verbose", "v", false, "whether to print the device with verbose output")
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runWatchDeviceCommand(cmd *cobra.Command, args []string) error {
	var id string
	if len(args) > 0 {
		id = args[0]
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := topoapi.CreateTopoClient(conn)

	stream, err := client.Watch(context.Background(), &topoapi.WatchRequest{})
	if err != nil {
		return err
	}

	writer := new(tabwriter.Writer)
	writer.Init(cli.GetOutput(), 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		_, _ = fmt.Fprintln(writer, "EVENT\t")
		outputDeviceHeader(writer, verbose)
		_ = writer.Flush()
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		event := response.Event
		if id != "" && event.Object.ID != topoapi.ID(id) {
			continue
		}

		_, _ = fmt.Fprintf(writer, "%s\t", event.Type)
		outputDevice(event.Object, writer, verbose)
		_ = writer.Flush()
	}
}
