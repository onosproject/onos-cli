# Deploying onos-cli with Helm

This guide deploys `onos-cli` through it's [Helm] chart assumes you have a [Kubernets] cluster running 
with an atomix controller deployed in a namespace. if you dont' specify the `--namespace` in the commands 
below atomix controller must be deployed in the `default`
`onos-cli` Helm chart is based on Helm 3.0 version, with no need for the Tiller pod to be present. 
If you don't have a cluster running and want to try on your local machine please follow first 
the [Kubernetes] setup steps outlined in [deploy with Helm](https://docs.onosproject.org/developers/deploy_with_helm/).

## Installing the Chart

To install the chart, simply run `helm install deployments/helm/onos-cli` from
the root directory of this project:

```bash
helm install deployments/helm/onos-cli
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

To install the chart in a different namespace please modify the `default` occurances in the `values.yaml` file. 
Please be aware to change also `atomix-controller.default.svc.cluster.local:5679` 
to `atomix-controller.<your_name_space_here>.svc.cluster.local:5679`.
Then issue the `helm install` command
```bash
helm install onos-cli --namespace <your_name_space> deployments/helm/onos-cli
```
### Installing the chart with debug. 
`onos-cli` offers the capability to open a debug port (4000) to the image.
To enable the debug capabilities please set the debug flag to true in `values.yaml` or pass it to `helm install`
```bash
helm install onos-cli deployments/helm/onos-cli --set debug=true
```
### Troubleshoot

If your chart does not install or the pod is not running for some reason and/or you modified values Helm offers two flags to help you
debug your chart:  
- `--dry-run` check the chart without actually installing the pod. 
- `--debug` prints out more information about your chart
```bash
helm install onos-cli --debug --dry-run ./deployments/helm/onos-cli/
```
## Uninstalling the chart.

To remove the `onos-cli` pod issue
```bash
 helm delete onos-cli
```
## Pod Information

To view the pods that are deployed, run `kubectl get pods`.

## Getting access to the onos-cli console

To gain acess to the `onos-cli` console and be able of issuing the different cli commands the following command is need:
```bash
> kubectl -n default exec <cli_pod_name> -it -- /bin/sh
```

At this point you can execute `topo`, `config` and all the other commands. For example:
 ```bash
 > onos topo get devices
 ```

[Helm]: https://helm.sh/
[Kubernetes]: https://kubernetes.io/
[kind]: https://kind.sigs.k8s.io

