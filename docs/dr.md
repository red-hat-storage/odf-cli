# DR Commands

The dr command supports the following sub-commands:

* [init](#init)
* [test](#test)
* [validate](#validate)
* [gather](#gather)

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

## validate

The validate commands help to troubleshoot disaster recovery problems. They
gathers data from the clusters and detects problems in configuration and the
current status of the clusters or protected applications.

```console
$ odf dr validate -h
Detect disaster recovery problems

Usage:
  odf dr validate [command]

Available Commands:
  application Detect problems in disaster recovery protected application

Flags:
  -h, --help            help for validate
  -o, --output string   output directory

Global Flags:
  -c, --config string   configuration file (default "config.yaml")

Use "odf dr validate [command] --help" for more information about a command.
```

The command supports the following sub-commands:

* [application](#validate-application)

### validate application

The validate application command validates a specific DR-protected application
by gathering related namespaces from all clusters and inspecting the gathered
resources.

#### Looking up applications

To run the validate application command, we need to find the protected
application name and namespace. Run the following command:

```console
$ oc get drpc -A --kubeconfig hub
NAMESPACE           NAME                   AGE   PREFERREDCLUSTER   FAILOVERCLUSTER   DESIREDSTATE   CURRENTSTATE
openshift-gitops    appset-deploy-rbd      69m   dr1                dr2               Relocate       Relocated
```

#### Validating an application

To validate the application `appset-deploy-rbd` in namespace `openshift-gitops` run the
following command:

```console
$ odf dr validate application --name appset-deploy-rbd --namespace openshift-gitops -o out
⭐ Using config "config.yaml"
⭐ Using report "out"

🔎 Validate config ...
   ✅ Config validated

🔎 Validate application ...
   ✅ Inspected application
   ✅ Gathered data from cluster "dr2"
   ✅ Gathered data from cluster "dr1"
   ✅ Gathered data from cluster "hub"
   ✅ Application validated

✅ Validation completed (21 ok, 0 stale, 0 problem)
```

The command gathered related namespaces from all clusters, inspected the
resources, and stored output files in the specified output directory:

```console
$ tree -L1 out
out
├── validate-application.data
├── validate-application.log
└── validate-application.yaml
```

> [!IMPORTANT]
> When reporting DR related issues, please create an archive with the output
> directory and upload it to the issue tracker.

#### The validate-application.yaml

The `validate-application.yaml` report is a machine and human readable
description of the command and the application status.

The most important part of the report is the applicationStatus:

```yaml
applicationStatus:
  hub:
    drpc:
      action:
        state: ok ✅
        value: Relocate
      conditions:
      - state: ok ✅
        type: Available
      - state: ok ✅
        type: PeerReady
      - state: ok ✅
        type: Protected
      deleted:
        state: ok ✅
      drPolicy: dr-policy
      name: appset-deploy-rbd
      namespace: openshift-gitops
      phase:
        state: ok ✅
        value: Relocated
      progression:
        state: ok ✅
        value: Completed
  primaryCluster:
    name: dr1
    vrg:
      conditions:
      - state: ok ✅
        type: DataReady
      - state: ok ✅
        type: ClusterDataReady
      - state: ok ✅
        type: ClusterDataProtected
      - state: ok ✅
        type: KubeObjectsReady
      - state: ok ✅
        type: NoClusterDataConflict
      deleted:
        state: ok ✅
      name: appset-deploy-rbd
      namespace: e2e-appset-deploy-rbd
      protectedPVCs:
      - conditions:
        - state: ok ✅
          type: DataReady
        - state: ok ✅
          type: ClusterDataProtected
        deleted:
          state: ok ✅
        name: busybox-pvc
        namespace: e2e-appset-deploy-rbd
        phase:
          state: ok ✅
          value: Bound
        replication: volrep
      state:
        state: ok ✅
        value: Primary
  secondaryCluster:
    name: dr2
    vrg:
      conditions:
      - state: ok ✅
        type: NoClusterDataConflict
      deleted:
        state: ok ✅
      name: appset-deploy-rbd
      namespace: e2e-appset-deploy-rbd
      state:
        state: ok ✅
        value: Secondary
```

#### The validate-application.data directory

This directory contains all data gathered during validation. The data depend on
the application deployment type. Use the gathered data to investigate the
problems reported in the `validate-application.yaml` report.

```console
$ tree -L3 out/validate-application.data
out/validate-application.data
├── dr1
│   ├── cluster
│   │   ├── namespaces
│   │   ├── persistentvolumes
│   │   └── storage.k8s.io
│   └── namespaces
│       └── e2e-appset-deploy-rbd
├── dr2
│   ├── cluster
│   │   └── namespaces
│   └── namespaces
│       └── e2e-appset-deploy-rbd
└── hub
    ├── cluster
    │   └── namespaces
    └── namespaces
        └── openshift-gitops
```

#### The validate-application.log

This log includes detailed information that may help to troubleshoot the
validate application command. If the command failed, check the error details in
the log.

## gather

The gather command helps to troubleshoot disaster recovery issues by gathering
data about a protected application.

```console
$ odf dr gather -h
Collect diagnostic data from your clusters

Usage:
  odf dr gather [command]

Available Commands:
  application Collect data for a protected application

Flags:
  -h, --help            help for gather
  -o, --output string   output directory

Global Flags:
  -c, --config string   configuration file (default "config.yaml")

Use "odf dr gather [command] --help" for more information about a command.
```

> [!IMPORTANT]
> The gather command requires a configuration file. See [init](docs/init.md) to
> learn how to create one.

### gather application

The gather application command gathers data for a specific disaster recover
protected application. It gathers entire namespaces related to the protected
application across the hub and the managed clusters.

#### Looking up applications

To run the gather application command, we need to find the protected
application name and namespace. Run the following command:

```console
$ kubectl get drpc -A --context hub
NAMESPACE   NAME                AGE     PREFERREDCLUSTER   FAILOVERCLUSTER   DESIREDSTATE   CURRENTSTATE
argocd      appset-deploy-rbd   6m16s   dr1                                                 Deployed
```

#### Gathering application data

To gather data for the application `appset-deploy-rbd` in namespace `argocd`
run the following command:

```console
$ odf dr gather application --name appset-deploy-rbd --namespace argocd -o out
⭐ Using config "config.yaml"
⭐ Using report "out"

🔎 Validate config ...
   ✅ Config validated

🔎 Gather application data ...
   ✅ Inspected application
   ✅ Gathered data from cluster "hub"
   ✅ Gathered data from cluster "dr1"
   ✅ Gathered data from cluster "dr2"

✅ Gather completed
```

The command gathered related namespaces from all clusters and stored output
files in the specified output directory:

```console
$ tree -L1 out
out
├── gather-application.data
├── gather-application.log
└── gather-application.yaml
```

#### The gather-application.data directory

This directory contains the namespaces and cluster scope resources related to
the protected application. The data depend on the application deployment type.

```console
$ tree -L3 out/gather-application.data
out/gather-application.data
├── dr1
│   ├── cluster
│   │   ├── namespaces
│   │   ├── persistentvolumes
│   │   └── storage.k8s.io
│   └── namespaces
│       ├── e2e-appset-deploy-rbd
│       └── ramen-system
├── dr2
│   ├── cluster
│   │   └── namespaces
│   └── namespaces
│       ├── e2e-appset-deploy-rbd
│       └── ramen-system
└── hub
    ├── cluster
    │   └── namespaces
    └── namespaces
        ├── argocd
        └── ramen-system
```

You can use standard tools to inspect the resources:

```console
$ yq '.status.protectedPVCs[0].conditions' < out/gather-application.data/dr1/namespaces/e2e-appset-deploy-rbd/ramendr.openshift.io/volumereplicationgroups/appset-deploy-rbd.yaml
- lastTransitionTime: "2025-08-17T17:45:41Z"
  message: PVC in the VolumeReplicationGroup is ready for use
  observedGeneration: 1
  reason: Ready
  status: "True"
  type: DataReady
- lastTransitionTime: "2025-08-17T17:45:40Z"
  message: PV cluster data already protected for PVC busybox-pvc
  observedGeneration: 1
  reason: Uploaded
  status: "True"
  type: ClusterDataProtected
- lastTransitionTime: "2025-08-17T17:45:41Z"
  message: PVC in the VolumeReplicationGroup is ready for use
  observedGeneration: 1
  reason: Replicating
  status: "False"
  type: DataProtected
```

You can also inspect ramen logs in all clusters:

```console
$ grep -E 'ERROR.+appset-deploy-rbd' out/gather-application.data/dr1/namespaces/ramen-system/pods/ramen-dr-cluster-operator-67dff877f5-k4gjm/manager/current.log
2025-08-17T17:45:40.644Z	ERROR	vrg	controller/vrg_volrep.go:122	Requeuing due to failure to upload PV object to S3 store(s)	{"vrg": {"name":"appset-deploy-rbd","namespace":"e2e-appset-deploy-rbd"}, "rid": "1c5b6d55", "State": "primary", "pvc": "e2e-appset-deploy-rbd/busybox-pvc", "error": "failed to add archived annotation for PVC (e2e-appset-deploy-rbd/busybox-pvc): failed to update PersistentVolumeClaim (e2e-appset-deploy-rbd/busybox-pvc) annotation (volumereplicationgroups.ramendr.openshift.io/vr-archived) belonging to VolumeReplicationGroup (e2e-appset-deploy-rbd/appset-deploy-rbd), Operation cannot be fulfilled on persistentvolumeclaims \"busybox-pvc\": the object has been modified; please apply your changes to the latest version and try again"}
```

#### The gather-application.yaml

The `gather-application.yaml` report is a machine and human readable description
of the command. It can be useful to troubleshoot the gather application command.

#### The gather-application.log

This log includes detailed information that may help to troubleshoot the gather
application command. If the command failed, check the error details in the log.
