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

package controllers

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	clusterv1beta2 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/conditions/deprecated/v1beta1"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1beta1 "github.com/moondev/cluster-api-provider-maas/api/v1beta1"
	"github.com/moondev/cluster-api-provider-maas/pkg/maas/dns"
	"github.com/moondev/cluster-api-provider-maas/pkg/maas/scope"
	infrautil "github.com/moondev/cluster-api-provider-maas/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

const (
	apiServerServiceNameAnnotation = "capmaas.spectrocloud.io/apiserver-service-name"
	apiServerIngressNameAnnotation = "capmaas.spectrocloud.io/apiserver-ingress-name"
	enableExternalControlPlaneAnno = "capmaas.spectrocloud.io/enable-external-control-plane-endpoint"
)

// MaasClusterReconciler reconciles a MaasCluster object
type MaasClusterReconciler struct {
	client.Client
	Log                 logr.Logger
	Scheme              *runtime.Scheme
	Recorder            record.EventRecorder
	GenericEventChannel chan event.GenericEvent
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=maasclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=maasclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups="",resources=services;services/status,verbs=get;list;watch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses;ingresses/status,verbs=get;list;watch

// Reconcile reads that state of the cluster for a MaasCluster object and makes changes based on the state read
// and what is in the MaasCluster.Spec
func (r *MaasClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, rerr error) {
	log := r.Log.WithValues("maascluster", req.Name)

	// Fetch the MaasCluster instance
	maasCluster := &infrav1beta1.MaasCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, maasCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, maasCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Waiting for Cluster Controller to set OwnerRef on MaasCluster")
		return ctrl.Result{}, nil
	}

	// Create the scope.
	clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:              r.Client,
		Logger:              log,
		Cluster:             cluster,
		MaasCluster:         maasCluster,
		ClusterEventChannel: r.GenericEventChannel,
		ControllerName:      "maascluster",
	})
	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	// Always close the scope when exiting this function so we can persist any MAAS Cluster changes.
	defer func() {
		if err := clusterScope.Close(); err != nil && rerr == nil {
			rerr = err
		}
	}()

	// Support FailureDomains
	// In cloud providers this would likely look up which failure domains are supported and set the status appropriately.
	// so kCP will distribute the CPs across multiple failure domains
	failureDomains := make(clusterv1beta1.FailureDomains)
	for _, az := range maasCluster.Spec.FailureDomains {
		failureDomains[az] = clusterv1beta1.FailureDomainSpec{
			ControlPlane: true,
		}
	}
	maasCluster.Status.FailureDomains = failureDomains

	// Handle deleted clusters
	if !maasCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, clusterScope)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, clusterScope)
}

func (r *MaasClusterReconciler) reconcileDelete(ctx context.Context, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	clusterScope.Info("Reconciling MaasCluster delete")

	maasCluster := clusterScope.MaasCluster

	maasMachines, err := infrautil.GetMAASMachinesInCluster(ctx, r.Client, clusterScope.Cluster.Namespace, clusterScope.Cluster.Name)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err,
			"unable to list MAASMachines part of MAASCluster %s/%s", clusterScope.Cluster.Namespace, clusterScope.Cluster.Name)
	}

	if len(maasMachines) > 0 {
		r.Log.Info("Waiting for MAASMachines to be deleted", "count", len(maasMachines))
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(maasCluster, infrav1beta1.ClusterFinalizer)

	// TODO(saamalik) implement the recorder stuff (look at aws)

	return reconcile.Result{}, nil
}

