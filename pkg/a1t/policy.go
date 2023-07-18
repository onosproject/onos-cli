// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package a1t

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	a1 "github.com/onosproject/onos-api/go/onos/a1t/admin"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"io"
	"text/tabwriter"
)

const (
	policyTypeFormat   = "%-50s\t%s\n"
	policyObjectFormat = "%-50s\t%-30s\n"
	policyStatusFormat = "%-50s\t%-30s\t%s\n"
	verboseFormat      = "%s:\t%s\n"
	jsonFormat         = "%s:\n%s\n"
)

func addIndentJSONString(jsonStringObj string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(jsonStringObj), "", "\t")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return out.String()
}

func displayPolicyTypeListHeader(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, policyTypeFormat,
		"PolicyTypeID", "List(PolicyObjectID)")
}

func displayPolicyObjectListHeader(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, policyObjectFormat,
		"PolicyTypeID", "PolicyObjectID")
}

func displayPolicyStatusListHeader(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, policyStatusFormat,
		"PolicyTypeID", "PolicyObjectID", "Status")
}

func displayPolicyTypeListElement(writer io.Writer, resp *a1.GetPolicyTypeObjectResponse) {
	_, _ = fmt.Fprintf(writer, policyTypeFormat,
		resp.PolicyTypeId, resp.PolicyIds)
}

func displayPolicyType(writer io.Writer, resp *a1.GetPolicyTypeObjectResponse) {
	_, _ = fmt.Fprintf(writer, verboseFormat, "PolicyTypeID", resp.PolicyTypeId)
	_, _ = fmt.Fprintf(writer, verboseFormat, "PolicyObjectIDs", resp.PolicyIds)
	_, _ = fmt.Fprintf(writer, jsonFormat, "PolicyTypeObject", addIndentJSONString(resp.PolicyTypeObject))
}

func displayPolicyObjectListElement(writer io.Writer, resp *a1.GetPolicyObjectResponse) {
	_, _ = fmt.Fprintf(writer, policyObjectFormat,
		resp.PolicyTypeId, resp.PolicyObjectId)
}

func displayPolicyObject(writer io.Writer, resp *a1.GetPolicyObjectResponse) {
	_, _ = fmt.Fprintf(writer, verboseFormat, "PolicyTypeID", resp.PolicyTypeId)
	_, _ = fmt.Fprintf(writer, verboseFormat, "PolicyObjectIDs", resp.PolicyObjectId)
	_, _ = fmt.Fprintf(writer, jsonFormat, "PolicyTypeObject", addIndentJSONString(resp.PolicyObject))
}

func displayPolicyStatusListElement(writer io.Writer, resp *a1.GetPolicyObjectStatusResponse) {
	_, _ = fmt.Fprintf(writer, policyStatusFormat,
		resp.PolicyTypeId, resp.PolicyObjectId, resp.PolicyObjectStatus)
}

func displayPolicyStatus(writer io.Writer, resp *a1.GetPolicyObjectStatusResponse) {
	_, _ = fmt.Fprintf(writer, verboseFormat, "PolicyTypeID", resp.PolicyTypeId)
	_, _ = fmt.Fprintf(writer, verboseFormat, "PolicyObjectIDs", resp.PolicyObjectId)
	_, _ = fmt.Fprintf(writer, jsonFormat, "PolicyTypeObject", addIndentJSONString(resp.PolicyObjectStatus))
}

func getPolicyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy {type/object/status} [args]",
		Short: "Policy command",
	}

	cmd.AddCommand(getGetPolicyType())
	cmd.AddCommand(getGetPolicyObject())
	cmd.AddCommand(getGetPolicyStatus())
	return cmd
}

func getGetPolicyType() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "type [args]",
		Short: "Get policy type",
		RunE:  runGetPolicyType,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().String("policyTypeID", "", "Policy Type ID")
	return cmd
}

func getGetPolicyObject() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "object [args]",
		Short: "Get policy object",
		RunE:  runGetPolicyObject,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().String("policyTypeID", "", "Policy Type ID")
	cmd.Flags().String("policyObjectID", "", "Policy Object ID")
	return cmd
}

func getGetPolicyStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [args]",
		Short: "Get policy status",
		RunE:  runGetPolicyStatus,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().String("policyTypeID", "", "Policy Type ID")
	cmd.Flags().String("policyObjectID", "", "Policy Object ID")
	return cmd
}

