// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package kpimon

import "github.com/spf13/cobra"

func getListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list {metrics} [args]",
		Short: "List KPIMON resources",
	}
	cmd.AddCommand(getListMetricsCommand())
	return cmd
}

func getWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch {metrics} [args]",
		Short: "Watch KPIMON resources",
	}
	cmd.AddCommand(getWatchMetricsCommand())
	return cmd
}

func getSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {report_interval} [args]",
		Short: "Set KPIMON parameters",
	}
	cmd.AddCommand(setReportIntervalCommand())
	return cmd
}