func (r *MaasClusterReconciler) reconcileDNSAttachments(clusterScope *scope.ClusterScope, dnssvc *dns.Service) error {
	machines, err := clusterScope.GetClusterMaasMachines()
	if err != nil {
		return errors.Wrapf(err, "Unable to list all maas machines")
	}

	var runningIpAddresses []string

	currentIPs, err := dnssvc.GetAPIServerDNSRecords()
	if err != nil {
		return errors.Wrap(err, "Unable to get the dns resources")
	}

	machinesPendingAttachment := make([]*infrav1beta1.MaasMachine, 0)
	machinesPendingDetachment := make([]*infrav1beta1.MaasMachine, 0)

	for _, m := range machines {
		if !IsControlPlaneMachine(m) {
			continue
		}

		machineIP := getExternalMachineIP(m)
		attached := currentIPs.Has(machineIP)
		isRunningHealthy := IsRunning(m)

		if !m.DeletionTimestamp.IsZero() || !isRunningHealthy {
			if attached {
				clusterScope.Info("Cleaning up IP on unhealthy machine", "machine", m.Name)
				machinesPendingDetachment = append(machinesPendingDetachment, m)
			}
		} else if IsRunning(m) {
			if !attached {
				clusterScope.Info("Healthy machine without DNS attachment; attaching.", "machine", m.Name)
				machinesPendingAttachment = append(machinesPendingAttachment, m)
			}

			runningIpAddresses = append(runningIpAddresses, machineIP)
		}
		//r.Recorder.Eventf(machineScope.MaasMachine, corev1.EventTypeNormal, "SuccessfulDetachControlPlaneDNS",
		//	"Control plane instance %q is de-registered from load balancer", i.ID)
		//runningIpAddresses = append(runningIpAddresses, m.)
	}

	if err := dnssvc.UpdateDNSAttachments(runningIpAddresses); err != nil {
		return err
	} else if len(machinesPendingAttachment) > 0 || len(machinesPendingDetachment) > 0 {
		clusterScope.Info("Pending DNS attachments or detachments; will retry again")
		return ErrRequeueDNS
	}

	return nil
}

func getAnnotation(cluster *clusterv1.Cluster, maasCluster *infrav1beta1.MaasCluster, key string) string {
	if cluster != nil {
		if v, ok := cluster.Annotations[key]; ok {
			return v
		}
	}
	if maasCluster != nil {
		if v, ok := maasCluster.Annotations[key]; ok {
			return v
		}
	}
	return ""
}

func isExternalEndpointEnabled(cluster *clusterv1.Cluster, maasCluster *infrav1beta1.MaasCluster) bool {
	v := getAnnotation(cluster, maasCluster, enableExternalControlPlaneAnno)
	return v == "true" || v == "1" || v == "yes"
}

func resolveServiceEndpoint(ctx context.Context, c client.Client, namespace, name string) (string, int, bool, error) {
	if name == "" {
		return "", 0, false, nil
	}
	var svc corev1.Service
	if err := c.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &svc); err != nil {
		return "", 0, false, err
	}
	if svc.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return "", 0, false, nil
	}
	if len(svc.Status.LoadBalancer.Ingress) == 0 {
		return "", 0, false, nil
	}
	addr := svc.Status.LoadBalancer.Ingress[0]
	host := addr.Hostname
	if host == "" {
		host = addr.IP
	}
	if host == "" {
		return "", 0, false, nil
	}
	port := 0
	for _, p := range svc.Spec.Ports {
		if p.Name == "https" || p.Port == 6443 {
			port = int(p.Port)
			break
		}
	}
	if port == 0 && len(svc.Spec.Ports) > 0 {
		port = int(svc.Spec.Ports[0].Port)
	}
	if port == 0 {
		port = 6443
	}
	return host, port, true, nil
}

func resolveIngressEndpoint(ctx context.Context, c client.Client, namespace, name string) (string, int, bool, error) {
	if name == "" {
		return "", 0, false, nil
	}
	var ing networkingv1.Ingress
	if err := c.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &ing); err != nil {
		return "", 0, false, err
	}
	host := ""
	if len(ing.Status.LoadBalancer.Ingress) > 0 {
		addr := ing.Status.LoadBalancer.Ingress[0]
		host = addr.Hostname
		if host == "" {
			host = addr.IP
		}
	}
	if host == "" && len(ing.Spec.Rules) > 0 {
		host = ing.Spec.Rules[0].Host
	}
	if host == "" {
		return "", 0, false, nil
	}
	port := 443
	if len(ing.Spec.TLS) == 0 {
		port = 80
	}
	return host, port, true, nil
}

// IsControlPlaneMachine checks machine is a control plane node.
func IsControlPlaneMachine(m *infrav1beta1.MaasMachine) bool {
	_, ok := m.ObjectMeta.Labels[clusterv1beta1.MachineControlPlaneLabel]
	return ok
}