func runGetPolicyType(cmd *cobra.Command, _ []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := a1.NewA1TAdminServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutTimer)
	defer cancel()

	policyTypeID := ""

	if cmd.Flags().Changed("policyTypeID") {
		policyTypeID, err = cmd.Flags().GetString("policyTypeID")
		if err != nil {
			return err
		}
	}

	if policyTypeID != "" {
		noHeaders = true
	}

	if !noHeaders {
		displayPolicyTypeListHeader(writer)
		_ = writer.Flush()
	}

	req := &a1.GetPolicyTypeObjectRequest{
		PolicyTypeId: policyTypeID,
	}

	stream, err := client.GetPolicyTypeObject(ctx, req)
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			cli.Output("Error receiving notification: %v", err)
			return err
		}
		if policyTypeID == "" {
			displayPolicyTypeListElement(writer, resp)
			_ = writer.Flush()
		} else {
			displayPolicyType(writer, resp)
			_ = writer.Flush()
		}
	}

	return nil
}

func runGetPolicyObject(cmd *cobra.Command, _ []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := a1.NewA1TAdminServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutTimer)
	defer cancel()

	policyTypeID := ""
	if cmd.Flags().Changed("policyTypeID") {
		policyTypeID, err = cmd.Flags().GetString("policyTypeID")
		if err != nil {
			return err
		}
	}

	policyObjectID := ""
	if cmd.Flags().Changed("policyObjectID") {
		policyObjectID, err = cmd.Flags().GetString("policyObjectID")
		if err != nil {
			return err
		}
	}

	if (policyTypeID != "" && policyObjectID == "") || (policyTypeID == "" && policyObjectID != "") {
		cli.Output("To show all policyObjects, policyObjectID and policyTypeID should be blank\n")
		cli.Output("To show all specific policy object, both policyObjectID and policyTypeID should not be blank\n")
		_ = writer.Flush()
	}

	if policyTypeID != "" && policyObjectID != "" {
		noHeaders = true
	}

	if !noHeaders {
		displayPolicyObjectListHeader(writer)
		_ = writer.Flush()
	}

	req := &a1.GetPolicyObjectRequest{
		PolicyTypeId:   policyTypeID,
		PolicyObjectId: policyObjectID,
	}

	stream, err := client.GetPolicyObject(ctx, req)
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			cli.Output("Error receiving notification: %v", err)
			return err
		}
		if policyTypeID == "" {
			displayPolicyObjectListElement(writer, resp)
			_ = writer.Flush()
		} else {
			displayPolicyObject(writer, resp)
			_ = writer.Flush()
		}
	}

	return nil
}

func runGetPolicyStatus(cmd *cobra.Command, _ []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	client := a1.NewA1TAdminServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutTimer)
	defer cancel()

	policyTypeID := ""
	if cmd.Flags().Changed("policyTypeID") {
		policyTypeID, err = cmd.Flags().GetString("policyTypeID")
		if err != nil {
			return err
		}
	}

	policyObjectID := ""
	if cmd.Flags().Changed("policyObjectID") {
		policyObjectID, err = cmd.Flags().GetString("policyObjectID")
		if err != nil {
			return err
		}
	}

	if (policyTypeID != "" && policyObjectID == "") || (policyTypeID == "" && policyObjectID != "") {
		cli.Output("To show all policyObjects, policyObjectID and policyTypeID should be blank\n")
		cli.Output("To show all specific policy object, both policyObjectID and policyTypeID should not be blank\n")
		_ = writer.Flush()
	}

	if policyTypeID != "" && policyObjectID != "" {
		noHeaders = true
	}

	if !noHeaders {
		displayPolicyStatusListHeader(writer)
		_ = writer.Flush()
	}

	req := &a1.GetPolicyObjectStatusRequest{
		PolicyTypeId:   policyTypeID,
		PolicyObjectId: policyObjectID,
	}

	stream, err := client.GetPolicyObjectStatus(ctx, req)
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			cli.Output("Error receiving notification: %v", err)
			return err
		}
		if policyTypeID == "" {
			displayPolicyStatusListElement(writer, resp)
			_ = writer.Flush()
		} else {
			displayPolicyStatus(writer, resp)
			_ = writer.Flush()
		}
	}

	return nil
}
