# Installation and Setup of Gardener Syncer

## Context

The Gardener Syncer is an application that periodically fetches available Gardener Seed data and stores a list of Seed regions in a ConfigMap within the Kyma Control Plane (KCP).
This information is essential for the KCP services to validate Kyma runtime provisioning parameters.

When the user enables the "shoot in the same region as seed" feature, the Kyma Environment Broker (KEB) service uses the Gardener Syncer ConfigMap to check if a Seed is available in the same region as the Shoot cluster.
If there is no Seed available, the KEB service interrupts the provisioning process and returns an error.

## Installation   

The Gardener Syncer is deployed as a Kubernetes CronJob on KCP. It is configured to run periodically with user user-configured schedule.
Gardener Syncer is installed on KCP using the Helm chart provided in a separate repository.

This Helm chart includes all necessary configurations to deploy the Gardener Syncer job and contains some values that can be customized to fit the target environment.
The chart is prepared as an ArgoCD application and can be installed using the [ArgoCD](https://argoproj.github.io/) UI or CLI.


![Deployment](./assets/syncer-deployment.png)

See the [Configuration documentation](./configuration.md) for the complete list of Helm Chart parameters.

## Dependencies for Loading Taint Toleration Data

The Gardener Syncer job requires an additional configuration file to get the information on [Taint toleration](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/). 
When processed, Seed has a Taint, and the toleration for it is not present in the configuration file. Such a Seed is not processed.

The toleration configuration file is mounted from a ConfigMap named `gardener-shoot-converter-config` in the `kcp-system` namespace.
Because this ConfigMap is a part of the Kyma Infrastructure Manager (KIM) Helm ArgoCD application, KIM must always be present in the environment before Gardener Syncer Job is executed.
If the configuration file is detected as missing during job execution, the Gardener Syncer Job fails.

## Runtime Environment Prerequisites

- Access to KCP where the Gardener Syncer runs.
- Access to the Gardener Cluster API with a valid kubeconfig file.
- Access to the KIM converter configuration file. 
