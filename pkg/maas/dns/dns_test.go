package dns

import (
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2/klogr"
	clusterv1beta2 "sigs.k8s.io/cluster-api/api/core/v1beta2"

	maasclient "github.com/spectrocloud/maas-client-go/maasclient"
	infrav1beta2 "github.com/moondev/cluster-api-provider-maas/api/v1beta2"
	mockclientset "github.com/moondev/cluster-api-provider-maas/pkg/maas/client/mock"
	"github.com/moondev/cluster-api-provider-maas/pkg/maas/scope"
)

func TestDNS(t *testing.T) {
	log := klogr.New()
	cluster := &clusterv1beta2.Cluster{
		ObjectMeta: v1.ObjectMeta{
			Name: "a",
		},
	}
	maasCluster := &infrav1beta2.MaasCluster{
		Spec: infrav1beta2.MaasClusterSpec{
			DNSDomain: "b.com",
		},
	}

	t.Run("reconcile dns", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		mockClientSetInterface := mockclientset.NewMockClientSetInterface(ctrl)
		mockDNSResources := mockclientset.NewMockDNSResources(ctrl)
		mockDNSResourceBuilder := mockclientset.NewMockDNSResourceBuilder(ctrl)
		s := &Service{
			scope: &scope.ClusterScope{
				Logger:      log,
				Cluster:     cluster,
				MaasCluster: maasCluster,
			},
			maasClient: mockClientSetInterface,
		}
		mockClientSetInterface.EXPECT().DNSResources().Return(mockDNSResources)
		mockDNSResources.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
		mockClientSetInterface.EXPECT().DNSResources().Return(mockDNSResources)
		mockDNSResources.EXPECT().Builder().Return(mockDNSResourceBuilder)
		mockDNSResourceBuilder.EXPECT().WithFQDN(gomock.Any()).Return(mockDNSResourceBuilder)
		mockDNSResourceBuilder.EXPECT().WithAddressTTL("10").Return(mockDNSResourceBuilder)
		mockDNSResourceBuilder.EXPECT().WithIPAddresses(nil).Return(mockDNSResourceBuilder)
		mockDNSResourceBuilder.EXPECT().Create(gomock.Any())
		err := s.ReconcileDNS()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(s.scope.GetDNSName()).To(ContainSubstring(cluster.Name))
		g.Expect(s.scope.GetDNSName()).To(ContainSubstring(maasCluster.Spec.DNSDomain))
	})

	t.Run("update dns attachment", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		mockClientSetInterface := mockclientset.NewMockClientSetInterface(ctrl)
		mockDNSResources := mockclientset.NewMockDNSResources(ctrl)
		mockDNSResource := mockclientset.NewMockDNSResource(ctrl)
		mockDNSResourceModifier := mockclientset.NewMockDNSResourceModifier(ctrl)
		s := &Service{
			scope: &scope.ClusterScope{
				Logger:      log,
				Cluster:     cluster,
				MaasCluster: maasCluster,
			},
			maasClient: mockClientSetInterface,
		}

		mockClientSetInterface.EXPECT().DNSResources().Return(mockDNSResources)
		mockDNSResources.EXPECT().List(gomock.Any(), gomock.Any()).Return([]maasclient.DNSResource{mockDNSResource}, nil)
		mockDNSResource.EXPECT().Modifier().Return(mockDNSResourceModifier)
		mockDNSResourceModifier.EXPECT().SetIPAddresses([]string{"1.1.1.1", "8.8.8.8"}).Return(mockDNSResourceModifier)
		mockDNSResourceModifier.EXPECT().Modify(gomock.Any()).Return(mockDNSResource, nil)

		err := s.UpdateDNSAttachments([]string{"1.1.1.1", "8.8.8.8"})

		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("machine is registered", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		mockClientSetInterface := mockclientset.NewMockClientSetInterface(ctrl)
		mockDNSResources := mockclientset.NewMockDNSResources(ctrl)
		mockDNSResource := mockclientset.NewMockDNSResource(ctrl)
		mockIP1 := mockclientset.NewMockIPAddress(ctrl)
		mockIP2 := mockclientset.NewMockIPAddress(ctrl)
		s := &Service{
			scope: &scope.ClusterScope{
				Logger:      log,
				Cluster:     cluster,
				MaasCluster: maasCluster,
			},
			maasClient: mockClientSetInterface,
		}
		mockClientSetInterface.EXPECT().DNSResources().Return(mockDNSResources)
		mockDNSResources.EXPECT().List(gomock.Any(), gomock.Any()).Return([]maasclient.DNSResource{mockDNSResource}, nil)
		mockDNSResource.EXPECT().IPAddresses().Return([]maasclient.IPAddress{mockIP1, mockIP2})
		mockIP1.EXPECT().IP().Return(net.ParseIP("1.1.1.1")).AnyTimes()
		mockIP2.EXPECT().IP().Return(net.ParseIP("8.8.8.8")).AnyTimes()

		res, err := s.MachineIsRegisteredWithAPIServerDNS(&infrav1beta2.Machine{
			Addresses: []clusterv1beta2.MachineAddress{
				{
					Type:    clusterv1beta2.MachineInternalIP,
					Address: "1.1.1.1",
				},
				{
					Type:    clusterv1beta2.MachineInternalIP,
					Address: "8.8.8.8",
				},
			},
		})

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(res).To(BeTrue())
	})
}
