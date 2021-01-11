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

package topo

const bashCompletion = `
__onos_topo_get_devices() {
    local onos_output out
    if onos_output=$(onos topo get devices --no-headers 2>/dev/null); then
        out=($(echo "${onos_output}" | awk '{print $1}'))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__onit_custom_func() {
    case ${last_command} in
        onos_topo_get_device | onos_topo_update_device | onos_topo_remove_device)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __onit_get_clusters
            fi
            return
            ;;
        *)
            ;;
    esac
}
`

// GetBashCompletion returns the bash completion script for topo
func GetBashCompletion() string {
	return bashCompletion
}
