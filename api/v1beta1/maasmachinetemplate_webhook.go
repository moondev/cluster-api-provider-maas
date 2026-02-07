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
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// log is for logging in this package.
var maasmachinetemplatelog = logf.Log.WithName("maasmachinetemplate-resource")

func (r *MaasMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-maasmachinetemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasmachinetemplates,verbs=create;update,versions=v1beta1,name=mmaasmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-maasmachinetemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasmachinetemplates,versions=v1beta1,name=vmaasmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1

var (
	_ admission.CustomDefaulter = &MaasMachineTemplate{}
	_ admission.CustomValidator = &MaasMachineTemplate{}
)

// Default implements admission.CustomDefaulter so a webhook will be registered for the type
func (r *MaasMachineTemplate) Default(_ context.Context, obj runtime.Object) error {
	maasMachineTemplate, ok := obj.(*MaasMachineTemplate)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a MaasMachineTemplate but got a %T", obj))
	}
	maasmachinetemplatelog.Info("default", "name", maasMachineTemplate.Name)
	return nil
}

// ValidateCreate implements admission.CustomValidator so a webhook will be registered for the type
func (r *MaasMachineTemplate) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	maasMachineTemplate, ok := obj.(*MaasMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a MaasMachineTemplate but got a %T", obj))
	}
	maasmachinetemplatelog.Info("validate create", "name", maasMachineTemplate.Name)

	// ClusterClass/managed topology validation
	if isManagedTopologyTemplate(maasMachineTemplate) {
		flds := field.ErrorList{}
		spec := maasMachineTemplate.Spec.Template.Spec
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
				maasMachineTemplate.Name,
				flds,
			)
		}
	}
	return nil, nil
}

// ValidateUpdate implements admission.CustomValidator so a webhook will be registered for the type
func (r *MaasMachineTemplate) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	oldM, ok := oldObj.(*MaasMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a MaasMachineTemplate but got a %T", oldObj))
	}
	maasMachineTemplate, ok := newObj.(*MaasMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a MaasMachineTemplate but got a %T", newObj))
	}
	maasmachinetemplatelog.Info("validate update", "name", maasMachineTemplate.Name)

	if maasMachineTemplate.Spec.Template.Spec.Image != oldM.Spec.Template.Spec.Image {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("maas machine template image change is not allowed, old=%s, new=%s", oldM.Spec.Template.Spec.Image, maasMachineTemplate.Spec.Template.Spec.Image))
	}

	if *maasMachineTemplate.Spec.Template.Spec.MinCPU != *oldM.Spec.Template.Spec.MinCPU {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("maas machine template min cpu count change is not allowed, old=%d, new=%d", oldM.Spec.Template.Spec.MinCPU, maasMachineTemplate.Spec.Template.Spec.MinCPU))
	}

	if *maasMachineTemplate.Spec.Template.Spec.MinMemoryInMB != *oldM.Spec.Template.Spec.MinMemoryInMB {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("maas machine template min memory change is not allowed, old=%d MB, new=%d MB", oldM.Spec.Template.Spec.MinMemoryInMB, maasMachineTemplate.Spec.Template.Spec.MinMemoryInMB))
	}

	// ClusterClass/managed topology validation
	if isManagedTopologyTemplate(maasMachineTemplate) {
		flds := field.ErrorList{}
		spec := maasMachineTemplate.Spec.Template.Spec
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
				maasMachineTemplate.Name,
				flds,
			)
		}
	}
	return nil, nil
}

// ValidateDelete implements admission.CustomValidator so a webhook will be registered for the type
func (r *MaasMachineTemplate) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	maasMachineTemplate, ok := obj.(*MaasMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a MaasMachineTemplate but got a %T", obj))
	}
	maasmachinetemplatelog.Info("validate delete", "name", maasMachineTemplate.Name)
	return nil, nil
}

// isManagedTopologyTemplate returns true if the MaasMachineTemplate is being used in a managed topology (ClusterClass) context.
func isManagedTopologyTemplate(r *MaasMachineTemplate) bool {
	_, ok := r.Labels["topology.cluster.x-k8s.io/owned"]
	return ok
}
