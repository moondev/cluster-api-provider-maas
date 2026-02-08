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
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var maasmachinetemplatelog = logf.Log.WithName("maasmachinetemplate-resource")

func (r *MaasMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-maasmachinetemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasmachinetemplates,verbs=create;update,versions=v1beta2,name=mmaasmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-maasmachinetemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasmachinetemplates,versions=v1beta2,name=vmaasmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1

var (
	_ webhook.CustomDefaulter = &MaasMachineTemplate{}
	_ webhook.CustomValidator = &MaasMachineTemplate{}
)

func (r *MaasMachineTemplate) Default(_ context.Context, obj runtime.Object) error {
	template := obj.(*MaasMachineTemplate)
	maasmachinetemplatelog.Info("default", "name", template.Name)
	return nil
}

func (r *MaasMachineTemplate) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	template := obj.(*MaasMachineTemplate)
	maasmachinetemplatelog.Info("validate create", "name", template.Name)
	if isManagedTopologyTemplate(template) {
		flds := field.ErrorList{}
		spec := template.Spec.Template.Spec
		if spec.Image == "" {
			flds = append(flds, field.Required(field.NewPath("spec").Child("template").Child("spec").Child("image"), "spec.template.spec.image is required for managed topology (ClusterClass)"))
		}
		if spec.MinCPU == nil {
			flds = append(flds, field.Required(field.NewPath("spec").Child("template").Child("spec").Child("minCPU"), "spec.template.spec.minCPU is required for managed topology (ClusterClass)"))
		}
		if spec.MinMemoryInMB == nil {
			flds = append(flds, field.Required(field.NewPath("spec").Child("template").Child("spec").Child("minMemory"), "spec.template.spec.minMemory is required for managed topology (ClusterClass)"))
		}
		if len(flds) > 0 {
			return nil, apierrors.NewInvalid(
				GroupVersion.WithKind("MaasMachineTemplate").GroupKind(),
				template.Name,
				flds,
			)
		}
	}
	return nil, nil
}

func (r *MaasMachineTemplate) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	template := newObj.(*MaasMachineTemplate)
	maasmachinetemplatelog.Info("validate update", "name", template.Name)
	oldM := oldObj.(*MaasMachineTemplate)
	if template.Spec.Template.Spec.Image != oldM.Spec.Template.Spec.Image {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("maas machine template image change is not allowed, old=%s, new=%s", oldM.Spec.Template.Spec.Image, template.Spec.Template.Spec.Image))
	}
	if (template.Spec.Template.Spec.MinCPU != nil) != (oldM.Spec.Template.Spec.MinCPU != nil) || (template.Spec.Template.Spec.MinCPU != nil && *template.Spec.Template.Spec.MinCPU != *oldM.Spec.Template.Spec.MinCPU) {
		return nil, apierrors.NewBadRequest("maas machine template min cpu count change is not allowed")
	}
	if (template.Spec.Template.Spec.MinMemoryInMB != nil) != (oldM.Spec.Template.Spec.MinMemoryInMB != nil) || (template.Spec.Template.Spec.MinMemoryInMB != nil && *template.Spec.Template.Spec.MinMemoryInMB != *oldM.Spec.Template.Spec.MinMemoryInMB) {
		return nil, apierrors.NewBadRequest("maas machine template min memory change is not allowed")
	}
	if isManagedTopologyTemplate(template) {
		flds := field.ErrorList{}
		spec := template.Spec.Template.Spec
		if spec.Image == "" {
			flds = append(flds, field.Required(field.NewPath("spec").Child("template").Child("spec").Child("image"), "spec.template.spec.image is required for managed topology (ClusterClass)"))
		}
		if spec.MinCPU == nil {
			flds = append(flds, field.Required(field.NewPath("spec").Child("template").Child("spec").Child("minCPU"), "spec.template.spec.minCPU is required for managed topology (ClusterClass)"))
		}
		if spec.MinMemoryInMB == nil {
			flds = append(flds, field.Required(field.NewPath("spec").Child("template").Child("spec").Child("minMemory"), "spec.template.spec.minMemory is required for managed topology (ClusterClass)"))
		}
		if len(flds) > 0 {
			return nil, apierrors.NewInvalid(
				GroupVersion.WithKind("MaasMachineTemplate").GroupKind(),
				template.Name,
				flds,
			)
		}
	}
	return nil, nil
}

func (r *MaasMachineTemplate) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	template := obj.(*MaasMachineTemplate)
	maasmachinetemplatelog.Info("validate delete", "name", template.Name)
	return nil, nil
}

func isManagedTopologyTemplate(r *MaasMachineTemplate) bool {
	_, ok := r.Labels["topology.cluster.x-k8s.io/owned"]
	return ok
}
