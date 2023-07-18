// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mlb

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"sort"
	"strings"
	"text/tabwriter"

	mlbapi "github.com/onosproject/onos-api/go/onos/mlb"
)

type qOffsetRange int

var strListQOffsetRange = []string{
	"QOffsetMinus24dB",
	"QOffsetMinus22dB",
	"QOffsetMinus20dB",
	"QOffsetMinus18dB",
	"QOffsetMinus16dB",
	"QOffsetMinus14dB",
	"QOffsetMinus12dB",
	"QOffsetMinus10dB",
	"QOffsetMinus8dB",
	"QOffsetMinus6dB",
	"QOffsetMinus5dB",
	"QOffsetMinus4dB",
	"QOffsetMinus3dB",
	"QOffsetMinus2dB",
	"QOffsetMinus1dB",
	"QOffset0dB",
	"QOffset1dB",
	"QOffset2dB",
	"QOffset3dB",
	"QOffset4dB",
	"QOffset5dB",
	"QOffset6dB",
	"QOffset8dB",
	"QOffset10dB",
	"QOffset12dB",
	"QOffset14dB",
	"QOffset16dB",
	"QOffset18dB",
	"QOffset20dB",
	"QOffset22dB",
	"QOffset24dB",
}

var valueListQOffsetRange = []int{
	-24, //"QOffsetMinus24dB"
	-22, //"QOffsetMinus22dB"
	-20, //"QOffsetMinus20dB"
	-18, //"QOffsetMinus18dB"
	-16, //"QOffsetMinus16dB"
	-14, //"QOffsetMinus14dB"
	-12, //"QOffsetMinus12dB"
	-10, //"QOffsetMinus10dB"
	-8,  //"QOffsetMinus8dB"
	-6,  //"QOffsetMinus6dB"
	-5,  //"QOffsetMinus5dB"
	-4,  //"QOffsetMinus4dB"
	-3,  //"QOffsetMinus3dB"
	-2,  //"QOffsetMinus2dB"
	-1,  //"QOffsetMinus1dB"
	0,   //"QOffset0dB"
	1,   //"QOffset1dB"
	2,   //"QOffset2dB"
	3,   //"QOffset3dB"
	4,   //"QOffset4dB"
	5,   //"QOffset5dB"
	6,   //"QOffset6dB"
	8,   //"QOffset8dB"
	10,  //"QOffset10dB"
	12,  //"QOffset12dB"
	14,  //"QOffset14dB"
	16,  //"QOffset16dB"
	18,  //"QOffset18dB"
	20,  //"QOffset20dB"
	22,  //"QOffset22dB"
	24,  //"QOffset24dB"
}

// GetValue is the get function for value as interface type
func (q *qOffsetRange) GetValue() interface{} {
	return valueListQOffsetRange[*q]
}

// String returns value as string type
func (q *qOffsetRange) String() string {
	return strListQOffsetRange[*q]
}

const (
	// QOffsetMinus24dB is the Q-Offset value -24 dB
	QOffsetMinus24dB qOffsetRange = iota
	// QOffsetMinus22dB is the Q-Offset value -22 dB
	QOffsetMinus22dB
	// QOffsetMinus20dB is the Q-Offset value -20 dB
	QOffsetMinus20dB
	// QOffsetMinus18dB is the Q-Offset value -18 dB
	QOffsetMinus18dB
	// QOffsetMinus16dB is the Q-Offset value -16 dB
	QOffsetMinus16dB
	// QOffsetMinus14dB is the Q-Offset value -14 dB
	QOffsetMinus14dB
	// QOffsetMinus12dB is the Q-Offset value -12 dB
	QOffsetMinus12dB
	// QOffsetMinus10dB is the Q-Offset value -10 dB
	QOffsetMinus10dB
	// QOffsetMinus8dB is the Q-Offset value -8 dB
	QOffsetMinus8dB
	// QOffsetMinus6dB is the Q-Offset value -6 dB
	QOffsetMinus6dB
	// QOffsetMinus5dB is the Q-Offset value -5 dB
	QOffsetMinus5dB
	// QOffsetMinus4dB is the Q-Offset value -4 dB
	QOffsetMinus4dB
	// QOffsetMinus3dB is the Q-Offset value -3 dB
	QOffsetMinus3dB
	// QOffsetMinus2dB is the Q-Offset value -2 dB
	QOffsetMinus2dB
	// QOffsetMinus1dB is the Q-Offset value -1 dB
	QOffsetMinus1dB
	// QOffset0dB is the Q-Offset value 0 dB
	QOffset0dB
	// QOffset1dB is the Q-Offset value 1 dB
	QOffset1dB
	// QOffset2dB is the Q-Offset value 2 dB
	QOffset2dB
	// QOffset3dB is the Q-Offset value 3 dB
	QOffset3dB
	// QOffset4dB is the Q-Offset value 4 dB
	QOffset4dB
	// QOffset5dB is the Q-Offset value 5 dB
	QOffset5dB
	// QOffset6dB is the Q-Offset value 6 dB
	QOffset6dB
	// QOffset8dB is the Q-Offset value 8 dB
	QOffset8dB
	// QOffset10dB is the Q-Offset value 10 dB
	QOffset10dB
	// QOffset12dB is the Q-Offset value 12 dB
	QOffset12dB
	// QOffset14dB is the Q-Offset value 14 dB
	QOffset14dB
	// QOffset16dB is the Q-Offset value 16 dB
	QOffset16dB
	// QOffset18dB is the Q-Offset value 18 dB
	QOffset18dB
	// QOffset20dB is the Q-Offset value 20 dB
	QOffset20dB
	// QOffset22dB is the Q-Offset value 22 dB
	QOffset22dB
	// QOffset24dB is the Q-Offset value 24 dB
	QOffset24dB
)

