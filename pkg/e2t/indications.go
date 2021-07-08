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

package e2t

// TODO: Remove: deprecated
/*
func getWatchIndicationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indications",
		Short: "watch indications traffic",
		RunE:  runWatchIndicationsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().String("service-model", "oran-e2sm-rc-pre:v2", "service model as name:version")
	cmd.Flags().String("node", "", "E2 node ID")
	return cmd
}

func runWatchIndicationsCommand(cmd *cobra.Command, args []string) error {
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	sm, _ := cmd.Flags().GetString("service-model")
	smf := strings.Split(sm, ":")
	if len(smf) < 2 {
		return errors.New("invalid model format")
	}
	serviceModelName := e2client.ServiceModelName(smf[0])
	serviceModelVersion := e2client.ServiceModelVersion(smf[1])

	nodeID, _ := cmd.Flags().GetString("node")

	e2tHost := "onos-e2t"
	e2tPort := 5150

	client := e2client.NewClient(
		e2client.WithServiceModel(serviceModelName, serviceModelVersion),
		e2client.WithAppID("onos-cli"),
		e2client.WithE2TAddress(e2tHost, e2tPort))

	actions := []e2api.Action{
		{
			ID:   100,
			Type: e2api.ActionType_ACTION_TYPE_REPORT,
			SubsequentAction: &e2api.SubsequentAction{
				Type:       e2api.SubsequentActionType_SUBSEQUENT_ACTION_TYPE_CONTINUE,
				TimeToWait: e2api.TimeToWait_TIME_TO_WAIT_ZERO,
			},
		},
	}

	trigger := e2api.EventTrigger{
		Payload: []byte{},
	}

	ch := make(chan e2api.Indication)
	node := client.Node(e2client.NodeID(nodeID))
	subName := "onos-pci-subscription"
	subSpec := e2api.SubscriptionSpec{
		Actions:      actions,
		EventTrigger: trigger,
	}
	ctx := context.Background()
	channelID, err := node.Subscribe(ctx, subName, subSpec, ch)
	if err != nil {
		return err
	}

	// TODO: to be implemented
	fmt.Printf("%v", channelID)

	return nil
}
*/