// IsRunning returns if the machine is running
func IsRunning(m *infrav1beta1.MaasMachine) bool {
	if !m.Status.MachinePowered {
		return false
	}

	state := m.Status.MachineState
	return state != nil && infrav1beta1.MachineRunningStates.Has(string(*state))
}

func getExternalMachineIP(machine *infrav1beta1.MaasMachine) string {
	for _, i := range machine.Status.Addresses {
		if i.Type == clusterv1beta1.MachineExternalIP {
			return i.Address
		}
	}
	return ""
}

func (r *MaasClusterReconciler) reconcileNormal(_ context.Context, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	clusterScope.Info("Reconciling MaasCluster")

	maasCluster := clusterScope.MaasCluster
	cluster := clusterScope.Cluster

	// Add finalizer first if not exist to avoid the race condition between init and delete
	if !controllerutil.ContainsFinalizer(maasCluster, infrav1beta1.ClusterFinalizer) {
		controllerutil.AddFinalizer(maasCluster, infrav1beta1.ClusterFinalizer)
		return ctrl.Result{}, nil
	}

	// If Cluster.ControlPlaneEndpoint is already set, honor it first.
	if cluster.Spec.ControlPlaneEndpoint.IsValid() {
		maasCluster.Spec.ControlPlaneEndpoint = infrav1beta1.APIEndpoint{
			Host: cluster.Spec.ControlPlaneEndpoint.Host,
			Port: int(cluster.Spec.ControlPlaneEndpoint.Port),
		}
		maasCluster.Status.Network.DNSName = cluster.Spec.ControlPlaneEndpoint.Host
		v1beta1conditions.MarkTrue(maasCluster, infrav1beta1.DNSReadyCondition)
		return ctrl.Result{}, nil
	}

	// Service LoadBalancer endpoint (optional, via annotation)
	if svcName := getAnnotation(cluster, maasCluster, apiServerServiceNameAnnotation); svcName != "" {
		host, port, ok, err := resolveServiceEndpoint(context.TODO(), r.Client, maasCluster.Namespace, svcName)
		if err != nil {
			return ctrl.Result{}, err
		}
		if ok {
			maasCluster.Spec.ControlPlaneEndpoint = infrav1beta1.APIEndpoint{Host: host, Port: port}
			maasCluster.Status.Network.DNSName = host
			v1beta1conditions.MarkTrue(maasCluster, infrav1beta1.DNSReadyCondition)
			return ctrl.Result{}, nil
		}
	}

	// Ingress endpoint (optional, via annotation)
	if ingName := getAnnotation(cluster, maasCluster, apiServerIngressNameAnnotation); ingName != "" {
		host, port, ok, err := resolveIngressEndpoint(context.TODO(), r.Client, maasCluster.Namespace, ingName)
		if err != nil {
			return ctrl.Result{}, err
		}
		if ok {
			maasCluster.Spec.ControlPlaneEndpoint = infrav1beta1.APIEndpoint{Host: host, Port: port}
			maasCluster.Status.Network.DNSName = host
			v1beta1conditions.MarkTrue(maasCluster, infrav1beta1.DNSReadyCondition)
			return ctrl.Result{}, nil
		}
	}

	// For non-MaaS control planes, skip MaaS DNS unless explicitly enabled.
	if cluster.Spec.InfrastructureRef.Kind != "MaasCluster" && !isExternalEndpointEnabled(cluster, maasCluster) {
		return ctrl.Result{}, nil
	}

	dnsService := dns.NewService(clusterScope)

	if err := dnsService.ReconcileDNS(); err != nil {
		clusterScope.Error(err, "failed to reconcile load balancer")
		v1beta1conditions.MarkFalse(maasCluster, infrav1beta1.DNSReadyCondition, infrav1beta1.DNSFailedReason, clusterv1beta2.ConditionSeverityError, "%v", err)
		return reconcile.Result{}, err
	}

	if maasCluster.Status.Network.DNSName == "" {
		v1beta1conditions.MarkFalse(maasCluster, infrav1beta1.DNSReadyCondition, infrav1beta1.WaitForDNSNameReason, clusterv1beta2.ConditionSeverityInfo, "")
		clusterScope.Info("Waiting on API server DNS name")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}

	maasCluster.Spec.ControlPlaneEndpoint = infrav1beta1.APIEndpoint{
		Host: maasCluster.Status.Network.DNSName,
		Port: clusterScope.APIServerPort(),
	}

	maasCluster.Status.Ready = true

	// Mark the maasCluster ready
	v1beta1conditions.MarkTrue(maasCluster, infrav1beta1.DNSReadyCondition)

	if err := r.reconcileDNSAttachments(clusterScope, dnsService); err != nil {
		if errors.Is(err, ErrRequeueDNS) {
			return ctrl.Result{}, nil
			//return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		}

		clusterScope.Error(err, "failed to reconcile load balancer")
		return reconcile.Result{}, err

	}

	clusterScope.ReconcileMaasClusterWhenAPIServerIsOnline()
	if k, _ := clusterScope.IsAPIServerOnline(); !k {
		v1beta1conditions.MarkFalse(maasCluster, infrav1beta1.APIServerAvailableCondition, infrav1beta1.APIServerNotReadyReason, clusterv1beta2.ConditionSeverityWarning, "")
		return ctrl.Result{}, nil
	}

	v1beta1conditions.MarkTrue(maasCluster, infrav1beta1.APIServerAvailableCondition)
	clusterScope.Info("API Server is available")

	return ctrl.Result{}, nil
}

// SetupWithManager will add watches for this controller
func (r *MaasClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	recover := true
	if r.GenericEventChannel == nil {
		r.GenericEventChannel = make(chan event.GenericEvent)
	}

	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&infrav1beta1.MaasCluster{}).
		WithOptions(controller.Options{
			RecoverPanic: &recover,
		}).
		Watches(
			&infrav1beta1.MaasMachine{},
			handler.EnqueueRequestsFromMapFunc(r.controlPlaneMachineToCluster),
		).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(
				util.ClusterToInfrastructureMapFunc(context.Background(), infrav1beta1.GroupVersion.WithKind("MaasCluster"), mgr.GetClient(), &infrav1beta1.MaasCluster{}),
			),
			//predicates.ClusterUnpaused(mgr.GetScheme(), r.Log),
		).
		WithEventFilter(predicates.ResourceNotPaused(mgr.GetScheme(), r.Log)).
		Build(r)
	if err != nil {
		return err
	}

	if err := c.Watch(
		source.Channel(r.GenericEventChannel, &handler.EnqueueRequestForObject{}),
	); err != nil {
		return err
	}

	return err

}

