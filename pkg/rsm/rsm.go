// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package rsm

import (
	"context"
	"fmt"
	rsmapi "github.com/onosproject/onos-api/go/onos/rsm"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getSetAssociation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "association",
		Short: "Set UE-Slice association",
		RunE:  runSetAssociation,
	}
	cmd.Flags().Bool("no-headers", false, "disable output headers")
	cmd.Flags().String("e2NodeID", "", "E2 Node ID")
	cmd.Flags().String("DuUeF1apID", "", "DU-UE-F1AP-ID")
	cmd.Flags().String("CuUeF1apID", "", "CU-UE-F1AP-ID")
	cmd.Flags().String("RanUeNgapID", "", "RAN-UE-NGAP-ID")
	cmd.Flags().String("eNBUeS1apID", "", "ENB-UE-S1AP-ID")
	cmd.Flags().String("dlSliceID", "", "DL Slice ID")
	cmd.Flags().String("ulSliceID", "", "UL Slice ID")
	cmd.Flags().String("drbID", "", "DRB-ID")
	return cmd
}

func getCreateSlice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slice",
		Short: "Create slice",
		RunE:  runGetCreateSlice,
	}
	cmd.Flags().Bool("no-headers", false, "disable output headers")
	cmd.Flags().String("e2NodeID", "", "E2 Node ID")
	cmd.Flags().String("sliceID", "", "Slice ID")
	cmd.Flags().String("scheduler", "", "Scheduler Type {RR, PF, QoS}")
	cmd.Flags().String("weight", "", "Weight")
	cmd.Flags().String("sliceType", "", "Slice Type {DL, UL}")
	return cmd
}

func getUpdateSlice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slice",
		Short: "Update slice",
		RunE:  runGetUpdateSlice,
	}
	cmd.Flags().Bool("no-headers", false, "disable output headers")
	cmd.Flags().String("e2NodeID", "", "E2 Node ID")
	cmd.Flags().String("sliceID", "", "Slice ID")
	cmd.Flags().String("scheduler", "", "Scheduler Type {RR, PF, QoS}")
	cmd.Flags().String("weight", "", "Weight")
	cmd.Flags().String("sliceType", "", "Slice Type {DL, UL}")
	return cmd
}

func getDeleteSlice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slice",
		Short: "Delete slice",
		RunE:  runGetDeleteSlice,
	}
	cmd.Flags().Bool("no-headers", false, "disable output headers")
	cmd.Flags().String("e2NodeID", "", "E2 Node ID")
	cmd.Flags().String("sliceID", "", "Slice ID")
	cmd.Flags().String("sliceType", "", "Slice Type {DL, UL}")
	return cmd
}

