# Release Notes - v0.6.2

## 🚀 New Features

### Ephemeral Deployment Support
- **Added `ephemeralDeploy` field** to `MaasMachineSpec` and `MaasMachineTemplateSpec`
- **In-memory deployment**: Machines can now be deployed in ephemeral mode (in-memory) instead of on-disk
- **Performance benefits**: Faster deployment and reduced storage costs for temporary workloads
- **Use cases**: Ideal for testing, development environments, and stateless workloads

### Configuration
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

## 🔧 Technical Changes

### API Changes
- Added `ephemeralDeploy *bool` field to `MaasMachineSpec`
- Field is optional and defaults to `false` (persistent deployment)
- Updated CRDs to include the new field

### Controller Changes
- Updated machine deployment logic in `pkg/maas/machine/machine.go`
- Added conditional logic to set ephemeral deployment based on spec
- Enhanced logging to track deployment mode

### Dependencies
- Updated to use `github.com/moondev/maas-client-go` instead of `github.com/spectrocloud/maas-client-go`
- Updated module name to `github.com/moondev/cluster-api-provider-maas`

## 📋 Files Modified

### API Types
- `api/v1beta1/maasmachine_types.go` - Added ephemeralDeploy field

### Controller Logic
- `pkg/maas/machine/machine.go` - Updated deployment logic

### CRDs
- `config/crd/bases/infrastructure.cluster.x-k8s.io_maasmachines.yaml`
- `config/crd/bases/infrastructure.cluster.x-k8s.io_maasmachinetemplates.yaml`

### Samples
- `config/samples/infrastructure_v1beta1_maasmachine.yaml`
- `config/samples/infrastructure_v1beta1_maasmachine_ephemeral.yaml` (new)
- `config/samples/infrastructure_v1beta1_maasmachinetemplate.yaml`
- `config/samples/infrastructure_v1beta1_maasmachinetemplate_ephemeral.yaml` (new)

### Documentation
- `README.md` - Added comprehensive documentation

## 🧪 Testing

### Test Fixes
- Fixed mock generation to include all required interfaces
- Updated test context expectations to match implementation
- All tests now pass successfully

## 📦 Release Artifacts

The following files are generated for release:

### Development Build (`_build/dev/`)
- `infrastructure-components.yaml` - Controller deployment manifests
- `cluster-template.yaml` - Cluster template for clusterctl
- `clusterclass-maas.yaml` - ClusterClass template
- `metadata.yaml` - Provider metadata

### Production Manifests (`_build/manifests/`)
- `infrastructure-components.yaml` - Complete infrastructure components

## 🔄 Migration Notes

### From Previous Versions
- **Backward Compatible**: The new `ephemeralDeploy` field is optional
- **Default Behavior**: Existing deployments will continue to use persistent mode
- **No Breaking Changes**: All existing functionality remains unchanged

### Module Name Change
- **Old**: `github.com/spectrocloud/cluster-api-provider-maas`
- **New**: `github.com/moondev/cluster-api-provider-maas`
- **Impact**: Users should update their clusterctl configuration

## 🎯 Usage Examples

### Persistent Deployment (Default)
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

### Ephemeral Deployment
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

## 🚀 Benefits of Ephemeral Deployment

- **Faster Deployment**: In-memory deployment is typically faster than on-disk deployment
- **Resource Efficiency**: Ephemeral machines don't consume disk space for the OS image
- **Testing & Development**: Ideal for temporary workloads and development environments
- **Cost Optimization**: Reduced storage costs for temporary workloads

## ⚠️ Considerations

- **Data Persistence**: Ephemeral machines lose all data when powered off or rebooted
- **Performance**: May have different performance characteristics compared to persistent deployments
- **Use Cases**: Best suited for stateless workloads, testing, and development 