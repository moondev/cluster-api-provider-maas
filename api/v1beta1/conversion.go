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
	"github.com/moondev/cluster-api-provider-maas/api/v1beta2"
	clusterv1beta2 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	conversion "sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the receiver to the hub version (v1beta2).
func (src *MaasCluster) ConvertTo(dst conversion.Hub) error {
	return src.convertTo(dst.(*v1beta2.MaasCluster))
}

func (src *MaasCluster) convertTo(dst *v1beta2.MaasCluster) error {
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.DNSDomain = src.Spec.DNSDomain
	dst.Spec.ControlPlaneEndpoint = v1beta2.APIEndpoint{Host: src.Spec.ControlPlaneEndpoint.Host, Port: src.Spec.ControlPlaneEndpoint.Port}
	dst.Spec.FailureDomains = append([]string(nil), src.Spec.FailureDomains...)
	dst.Status.Ready = src.Status.Ready
	dst.Status.Network = v1beta2.Network(src.Status.Network)
	if src.Status.FailureDomains != nil {
		dst.Status.FailureDomains = make([]clusterv1beta2.FailureDomain, len(src.Status.FailureDomains))
		for i := range src.Status.FailureDomains {
			dst.Status.FailureDomains[i] = *(&src.Status.FailureDomains[i]).DeepCopy()
		}
	}
	dst.Status.Conditions = copyConditions(src.Status.Conditions)
	return nil
}

// ConvertFrom converts from the hub version (v1beta2) to the receiver.
func (dst *MaasCluster) ConvertFrom(src conversion.Hub) error {
	return dst.convertFrom(src.(*v1beta2.MaasCluster))
}

func (dst *MaasCluster) convertFrom(src *v1beta2.MaasCluster) error {
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.DNSDomain = src.Spec.DNSDomain
	dst.Spec.ControlPlaneEndpoint = APIEndpoint{Host: src.Spec.ControlPlaneEndpoint.Host, Port: src.Spec.ControlPlaneEndpoint.Port}
	dst.Spec.FailureDomains = append([]string(nil), src.Spec.FailureDomains...)
	dst.Status.Ready = src.Status.Ready
	dst.Status.Network = Network(src.Status.Network)
	if src.Status.FailureDomains != nil {
		dst.Status.FailureDomains = make([]clusterv1beta2.FailureDomain, len(src.Status.FailureDomains))
		for i := range src.Status.FailureDomains {
			dst.Status.FailureDomains[i] = *(&src.Status.FailureDomains[i]).DeepCopy()
		}
	}
	dst.Status.Conditions = copyConditions(src.Status.Conditions)
	return nil
}

// MaasClusterList
func (src *MaasClusterList) ConvertTo(dst conversion.Hub) error {
	return src.convertTo(dst.(*v1beta2.MaasClusterList))
}