func runSetAssociation(cmd *cobra.Command, _ []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := rsmapi.NewRsmClient(conn)
	e2NodeID := ""
	duUeF1apID := ""
	cuUeF1apID := ""
	ranUeNgapID := ""
	enbUeS1apID := ""
	dlSliceID := ""
	ulSliceID := ""
	drbID := ""

	if cmd.Flags().Changed("e2NodeID") {
		e2NodeID, err = cmd.Flags().GetString("e2NodeID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("DuUeF1apID") {
		duUeF1apID, err = cmd.Flags().GetString("DuUeF1apID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("CuUeF1apID") {
		cuUeF1apID, err = cmd.Flags().GetString("CuUeF1apID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("ranUeNgapID") {
		ranUeNgapID, err = cmd.Flags().GetString("ranUeNgapID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("eNBUeS1apID") {
		enbUeS1apID, err = cmd.Flags().GetString("eNBUeS1apID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("dlSliceID") {
		dlSliceID, err = cmd.Flags().GetString("dlSliceID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("ulSliceID") {
		ulSliceID, err = cmd.Flags().GetString("ulSliceID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("drbID") {
		drbID, err = cmd.Flags().GetString("drbID")
		if err != nil {
			return err
		}
	}

	ueIDList := make([]*rsmapi.UeId, 0)
	duUeF1apIDField := &rsmapi.UeId{
		UeId: duUeF1apID,
		Type: rsmapi.UeIdType_UE_ID_TYPE_DU_UE_F1_AP_ID,
	}
	cuUeF1apIDField := &rsmapi.UeId{
		UeId: cuUeF1apID,
		Type: rsmapi.UeIdType_UE_ID_TYPE_CU_UE_F1_AP_ID,
	}
	ranUeNgapIDField := &rsmapi.UeId{
		UeId: ranUeNgapID,
		Type: rsmapi.UeIdType_UE_ID_TYPE_RAN_UE_NGAP_ID,
	}
	enbUeS1apIDField := &rsmapi.UeId{
		UeId: enbUeS1apID,
		Type: rsmapi.UeIdType_UE_ID_TYPE_ENB_UE_S1_AP_ID,
	}
	ueIDList = append(ueIDList, duUeF1apIDField)
	ueIDList = append(ueIDList, cuUeF1apIDField)
	ueIDList = append(ueIDList, ranUeNgapIDField)
	ueIDList = append(ueIDList, enbUeS1apIDField)

	setRequest := rsmapi.SetUeSliceAssociationRequest{
		E2NodeId:  e2NodeID,
		UeId:      ueIDList,
		DlSliceId: dlSliceID,
		UlSliceId: ulSliceID,
		DrbId:     drbID,
	}
	setResponse, err := client.SetUeSliceAssociation(context.Background(), &setRequest)
	if err != nil {
		return err
	}

	if !setResponse.Ack.Success {
		return fmt.Errorf("%v", setResponse.Ack.Cause)
	}
	return nil
}

func runGetCreateSlice(cmd *cobra.Command, _ []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := rsmapi.NewRsmClient(conn)

	e2NodeID := ""
	sliceID := ""
	schedulerType := ""
	weight := ""
	sliceType := ""

	if cmd.Flags().Changed("e2NodeID") {
		e2NodeID, err = cmd.Flags().GetString("e2NodeID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("sliceID") {
		sliceID, err = cmd.Flags().GetString("sliceID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("scheduler") {
		schedulerType, err = cmd.Flags().GetString("scheduler")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("weight") {
		weight, err = cmd.Flags().GetString("weight")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("sliceType") {
		sliceType, err = cmd.Flags().GetString("sliceType")
		if err != nil {
			return err
		}
	}

	var schedulerTypeField rsmapi.SchedulerType
	switch schedulerType {
	case "RR":
		schedulerTypeField = rsmapi.SchedulerType_SCHEDULER_TYPE_ROUND_ROBIN
	case "PF":
		schedulerTypeField = rsmapi.SchedulerType_SCHEDULER_TYPE_PROPORTIONALLY_FAIR
	case "QoS":
		schedulerTypeField = rsmapi.SchedulerType_SCHEDULER_TYPE_QOS_BASED
	default:
		return fmt.Errorf("scheduler should be {RR, PF, QoS}")
	}

	var sliceTypeField rsmapi.SliceType
	switch sliceType {
	case "DL":
		sliceTypeField = rsmapi.SliceType_SLICE_TYPE_DL_SLICE
	case "UL":
		sliceTypeField = rsmapi.SliceType_SLICE_TYPE_UL_SLICE
	default:
		return fmt.Errorf("sliceType should be {DL, UL}")
	}

	setRequest := rsmapi.CreateSliceRequest{
		E2NodeId:      e2NodeID,
		SliceId:       sliceID,
		Weight:        weight,
		SchedulerType: schedulerTypeField,
		SliceType:     sliceTypeField,
	}
	setResponse, err := client.CreateSlice(context.Background(), &setRequest)
	if err != nil {
		return err
	}

	if !setResponse.Ack.Success {
		return fmt.Errorf("%v", setResponse.Ack.Cause)
	}

	return nil
}

func runGetUpdateSlice(cmd *cobra.Command, _ []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := rsmapi.NewRsmClient(conn)

	e2NodeID := ""
	sliceID := ""
	schedulerType := ""
	weight := ""
	sliceType := ""

	if cmd.Flags().Changed("e2NodeID") {
		e2NodeID, err = cmd.Flags().GetString("e2NodeID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("sliceID") {
		sliceID, err = cmd.Flags().GetString("sliceID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("scheduler") {
		schedulerType, err = cmd.Flags().GetString("scheduler")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("weight") {
		weight, err = cmd.Flags().GetString("weight")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("sliceType") {
		sliceType, err = cmd.Flags().GetString("sliceType")
		if err != nil {
			return err
		}
	}

	var schedulerTypeField rsmapi.SchedulerType
	switch schedulerType {
	case "RR":
		schedulerTypeField = rsmapi.SchedulerType_SCHEDULER_TYPE_ROUND_ROBIN
	case "PF":
		schedulerTypeField = rsmapi.SchedulerType_SCHEDULER_TYPE_PROPORTIONALLY_FAIR
	case "QoS":
		schedulerTypeField = rsmapi.SchedulerType_SCHEDULER_TYPE_QOS_BASED
	default:
		return fmt.Errorf("scheduler should be {RR, PF, QoS}")
	}

	var sliceTypeField rsmapi.SliceType
	switch sliceType {
	case "DL":
		sliceTypeField = rsmapi.SliceType_SLICE_TYPE_DL_SLICE
	case "UL":
		sliceTypeField = rsmapi.SliceType_SLICE_TYPE_UL_SLICE
	default:
		return fmt.Errorf("sliceType should be {DL, UL}")
	}

	setRequest := rsmapi.UpdateSliceRequest{
		E2NodeId:      e2NodeID,
		SliceId:       sliceID,
		Weight:        weight,
		SchedulerType: schedulerTypeField,
		SliceType:     sliceTypeField,
	}
	setResponse, err := client.UpdateSlice(context.Background(), &setRequest)
	if err != nil {
		return err
	}

	if !setResponse.Ack.Success {
		return fmt.Errorf("%v", setResponse.Ack.Cause)
	}

	return nil
}

func runGetDeleteSlice(cmd *cobra.Command, _ []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := rsmapi.NewRsmClient(conn)

	e2NodeID := ""
	sliceID := ""
	sliceType := ""

	if cmd.Flags().Changed("e2NodeID") {
		e2NodeID, err = cmd.Flags().GetString("e2NodeID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("sliceID") {
		sliceID, err = cmd.Flags().GetString("sliceID")
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("sliceType") {
		sliceType, err = cmd.Flags().GetString("sliceType")
		if err != nil {
			return err
		}
	}

	var sliceTypeField rsmapi.SliceType
	switch sliceType {
	case "DL":
		sliceTypeField = rsmapi.SliceType_SLICE_TYPE_DL_SLICE
	case "UL":
		sliceTypeField = rsmapi.SliceType_SLICE_TYPE_UL_SLICE
	default:
		return fmt.Errorf("sliceType should be {DL, UL}")
	}

	setRequest := rsmapi.DeleteSliceRequest{
		E2NodeId:  e2NodeID,
		SliceId:   sliceID,
		SliceType: sliceTypeField,
	}
	setResponse, err := client.DeleteSlice(context.Background(), &setRequest)
	if err != nil {
		return err
	}

	if !setResponse.Ack.Success {
		return fmt.Errorf("%v", setResponse.Ack.Cause)
	}

	return nil
}