func getListParameters() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parameters",
		Short: "Get all MLB parameters",
		RunE:  runListParameters,
	}
	cmd.Flags().Bool("no-headers", false, "disable output headers")
	return cmd
}

func getListOcns() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ocns",
		Short: "Get all Ocn for all cells",
		RunE:  runListOcns,
	}
	cmd.Flags().Bool("no-headers", false, "disable output headers")
	return cmd
}

func getSetParameters() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parameters",
		Short: "Set MLB parameters",
		RunE:  runSetParameters,
	}
	cmd.Flags().Int32("interval", int32(10), "MLB interval")
	cmd.Flags().Int32("delta-ocn", int32(3), "Delta Ocn per step")
	cmd.Flags().Int32("overload-threshold", int32(100), "Overload threshold [%]")
	cmd.Flags().Int32("target-threshold", int32(100), "Target threshold [%]")
	return cmd
}

func runListParameters(cmd *cobra.Command, _ []string) error {
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

func runListOcns(cmd *cobra.Command, _ []string) error {
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
	sortedOcnMap := getIDListSortedByString(func() []string {
		list := make([]string, 0)
		for k := range response.GetOcnMap() {
			list = append(list, k)
		}
		return list
	})
	for _, k := range sortedOcnMap {
		sCellIDs := strings.Split(k, ":")
		sCellNodeID := fmt.Sprintf("%s:%s", sCellIDs[0], sCellIDs[1])
		sCellPlmnID := sCellIDs[2]
		sCellCellID := sCellIDs[3]
		sCellObjID := sCellIDs[4]
		key := k // to avoide scopelint error
		sortedInnerOcnMap := getIDListSortedByString(func() []string {
			list := make([]string, 0)
			for fk := range response.GetOcnMap()[key].GetOcnRecord() {
				list = append(list, fk)
			}
			return list
		})
		for _, ik := range sortedInnerOcnMap {
			nCellIDs := strings.Split(ik, ":")
			nCellPlmnID := nCellIDs[1]
			nCellCellID := nCellIDs[2]
			ocn := qOffsetRange(response.GetOcnMap()[k].GetOcnRecord()[ik])
			_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%d\n",
				sCellNodeID, sCellPlmnID, sCellCellID, sCellObjID,
				nCellPlmnID, nCellCellID, ocn.GetValue())
		}
	}
	_ = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func runSetParameters(cmd *cobra.Command, _ []string) error {
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
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("delta-ocn") {
		deltaOcn, err = cmd.Flags().GetInt32("delta-ocn")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("overload-threshold") {
		overloadThr, err = cmd.Flags().GetInt32("overload-threshold")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("target-threshold") {
		targetThr, err = cmd.Flags().GetInt32("target-threshold")
		if err != nil {
			return err
		}
	}

	setRequest := mlbapi.SetMlbParamRequest{
		Interval:          interval,
		DeltaOcn:          deltaOcn,
		OverloadThreshold: overloadThr,
		TargetThreshold:   targetThr,
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
