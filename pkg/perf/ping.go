// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package perf

import (
	"context"
	"github.com/aybabtme/uniplot/histogram"
	"io"
	"math/rand"
	"time"

	perfapi "github.com/onosproject/onos-api/go/onos/perf"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getPingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping",
		Args:  cobra.ExactArgs(0),
		Short: "Ping request/response",
		RunE:  runPingCommand,
	}
	cmd.Flags().IntP("count", "c", 1000, "number of ping invocations")
	cmd.Flags().Uint32P("size", "s", 4096, "size of the payload")
	return cmd
}

func getStreamCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stream",
		Aliases: []string{"relations"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Ping request/response over a stream",
		RunE:    runPingStreamCommand,
	}
	cmd.Flags().IntP("count", "c", 1000, "number of ping invocations")
	cmd.Flags().Uint32P("size", "s", 4096, "size of the payload")
	cmd.Flags().Uint32P("responses", "r", 1, "number of responses for each request")
	return cmd
}

func runPingCommand(cmd *cobra.Command, _ []string) error {
	count, _ := cmd.Flags().GetInt("count")
	size, _ := cmd.Flags().GetUint32("size")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := perfapi.NewPerfServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	payload := createPayload(size)
	latencies := make([]float64, 0, count)
	for i := 0; i < count; i++ {
		response, err := client.Ping(ctx, &perfapi.PingRequest{Payload: payload, Timestamp: uint64(time.Now().UnixNano())})
		if err != nil {
			cli.Output("get error")
			return err
		}
		latencies = append(latencies, float64(uint64(time.Now().UnixNano())-response.Timestamp)/1000.0)
	}
	analyzeLatencies(latencies)
	return nil
}

func runPingStreamCommand(cmd *cobra.Command, _ []string) error {
	count, _ := cmd.Flags().GetInt("count")
	size, _ := cmd.Flags().GetUint32("size")
	responses, _ := cmd.Flags().GetUint32("responses")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := perfapi.NewPerfServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	latencies := make([]float64, 0, count*int(responses))
	stream, err := client.PingStream(ctx)
	if err != nil {
		return err
	}

	waitc := make(chan struct{})
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				cli.Output("get error: %+v", err)
				return
			}
			latencies = append(latencies, float64(uint64(time.Now().UnixNano())-response.Timestamp)/1000.0)
		}
	}()

	for i := 0; i < count; i++ {
		payload := createPayload(size)
		err = stream.Send(&perfapi.PingRequest{Payload: payload, Timestamp: uint64(time.Now().UnixNano()), RepeatCount: responses})
		if err != nil {
			cli.Output("get error: %+v", err)
			return err
		}
	}
	err = stream.CloseSend()
	<-waitc

	analyzeLatencies(latencies)
	return err
}

func createPayload(size uint32) *perfapi.Data {
	payload := &perfapi.Data{
		Length: size,
		Data:   make([]byte, size),
	}
	rand.Read(payload.Data)
	return payload
}

func analyzeLatencies(latencies []float64) {
	hist := histogram.Hist(20, latencies)
	_ = histogram.Fprintf(cli.GetOutput(), hist, histogram.Linear(50), func(v float64) string {
		return time.Duration(v).String()
	})
}
