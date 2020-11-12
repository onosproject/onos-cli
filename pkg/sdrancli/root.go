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

package sdrancli

import (
	"fmt"
	"github.com/onosproject/onos-cli/pkg/cli"
	// Needed to keep ran-sim happy for the mo
	_ "github.com/onosproject/onos-lib-go/pkg/cli"

	e2sub "github.com/onosproject/onos-e2sub/pkg/cli"
	e2t "github.com/onosproject/onos-e2t/pkg/cli"
	richo "github.com/onosproject/onos-ric/pkg/apps/onos-ric-ho/cli"
	ricmlb "github.com/onosproject/onos-ric/pkg/apps/onos-ric-mlb/cli"
	ric "github.com/onosproject/onos-ric/pkg/cli"
	ransim "github.com/onosproject/ran-simulator/pkg/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

// Execute runs the root command and any sub-commands.
func Execute() {
	rootCmd := GetRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// GenerateCliDocs generate markdown files for onos-cli commands
func GenerateCliDocs() {
	cmd := GetRootCommand()
	err := doc.GenMarkdownTree(cmd, "docs/cli")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// GetRootCommand returns the root onos command
func GetRootCommand() *cobra.Command {
	cmd := cli.GetRootCommand()
	// Do something with completions
	cmd.AddCommand(ric.GetCommand())
	cmd.AddCommand(richo.GetCommand())
	cmd.AddCommand(ricmlb.GetCommand())
	cmd.AddCommand(ransim.GetCommand())
	cmd.AddCommand(e2t.GetCommand())
	cmd.AddCommand(e2sub.GetCommand())
	return cmd
}
