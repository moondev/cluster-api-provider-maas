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
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var maasmachinelog = logf.Log.WithName("maasmachine-resource")

func (r *MaasMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-maasmachine,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasmachines,verbs=create;update,versions=v1beta2,name=mmaasmachine.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-maasmachine,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=maasmachines,versions=v1beta2,name=vmaasmachine.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1

var (
	_ webhook.CustomDefaulter = &MaasMachine{}
	_ webhook.CustomValidator = &MaasMachine{}
)

func (r *MaasMachine) Default(_ context.Context, obj runtime.Object) error {
	machine := obj.(*MaasMachine)
	maasmachinelog.Info("default", "name", machine.Name)
	return nil
}

func (r *MaasMachine) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	machine := obj.(*MaasMachine)
	maasmachinelog.Info("validate create", "name", machine.Name)
	return nil, nil
}

func (r *MaasMachine) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	machine := obj.(*MaasMachine)
	maasmachinelog.Info("validate delete", "name", machine.Name)
	return nil, nil
}

func (r *MaasMachine) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	machine := newObj.(*MaasMachine)
	maasmachinelog.Info("validate update", "name", machine.Name)
	oldM := oldObj.(*MaasMachine)
	if machine.Spec.Image != oldM.Spec.Image {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("maas machine image change is not allowed, old=%s, new=%s", oldM.Spec.Image, machine.Spec.Image))
	}
	if (machine.Spec.MinCPU != nil) != (oldM.Spec.MinCPU != nil) || (machine.Spec.MinCPU != nil && *machine.Spec.MinCPU != *oldM.Spec.MinCPU) {
		return nil, apierrors.NewBadRequest("maas machine min cpu count change is not allowed")
	}
	if (machine.Spec.MinMemoryInMB != nil) != (oldM.Spec.MinMemoryInMB != nil) || (machine.Spec.MinMemoryInMB != nil && *machine.Spec.MinMemoryInMB != *oldM.Spec.MinMemoryInMB) {
		return nil, apierrors.NewBadRequest("maas machine min memory change is not allowed")
	}
	return nil, nil
}
