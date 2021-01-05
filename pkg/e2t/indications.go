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

package e2t

import (
	"context"
	"errors"
	"fmt"
	"github.com/onosproject/onos-api/go/onos/e2sub/subscription"
	"github.com/onosproject/onos-api/go/onos/e2t/admin"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/creds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	e2client "github.com/onosproject/onos-ric-sdk-go/pkg/e2"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/indication"

	"text/tabwriter"

	"github.com/spf13/cobra"
)

func getWatchIndicationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indications",
		Short: "watch indications traffic",
		RunE:  runWatchIndicationsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.PersistentFlags().String("service-address", "onos-e2sub:5150", "the gRPC endpoint")
	cmd.Flags().Duration("timeout", time.Hour, "specifies maximum wait time for new indications")
	return cmd
}

// createSubscriptionRequest make a proto-encoded request for a subscription to indication data.
// TODO : revisit this when JSON encoding is supported, and make this more general
func createSubscriptionRequest(nodeID string) (subscription.SubscriptionDetails, error) {
	return subscription.SubscriptionDetails{
		E2NodeID: subscription.E2NodeID(nodeID),
		ServiceModel: subscription.ServiceModel{
			ID: subscription.ServiceModelID("test"),
		},
		EventTrigger: subscription.EventTrigger{
			Payload: subscription.Payload{
				Encoding: subscription.Encoding_ENCODING_PROTO,
				Data:     []byte{},
			},
		},
		Actions: []subscription.Action{
			{
				ID:   100,
				Type: subscription.ActionType_ACTION_TYPE_REPORT,
				SubsequentAction: &subscription.SubsequentAction{
					Type:       subscription.SubsequentActionType_SUBSEQUENT_ACTION_TYPE_CONTINUE,
					TimeToWait: subscription.TimeToWait_TIME_TO_WAIT_ZERO,
				},
			},
		},
	}, nil
}

const (
	onosE2TAddress = "onos-e2t:5150"
)

// GetNodeIDs get list of E2 node IDs
// TODO this function should be replaced with topology API
func getNodeIDs() ([]string, error) {
	tlsConfig, err := creds.GetClientCredentials()
	var nodeIDs []string
	if err != nil {
		return []string{}, err
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	conn, err := grpc.DialContext(context.Background(), onosE2TAddress, opts...)
	if err != nil {
		return []string{}, err
	}
	adminClient := admin.NewE2TAdminServiceClient(conn)
	connections, err := adminClient.ListE2NodeConnections(context.Background(), &admin.ListE2NodeConnectionsRequest{})

	if err != nil {
		return []string{}, err
	}

	for {
		connection, err := connections.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []string{}, err
		}
		if connection != nil {
			nodeID := connection.Id
			nodeIDs = append(nodeIDs, nodeID)
		}

	}

	return nodeIDs, nil
}

func runWatchIndicationsCommand(cmd *cobra.Command, args []string) error {
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	timeout, _ := cmd.Flags().GetDuration("timeout")
	address, _ := cmd.Flags().GetString("service-address")
	tokens := strings.Split(address, ":")
	if len(tokens) != 2 {
		return errors.New("service-address must be of the form host:port")
	}
	host := tokens[0]

	port, err := strconv.Atoi(tokens[1])
	if err != nil {
		return err
	}

	clientConfig := e2client.Config{
		AppID: "subscription-test",
		SubscriptionService: e2client.ServiceConfig{
			Host: host,
			Port: port,
		},
	}

	client, err := e2client.NewClient(clientConfig)
	if err != nil {
		return err
	}

	ch := make(chan indication.Indication)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nodeIDs, err := getNodeIDs()
	if err != nil {
		return err
	}

	subReq, err := createSubscriptionRequest(nodeIDs[0])
	if err != nil {
		return err
	}

	_, err = client.Subscribe(ctx, subReq, ch)
	if err != nil {
		return err
	}

	done := false
	for !done {
		select {
		case indicationMsg := <-ch:
			_, _ = fmt.Fprintf(writer, "%v\n\n", indicationMsg)
			_ = writer.Flush()
		case <-time.After(timeout):
			done = true
		}
	}

	return nil
}