func (src *MaasClusterList) convertTo(dst *v1beta2.MaasClusterList) error {
	dst.ListMeta = src.ListMeta
	dst.Items = make([]v1beta2.MaasCluster, len(src.Items))
	for i := range src.Items {
		if err := src.Items[i].convertTo(&dst.Items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (dst *MaasClusterList) ConvertFrom(src conversion.Hub) error {
	return dst.convertFrom(src.(*v1beta2.MaasClusterList))
}

func (dst *MaasClusterList) convertFrom(src *v1beta2.MaasClusterList) error {
	dst.ListMeta = src.ListMeta
	dst.Items = make([]MaasCluster, len(src.Items))
	for i := range src.Items {
		if err := dst.Items[i].convertFrom(&src.Items[i]); err != nil {
			return err
		}
	}
	return nil
}

// MaasMachine
func (src *MaasMachine) ConvertTo(dst conversion.Hub) error {
	return src.convertTo(dst.(*v1beta2.MaasMachine))
}

func (src *MaasMachine) convertTo(dst *v1beta2.MaasMachine) error {
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = v1beta2.MaasMachineSpec{
		FailureDomain: src.Spec.FailureDomain,
		SystemID:      src.Spec.SystemID,
		ProviderID:    src.Spec.ProviderID,
		ResourcePool:  src.Spec.ResourcePool,
		MinCPU:        src.Spec.MinCPU,
		MinMemoryInMB: src.Spec.MinMemoryInMB,
		Tags:          append([]string(nil), src.Spec.Tags...),
		Image:         src.Spec.Image,
		Ephemeral:     src.Spec.Ephemeral,
	}
	dst.Status.Ready = src.Status.Ready
	dst.Status.MachineState = (*v1beta2.MachineState)(src.Status.MachineState)
	dst.Status.MachinePowered = src.Status.MachinePowered
	dst.Status.Hostname = src.Status.Hostname
	dst.Status.DNSAttached = src.Status.DNSAttached
	dst.Status.Addresses = append([]clusterv1beta2.MachineAddress(nil), src.Status.Addresses...)
	dst.Status.FailureReason = src.Status.FailureReason
	dst.Status.FailureMessage = src.Status.FailureMessage
	dst.Status.Conditions = copyConditions(src.Status.Conditions)
	return nil
}

func (dst *MaasMachine) ConvertFrom(src conversion.Hub) error {
	return dst.convertFrom(src.(*v1beta2.MaasMachine))
}

func (dst *MaasMachine) convertFrom(src *v1beta2.MaasMachine) error {
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec = MaasMachineSpec{
		FailureDomain: src.Spec.FailureDomain,
		SystemID:      src.Spec.SystemID,
		ProviderID:    src.Spec.ProviderID,
		ResourcePool:  src.Spec.ResourcePool,
		MinCPU:        src.Spec.MinCPU,
		MinMemoryInMB: src.Spec.MinMemoryInMB,
		Tags:          append([]string(nil), src.Spec.Tags...),
		Image:         src.Spec.Image,
		Ephemeral:     src.Spec.Ephemeral,
	}
	dst.Status.Ready = src.Status.Ready
	dst.Status.MachineState = (*MachineState)(src.Status.MachineState)
	dst.Status.MachinePowered = src.Status.MachinePowered
	dst.Status.Hostname = src.Status.Hostname
	dst.Status.DNSAttached = src.Status.DNSAttached
	dst.Status.Addresses = append([]clusterv1beta2.MachineAddress(nil), src.Status.Addresses...)
	dst.Status.FailureReason = src.Status.FailureReason
	dst.Status.FailureMessage = src.Status.FailureMessage
	dst.Status.Conditions = copyConditions(src.Status.Conditions)
	return nil
}

// MaasMachineList
func (src *MaasMachineList) ConvertTo(dst conversion.Hub) error {
	return src.convertTo(dst.(*v1beta2.MaasMachineList))
}

func (src *MaasMachineList) convertTo(dst *v1beta2.MaasMachineList) error {
	dst.ListMeta = src.ListMeta
	dst.Items = make([]v1beta2.MaasMachine, len(src.Items))
	for i := range src.Items {
		if err := src.Items[i].convertTo(&dst.Items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (dst *MaasMachineList) ConvertFrom(src conversion.Hub) error {
	return dst.convertFrom(src.(*v1beta2.MaasMachineList))
}

func (dst *MaasMachineList) convertFrom(src *v1beta2.MaasMachineList) error {
	dst.ListMeta = src.ListMeta
	dst.Items = make([]MaasMachine, len(src.Items))
	for i := range src.Items {
		if err := dst.Items[i].convertFrom(&src.Items[i]); err != nil {
			return err
		}
	}
	return nil
}

// MaasMachineTemplate
func (src *MaasMachineTemplate) ConvertTo(dst conversion.Hub) error {
	return src.convertTo(dst.(*v1beta2.MaasMachineTemplate))
}

func (src *MaasMachineTemplate) convertTo(dst *v1beta2.MaasMachineTemplate) error {
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Template.Spec = v1beta2.MaasMachineSpec{
		FailureDomain: src.Spec.Template.Spec.FailureDomain,
		SystemID:      src.Spec.Template.Spec.SystemID,
		ProviderID:    src.Spec.Template.Spec.ProviderID,
		ResourcePool:  src.Spec.Template.Spec.ResourcePool,
		MinCPU:        src.Spec.Template.Spec.MinCPU,
		MinMemoryInMB: src.Spec.Template.Spec.MinMemoryInMB,
		Tags:          append([]string(nil), src.Spec.Template.Spec.Tags...),
		Image:         src.Spec.Template.Spec.Image,
		Ephemeral:     src.Spec.Template.Spec.Ephemeral,
	}
	return nil
}

func (dst *MaasMachineTemplate) ConvertFrom(src conversion.Hub) error {
	return dst.convertFrom(src.(*v1beta2.MaasMachineTemplate))
}

func (dst *MaasMachineTemplate) convertFrom(src *v1beta2.MaasMachineTemplate) error {
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Template.Spec = MaasMachineSpec{
		FailureDomain: src.Spec.Template.Spec.FailureDomain,
		SystemID:      src.Spec.Template.Spec.SystemID,
		ProviderID:    src.Spec.Template.Spec.ProviderID,
		ResourcePool:  src.Spec.Template.Spec.ResourcePool,
		MinCPU:        src.Spec.Template.Spec.MinCPU,
		MinMemoryInMB: src.Spec.Template.Spec.MinMemoryInMB,
		Tags:          append([]string(nil), src.Spec.Template.Spec.Tags...),
		Image:         src.Spec.Template.Spec.Image,
		Ephemeral:     src.Spec.Template.Spec.Ephemeral,
	}
	return nil
}

// MaasMachineTemplateList
func (src *MaasMachineTemplateList) ConvertTo(dst conversion.Hub) error {
	return src.convertTo(dst.(*v1beta2.MaasMachineTemplateList))
}

func (src *MaasMachineTemplateList) convertTo(dst *v1beta2.MaasMachineTemplateList) error {
	dst.ListMeta = src.ListMeta
	dst.Items = make([]v1beta2.MaasMachineTemplate, len(src.Items))
	for i := range src.Items {
		if err := src.Items[i].convertTo(&dst.Items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (dst *MaasMachineTemplateList) ConvertFrom(src conversion.Hub) error {
	return dst.convertFrom(src.(*v1beta2.MaasMachineTemplateList))
}

func (dst *MaasMachineTemplateList) convertFrom(src *v1beta2.MaasMachineTemplateList) error {
	dst.ListMeta = src.ListMeta
	dst.Items = make([]MaasMachineTemplate, len(src.Items))
	for i := range src.Items {
		if err := dst.Items[i].convertFrom(&src.Items[i]); err != nil {
			return err
		}
	}
	return nil
}

// MaasClusterTemplate
func (src *MaasClusterTemplate) ConvertTo(dst conversion.Hub) error {
	return src.convertTo(dst.(*v1beta2.MaasClusterTemplate))
}

func (src *MaasClusterTemplate) convertTo(dst *v1beta2.MaasClusterTemplate) error {
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Template.Spec.DNSDomain = src.Spec.Template.Spec.DNSDomain
	dst.Spec.Template.Spec.ControlPlaneEndpoint = v1beta2.APIEndpoint{Host: src.Spec.Template.Spec.ControlPlaneEndpoint.Host, Port: src.Spec.Template.Spec.ControlPlaneEndpoint.Port}
	dst.Spec.Template.Spec.FailureDomains = append([]string(nil), src.Spec.Template.Spec.FailureDomains...)
	dst.Status.Ready = src.Status.Ready
	dst.Status.Network = v1beta2.Network(src.Status.Network)
	if src.Status.FailureDomains != nil {
		dst.Status.FailureDomains = make([]clusterv1beta2.FailureDomain, len(src.Status.FailureDomains))
		for i := range src.Status.FailureDomains {
			dst.Status.FailureDomains[i] = *(&src.Status.FailureDomains[i]).DeepCopy()
		}
	}
	dst.Status.Conditions = copyConditions(src.Status.Conditions)
	return nil
}

func (dst *MaasClusterTemplate) ConvertFrom(src conversion.Hub) error {
	return dst.convertFrom(src.(*v1beta2.MaasClusterTemplate))
}

func (dst *MaasClusterTemplate) convertFrom(src *v1beta2.MaasClusterTemplate) error {
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Template.Spec.DNSDomain = src.Spec.Template.Spec.DNSDomain
	dst.Spec.Template.Spec.ControlPlaneEndpoint = APIEndpoint{Host: src.Spec.Template.Spec.ControlPlaneEndpoint.Host, Port: src.Spec.Template.Spec.ControlPlaneEndpoint.Port}
	dst.Spec.Template.Spec.FailureDomains = append([]string(nil), src.Spec.Template.Spec.FailureDomains...)
	dst.Status.Ready = src.Status.Ready
	dst.Status.Network = Network(src.Status.Network)
	if src.Status.FailureDomains != nil {
		dst.Status.FailureDomains = make([]clusterv1beta2.FailureDomain, len(src.Status.FailureDomains))
		for i := range src.Status.FailureDomains {
			dst.Status.FailureDomains[i] = *(&src.Status.FailureDomains[i]).DeepCopy()
		}
	}
	dst.Status.Conditions = copyConditions(src.Status.Conditions)
	return nil
}

// MaasClusterTemplateList
func (src *MaasClusterTemplateList) ConvertTo(dst conversion.Hub) error {
	return src.convertTo(dst.(*v1beta2.MaasClusterTemplateList))
}

func (src *MaasClusterTemplateList) convertTo(dst *v1beta2.MaasClusterTemplateList) error {
	dst.ListMeta = src.ListMeta
	dst.Items = make([]v1beta2.MaasClusterTemplate, len(src.Items))
	for i := range src.Items {
		if err := src.Items[i].convertTo(&dst.Items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (dst *MaasClusterTemplateList) ConvertFrom(src conversion.Hub) error {
	return dst.convertFrom(src.(*v1beta2.MaasClusterTemplateList))
}

func (dst *MaasClusterTemplateList) convertFrom(src *v1beta2.MaasClusterTemplateList) error {
	dst.ListMeta = src.ListMeta
	dst.Items = make([]MaasClusterTemplate, len(src.Items))
	for i := range src.Items {
		if err := dst.Items[i].convertFrom(&src.Items[i]); err != nil {
			return err
		}
	}
	return nil
}

// copyConditions copies conditions from source to a new slice.
func copyConditions(in clusterv1beta2.Conditions) clusterv1beta2.Conditions {
	if in == nil {
		return nil
	}
	out := make(clusterv1beta2.Conditions, len(in))
	for i := range in {
		in[i].DeepCopyInto(&out[i])
	}
	return out
}