// controlPlaneMachineToCluster is a handler.ToRequestsFunc to be used
// to enqueue requests for reconciliation for MaasCluster to update
// its status.apiEndpoints field.
func (r *MaasClusterReconciler) controlPlaneMachineToCluster(_ context.Context, o client.Object) []ctrl.Request {
	maasMachine, ok := o.(*infrav1beta1.MaasMachine)
	if !ok {
		r.Log.Error(nil, fmt.Sprintf("expected a MaasMachine but got a %T", o))
		return nil
	}
	if !IsControlPlaneMachine(maasMachine) {
		return nil
	}

	ctx := context.TODO()

	// Fetch the CAPI Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, maasMachine.ObjectMeta)
	if err != nil {
		r.Log.Error(err, "MaasMachine is missing cluster label or cluster does not exist",
			"namespace", maasMachine.Namespace, "name", maasMachine.Name)
		return nil
	}

	// Fetch the MaasCluster
	maasCluster := &infrav1beta1.MaasCluster{}
	maasClusterKey := client.ObjectKey{
		Namespace: maasMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, maasClusterKey, maasCluster); err != nil {
		r.Log.Error(err, "failed to get MaasCluster",
			"namespace", maasClusterKey.Namespace, "name", maasClusterKey.Name)
		return nil
	}

	return []ctrl.Request{{
		NamespacedName: types.NamespacedName{
			Namespace: maasClusterKey.Namespace,
			Name:      maasClusterKey.Name,
		},
	}}
}
