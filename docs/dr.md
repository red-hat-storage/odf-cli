# DR Commands

The dr command supports the following sub-commands:

* [init](#init)
* [test](#test)

## init

The init command crates a configuration file required for all other dr commands.

```bash
$ odf dr init

✅ Created config file "config.yaml" - please modify for your clusters
```

> [!IMPORTANT]
> Before using the config file, you need to edit it and configure your
> clusters and storage.

### Sample configuration file

```yaml
## odf dr configuration file

## Clusters configuration.
# - Modify clusters "kubeconfig" to match your hub and managed clusters
#   kubeconfig files.
clusters:
  hub:
    kubeconfig: hub/config
  c1:
    kubeconfig: primary/config
  c2:
    kubeconfig: secondary/config

## Git repository for test command.
# - Modify "url" to use your own Git repository.
# - Modify "branch" to test a different branch.
repo:
  url: https://github.com/RamenDR/ocm-ramen-samples.git
  branch: main

## DRPolicy for test command.
# - Modify to match actual DRPolicy in the hub cluster.
drPolicy: dr-policy

## ClusterSet for test command.
# - Modify to match your Open Cluster Management configuration.
clusterSet: default

## PVC specifications for test command.
# - Modify items "storageclassname" to match the actual storage classes in the
#   managed clusters.
# - Add new items for testing more storage types.
PVCSpecs:
- name: rbd
  storageClassName: rook-ceph-block
  accessModes: ReadWriteOnce
- name: cephfs
  storageClassName: rook-cephfs-fs1
  accessModes: ReadWriteMany

# Deployer specifications for test command.
# - Modify items "name" and "type" to match your deployer configurations.
# - Add new items for testing more deployers.
# - Available types: appset, subscr, disapp
deployers:
- name: appset
  type: appset
  description: ApplicationSet deployer for ArgoCD
- name: subscr
  type: subscr
  description: Subscription deployer for OCM subscriptions
- name: disapp
  type: disapp
  description: Discovered Application deployer for discovered applications

## Tests cases for test command.
# - Modify the test for your preferred workload or deployment type.
# - Add new tests for testing more combinations in parallel.
# - Available workloads: deploy
# - Available deployers: appset, subscr, disapp
tests:
- workload: deploy
  deployer: appset
  pvcSpec: rbd
```

## test

The command supports the following sub-commands:

* [run](#test-run)
* [clean](#test-clean)

### test run

The test command tests complete disaster recovery flow with a tiny application.

```bash
$ odf dr test run -o odf-dr-test
⭐ Using report "odf-dr-test"
⭐ Using config "config.yaml"

🔎 Validate config ...
   ✅ Config validated

🔎 Setup environment ...
   ✅ Environment setup

🔎 Run tests ...
   ✅ Application "appset-deploy-rbd" deployed
   ✅ Application "appset-deploy-rbd" protected
   ✅ Application "appset-deploy-rbd" failed over
   ✅ Application "appset-deploy-rbd" relocated
   ✅ Application "appset-deploy-rbd" unprotected
   ✅ Application "appset-deploy-rbd" undeployed

✅ passed (1 passed, 0 failed, 0 skipped, 0 canceled)
```

The command stores `test-run.yaml` and `test-run.log` in the specified output
directory:

```bash
$ tree odf-dr-test
odf-dr-test
├── test-run.log
└── test-run.yaml
```

> [!IMPORTANT]
> When reporting DR related issues, please create an archive with the output
> directory and upload it to the issue tracker.

See [The test-run.yaml](#the-test-run.yaml) for more info on the test report.

To clean up after the test use the [clean](#test-clean) command.

### test clean

The clean command delete resources created by the [run](#test-run) command.

```bash
$ odf dr test clean -o odf-dr-test
⭐ Using report "odf-dr-test"
⭐ Using config "config.yaml"

🔎 Validate config ...
   ✅ Config validated

🔎 Clean tests ...
   ✅ Application "appset-deploy-rbd" cleaned up

🔎 Clean environment ...
   ✅ Environment cleaned

✅ passed (1 passed, 0 failed, 0 skipped, 0 canceled)
```

The command stores `test-clean.yaml` and `test-clean.log` in the specified
output directory:

```bash
$ tree odf-dr-test
odf-dr-test
├── test-clean.log
├── test-clean.yaml
├── test-run.log
└── test-run.yaml
```

### The test-run.yaml

The `test-run.yaml` is a machine and human readable description of the test run.

```yaml
config:
  channel:
    name: https-github-com-ramendr-ocm-ramen-samples-git
    namespace: test-gitops
  clusterSet: clusterset-submariner-52bbff94cfe4421185
  clusters:
    c1:
      kubeconfig: /Users/nir/envs/ocp/c1
    c2:
      kubeconfig: /Users/nir/envs/ocp/c2
    hub:
      kubeconfig: /Users/nir/envs/ocp/hub
    passive-hub:
      kubeconfig: ""
  deployers:
  - description: ApplicationSet deployer for ArgoCD
    name: appset
    type: appset
  - description: Subscription deployer for OCM subscriptions
    name: subscr
    type: subscr
  - description: Discovered Application deployer for discovered applications
    name: disapp
    type: disapp
  distro: ocp
  drPolicy: dr-policy-1m
  namespaces:
    argocdNamespace: openshift-gitops
    ramenDRClusterNamespace: openshift-dr-system
    ramenHubNamespace: openshift-operators
    ramenOpsNamespace: openshift-dr-ops
  pvcSpecs:
  - accessModes: ReadWriteOnce
    name: rbd
    storageClassName: ocs-storagecluster-ceph-rbd
  - accessModes: ReadWriteMany
    name: cephfs
    storageClassName: ocs-storagecluster-cephfs
  repo:
    branch: main
    url: https://github.com/RamenDR/ocm-ramen-samples.git
  tests:
  - deployer: appset
    pvcSpec: rbd
    workload: deploy
created: "2025-08-19T15:54:21.917077+03:00"
duration: 929.939373665
host:
  arch: arm64
  cpus: 12
  os: darwin
name: test-run
status: passed
steps:
- duration: 2.989362291
  name: validate
  status: passed
- duration: 0.420524791
  name: setup
  status: passed
- duration: 926.529486583
  items:
  - duration: 926.5290810419999
    items:
    - duration: 21.914995542
      name: deploy
      status: passed
    - duration: 74.045327833
      name: protect
      status: passed
    - duration: 423.177328667
      name: failover
      status: passed
    - duration: 350.520063166
      name: relocate
      status: passed
    - duration: 35.47629975
      name: unprotect
      status: passed
    - duration: 21.395066084
      name: undeploy
      status: passed
    name: appset-deploy-rbd
    status: passed
  name: tests
  status: passed
summary:
  canceled: 0
  failed: 0
  passed: 1
  skipped: 0
```

You can query it with tools like `yq`:

```bash
$ yq .status < odf-dr-test/test-run.yaml
passed
```
