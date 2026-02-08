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
)

const (
	ClusterFinalizer = "maascluster.infrastructure.cluster.x-k8s.io"
)

// MaasClusterSpec defines the desired state of MaasCluster
type MaasClusterSpec struct {
	DNSDomain string `json:"dnsDomain"`

	// +optional
	ControlPlaneEndpoint APIEndpoint `json:"controlPlaneEndpoint,omitempty"`

	// +optional
	FailureDomains []string `json:"failureDomains,omitempty"`
}

// MaasClusterStatus defines the observed state of MaasCluster
type MaasClusterStatus struct {
	Ready bool `json:"ready"`

	Network         Network                         `json:"network,omitempty"`
	FailureDomains  []clusterv1beta2.FailureDomain  `json:"failureDomains,omitempty"`
	Conditions      clusterv1beta2.Conditions       `json:"conditions,omitempty"`
}

// Network encapsulates the Cluster Network
type Network struct {
	DNSName string `json:"dnsName,omitempty"`
}

// APIEndpoint represents a reachable Kubernetes API endpoint (Host and Port optional per CAPI v1beta2).
type APIEndpoint struct {
	// +optional
	Host string `json:"host,omitempty"`
	// +optional
	Port int `json:"port,omitempty"`
}

// IsZero returns true if both host and port are zero values.
func (in APIEndpoint) IsZero() bool {
	return in.Host == "" && in.Port == 0
}

// +kubebuilder:resource:path=maasclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// MaasCluster is the Schema for the maasclusters API
type MaasCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaasClusterSpec   `json:"spec,omitempty"`
	Status MaasClusterStatus `json:"status,omitempty"`
}

func (in *MaasCluster) GetConditions() clusterv1beta2.Conditions {
	return in.Status.Conditions
}

func (in *MaasCluster) SetConditions(conditions clusterv1beta2.Conditions) {
	in.Status.Conditions = conditions
}

// GetV1Beta1Conditions returns the list of conditions for the deprecated v1beta1 conditions API.
func (in *MaasCluster) GetV1Beta1Conditions() clusterv1beta2.Conditions {
	return in.Status.Conditions
}

// SetV1Beta1Conditions sets the list of conditions for the deprecated v1beta1 conditions API.
func (in *MaasCluster) SetV1Beta1Conditions(conditions clusterv1beta2.Conditions) {
	in.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// MaasClusterList contains a list of MaasCluster
type MaasClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaasCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MaasCluster{}, &MaasClusterList{})
}
