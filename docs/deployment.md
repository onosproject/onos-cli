<!--
SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>

SPDX-License-Identifier: Apache-2.0
-->

# Deploying onos-cli with Helm

This guide deploys `onos-cli` through it's [Helm] chart assumes you have a [Kubernetes] cluster running 
with an atomix controller deployed in a namespace.  
`onos-cli` Helm chart is based on Helm 3.0 version, with no need for the Tiller pod to be present. 
If you don't have a cluster running and want to try on your local machine please follow first 
the [Kubernetes] setup steps outlined in [deploy with Helm](https://docs.onosproject.org/developers/deploy_with_helm/).
The following steps assume you have the setup outlined in that page, including the `micro-onos` namespace configured. 
## Installing the Chart

To install the chart in the `micro-onos` namespace run from the root directory of the `onos-helm-charts` repo the command:
```bash
helm install -n micro-onos onos-cli onos-cli
```
The output should be:
```bash
NAME: onos-cli
LAST DEPLOYED: Tue Nov 26 13:31:42 2019
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

`helm install` assigns a unique name to the chart and displays all the k8s resources that were
created by it. To list the charts that are installed and view their statuses, run `helm ls`:

```bash
helm ls
NAME          	REVISION	UPDATED                 	STATUS  	CHART                    	APP VERSION	NAMESPACE
...
onos-cli	1       	Tue May 14 18:56:39 2019	DEPLOYED	onos-cli-0.0.1	        0.0.1      	default
```

### Installing the chart in a different namespace.

Issue the `helm install` command substituting `micro-onos` with your namespace.
```bash
helm install -n <your_name_space> onos-cli onos-cli
```
### Installing the chart with debug. 
`onos-cli` offers the capability to open a debug port (4000) to the image.
To enable the debug capabilities please set the debug flag to true in `values.yaml` or pass it to `helm install`
```bash
helm install -n micro-onos onos-cli onos-cli --set debug=true
```
### Troubleshoot

If your chart does not install or the pod is not running for some reason and/or you modified values Helm offers two flags to help you
debug your chart:

* `--dry-run` check the chart without actually installing the pod. 
* `--debug` prints out more information about your chart

```bash
helm install -n micro-onos onos-cli --debug --dry-run onos-cli/
```
Also to verify how template values are expanded, run:
```bash
helm install template onos-gui
```

## Uninstalling the chart.

To remove the `onos-cli` pod issue
```bash
helm delete -n micro-onos onos-cli
```
## Pod Information

To view the pods that are deployed, run `kubectl -n micro-onos get pods`.

## Getting access to the onos-cli console

To gain acess to the `onos-cli` console and be able of issuing the different cli commands the following command is need:
```bash
kubectl -n micro-onos exec -it $(kubectl -n micro-onos get pods -l type=cli -o name) -- /bin/sh
```

At this point you can execute `topo`, `config` and all the other commands. For example:
```bash
 onos topo get devices
```

[Helm]: https://helm.sh/
[Kubernetes]: https://kubernetes.io/

