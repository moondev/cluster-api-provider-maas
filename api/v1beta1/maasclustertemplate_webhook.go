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
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var maasclustertemplatelog = logf.Log.WithName("maasclustertemplate-resource")

func (r *MaasClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-maasclustertemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasclustertemplates,verbs=create;update,versions=v1beta1,name=mmaasclustertemplate.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-maasclustertemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasclustertemplates,versions=v1beta1,name=vmaasclustertemplate.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1

var (
	_ webhook.Defaulter = &MaasClusterTemplate{}
	_ webhook.Validator = &MaasClusterTemplate{}
)

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MaasClusterTemplate) Default() {
	maasclustertemplatelog.Info("default", "name", r.Name)
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MaasClusterTemplate) ValidateCreate() (admission.Warnings, error) {
	maasclustertemplatelog.Info("validate create", "name", r.Name)

	// ClusterClass/managed topology validation
	if isManagedTopologyClusterTemplate(r) {
		flds := field.ErrorList{}
		spec := r.Spec.Template.Spec
		if spec.DNSDomain == "" {
			flds = append(flds, field.Required(field.NewPath("spec").Child("template").Child("spec").Child("dnsDomain"), "spec.template.spec.dnsDomain is required for managed topology (ClusterClass)"))
		}
		if len(flds) > 0 {
			return nil, apierrors.NewInvalid(
				GroupVersion.WithKind("MaasClusterTemplate").GroupKind(),
				r.Name,
				flds,
			)
		}
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MaasClusterTemplate) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	maasclustertemplatelog.Info("validate update", "name", r.Name)
	oldT, ok := old.(*MaasClusterTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a MaasClusterTemplate but got a %T", old))
	}

	// Prevent changing dnsDomain
	if r.Spec.Template.Spec.DNSDomain != oldT.Spec.Template.Spec.DNSDomain {
		return nil, apierrors.NewBadRequest("changing cluster DNS Domain not allowed")
	}

	// ClusterClass/managed topology validation
	if isManagedTopologyClusterTemplate(r) {
		flds := field.ErrorList{}
		spec := r.Spec.Template.Spec
		if spec.DNSDomain == "" {
			flds = append(flds, field.Required(field.NewPath("spec").Child("template").Child("spec").Child("dnsDomain"), "spec.template.spec.dnsDomain is required for managed topology (ClusterClass)"))
		}
		if len(flds) > 0 {
			return nil, apierrors.NewInvalid(
				GroupVersion.WithKind("MaasClusterTemplate").GroupKind(),
				r.Name,
				flds,
			)
		}
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MaasClusterTemplate) ValidateDelete() (admission.Warnings, error) {
	maasclustertemplatelog.Info("validate delete", "name", r.Name)
	return nil, nil
}

// isManagedTopologyClusterTemplate returns true if the MaasClusterTemplate is being used in a managed topology (ClusterClass) context.
func isManagedTopologyClusterTemplate(r *MaasClusterTemplate) bool {
	_, ok := r.Labels["topology.cluster.x-k8s.io/owned"]
	return ok
}
