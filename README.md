[![REUSE status](https://api.reuse.software/badge/github.com/kyma-project/gardener-syncer)](https://api.reuse.software/info/github.com/kyma-project/gardener-syncer)

# Gardener Syncer Job


## Overview

Cronjob which is regularly synchronizing the Seed-data from Gardener System to Kyma Control Plane.
The fetched Gardener Seed data for available cloud providers regions is stored inside a ConfigMap.

## Prerequisites

- Gardener Syncer Job is running in the Kyma Control Plane.
- Gardener Syncer Job is configured to connect to the Gardener System with kubeconfig file. 

## Installation

The Gardener Syncer Job is installed and run as an ArgoCD application in the Kyma Control Plane.

### Program arguments:
* **gardener-timeout** - Timeout of the Gardener API call, default is 5 seconds
* **gardener-kubeconfig-path** - Path to the kubeconfig file for the Gardener API mounted inside the container, default is `/gardener/kubeconfig/kubeconfig`
* **gardener-seed-map-name** - Name of the ConfigMap where the Seed data will be stored, default is `gardener-seeds-cache`
* **gardener-seed-map-namespace** - Namespace of the ConfigMap where the Seed data will be stored, default is `kcp-system`
* **log-level** - Logging level, possible values are `INFO`, `DEBUG`, default is `INFO`


## Usage

The Gardener Syncer Job runs periodically and fetches the current Seed data from the Gardener System. \
The fetched data is stored in a ConfigMap named `gardener-seeds-cache` in the `kcp-system` namespace. 

The seed information, grouped as regions for all available cloud providers, is stored in the config map structured as follows:


```yaml
apiVersion: v1
data:
  alicloud: |-
    seedRegions:
    - eu-central-1
  aws: |-
    seedRegions:
    - eu-west-1
    - eu-central-1
    - us-east-1
  azure: |-
    seedRegions:
    - westeurope
    - northeurope
    - westus2
    - eastus
    - eastus2
  gcp: |-
    seedRegions:
    - europe-west1
    - us-central1
  openstack: |-
    seedRegions:
    - eu-de-1
kind: ConfigMap
metadata:
  name: gardener-seeds-cache
  namespace: kcp-system
```

You can check the status of the job by looking at the logs of the CronJob in the Kyma Control Plane. The job will log any errors encountered during the synchronization process.
You can also manually trigger the job by running the following command:

```bash 
kubectl create job --from=cronjob/gardener-syncer-job gardener-syncer-job-manual --namespace kcp-system
```

Finally, the list of regions with exiting seed information for each cloud provider is available for all interested KCP services like Kyma Infrastructure Manager (KIM) or Kyma Environment Broker (KEB)

## Contributing
<!--- mandatory section - do not change this! --->

See the [Contributing Rules](CONTRIBUTING.md) and [docs](./docs/contributor/kcp-sync-gardener-seed.md).

## Code of Conduct
<!--- mandatory section - do not change this! --->

See the [Code of Conduct](CODE_OF_CONDUCT.md) document.

## Licensing
<!--- mandatory section - do not change this! --->

See the [license](./LICENSE) file.
