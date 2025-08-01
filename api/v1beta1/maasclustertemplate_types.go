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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MaasClusterTemplateSpec defines the desired state of MaasClusterTemplate
// +kubebuilder:object:generate=true
// +kubebuilder:resource:path=maasclustertemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
type MaasClusterTemplateSpec struct {
	Template MaasClusterTemplateResource `json:"template"`
}

// MaasClusterTemplateResource describes the data needed to create a MaasCluster from a template
type MaasClusterTemplateResource struct {
	// Spec is the specification of the desired behavior of the cluster.
	Spec MaasClusterSpec `json:"spec"`
}

// MaasClusterTemplate is the Schema for the maasclustertemplates API
// +kubebuilder:object:root=true
type MaasClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaasClusterTemplateSpec `json:"spec,omitempty"`
	Status MaasClusterStatus       `json:"status,omitempty"`
}

// MaasClusterTemplateList contains a list of MaasClusterTemplate
// +kubebuilder:object:root=true
type MaasClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaasClusterTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MaasClusterTemplate{}, &MaasClusterTemplateList{})
}
