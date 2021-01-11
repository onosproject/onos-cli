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
			if verbose {
				_, _ = fmt.Fprintln(writer, "ID\tDISPLAYNAME\tADDRESS\tVERSION\tTYPE\tSTATE\tUSER\tPASSWORD\tATTRIBUTES")
			} else {
				_, _ = fmt.Fprintln(writer, "ID\tDISPLAYNAME\tADDRESS\tVERSION\tTYPE\tSTATE")
			}
		}

		for _, obj := range resp.Objects {
			ent := obj.GetEntity()
			if ent != nil {
				state := stateString(ent)
				attrs := obj.Attributes
				if verbose {
					attributesBuf := bytes.Buffer{}
					for key, attribute := range attrs {
						attributesBuf.WriteString(key)
						attributesBuf.WriteString(": ")
						attributesBuf.WriteString(attribute)
						attributesBuf.WriteString(", ")
					}
					_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", obj.ID,
						attrs[topoapi.Displayname], attrs[topoapi.Address], attrs[topoapi.Version], attrs[topoapi.Type], state,
						attrs[topoapi.User], attrs[topoapi.Password], attributesBuf.String())
				} else {
					_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n", obj.ID,
						attrs[topoapi.Displayname], attrs[topoapi.Address], attrs[topoapi.Version], attrs[topoapi.Type], state)
				}
			}
		}
	} else {
		response, err := client.Get(ctx, &topoapi.GetRequest{
			ID: topoapi.ID(args[0]),
		})
		if err != nil {
			cli.Output("get error")
			return err
		}

		obj := response.Object
		ent := obj.GetEntity()
		if ent != nil {
			state := stateString(ent)
			attrs := obj.Attributes

			_, _ = fmt.Fprintf(writer, "ID\t%s\n", obj.ID)
			_, _ = fmt.Fprintf(writer, "DisplayName\t%s\n", attrs[topoapi.Displayname])
			_, _ = fmt.Fprintf(writer, "ADDRESS\t%s\n", attrs[topoapi.Address])
			_, _ = fmt.Fprintf(writer, "VERSION\t%s\n", attrs[topoapi.Version])
			_, _ = fmt.Fprintf(writer, "TYPE\t%s\n", attrs[topoapi.Type])
			_, _ = fmt.Fprintf(writer, "STATE\t%s\n", state)
			if verbose {
				_, _ = fmt.Fprintf(writer, "USER\t%s\n", attrs[topoapi.User])
				_, _ = fmt.Fprintf(writer, "PASSWORD\t%s\n", attrs[topoapi.Password])
				for key, attribute := range attrs {
					_, _ = fmt.Fprintf(writer, "%s\t%s\n", strings.ToUpper(key), attribute)
				}
			}
		}
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

	if attributes != nil {
		obj.Attributes = attributes
	} else {
		obj.Attributes = make(map[string]string)
	}
	setAttribute(obj, topoapi.Type, deviceType)
	setAttribute(obj, topoapi.Role, deviceRole)
	setAttribute(obj, topoapi.Target, deviceTarget)
	setAttribute(obj, topoapi.Address, address)
	setAttribute(obj, topoapi.User, user)
	setAttribute(obj, topoapi.Password, password)
	setAttribute(obj, topoapi.Version, version)
	setAttribute(obj, topoapi.Displayname, displayName)
	setAttribute(obj, topoapi.TLSKey, key)
	setAttribute(obj, topoapi.TLSCert, cert)
	setAttribute(obj, topoapi.TLSCaCert, caCert)
	setAttribute(obj, topoapi.TLSPlain, fmt.Sprintf("%t", plain))
	setAttribute(obj, topoapi.TLSInsecure, fmt.Sprintf("%t", insecure))
	setAttribute(obj, topoapi.Timeout, fmt.Sprintf("%d", timeout))

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

	if cmd.Flags().Changed("attributes") {
		attributes, _ := cmd.Flags().GetStringToString("attributes")
		obj.Attributes = attributes
	} else if obj.Attributes == nil {
		obj.Attributes = make(map[string]string)
	}
	attrs := obj.Attributes

	if cmd.Flags().Changed("type") {
		deviceType, _ := cmd.Flags().GetString("type")
		attrs[topoapi.Type] = deviceType
	}
	if cmd.Flags().Changed("target") {
		deviceTarget, _ := cmd.Flags().GetString("target")
		attrs[topoapi.Target] = deviceTarget
	}
	if cmd.Flags().Changed("role") {
		deviceRole, _ := cmd.Flags().GetString("role")
		attrs[topoapi.Role] = deviceRole
	}
	if cmd.Flags().Changed("address") {
		address, _ := cmd.Flags().GetString("address")
		attrs[topoapi.Address] = address
	}
	if cmd.Flags().Changed("user") {
		user, _ := cmd.Flags().GetString("user")
		attrs[topoapi.User] = user
	}
	if cmd.Flags().Changed("password") {
		password, _ := cmd.Flags().GetString("password")
		attrs[topoapi.Password] = password
	}
	if cmd.Flags().Changed("version") {
		version, _ := cmd.Flags().GetString("version")
		attrs[topoapi.Version] = version
	}
	if cmd.Flags().Changed("displayname") {
		displayName, _ := cmd.Flags().GetString("displayname")
		attrs[topoapi.Displayname] = displayName
	}
	if cmd.Flags().Changed("key") {
		key, _ := cmd.Flags().GetString("key")
		attrs[topoapi.TLSKey] = key
	}
	if cmd.Flags().Changed("cert") {
		cert, _ := cmd.Flags().GetString("cert")
		attrs[topoapi.TLSCert] = cert
	}
	if cmd.Flags().Changed("ca-cert") {
		caCert, _ := cmd.Flags().GetString("ca-cert")
		attrs[topoapi.TLSCaCert] = caCert
	}
	if cmd.Flags().Changed("plain") {
		plain, _ := cmd.Flags().GetBool("plain")
		attrs[topoapi.TLSPlain] = fmt.Sprintf("%t", plain)
	}
	if cmd.Flags().Changed("insecure") {
		insecure, _ := cmd.Flags().GetBool("insecure")
		attrs[topoapi.TLSInsecure] = fmt.Sprintf("%t", insecure)
	}
	if cmd.Flags().Changed("timeout") {
		timeout, _ := cmd.Flags().GetDuration("timeout")
		attrs[topoapi.Timeout] = fmt.Sprintf("%d", timeout)
	}

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
		if verbose {
			_, _ = fmt.Fprintln(writer, "EVENT\tID\tADDRESS\tVERSION\tUSER\tPASSWORD")
		} else {
			_, _ = fmt.Fprintln(writer, "EVENT\tID\tADDRESS\tVERSION")
		}
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

		obj := event.Object
		attrs := obj.Attributes
		if verbose {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n", event.Type, obj.ID,
				attrs[topoapi.Address], attrs[topoapi.Version], attrs[topoapi.User], attrs[topoapi.Password])
		} else {
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", event.Type, obj.ID,
				attrs[topoapi.Address], attrs[topoapi.Version])
		}
		_ = writer.Flush()
	}
}
