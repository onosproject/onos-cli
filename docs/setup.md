<!--
SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>

SPDX-License-Identifier: Apache-2.0
-->

# How to install and run onos-cli?

To install onos command line client, run the following command:

```bash
> export GO111MODULE=on
> go get github.com/onosproject/onos-cli/cmd/onos
```

# ONOS-CLI Auto Completion
The onos client supports shell auto-completion for its various commands, sub-commands and flags. 
You can enable this feature for bash or zsh as follows:

## Bash Auto-Completion
```bash
> eval "$(onos completion bash)"
```
After that, you should be able to use the `TAB` key to obtain suggestions for 
valid options.

## Zsh Auto-Completion

```bash
> source <(onos completion zsh)
```

**Note**: Note: We also recommend to add the output of the above commands to *.bashrc* or *.zshrc*.

## How to run onos client?

After the above steps, if you run *onos* from command line interface,
you should be able to see an output like this:

```bash
> onos
ONOS command line client

Usage:
  onos [command]

Available Commands:
  changes     Lists records of configuration changes
  completion  Generated bash or zsh auto-completion script
  config      Read and update CLI configuration options
  configs     Lists details of device configuration changes
  devices     Manages inventory of network devices
  devicetree  Lists devices and their configuration in tree format
  help        Help about any command
  init        Initialize the ONOS CLI configuration
  models      Manages model plugins
  net-changes Lists network configuration changes
  rollback    Rolls-back a network configuration change

Flags:
  -a, --address string    the controller address (default ":5150")
  -c, --certPath string   path to client certificate (default "client1.crt")
      --config string     config file (default: $HOME/.onos/config.yaml)
  -h, --help              help for onos
  -k, --keyPath string    path to client private key (default "client1.key")

Use "onos [command] --help" for more information about a command.
```
