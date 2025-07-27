# Ephemeral Deployment Feature

This document describes the ephemeral deployment feature for the MAAS Cluster API Provider, which allows machines to be deployed in memory instead of on disk.

## Overview

Ephemeral deployment is a MAAS feature that enables machines to run entirely in memory without persistent disk storage. This is useful for:

- **Testing and Development**: Quick deployment and teardown of test environments
- **CI/CD Pipelines**: Temporary environments for build and test processes
- **Demo Environments**: Short-lived demonstration clusters
- **Resource Optimization**: Reduced storage requirements for temporary workloads

## Configuration

### Enabling Ephemeral Deployment

To enable ephemeral deployment, set the `ephemeral` field to `true` in your `MaasMachineTemplate`:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: MaasMachineTemplate
metadata:
  name: my-cluster-control-plane
spec:
  template:
    spec:
      minCPU: 2
      minMemory: 4096
      image: ubuntu/focal
      resourcePool: default
      tags: [control-plane]
      # Enable ephemeral deployment
      ephemeral: true
```

### API Specification

The `ephemeral` field has been added to the `MaasMachineSpec`:

```go
type MaasMachineSpec struct {
    // ... existing fields ...
    
    // Ephemeral, if true, deploys the machine in memory instead of disk
    // +optional
    Ephemeral bool `json:"ephemeral,omitempty"`
}
```

## Usage with clusterctl

### Using the Ephemeral Template

1. **Set environment variables** for your MAAS configuration:

```bash
export MAAS_ENDPOINT="http://your-maas-server:5240/MAAS"
export MAAS_API_KEY="your-api-key"
export CLUSTER_NAME="ephemeral-cluster"
export NAMESPACE="default"
export KUBERNETES_VERSION="v1.28.0"
export CONTROL_PLANE_MACHINE_COUNT=1
export WORKER_MACHINE_COUNT=2
export DNS_DOMAIN="example.com"
export CONTROL_PLANE_MACHINE_MINCPU=2
export CONTROL_PLANE_MACHINE_MINMEMORY=4096
export CONTROL_PLANE_MACHINE_IMAGE="ubuntu/focal"
export CONTROL_PLANE_MACHINE_RESOURCEPOOL="default"
export CONTROL_PLANE_MACHINE_TAG="control-plane"
export WORKER_MACHINE_MINCPU=2
export WORKER_MACHINE_MINMEMORY=4096
export WORKER_MACHINE_IMAGE="ubuntu/focal"
export WORKER_MACHINE_RESOURCEPOOL="default"
export WORKER_MACHINE_TAG="worker"
```

2. **Generate the cluster configuration**:

```bash
clusterctl generate cluster $CLUSTER_NAME \
  --from templates/cluster-template-ephemeral.yaml \
  --target-namespace $NAMESPACE > ephemeral-cluster.yaml
```

3. **Apply the configuration**:

```bash
kubectl apply -f ephemeral-cluster.yaml
```

### Customizing Ephemeral Deployment

You can customize which machines use ephemeral deployment by modifying the template:

```yaml
# Control plane with ephemeral deployment
spec:
  template:
    spec:
      ephemeral: true  # Control plane in memory

---
# Worker nodes with disk deployment
spec:
  template:
    spec:
      ephemeral: false  # Workers on disk (default)
```

## Technical Details

### Implementation

The ephemeral deployment feature is implemented in the machine service:

1. **API Level**: The `ephemeral` field is added to `MaasMachineSpec`
2. **Deployment Logic**: When `ephemeral: true`, the deployment parameters include `ephemeral_deploy: true`
3. **MAAS Integration**: Uses the canonical MAAS client's `MachineDeployParams.EphemeralDeploy` field

### Client Migration

The implementation uses a hybrid approach:

- **Machine Service**: Uses the canonical MAAS client (`github.com/canonical/gomaasclient`) for machine operations
- **DNS Service**: Uses the spectrocloud MAAS client for DNS operations (maintains compatibility)
- **Scope**: Provides both client types through separate functions

### Code Example

```go
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
```

## Use Cases

### 1. Development and Testing

```yaml
# Development cluster template
spec:
  template:
    spec:
      ephemeral: true  # Fast deployment/teardown
      minCPU: 1        # Minimal resources
      minMemory: 2048
```

### 2. CI/CD Pipelines

```yaml
# CI/CD pipeline cluster
spec:
  template:
    spec:
      ephemeral: true  # No persistent storage needed
      tags: [ci, pipeline]
```

### 3. Demo Environments

```yaml
# Demo cluster
spec:
  template:
    spec:
      ephemeral: true  # Easy cleanup
      minCPU: 2
      minMemory: 4096
```

## Limitations and Considerations

### Limitations

1. **No Persistence**: Data is lost when machines are powered off
2. **Memory Requirements**: Machines must have sufficient RAM to hold the entire OS
3. **Performance**: May have different performance characteristics than disk-based deployment
4. **MAAS Version**: Requires MAAS version that supports ephemeral deployment

### Considerations

1. **Resource Planning**: Ensure adequate memory for ephemeral machines
2. **Data Management**: Don't use ephemeral deployment for persistent workloads
3. **Monitoring**: Monitor memory usage of ephemeral machines
4. **Backup Strategy**: Implement appropriate backup strategies for non-ephemeral components

## Troubleshooting

### Common Issues

1. **Deployment Fails**: Check MAAS version compatibility
2. **Memory Errors**: Ensure sufficient RAM for ephemeral deployment
3. **Client Errors**: Verify MAAS endpoint and API key configuration

### Debugging

Enable verbose logging to troubleshoot deployment issues:

```bash
export LOG_LEVEL=debug
clusterctl generate cluster $CLUSTER_NAME --from templates/cluster-template-ephemeral.yaml
```

### Support

For issues with ephemeral deployment:

1. Check MAAS server logs
2. Verify machine allocation and deployment parameters
3. Ensure MAAS supports ephemeral deployment
4. Review cluster API provider logs

## Migration Guide

### From Disk to Ephemeral

To migrate existing clusters to use ephemeral deployment:

1. **Update Machine Templates**:
   ```yaml
   spec:
     template:
       spec:
         ephemeral: true  # Add this field
   ```

2. **Rolling Update**:
   ```bash
   kubectl patch machinedeployment my-cluster-md-0 \
     --type='merge' \
     -p='{"spec":{"template":{"spec":{"infrastructureRef":{"name":"new-template-with-ephemeral"}}}}}'
   ```

### From Ephemeral to Disk

To migrate back to disk-based deployment:

1. **Update Machine Templates**:
   ```yaml
   spec:
     template:
       spec:
         ephemeral: false  # or remove the field entirely
   ```

2. **Rolling Update**: Same process as above

## Future Enhancements

Planned improvements for the ephemeral deployment feature:

1. **Selective Ephemeral**: Per-machine ephemeral configuration
2. **Hybrid Deployments**: Mix of ephemeral and disk-based machines
3. **Auto-scaling**: Automatic ephemeral deployment for scaling events
4. **Monitoring**: Enhanced monitoring for ephemeral machines
5. **Persistence Options**: Configurable persistence levels 