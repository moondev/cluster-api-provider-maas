/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta2

import (
	"k8s.io/apimachinery/pkg/util/sets"
	clusterv1beta2 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

// MachineState describes the state of an MAAS Machine.
type MachineState string

var (
	MachineStateAllocated  = MachineState("Allocated")
	MachineStateDeploying = MachineState("Deploying")
	MachineStateDeployed  = MachineState("Deployed")
	MachineStateReady      = MachineState("Ready")
	MachineStateDiskErasing = MachineState("Disk erasing")
	MachineStateReleasing  = MachineState("Releasing")
	MachineStateNew        = MachineState("New")

	MachineRunningStates = sets.NewString(
		string(MachineStateDeploying),
		string(MachineStateDeployed),
	)

	MachineOperationalStates = MachineRunningStates.Union(
		sets.NewString(string(MachineStateAllocated)),
	)

	MachineKnownStates = MachineOperationalStates.Union(
		sets.NewString(
			string(MachineStateDiskErasing),
			string(MachineStateReleasing),
			string(MachineStateReady),
			string(MachineStateNew),
		),
	)
)

// Machine describes an MAAS Machine (internal representation).
type Machine struct {
	ID                string
	Hostname          string
	State             MachineState
	Powered           bool
	AvailabilityZone  string
	Addresses         []clusterv1beta2.MachineAddress
}
