// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"github.com/onosproject/onos-cli/pkg/a1t"
	"github.com/onosproject/onos-cli/pkg/discovery"
	"github.com/onosproject/onos-cli/pkg/fabricsim"
	"github.com/onosproject/onos-cli/pkg/mlb"
	"github.com/onosproject/onos-cli/pkg/perf"
	"github.com/onosproject/onos-cli/pkg/provisioner"
	"github.com/onosproject/onos-cli/pkg/rsm"

	"github.com/onosproject/onos-cli/pkg/config"
	"github.com/onosproject/onos-cli/pkg/e2t"
	"github.com/onosproject/onos-cli/pkg/kpimon"
	"github.com/onosproject/onos-cli/pkg/mho"
	"github.com/onosproject/onos-cli/pkg/pci"
	"github.com/onosproject/onos-cli/pkg/topo"
	"github.com/onosproject/onos-cli/pkg/uenib"

	// Needed to keep ran-sim happy for the mo
	_ "github.com/onosproject/onos-lib-go/pkg/cli"

	"os"

	"github.com/onosproject/onos-cli/pkg/ransim"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Execute runs the root command and any sub-commands.
func Execute() {
	rootCmd := GetRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const apacheLicense = "<!--\nSP" + "DX-FileCopyright" +
	"Text: 2019-present Open Networking Foundation <info@opennetworking.org>\n\nSP" +
	"DX-License-Identifier: Apache-2.0\n-->\n\n"

// GenerateCliDocs generate markdown files for onos-cli commands
func GenerateCliDocs() {
	cmd := GetRootCommand()
	identity := func(s string) string { return s }
	licenseStr := func(s string) string { return apacheLicense }
	err := doc.GenMarkdownTreeCustom(cmd, "docs/cli", licenseStr, identity)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// GetRootCommand returns the root onos command
func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "onos",
		Short:                  "ONOS command line client",
		BashCompletionFunction: getBashCompletions(),
		SilenceUsage:           true,
		SilenceErrors:          true,
	}
	cmd.AddCommand(topo.GetCommand())
	cmd.AddCommand(uenib.GetCommand())
	cmd.AddCommand(config.GetCommand())
	cmd.AddCommand(fabricsim.GetCommand())
	cmd.AddCommand(e2t.GetCommand())
	cmd.AddCommand(ransim.GetCommand())
	cmd.AddCommand(kpimon.GetCommand())
	cmd.AddCommand(mho.GetCommand())
	cmd.AddCommand(pci.GetCommand())
	cmd.AddCommand(mlb.GetCommand())
	cmd.AddCommand(rsm.GetCommand())
	cmd.AddCommand(perf.GetCommand())
	cmd.AddCommand(a1t.GetCommand())
	cmd.AddCommand(provisioner.GetCommand())
	cmd.AddCommand(discovery.GetCommand())
	cmd.AddCommand()

	cmd.AddCommand(getCompletionCommand())

	return cmd
}
