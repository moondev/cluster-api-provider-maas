package machine

import (
	"github.com/canonical/gomaasclient/client"
	"github.com/canonical/gomaasclient/entity"
	"github.com/pkg/errors"
	"k8s.io/klog/v2/textlogger"

	infrav1beta1 "github.com/moondev/cluster-api-provider-maas/api/v1beta1"
	"github.com/moondev/cluster-api-provider-maas/pkg/maas/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// Service manages the MaaS machine
type Service struct {
	scope      *scope.MachineScope
	maasClient *client.Client
}

// DNS service returns a new helper for managing a MaaS "DNS" (DNS client loadbalancing)
func NewService(machineScope *scope.MachineScope) *Service {
	return &Service{
		scope:      machineScope,
		maasClient: scope.NewMaasClient(machineScope.ClusterScope),
	}
}

func (s *Service) GetMachine(systemID string) (*infrav1beta1.Machine, error) {

	if systemID == "" {
		return nil, nil
	}

	params := &entity.MachinesParams{
		ID: []string{systemID},
	}

	machines, err := s.maasClient.Machines.Get(params)
	if err != nil {
		return nil, err
	}

	if len(machines) == 0 {
		return nil, nil
	}

	machine := fromSDKTypeToMachine(&machines[0])

	return machine, nil
}

func (s *Service) ReleaseMachine(systemID string) error {
	err := s.maasClient.Machines.Release([]string{systemID}, "released by cluster-api-provider-maas")
	if err != nil {
		return errors.Wrapf(err, "Unable to release machine")
	}

	return nil
}

func (s *Service) DeployMachine(userDataB64 string) (_ *infrav1beta1.Machine, rerr error) {
	log := textlogger.NewLogger(textlogger.NewConfig())

	mm := s.scope.MaasMachine

	failureDomain := mm.Spec.FailureDomain
	if failureDomain == nil {
		failureDomain = s.scope.Machine.Spec.FailureDomain
	}

	var m *entity.Machine
	var err error

	if s.scope.GetProviderID() == "" {
		// Allocate a new machine
		allocateParams := &entity.MachineAllocateParams{
			CPUCount: *mm.Spec.MinCPU,
			Mem:      int64(*mm.Spec.MinMemoryInMB),
		}

		if failureDomain != nil {
			allocateParams.Zone = *failureDomain
		}

		if mm.Spec.ResourcePool != nil {
			allocateParams.Pool = *mm.Spec.ResourcePool
		}

		if len(mm.Spec.Tags) > 0 {
			allocateParams.Tags = mm.Spec.Tags
		}

		m, err = s.maasClient.Machines.Allocate(allocateParams)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to allocate machine")
		}

		s.scope.SetProviderID(m.SystemID, m.Zone.Name)
		err = s.scope.PatchObject()
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to patch object")
		}
	} else {
		// Get existing machine
		params := &entity.MachinesParams{
			ID: []string{s.scope.GetProviderID()},
		}
		machines, err := s.maasClient.Machines.Get(params)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to get machine")
		}
		if len(machines) == 0 {
			return nil, errors.New("Machine not found")
		}
		m = &machines[0]
	}

	// Deploy the machine with ephemeral support
	deployParams := &entity.MachineDeployParams{
		UserData:     userDataB64,
		DistroSeries: mm.Spec.Image,
	}

	// Add ephemeral deployment if specified
	if mm.Spec.Ephemeral {
		deployParams.EphemeralDeploy = true
	}

	deployingM, err := s.maasClient.Machine.Deploy(m.SystemID, deployParams)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to deploy machine")
	}

	log.Info("Machine deployed", "systemID", deployingM.SystemID, "hostname", deployingM.Hostname)

	machine := fromSDKTypeToMachine(deployingM)

	return machine, nil
}

func fromSDKTypeToMachine(m *entity.Machine) *infrav1beta1.Machine {
	machine := &infrav1beta1.Machine{
		ID:               m.SystemID,
		Hostname:         m.Hostname,
		State:            infrav1beta1.MachineState(m.StatusName),
		Powered:          m.PowerState == "on",
		AvailabilityZone: m.Zone.Name,
	}

	// Add IP addresses if available
	if m.IPAddresses != nil {
		for _, addr := range m.IPAddresses {
			machine.Addresses = append(machine.Addresses, clusterv1.MachineAddress{
				Type:    clusterv1.MachineExternalIP,
				Address: addr.String(),
			})
		}
	}

	return machine
}

func (s *Service) PowerOnMachine() error {
	// For the canonical client, we need to use the Machine API to power on
	// This would typically involve calling a power action on the machine
	// For now, we'll use a simple approach - the machine should be powered on during deployment
	return nil
}

//// ReconcileDNS reconciles the load balancers for the given cluster.
//func (s *Service) ReconcileDNS() error {
//	s.scope.V(2).Info("Reconciling DNS")
//
//	s.scope.SetDNSName("cluster1.maas")
//	return nil
//}
//
