/*
Copyright 2020 The Kubernetes Authors.

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

import clusterv1beta2 "sigs.k8s.io/cluster-api/api/core/v1beta2"

const (
	MachineDeployedCondition clusterv1beta2.ConditionType = "MachineDeployed"

	WaitingForClusterInfrastructureReason = "WaitingForClusterInfrastructure"
	WaitingForBootstrapDataReason         = "WaitingForBootstrapData"
	MachineDeployingReason                = "MachineDeploying"
	MachineTerminatedReason               = "MachineTerminatedReason"
	MachinePoweredOffReason               = "MachinePoweredOff"
	MachineNotFoundReason                 = "MachineNotFound"
	MachineDeployFailedReason             = "MachineDeployFailed"
	MachineDeployStartedReason            = "MachineDeployStartedReason"
)

const (
	DNSAttachedCondition clusterv1beta2.ConditionType = "DNSAttached"
	DNSDetachPending     = "DNSDetachPending"
	DNSAttachPending     = "DNSAttachPending"
)

const (
	DNSReadyCondition    clusterv1beta2.ConditionType = "LoadBalancerReady"
	DNSFailedReason      = "LoadBalancerFailed"
	WaitForDNSNameReason = "WaitForDNSName"
)

const (
	APIServerAvailableCondition clusterv1beta2.ConditionType = "APIServerAvailable"
	APIServerNotReadyReason     = "APIServerNotReady"
)
