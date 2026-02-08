package machine

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spectrocloud/maas-client-go/maasclient"

	"k8s.io/klog/v2/textlogger"

	infrav1beta1 "github.com/moondev/cluster-api-provider-maas/api/v1beta1"
	"github.com/moondev/cluster-api-provider-maas/pkg/maas/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// Service manages the MaaS machine
type Service struct {
	scope      *scope.MachineScope
	maasClient maasclient.ClientSetInterface
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

	machine := s.maasClient.Machines().Machine(systemID)
	// We need to fetch the machine details (GET)
	machine, err := machine.Get(context.TODO())
	if err != nil {
		return nil, err
	}

	return fromSDKTypeToMachine(machine), nil
}

func (s *Service) ReleaseMachine(systemID string) error {
	_, err := s.maasClient.Machines().Machine(systemID).Releaser().WithComment("released by cluster-api-provider-maas").Release(context.TODO())
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

	var m maasclient.Machine
	var err error

	if s.scope.GetProviderID() == "" {
		// Allocate a new machine
		allocator := s.maasClient.Machines().Allocator().
			WithCPUCount(*mm.Spec.MinCPU).
			WithMemory(*mm.Spec.MinMemoryInMB)

		if failureDomain != nil {
			allocator = allocator.WithZone(*failureDomain)
		}

		if mm.Spec.ResourcePool != nil {
			allocator = allocator.WithResourcePool(*mm.Spec.ResourcePool)
		}

		if len(mm.Spec.Tags) > 0 {
			allocator = allocator.WithTags(mm.Spec.Tags)
		}

		m, err = allocator.Allocate(context.TODO())
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to allocate machine")
		}

		s.scope.SetProviderID(m.SystemID(), m.ZoneName())
		err = s.scope.PatchObject()
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to patch object")
		}
	} else {
		// Get existing machine
		// But s.scope.GetProviderID() returns maas://... so we need systemID.
		// machineScope.GetInstanceID() parses it.
		// Let's use GetInstanceID()
		systemID := s.scope.GetInstanceID()
		if systemID == nil {
			return nil, errors.New("Unable to get systemID from providerID")
		}

		m = s.maasClient.Machines().Machine(*systemID)
		m, err = m.Get(context.TODO())
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to get machine")
		}
	}

	// Deploy the machine with ephemeral support
	deployer := m.Deployer().
		SetUserData(userDataB64).
		SetDistroSeries(mm.Spec.Image)

	// Add ephemeral deployment if specified
	if mm.Spec.Ephemeral {
		deployer = deployer.SetEphemeralDeploy(true)
	}

	// Use SetEphemeralDeploy for InMemory as well
	if mm.Spec.InMemory != nil && *mm.Spec.InMemory {
		deployer = deployer.SetEphemeralDeploy(true)
	}

	deployingM, err := deployer.Deploy(context.TODO())
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to deploy machine")
	}

	log.Info("Machine deployed", "systemID", deployingM.SystemID(), "hostname", deployingM.Hostname())

	machine := fromSDKTypeToMachine(deployingM)

	return machine, nil
}

func fromSDKTypeToMachine(m maasclient.Machine) *infrav1beta1.Machine {

	machine := &infrav1beta1.Machine{
		ID:               m.SystemID(),
		Hostname:         m.Hostname(),
		State:            infrav1beta1.MachineState(m.State()),
		Powered:          m.PowerState() == "on",
		AvailabilityZone: m.ZoneName(),
	}

	// Add IP addresses if available
	ips := m.IPAddresses()
	if len(ips) > 0 {
		for _, addr := range ips {
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
