# Ephemeral Deployment for MAAS Cluster API Provider

This document explains how to use ephemeral deployment with the MAAS Cluster API Provider.

## What is Ephemeral Deployment?

Ephemeral deployment deploys machines in memory instead of disk. This means:
- Machines run entirely in RAM
- No persistent storage is used
- Faster deployment and startup times
- Machines are automatically cleaned up when stopped
- Useful for testing, development, or temporary workloads

## Usage

### 1. Using the Ephemeral Template

The `cluster-template-ephemeral.yaml` template demonstrates how to enable ephemeral deployment for both control plane and worker nodes:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MaasMachineTemplate
metadata:
  name: "${CLUSTER_NAME}-control-plane"
spec:
  template:
    spec:
      # ... other configuration ...
      ephemeral: true  # Enable ephemeral deployment
```

### 2. Deploying with clusterctl

```bash
# Generate cluster configuration using the ephemeral template
clusterctl generate cluster my-ephemeral-cluster \
  --from templates/cluster-template-ephemeral.yaml \
  --kubernetes-version v1.28.0 \
  --control-plane-machine-count 1 \
  --worker-machine-count 1 \
  --target-namespace default \
  --output-dir ./my-ephemeral-cluster

# Apply the configuration
kubectl apply -f ./my-ephemeral-cluster/
```

### 3. Environment Variables

Set the following environment variables before generating the cluster:

```bash
export CLUSTER_NAME="my-ephemeral-cluster"
export MAAS_DNS_DOMAIN="maas.local"
export KUBERNETES_VERSION="v1.28.0"
export CONTROL_PLANE_MACHINE_COUNT="1"
export CONTROL_PLANE_MACHINE_MINCPU="2"
export CONTROL_PLANE_MACHINE_MINMEMORY="4096"
export CONTROL_PLANE_MACHINE_IMAGE="ubuntu/jammy"
export CONTROL_PLANE_MACHINE_RESOURCEPOOL=""
export CONTROL_PLANE_MACHINE_TAG="control-plane"
export WORKER_MACHINE_MINCPU="2"
export WORKER_MACHINE_MINMEMORY="4096"
export WORKER_MACHINE_IMAGE="ubuntu/jammy"
export WORKER_MACHINE_RESOURCEPOOL=""
export WORKER_MACHINE_TAG="worker"
```

## Configuration Options

### Ephemeral Field

The `ephemeral` field in `MaasMachineSpec` controls whether a machine is deployed in ephemeral mode:

```yaml
spec:
  ephemeral: true  # Deploy in memory
  ephemeral: false # Deploy to disk (default)
  # ephemeral: omit  # Deploy to disk (default)
```

### Use Cases

1. **Development and Testing**: Quick cluster setup for development work
2. **CI/CD Pipelines**: Temporary clusters for testing
3. **Demo Environments**: Fast deployment for demonstrations
4. **Resource Optimization**: When disk space is limited

### Limitations

- **No Persistence**: All data is lost when machines are stopped
- **Memory Requirements**: Machines need sufficient RAM to run entirely in memory
- **Performance**: May have different performance characteristics than disk-based deployment
- **Compatibility**: Not all workloads are suitable for ephemeral deployment

## Examples

### Control Plane Only (Ephemeral)

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MaasMachineTemplate
metadata:
  name: "${CLUSTER_NAME}-control-plane"
spec:
  template:
    spec:
      minCPU: 2
      minMemory: 4096
      image: "ubuntu/jammy"
      ephemeral: true
```

### Mixed Deployment

You can mix ephemeral and non-ephemeral machines in the same cluster:

```yaml
# Ephemeral control plane
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MaasMachineTemplate
metadata:
  name: "${CLUSTER_NAME}-control-plane"
spec:
  template:
    spec:
      ephemeral: true

---
# Non-ephemeral workers
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MaasMachineTemplate
metadata:
  name: "${CLUSTER_NAME}-md-0"
spec:
  template:
    spec:
      ephemeral: false  # or omit this field
```

## Troubleshooting

### Common Issues

1. **Insufficient Memory**: Ensure machines have enough RAM for ephemeral deployment
2. **MAAS Configuration**: Verify that your MAAS instance supports ephemeral deployment
3. **Image Compatibility**: Some images may not work well with ephemeral deployment

### Verification

To verify that machines are deployed in ephemeral mode:

1. Check the MAAS web interface for machine deployment status
2. Look for ephemeral deployment indicators in MAAS logs
3. Verify that machines are running in memory (no disk usage for OS)

## Related Documentation

- [MAAS Ephemeral Deployment](https://maas.io/docs/advanced-deployment)
- [Cluster API Documentation](https://cluster-api.sigs.k8s.io/)
- [MAAS Cluster API Provider](https://github.com/moondev/cluster-api-provider-maas) 