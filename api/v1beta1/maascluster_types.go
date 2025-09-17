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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

const (
	// ClusterFinalizer allows MaasClusterReconciler to clean up resources associated with MaasCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "maascluster.infrastructure.cluster.x-k8s.io"
)

// MaasClusterSpec defines the desired state of MaasCluster
type MaasClusterSpec struct {
	// DNSDomain configures the MaaS domain to create the cluster on (e.g maas)
	// +kubebuilder:validation:MinLength=1
	DNSDomain string `json:"dnsDomain"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint APIEndpoint `json:"controlPlaneEndpoint"`

	// FailureDomains are not usually defined on the spec.
	// but useful for MaaS since we can limit the domains to these
	// +optional
	FailureDomains []string `json:"failureDomains,omitempty"`
}

// MaasClusterStatus defines the observed state of MaasCluster
type MaasClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Ready denotes that the maas cluster (infrastructure) is ready.
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// Network represents the network
	Network Network `json:"network,omitempty"`

	// FailureDomains don't mean much in CAPMAAS since it's all local, but we can see how the rest of cluster API
	// will use this if we populate it.
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`

	// Conditions defines current service state of the MaasCluster.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// Network encapsulates the Cluster Network
type Network struct {
	// DNSName is the Kubernetes api server name
	DNSName string `json:"dnsName,omitempty"`
}

// APIEndpoint represents a reachable Kubernetes API endpoint.
type APIEndpoint struct {

	// Host is the hostname on which the API server is serving.
	Host string `json:"host"`

	// Port is the port on which the API server is serving.
	Port int `json:"port"`
}

// IsZero returns true if both host and port are zero values.
func (in APIEndpoint) IsZero() bool {
	return in.Host == "" && in.Port == 0
}

// +kubebuilder:resource:path=maasclusters,scope=Namespaced,categories=cluster-api
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// MaasCluster is the Schema for the maasclusters API
type MaasCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaasClusterSpec   `json:"spec,omitempty"`
	Status MaasClusterStatus `json:"status,omitempty"`
}

func (in *MaasCluster) GetConditions() []metav1.Condition {
	out := make([]metav1.Condition, len(in.Status.Conditions))
	for i := range in.Status.Conditions {
		c := in.Status.Conditions[i]
		out[i] = metav1.Condition{
			Type:               string(c.Type),
			Status:             metav1.ConditionStatus(c.Status),
			LastTransitionTime: c.LastTransitionTime,
			Reason:             c.Reason,
			Message:            c.Message,
		}
	}
	return out
}

func (in *MaasCluster) SetConditions(conditions []metav1.Condition) {
	out := make(clusterv1.Conditions, len(conditions))
	for i := range conditions {
		c := conditions[i]
		out[i] = clusterv1.Condition{
			Type:               clusterv1.ConditionType(c.Type),
			Status:             corev1.ConditionStatus(c.Status),
			LastTransitionTime: c.LastTransitionTime,
			Reason:             c.Reason,
			Message:            c.Message,
		}
	}
	in.Status.Conditions = out
}

//+kubebuilder:object:root=true

// MaasClusterList contains a list of MaasCluster
type MaasClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaasCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MaasCluster{}, &MaasClusterList{})
}
