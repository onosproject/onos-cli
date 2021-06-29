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

package mlb

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	mlbapi "github.com/onosproject/onos-api/go/onos/mlb"
	meastype "github.com/onosproject/rrm-son-lib/pkg/model/measurement/type"
)

func getListParameters() *cobra.Command {
	cmd := &cobra.Command{
		Use: "parameters",
		Short: "Get all MLB parameters",
		RunE: runListParameters,
	}
	cmd.Flags().Bool("no-headers", false, "disable output headers")
	return cmd
}

func getListOcns() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ocns",
		Short: "Get all Ocn for all cells",
		RunE: runListOcns,
	}
	cmd.Flags().Bool("no-headers", false, "disable output headers")
	return cmd
}

func getSetParameters() *cobra.Command {
	cmd := &cobra.Command{
		Use: "parameters",
		Short: "Set MLB parameters",
		RunE: runSetParameters,
	}
	cmd.Flags().Int32("interval", int32(10), "MLB interval")
	cmd.Flags().Int32("delta-ocn", int32(3), "Delta Ocn per step")
	cmd.Flags().Int32("overload-threshold", int32(100), "Overload threshold [%]")
	cmd.Flags().Int32("target-threshold", int32(100), "Target threshold [%]")
	return cmd
}

func runListParameters(cmd *cobra.Command, args []string) error {
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
		_, _ = fmt.Fprint(writer, "Name\tValue\n")
	}

	request := mlbapi.GetMlbParamRequest{}
	client := mlbapi.NewMlbClient(conn)
	response, err := client.GetMlbParams(context.Background(), &request)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(writer, "%s\t%d\n", "interval [sec]", response.GetInterval())
	_, _ = fmt.Fprintf(writer, "%s\t%d\n", "Delta Ocn per step", response.GetDeltaOcn())
	_, _ = fmt.Fprintf(writer, "%s\t%d\n", "Overload threshold [%]", response.GetOverloadThreshold())
	_, _ = fmt.Fprintf(writer, "%s\t%d\n", "Target threshold [%]", response.GetTargetThreshold())

	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func runListOcns(cmd *cobra.Command, args []string) error {
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
		_, _ = fmt.Fprintf(writer, "sCell node ID\tsCell PLMN ID\tsCell cell ID\tsCell object ID\tnCell PLMN ID\tnCell cell ID\tOcn [dB]\n")
	}

	request := mlbapi.GetOcnRequest{}
	client := mlbapi.NewMlbClient(conn)
	response, err := client.GetOcn(context.Background(), &request)
	if err != nil {
		return err
	}

	// need to sort keys
	sortedOcnMap := getIDListSortedByString(func() []string{
		list := make([]string, 0)
		for k := range response.GetOcnMap() {
			list = append(list, k)
		}
		return list
	})
	for _, k := range sortedOcnMap {
		sCellIDs := strings.Split(k, ":")
		sCellNodeID := sCellIDs[0]
		sCellPlmnID := sCellIDs[1]
		sCellCellID := sCellIDs[2]
		sCellObjID := sCellIDs[3]
		sortedInnerOcnMap := getIDListSortedByString(func() []string{
			list := make([]string, 0)
			for fk := range response.GetOcnMap()[k].GetOcnRecord() {
				list = append(list, fk)
			}
			return list
		})
		for _, ik := range sortedInnerOcnMap {
			nCellIDs := strings.Split(ik, ":")
			nCellPlmnIDDec := nCellIDs[1]
			nCellPlmnIDUint64, err := strconv.ParseUint(nCellPlmnIDDec, 10, 64)
			if err != nil {
				return err
			}

			nCellCellIDDec := nCellIDs[2]
			nCellCellIDUint64, err := strconv.ParseUint(nCellCellIDDec, 10, 64)
			if err != nil {
				return err
			}
			ocn := meastype.QOffsetRange(response.GetOcnMap()[k].GetOcnRecord()[ik])
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%x\t%x\t%d\n",
				sCellNodeID, sCellPlmnID, sCellCellID, sCellObjID,
				nCellPlmnIDUint64, nCellCellIDUint64, ocn.GetValue())
		}
	}
	_ = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func runSetParameters(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	// Get current parameters first
	request := mlbapi.GetMlbParamRequest{}
	client := mlbapi.NewMlbClient(conn)
	response, err := client.GetMlbParams(context.Background(), &request)
	if err != nil {
		return err
	}

	interval := response.GetInterval()
	deltaOcn := response.GetDeltaOcn()
	overloadThr := response.GetOverloadThreshold()
	targetThr := response.GetTargetThreshold()

	if cmd.Flags().Changed("interval") {
		interval, err = cmd.Flags().GetInt32("interval")
	}
	if cmd.Flags().Changed("delta-ocn") {
		deltaOcn, err = cmd.Flags().GetInt32("delta-ocn")
	}
	if cmd.Flags().Changed("overload-threshold") {
		overloadThr, err = cmd.Flags().GetInt32("overload-threshold")
	}
	if cmd.Flags().Changed("target-threshold") {
		targetThr, err = cmd.Flags().GetInt32("target-threshold")
	}

	setRequest := mlbapi.SetMlbParamRequest{
		Interval: interval,
		DeltaOcn: deltaOcn,
		OverloadThreshold: overloadThr,
		TargetThreshold: targetThr,
	}
	setResponse, err := client.SetMlbParams(context.Background(), &setRequest)
	if err != nil {
		return err
	} else if !setResponse.Success {
		return fmt.Errorf("failed to set MLB parameters")
	}
	return nil
}
func getIDListSortedByString(getList func() []string) []string {
	result := getList()
	sort.Strings(result)
	return result
}