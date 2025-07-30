```yaml
providers:
- name: maas
  type: InfrastructureProvider
  url: https://github.com/moondev/cluster-api-provider-maas/releases/latest/infrastructure-components.yaml

```
init
```shell
CLUSTER_TOPOLOGY=true MAAS_API_KEY=key MAAS_ENDPOINT=http://:5240/MAAS clusterctl init --infrastructure maas --addon helm -v=9
kubectl set image deployments.apps -n capmaas-system capmaas-controller-manager manager=docker.io/chadmoon/cluster-api-provider-maas-controller:v0.6.1-clusterclass

```

clusterclass
```yaml

---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MaasClusterTemplate
metadata:
  name: maas-cluster-template
  namespace: default
spec:
  template:
    spec:
      dnsDomain: maas


---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MaasMachineTemplate
metadata:
  name: mgmt-control-plane
  namespace: default
spec:
  template:
    spec:
      image: custom/u-2204-0-k-1316-0
      minCPU: 1
      minMemory: 2192
      resourcePool: default
      tags: []


---
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
kind: KubeadmControlPlaneTemplate
metadata:
  name: maas-control-plane-template
  namespace: default
spec:
  template:
    spec:
      kubeadmConfigSpec:
        clusterConfiguration:
          apiServer:
            extraArgs:
              anonymous-auth: "true"
              authorization-mode: RBAC,Node
              default-not-ready-toleration-seconds: "60"
              default-unreachable-toleration-seconds: "60"
              disable-admission-plugins: AlwaysAdmit
              enable-admission-plugins: AlwaysPullImages,NamespaceLifecycle,ServiceAccount,NodeRestriction
            timeoutForControlPlane: 10m0s
          controllerManager:
            extraArgs:
              feature-gates: RotateKubeletServerCertificate=true
              terminated-pod-gc-threshold: "25"
              use-service-account-credentials: "true"
          dns: {}
          etcd: {}
          networking: {}
          scheduler:
            extraArgs: null
        initConfiguration:
          localAPIEndpoint:
            advertiseAddress: ""
            bindPort: 0
          nodeRegistration:
            kubeletExtraArgs:
              event-qps: "0"
              feature-gates: RotateKubeletServerCertificate=true
              read-only-port: "0"
            taints: []
            name: '{{ v1.local_hostname }}'
        joinConfiguration:
          controlPlane:
            localAPIEndpoint:
              advertiseAddress: ""
              bindPort: 0
          discovery: {}
          nodeRegistration:
            kubeletExtraArgs:
              event-qps: "0"
              feature-gates: RotateKubeletServerCertificate=true
              read-only-port: "0"
            name: '{{ v1.local_hostname }}'
            taints: []
        preKubeadmCommands:
        - while [ ! -S /var/run/containerd/containerd.sock ]; do echo 'Waiting for containerd...';
          sleep 1; done
        - sed -ri '/\sswap\s/s/^#?/#/' /etc/fstab
        - swapoff -a
        useExperimentalRetryJoin: true

---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MaasMachineTemplate
metadata:
  name: mgmt-md-0
  namespace: default
spec:
  template:
    spec:
      image: custom/u-2204-0-k-1316-0
      minCPU: 1
      minMemory: 2192
      resourcePool: default
      tags: []

---
apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
kind: KubeadmConfigTemplate
metadata:
  name: mgmt-md-0
  namespace: default
spec:
  template:
    spec:
      joinConfiguration:
        nodeRegistration:
          kubeletExtraArgs:
            event-qps: "0"
            feature-gates: RotateKubeletServerCertificate=true
            read-only-port: "0"
          name: '{{ v1.local_hostname }}'
      preKubeadmCommands:
      - while [ ! -S /var/run/containerd/containerd.sock ]; do echo 'Waiting for containerd...';
        sleep 1; done
      - sed -ri '/\sswap\s/s/^#?/#/' /etc/fstab
      - swapoff -a
      useExperimentalRetryJoin: true


---
apiVersion: cluster.x-k8s.io/v1beta1
kind: ClusterClass
metadata:
  name: maas-clusterclass
spec:
  infrastructure:
    ref:
      apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
      kind: MaasClusterTemplate
      name: maas-cluster-template
      namespace: default
  controlPlane:
    ref:
      apiVersion: controlplane.cluster.x-k8s.io/v1beta1
      kind: KubeadmControlPlaneTemplate
      name: maas-control-plane-template
      namespace: default
    machineInfrastructure:
      ref:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: MaasMachineTemplate
        name: mgmt-control-plane
  workers:
    machineDeployments:
      - class: default-worker
        template:
          bootstrap:
            ref:
              apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
              kind: KubeadmConfigTemplate
              name: mgmt-md-0
          infrastructure:
            ref:
              apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
              kind: MaasMachineTemplate
              name: mgmt-md-0


---
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: my-cluster
  labels:
    cluster.x-k8s.io/cluster-name: my-cluster
spec:
  topology:
    class: maas-clusterclass
    version: v1.31.6
    controlPlane:
      metadata: {}
      replicas: 1
    workers:
      machineDeployments:
        - class: default-worker
          name: md-0
          replicas: 2
    # variables: {} 

```

