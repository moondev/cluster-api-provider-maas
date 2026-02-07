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
var maasclusterlog = logf.Log.WithName("maascluster-resource")

func (r *MaasCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-maascluster,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasclusters,verbs=create;update,versions=v1beta1,name=mmaascluster.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-maascluster,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasclusters,versions=v1beta1,name=vmaascluster.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1

var (
	_ admission.CustomDefaulter = &MaasCluster{}
	_ admission.CustomValidator = &MaasCluster{}
)

// Default implements admission.CustomDefaulter so a webhook will be registered for the type
func (r *MaasCluster) Default(_ context.Context, obj runtime.Object) error {
	maasCluster, ok := obj.(*MaasCluster)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a MaasCluster but got a %T", obj))
	}
	maasclusterlog.Info("default", "name", maasCluster.Name)
	return nil
}

// ValidateCreate implements admission.CustomValidator so a webhook will be registered for the type
func (r *MaasCluster) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	maasCluster, ok := obj.(*MaasCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a MaasCluster but got a %T", obj))
	}
	maasclusterlog.Info("validate create", "name", maasCluster.Name)

	// ClusterClass/managed topology validation
	if isManagedTopology(maasCluster) {
		if maasCluster.Spec.DNSDomain == "" {
			return nil, apierrors.NewInvalid(
				GroupVersion.WithKind("MaasCluster").GroupKind(),
				maasCluster.Name,
				field.ErrorList{
					field.Required(field.NewPath("spec").Child("dnsDomain"), "spec.dnsDomain is required for managed topology (ClusterClass)"),
				},
			)
		}
	}
	return nil, nil
}

// ValidateUpdate implements admission.CustomValidator so a webhook will be registered for the type
func (r *MaasCluster) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	oldC, ok := oldObj.(*MaasCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a MaasCluster but got a %T", oldObj))
	}
	maasCluster, ok := newObj.(*MaasCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a MaasCluster but got a %T", newObj))
	}

	maasclusterlog.Info("validate update", "name", maasCluster.Name)

	if maasCluster.Spec.DNSDomain != oldC.Spec.DNSDomain {
		return nil, apierrors.NewBadRequest("changing cluster DNS Domain not allowed")
	}

	// ClusterClass/managed topology validation
	if isManagedTopology(maasCluster) {
		if maasCluster.Spec.DNSDomain == "" {
			return nil, apierrors.NewInvalid(
				GroupVersion.WithKind("MaasCluster").GroupKind(),
				maasCluster.Name,
				field.ErrorList{
					field.Required(field.NewPath("spec").Child("dnsDomain"), "spec.dnsDomain is required for managed topology (ClusterClass)"),
				},
			)
		}
	}
	return nil, nil
}

// ValidateDelete implements admission.CustomValidator so a webhook will be registered for the type
func (r *MaasCluster) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	maasCluster, ok := obj.(*MaasCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a MaasCluster but got a %T", obj))
	}
	maasclusterlog.Info("validate delete", "name", maasCluster.Name)
	return nil, nil
}

// isManagedTopology returns true if the MaasCluster is being used in a managed topology (ClusterClass) context.
func isManagedTopology(r *MaasCluster) bool {
	// Managed topology clusters have the label: "topology.cluster.x-k8s.io/owned: ""
	_, ok := r.Labels["topology.cluster.x-k8s.io/owned"]
	return ok
}
