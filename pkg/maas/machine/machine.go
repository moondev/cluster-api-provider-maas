package machine

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/klog/v2/textlogger"

	maasclient "github.com/spectrocloud/maas-client-go/maasclient"

	infrav1beta2 "github.com/moondev/cluster-api-provider-maas/api/v1beta2"
	"github.com/moondev/cluster-api-provider-maas/pkg/maas/scope"
	clusterv1beta2 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

// Service manages the MaaS machine using the Spectro MaaS client (v0.1.2-beta1+).
type Service struct {
	scope      *scope.MachineScope
	maasClient maasclient.ClientSetInterface
}

// NewService returns a new machine service for the given scope using the Spectro MaaS client.
func NewService(machineScope *scope.MachineScope) *Service {
	return &Service{
		scope:      machineScope,
		maasClient: scope.NewSpectroMaasClient(machineScope.ClusterScope),
	}
}

func (s *Service) GetMachine(systemID string) (*infrav1beta2.Machine, error) {
	if systemID == "" {
		return nil, nil
	}

	ctx := context.Background()
	m, err := s.maasClient.Machines().Machine(systemID).Get(ctx)
	if err != nil {
		return nil, err
	}

	return fromMaasClientMachine(m), nil
}

func (s *Service) ReleaseMachine(systemID string) error {
	ctx := context.Background()
	_, err := s.maasClient.Machines().Machine(systemID).Releaser().Release(ctx)
	if err != nil {
		return errors.Wrapf(err, "Unable to release machine")
	}
	return nil
}

func (s *Service) DeployMachine(userDataB64 string) (_ *infrav1beta2.Machine, rerr error) {
	log := textlogger.NewLogger(textlogger.NewConfig())
	mm := s.scope.MaasMachine

	failureDomain := mm.Spec.FailureDomain
	if failureDomain == nil && s.scope.Machine.Spec.FailureDomain != "" {
		fd := s.scope.Machine.Spec.FailureDomain
		failureDomain = &fd
	}

	ctx := context.Background()
	var deployed maasclient.Machine

	if s.scope.GetProviderID() == "" {
		// Allocate a new machine
		alloc := s.maasClient.Machines().Allocator().
			WithCPUCount(*mm.Spec.MinCPU).
			WithMemory(*mm.Spec.MinMemoryInMB)
		if failureDomain != nil {
			alloc = alloc.WithZone(*failureDomain)
		}
		if mm.Spec.ResourcePool != nil {
			alloc = alloc.WithResourcePool(*mm.Spec.ResourcePool)
		}
		if len(mm.Spec.Tags) > 0 {
			alloc = alloc.WithTags(mm.Spec.Tags)
		}

		m, err := alloc.Allocate(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to allocate machine")
		}

		s.scope.SetProviderID(m.SystemID(), m.Zone().Name())
		if err := s.scope.PatchObject(); err != nil {
			return nil, errors.Wrapf(err, "Unable to patch object")
		}
	}

	// Deploy the machine (existing or just allocated)
	systemID := s.scope.GetProviderID()
	deployer := s.maasClient.Machines().Machine(systemID).Deployer().
		SetUserData(userDataB64).
		SetDistroSeries(mm.Spec.Image).
		SetEphemeralDeploy(mm.Spec.Ephemeral)

	deployed, err := deployer.Deploy(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to deploy machine")
	}

	log.Info("Machine deployed", "systemID", deployed.SystemID(), "hostname", deployed.Hostname())

	return fromMaasClientMachine(deployed), nil
}

// fromMaasClientMachine builds the provider's Machine from the Spectro client Machine interface.
func fromMaasClientMachine(m maasclient.Machine) *infrav1beta2.Machine {
	machine := &infrav1beta2.Machine{
		ID:               m.SystemID(),
		Hostname:         m.Hostname(),
		State:            infrav1beta2.MachineState(m.State()),
		Powered:          m.PowerState() == "on",
		AvailabilityZone: "",
	}
	if m.Zone() != nil {
		machine.AvailabilityZone = m.Zone().Name()
	}
	for _, ip := range m.IPAddresses() {
		if len(ip) > 0 && ip.String() != "" {
			machine.Addresses = append(machine.Addresses, clusterv1beta2.MachineAddress{
				Type:    clusterv1beta2.MachineExternalIP,
				Address: ip.String(),
			})
		}
	}
	return machine
}

func (s *Service) PowerOnMachine() error {
	// Power-on is typically part of deploy in MaaS; no separate call needed.
	return nil
}
