// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package a1t

import (
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"time"
)

const (
	configName     = "a1t"
	defaultAddress = "onos-a1t:5150"

	// TimeoutTimer is a timer time
	TimeoutTimer = time.Second * 5
)

func init() {
	cli.InitConfig(configName)
}

// Init is a hook called after cobra initialization
func Init() {
	// noop for now
}

// GetCommand returns the root command for the RAN service
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "a1t {get} [args]",
		Short: "ONOS a1t subsystem commands",
	}

	cli.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(getGetCommand())
	return cmd
}

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get {subscriptions/subscription/policy} [args]",
		Short:   "Get command",
		Aliases: []string{"list"},
	}

	cmd.AddCommand(getGetSubscriptionCommand())
	cmd.AddCommand(getPolicyCommand())

	return cmd
}
