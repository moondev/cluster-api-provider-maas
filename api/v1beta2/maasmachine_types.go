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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1beta2 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/errors"
)

const (
	MachineFinalizer = "maasmachine.infrastructure.cluster.x-k8s.io"
)

// MaasMachineSpec defines the desired state of MaasMachine
type MaasMachineSpec struct {
	FailureDomain *string `json:"failureDomain,omitempty"`
	SystemID      *string `json:"systemID,omitempty"`
	ProviderID    *string `json:"providerID,omitempty"`
	ResourcePool  *string `json:"resourcePool,omitempty"`
	MinCPU        *int    `json:"minCPU,omitempty"`
	MinMemoryInMB *int    `json:"minMemory,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Image     string  `json:"image"`
	// Ephemeral enables in-memory (ephemeral) deploy for this machine when true.
	// See https://github.com/spectrocloud/maas-client-go/releases (v0.1.2-beta1+).
	Ephemeral bool `json:"ephemeral,omitempty"`
}

// MaasMachineStatus defines the observed state of MaasMachine
type MaasMachineStatus struct {
	Ready          bool                           `json:"ready"`
	MachineState   *MachineState                   `json:"machineState,omitempty"`
	MachinePowered bool                           `json:"machinePowered,omitempty"`
	Hostname       *string                         `json:"hostname,omitempty"`
	DNSAttached    bool                            `json:"dnsAttached,omitempty"`
	Addresses      []clusterv1beta2.MachineAddress `json:"addresses,omitempty"`
	Conditions     clusterv1beta2.Conditions       `json:"conditions,omitempty"`
	FailureReason  *errors.MachineStatusError      `json:"failureReason,omitempty"`
	FailureMessage *string                         `json:"failureMessage,omitempty"`
}

// +kubebuilder:resource:path=maasmachines,scope=Namespaced,categories=cluster-api
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// MaasMachine is the Schema for the maasmachines API
type MaasMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaasMachineSpec   `json:"spec,omitempty"`
	Status MaasMachineStatus `json:"status,omitempty"`
}

func (c *MaasMachine) GetConditions() clusterv1beta2.Conditions {
	return c.Status.Conditions
}

func (c *MaasMachine) SetConditions(conditions clusterv1beta2.Conditions) {
	c.Status.Conditions = conditions
}

// GetV1Beta1Conditions returns the list of conditions for the deprecated v1beta1 conditions API.
func (c *MaasMachine) GetV1Beta1Conditions() clusterv1beta2.Conditions {
	return c.Status.Conditions
}

// SetV1Beta1Conditions sets the list of conditions for the deprecated v1beta1 conditions API.
func (c *MaasMachine) SetV1Beta1Conditions(conditions clusterv1beta2.Conditions) {
	c.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// MaasMachineList contains a list of MaasMachine
type MaasMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaasMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MaasMachine{}, &MaasMachineList{})
}
