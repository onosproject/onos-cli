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

package kpimon

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/onosproject/onos-ric-sdk-go/pkg/config/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/openconfig/gnmi/proto/gnmi_ext"

	"github.com/openconfig/gnmi/client"

	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/spf13/cobra"
)

const (
	onosConfigAddress  = "onos-config:5150"
	modelName          = "ric"
	modelVersion       = "1.0.0"
	reportIntervalPath = "/report_period/interval"
	target             = "onos-config"
)

func setReportIntervalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report-interval <interval>",
		Short: "Set report period interval",
		Args:  cobra.ExactArgs(1),
		RunE:  runSetReportIntervalCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd

}

// getClientCredentials returns the credentials for a service client
func getClientCredentials() (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(certs.DefaultClientCrt), []byte(certs.DefaultClientKey))
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}, nil
}

func runSetReportIntervalCommand(cmd *cobra.Command, args []string) error {
	creds, err := getClientCredentials()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dest := client.Destination{
		Addrs:   []string{onosConfigAddress},
		Target:  target,
		TLS:     creds,
		Timeout: 10 * time.Second,
	}

	conn, _ := grpc.Dial(onosConfigAddress, grpc.WithTransportCredentials(credentials.NewTLS(creds)))

	c, err := gclient.NewFromConn(ctx, conn, dest)
	if err != nil {
		return err
	}

	pbPath, err := utils.ToGNMIPath(reportIntervalPath)

	if err != nil {
		return err
	}

	interval, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return err
	}

	val := pb.TypedValue_UintVal{interval}

	update := &pb.Update{}
	update.Path = pbPath
	update.Path.Target = "onos-kpimon"
	update.Val = &pb.TypedValue{
		Value: &val,
	}
	extVersion := gnmi_ext.Extension_RegisteredExt{
		RegisteredExt: &gnmi_ext.RegisteredExtension{
			Id:  101,
			Msg: []byte(modelVersion),
		},
	}
	extModel := gnmi_ext.Extension_RegisteredExt{
		RegisteredExt: &gnmi_ext.RegisteredExtension{
			Id:  102,
			Msg: []byte(modelName),
		},
	}

	extensions := []*gnmi_ext.Extension{{Ext: &extVersion}, {Ext: &extModel}}

	request := &pb.SetRequest{
		Update:    []*pb.Update{update},
		Extension: extensions,
	}

	_, err = c.Set(context.Background(), request)
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)
	if err == nil {
		_, _ = fmt.Fprintf(writer, "Report period interval is set to %d successfully\n", interval)
	} else {
		_, _ = fmt.Fprintf(writer, "%v\n", err)
	}
	_ = writer.Flush()

	return nil

}