---
# Cluster-API-Provider-MAAS
Cluster API Provider for Canonical Metal-As-A-Service [maas.io](https://maas.io/)

You're welcome to join the upcoming [webinar](https://www.spectrocloud.com/webinars/managing-bare-metal-k8s-like-any-other-cluster/) for capmaas!


# Getting Started

## Public Images
Spectro Cloud public images

| Kubernetes Version | URL                                                                        |
|--------------------|----------------------------------------------------------------------------|
| 1.25.6             | https://maas-images-public.s3.amazonaws.com/u-2204-0-k-1256-0.tar.gz       |
| 1.26.1             | https://maas-images-public.s3.amazonaws.com/u-2204-0-k-1261-0.tar.gz       |



## Custom Image Generation
Refer [image-generation/](image-generation/README.md)

## Hello world

- Create kind cluster
```bash
kind create cluster --name=maas-cluster
```

- Install clusterctl v1beta1
https://release-1-1.cluster-api.sigs.k8s.io/user/quick-start.html

- Setup clusterctl configuration `~/.cluster-api/clusterctl.yaml`
```
# MAAS access endpoint and key
MAAS_API_KEY: <maas-api-key>
MAAS_ENDPOINT: http://<maas-endpoint>/MAAS
MAAS_DNS_DOMAIN: maas.domain

# Cluster configuration
KUBERNETES_VERSION: v1.26.4
CONTROL_PLANE_MACHINE_IMAGE: custom/u-2204-0-k-1264-0
CONTROL_PLANE_MACHINE_MINCPU: 4
CONTROL_PLANE_MACHINE_MINMEMORY: 8192
WORKER_MACHINE_IMAGE: custom/u-2204-0-k-1264-0
WORKER_MACHINE_MINCPU: 4
WORKER_MACHINE_MINMEMORY: 8192

# Selecting machine based on resourcepool (optional) and machine tag (optional)
CONTROL_PLANE_MACHINE_RESOURCEPOOL: resorcepool-controller
CONTROL_PLANE_MACHINE_TAG: hello-world
WORKER_MACHINE_RESOURCEPOOL: resourcepool-worker
WORKER_MACHINE_TAG: hello-world
```
- Initialize infrastructure
```bash
clusterctl init --infrastructure maas:v0.5.0
```
- Generate and create cluster
```
clusterctl generate cluster t-cluster --infrastructure=maas:v0.5.0 --kubernetes-version v1.26.4 --control-plane-machine-count=1 --worker-machine-count=3 | kubectl apply -f -
```

## ClusterClass Support (Managed Topology)

This provider supports ClusterClass and managed topology, enabling declarative, reusable cluster definitions.

- See `templates/clusterclass-maas.yaml` for a full example of ClusterClass, Cluster, and template resources for MAAS.
- To use ClusterClass:
  1. Apply the ClusterClass, MaasClusterTemplate, MaasMachineTemplate, and related resources from the example.
  2. Create a Cluster with a `spec.topology` referencing your ClusterClass.
  3. Customize the templates and variables as needed for your environment.

**Note:**
- The provider's webhooks validate required fields for managed topology (e.g., `spec.dnsDomain` for MaasCluster, and `image`, `minCPU`, `minMemory` for MaasMachineTemplate).
- All ClusterClass-managed resources must have the label `topology.cluster.x-k8s.io/owned` set by Cluster API.

## Developer Guide
- Create kind cluster
```shell
kind create cluster --name=maas-cluster
```

- Install clusterctl v1 depending on the version you are working with

- Makefile set IMG=<your docker repo>
- Run 
```shell
make docker-build && make docker-push
```
    
- Generate dev manifests
```shell
make dev-manifests
```

- Move _build/dev/ directory contents to ~/.clusterapi/overrides v0.5.0 depending on version you are working with

```text
.
├── clusterctl.yaml
├── overrides
│   ├── infrastructure-maas
│       └── v0.5.0
│           ├── cluster-template.yaml
│           ├── infrastructure-components.yaml
│           └── metadata.yaml
└── version.yaml

```

- Run
```shell
clusterctl init --infrastructure maas:v0.5.0
```


## Install CRDs
### v1beta1 v0.5.0 release
- Generate cluster using
```shell
clusterctl generate cluster t-cluster  --infrastructure=maas:v0.5.0 | kubectl apply -f -
```
or
```shell
clusterctl generate cluster t-cluster --infrastructure=maas:v0.5.0 --kubernetes-version v1.26.4 > my_cluster.yaml
kubectl apply -f my_cluster.yaml
```
or
```shell
clusterctl generate cluster t-cluster --infrastructure=maas:v0.5.0 --kubernetes-version v1.26.4 --control-plane-machine-count=1 --worker-machine-count=3 > my_cluster.yaml
kubectl apply -f my_cluster.yaml
```

## Ephemeral Deployment Feature

The MAAS provider supports ephemeral deployment, which allows machines to be deployed in-memory instead of on-disk. This feature is useful for temporary workloads, testing, or development environments.

### Configuration

To enable ephemeral deployment, add the `ephemeralDeploy` field to your MaasMachine or MaasMachineTemplate spec:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MaasMachine
metadata:
  name: my-ephemeral-machine
spec:
  minCPU: 2
  minMemory: 4096
  image: "ubuntu/focal"
  ephemeralDeploy: true  # Enable in-memory deployment
```

### Usage Examples

#### Persistent Deployment (Default)
```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MaasMachine
metadata:
  name: persistent-machine
spec:
  minCPU: 2
  minMemory: 4096
  image: "ubuntu/focal"
  ephemeralDeploy: false  # or omit this field
```

#### Ephemeral Deployment
```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MaasMachine
metadata:
  name: ephemeral-machine
spec:
  minCPU: 1
  minMemory: 2048
  image: "ubuntu/focal"
  ephemeralDeploy: true  # Enable in-memory deployment
```

### Benefits of Ephemeral Deployment

- **Faster Deployment**: In-memory deployment is typically faster than on-disk deployment
- **Resource Efficiency**: Ephemeral machines don't consume disk space for the OS image
- **Testing & Development**: Ideal for temporary workloads and development environments
- **Cost Optimization**: Reduced storage costs for temporary workloads

### Considerations

- **Data Persistence**: Ephemeral machines lose all data when powered off or rebooted
- **Performance**: May have different performance characteristics compared to persistent deployments
- **Use Cases**: Best suited for stateless workloads, testing, and development

### Sample Files

See the following sample files for complete examples:
- `config/samples/infrastructure_v1beta1_maasmachine.yaml` - Basic machine with persistent deployment
- `config/samples/infrastructure_v1beta1_maasmachine_ephemeral.yaml` - Machine with ephemeral deployment
- `config/samples/infrastructure_v1beta1_maasmachinetemplate.yaml` - Template with persistent deployment
- `config/samples/infrastructure_v1beta1_maasmachinetemplate_ephemeral.yaml` - Template with ephemeral deployment