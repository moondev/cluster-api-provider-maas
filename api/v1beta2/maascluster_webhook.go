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

var maasclusterlog = logf.Log.WithName("maascluster-resource")

func (r *MaasCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-maascluster,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasclusters,verbs=create;update,versions=v1beta2,name=mmaascluster.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-maascluster,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasclusters,versions=v1beta2,name=vmaascluster.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1

var (
	_ webhook.CustomDefaulter = &MaasCluster{}
	_ webhook.CustomValidator = &MaasCluster{}
)

func (r *MaasCluster) Default(_ context.Context, obj runtime.Object) error {
	cluster := obj.(*MaasCluster)
	maasclusterlog.Info("default", "name", cluster.Name)
	return nil
}

func (r *MaasCluster) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	cluster := obj.(*MaasCluster)
	maasclusterlog.Info("validate create", "name", cluster.Name)
	if isManagedTopology(cluster) {
		if cluster.Spec.DNSDomain == "" {
			return nil, apierrors.NewInvalid(
				GroupVersion.WithKind("MaasCluster").GroupKind(),
				cluster.Name,
				field.ErrorList{
					field.Required(field.NewPath("spec").Child("dnsDomain"), "spec.dnsDomain is required for managed topology (ClusterClass)"),
				},
			)
		}
	}
	return nil, nil
}

func (r *MaasCluster) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	cluster := newObj.(*MaasCluster)
	maasclusterlog.Info("validate update", "name", cluster.Name)
	oldC, ok := oldObj.(*MaasCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a MaasCluster but got a %T", oldObj))
	}
	if cluster.Spec.DNSDomain != oldC.Spec.DNSDomain {
		return nil, apierrors.NewBadRequest("changing cluster DNS Domain not allowed")
	}
	if isManagedTopology(cluster) && cluster.Spec.DNSDomain == "" {
		return nil, apierrors.NewInvalid(
			GroupVersion.WithKind("MaasCluster").GroupKind(),
			cluster.Name,
			field.ErrorList{
				field.Required(field.NewPath("spec").Child("dnsDomain"), "spec.dnsDomain is required for managed topology (ClusterClass)"),
			},
		)
	}
	return nil, nil
}

func (r *MaasCluster) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	cluster := obj.(*MaasCluster)
	maasclusterlog.Info("validate delete", "name", cluster.Name)
	return nil, nil
}

func isManagedTopology(r *MaasCluster) bool {
	_, ok := r.Labels["topology.cluster.x-k8s.io/owned"]
	return ok
}
