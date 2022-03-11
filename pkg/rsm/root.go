// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package rsm

import (
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	defaultAddress = "onos-rsm:5150"
)

// Init is a hook called after cobra initialization
func Init() {
	// noop for now
}

// GetCommand returns the root command for the RAN service
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rsm {set/create/update/delete} [args]",
		Short: "ONOS RSM subsystem commands",
	}

	cli.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(getSetCommand())
	cmd.AddCommand(getCreateCommand())
	cmd.AddCommand(getUpdateCommand())
	cmd.AddCommand(getDeleteCommand())
	return cmd
}
