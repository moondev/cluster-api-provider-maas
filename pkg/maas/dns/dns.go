package dns

import (
	"context"

	"github.com/pkg/errors"

	maasclient "github.com/spectrocloud/maas-client-go/maasclient"

	infrav1beta2 "github.com/moondev/cluster-api-provider-maas/api/v1beta2"
	"github.com/moondev/cluster-api-provider-maas/pkg/maas/scope"
	"k8s.io/apimachinery/pkg/util/sets"
)

type Service struct {
	scope      *scope.ClusterScope
	maasClient maasclient.ClientSetInterface
}

var ErrNotFound = errors.New("resource not found")

// NewService returns a new helper for managing MaaS DNS resources (e.g. API server load balancing).
func NewService(clusterScope *scope.ClusterScope) *Service {
	return &Service{
		scope:      clusterScope,
		maasClient: scope.NewSpectroMaasClient(clusterScope),
	}
}

// ReconcileDNS reconciles the DNS resource for the given cluster.
func (s *Service) ReconcileDNS() error {
	s.scope.V(2).Info("Reconciling DNS")
	ctx := context.TODO()

	dnsResource, err := s.GetDNSResource()
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	dnsName := s.scope.GetDNSName()

	if dnsResource == nil {
		_, err = s.maasClient.DNSResources().
			Builder().
			WithFQDN(s.scope.GetDNSName()).
			WithAddressTTL("10").
			WithIPAddresses(nil).
			Create(ctx)
		if err != nil {
			return errors.Wrapf(err, "Unable to create DNS Resources")
		}
	}

	s.scope.SetDNSName(dnsName)

	return nil
}

// UpdateDNSAttachments updates the DNS resource with the given IPs.
func (s *Service) UpdateDNSAttachments(IPs []string) error {
	s.scope.V(2).Info("Updating DNS Attachments")
	ctx := context.TODO()

	dnsResource, err := s.GetDNSResource()
	if err != nil {
		return err
	}

	_, err = dnsResource.Modifier().SetIPAddresses(IPs).Modify(ctx)
	if err != nil {
		return errors.Wrap(err, "Unable to update IPs")
	}

	return nil
}

// MachineIsRegisteredWithAPIServerDNS returns true if the machine's addresses are in the API server DNS records.
func (s *Service) MachineIsRegisteredWithAPIServerDNS(m *infrav1beta2.Machine) (bool, error) {
	ips, err := s.GetAPIServerDNSRecords()
	if err != nil {
		return false, err
	}

	for _, mAddr := range m.Addresses {
		if ips.Has(mAddr.Address) {
			return true, nil
		}
	}

	return false, nil
}

// GetAPIServerDNSRecords returns the set of IP addresses attached to the cluster's API server DNS resource.
func (s *Service) GetAPIServerDNSRecords() (sets.String, error) {
	dnsResource, err := s.GetDNSResource()
	if err != nil {
		return nil, err
	}

	ips := sets.NewString()
	for _, addr := range dnsResource.IPAddresses() {
		if addr.IP() != nil && addr.IP().String() != "" {
			ips.Insert(addr.IP().String())
		}
	}

	return ips, nil
}

// GetDNSResource returns the DNS resource for the cluster's DNS name, or ErrNotFound if none exists.
func (s *Service) GetDNSResource() (maasclient.DNSResource, error) {
	dnsName := s.scope.GetDNSName()
	if dnsName == "" {
		return nil, errors.New("No DNS on the cluster set!")
	}

	params := maasclient.ParamsBuilder().Set(maasclient.FQDNKey, dnsName)
	d, err := s.maasClient.DNSResources().List(context.Background(), params)
	if err != nil {
		return nil, errors.Wrapf(err, "error retrieving dns resources %q", dnsName)
	}
	if len(d) > 1 {
		return nil, errors.Errorf("expected 1 DNS Resource for %q, got %d", dnsName, len(d))
	}
	if len(d) == 0 {
		return nil, ErrNotFound
	}

	return d[0], nil
}
