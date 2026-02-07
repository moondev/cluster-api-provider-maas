# Templates

This directory contains example manifests for different deployment patterns.

## ClusterClass-less examples

- [cluster-template.yaml](cluster-template.yaml): baseline KubeadmControlPlane + MaaS workers.
- [cluster-template-ephemeral.yaml](cluster-template-ephemeral.yaml): ephemeral cluster workflow.
- [cluster-template-non-clusterclass-basic.yaml](cluster-template-non-clusterclass-basic.yaml): minimal CP + workers.
- [cluster-template-non-clusterclass-inmemory.yaml](cluster-template-non-clusterclass-inmemory.yaml): in-memory MaaS deploy for CP + workers.
- [cluster-template-non-clusterclass-control-plane-endpoint.yaml](cluster-template-non-clusterclass-control-plane-endpoint.yaml): explicit $\texttt{controlPlaneEndpoint}$ override.
- [cluster-template-non-clusterclass-multiple-worker-pools.yaml](cluster-template-non-clusterclass-multiple-worker-pools.yaml): multiple worker pools with distinct templates.

## External/hosted control plane examples

- [cluster-template-external-cp-service-lb.yaml](cluster-template-external-cp-service-lb.yaml): external CP endpoint from Service.
- [cluster-template-external-cp-ingress.yaml](cluster-template-external-cp-ingress.yaml): external CP endpoint from Ingress.
- [cluster-template-external-cp-mixed-infra.yaml](cluster-template-external-cp-mixed-infra.yaml): external CP + MaaS workers only.

## ClusterClass example

- [clusterclass-maas.yaml](clusterclass-maas.yaml): ClusterClass with MaaS infra.
