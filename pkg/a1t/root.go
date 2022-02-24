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

package a1t

import (
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
	"time"
)

const (
	configName = "a1t"
	defaultAddress = "onos-a1t:5150"
	TimeoutTimer = time.Second * 5
)

func init() {
	cli.InitConfig(configName)
}

func Init() {
	// noop for now
}

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "a1t {get} [args]",
		Short: "ONOS a1t subsystem commands",
	}

	cli.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(getGetCommand())
	return cmd
}

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command {
		Use: "get {subscriptions/subscription/policy} [args]",
		Short: "Get command",
		Aliases: []string{"list"},
	}

	cmd.AddCommand(getGetSubscriptionCommand())
	cmd.AddCommand(getPolicyCommand())

	return cmd
}
