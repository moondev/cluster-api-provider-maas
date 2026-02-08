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

var maasclustertemplatelog = logf.Log.WithName("maasclustertemplate-resource")

func (r *MaasClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-maasclustertemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasclustertemplates,verbs=create;update,versions=v1beta2,name=mmaasclustertemplate.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-maasclustertemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasclustertemplates,versions=v1beta2,name=vmaasclustertemplate.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1

var (
	_ webhook.CustomDefaulter = &MaasClusterTemplate{}
	_ webhook.CustomValidator = &MaasClusterTemplate{}
)

func (r *MaasClusterTemplate) Default(_ context.Context, obj runtime.Object) error {
	template := obj.(*MaasClusterTemplate)
	maasclustertemplatelog.Info("default", "name", template.Name)
	return nil
}

func (r *MaasClusterTemplate) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	template := obj.(*MaasClusterTemplate)
	maasclustertemplatelog.Info("validate create", "name", template.Name)
	if isManagedTopologyClusterTemplate(template) {
		flds := field.ErrorList{}
		spec := template.Spec.Template.Spec
		if spec.DNSDomain == "" {
			flds = append(flds, field.Required(field.NewPath("spec").Child("template").Child("spec").Child("dnsDomain"), "spec.template.spec.dnsDomain is required for managed topology (ClusterClass)"))
		}
		if len(flds) > 0 {
			return nil, apierrors.NewInvalid(
				GroupVersion.WithKind("MaasClusterTemplate").GroupKind(),
				template.Name,
				flds,
			)
		}
	}
	return nil, nil
}

func (r *MaasClusterTemplate) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	template := newObj.(*MaasClusterTemplate)
	maasclustertemplatelog.Info("validate update", "name", template.Name)
	oldT, ok := oldObj.(*MaasClusterTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a MaasClusterTemplate but got a %T", oldObj))
	}
	if template.Spec.Template.Spec.DNSDomain != oldT.Spec.Template.Spec.DNSDomain {
		return nil, apierrors.NewBadRequest("changing cluster DNS Domain not allowed")
	}
	if isManagedTopologyClusterTemplate(template) && template.Spec.Template.Spec.DNSDomain == "" {
		return nil, apierrors.NewInvalid(
			GroupVersion.WithKind("MaasClusterTemplate").GroupKind(),
			template.Name,
			field.ErrorList{
				field.Required(field.NewPath("spec").Child("template").Child("spec").Child("dnsDomain"), "spec.template.spec.dnsDomain is required for managed topology (ClusterClass)"),
			},
		)
	}
	return nil, nil
}

func (r *MaasClusterTemplate) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	template := obj.(*MaasClusterTemplate)
	maasclustertemplatelog.Info("validate delete", "name", template.Name)
	return nil, nil
}

func isManagedTopologyClusterTemplate(r *MaasClusterTemplate) bool {
	_, ok := r.Labels["topology.cluster.x-k8s.io/owned"]
	return ok
}
