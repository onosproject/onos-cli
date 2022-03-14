// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

// Package config holds ONOS command-line command implementations for onos-config
package config

import (
	clilib "github.com/onosproject/onos-lib-go/pkg/cli"
	loglib "github.com/onosproject/onos-lib-go/pkg/logging/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "config"
	defaultAddress = "onos-config:5150"
)

// init initializes the command line
func init() {
	clilib.InitConfig(configName)
}

// Init is a hook called after cobra initialization
func Init() {
	// noop for now
}

// GetCommand returns the root command for the config service.
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config {get,watch,rollback} [args]",
		Short: "ONOS configuration subsystem commands",
	}

	clilib.AddConfigFlags(cmd, defaultAddress)

	cmd.AddCommand(clilib.GetConfigCommand())
	cmd.AddCommand(getGetCommand())
	cmd.AddCommand(getRollbackCommand())
	cmd.AddCommand(getWatchCommand())
	cmd.AddCommand(loglib.GetCommand())
	return cmd
}
