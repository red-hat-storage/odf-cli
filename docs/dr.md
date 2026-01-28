# DR Commands

The dr command supports the following sub-commands:

* [init](#init)
* [test](#test)
* [validate](#validate)
* [gather](#gather)

> [!IMPORTANT]
> This command is a developer preview, unsupported and not fully tested.
> Please see the following document for more info on developer preview:
> https://access.redhat.com/support/offerings/devpreview.

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

The command supports the following sub-commands:

* [application](#validate-application)
* [clusters](#validate-clusters)

### validate application

The validate application command validates a specific DR-protected application
by gathering related namespaces from all clusters, S3 data and inspecting the
gathered resources.

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
   ✅ Inspected S3 profiles
   ✅ Gathered S3 profile "minio-on-dr1"
   ✅ Gathered S3 profile "minio-on-dr2"
   ✅ Application validated

✅ Validation completed (24 ok, 0 stale, 0 problem)
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

The most important part of the report is the `applicationStatus`:

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
  s3:
    profiles:
      state: ok ✅
      value:
      - gathered:
          state: ok ✅
          value: true
        name: minio-on-dr1
      - gathered:
          state: ok ✅
          value: true
        name: minio-on-dr2
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
├── hub
│   ├── cluster
│   │   └── namespaces
│   └── namespaces
│       └── argocd
└── s3
    ├── minio-on-dr1
    │   └── e2e-appset-deploy-rbd
    └── minio-on-dr2
        └── e2e-appset-deploy-rbd
```

#### The validate-application.log

This log includes detailed information that may help to troubleshoot the
validate application command. If the command failed, check the error details in
the log.

### validate clusters

The validate clusters command validates the disaster recovery clusters by
gathering cluster scoped and related ramen resources from all clusters,
and validates that configured S3 endpoints are accessible.

### Validating clusters

To validate the disaster recovery clusters, run the following command:

```console
$ odf dr validate clusters -o out
⭐ Using config "config.yaml"
⭐ Using report "out"

🔎 Validate config ...
   ✅ Config validated

🔎 Validate clusters ...
   ✅ Gathered data from cluster "hub"
   ✅ Gathered data from cluster "dr1"
   ✅ Gathered data from cluster "dr2"
   ✅ Inspected S3 profiles
   ✅ Checked S3 profile "minio-on-dr2"
   ✅ Checked S3 profile "minio-on-dr1"
   ✅ Clusters validated

✅ Validation completed (90 ok, 0 stale, 0 problem)
```

The command gathered cluster scoped and ramen resources from all clusters,
inspected the resources, and stored output files in the specified output
directory:

```console
$ tree -L1 out
out
├── validate-clusters.data
├── validate-clusters.log
└── validate-clusters.yaml
```

> [!IMPORTANT]
> When reporting DR related issues, please create an archive with the output
> directory and upload it to the issue tracker.

#### The validate-clusters.yaml

The `validate-clusters.yaml` report is a machine and human readable description
of the command and the clusters status.

The most important part of the report is the `clustersStatus`:

```yaml
clustersStatus:
  clusters:
  - name: dr1
    ramen:
      configmap:
        deleted:
          state: ok ✅
        name: ramen-dr-cluster-operator-config
        namespace: ramen-system
        ramenControllerType:
          state: ok ✅
          value: dr-cluster
        s3StoreProfiles:
          state: ok ✅
          value:
          - bucket:
              state: ok ✅
              value: bucket
            caCertificate:
              state: ok ✅
            endpoint:
              state: ok ✅
              value: http://dr1-endpoint:30000
            profileName: minio-on-dr1
            region:
              state: ok ✅
              value: us-west-1
            secret:
              awsAccessKeyID:
                state: ok ✅
                value: 57:91:B6:22:67:43:C5:6F:44:F8:27:4C:C6:7B:8B:EB:97:51:B6:40:3E:69:72:43:AD:00:FD:37:AF:56:35:E3
              awsSecretAccessKey:
                state: ok ✅
                value: 17:9C:07:6A:C5:22:15:9D:BB:D5:0D:4F:1D:84:FA:F0:55:51:FE:5B:59:D7:E5:82:4A:80:0D:46:55:9F:B1:D1
              deleted:
                state: ok ✅
              name:
                state: ok ✅
                value: ramen-s3-secret-dr1
              namespace:
                state: ok ✅
                value: ramen-system
          - bucket:
              state: ok ✅
              value: bucket
            caCertificate:
              state: ok ✅
            endpoint:
              state: ok ✅
              value: http://dr2-endpoint:30000
            profileName: minio-on-dr2
            region:
              state: ok ✅
              value: us-east-1
            secret:
              awsAccessKeyID:
                state: ok ✅
                value: 57:91:B6:22:67:43:C5:6F:44:F8:27:4C:C6:7B:8B:EB:97:51:B6:40:3E:69:72:43:AD:00:FD:37:AF:56:35:E3
              awsSecretAccessKey:
                state: ok ✅
                value: 17:9C:07:6A:C5:22:15:9D:BB:D5:0D:4F:1D:84:FA:F0:55:51:FE:5B:59:D7:E5:82:4A:80:0D:46:55:9F:B1:D1
              deleted:
                state: ok ✅
              name:
                state: ok ✅
                value: ramen-s3-secret-dr2
              namespace:
                state: ok ✅
                value: ramen-system
      deployment:
        conditions:
        - state: ok ✅
          type: Available
        - state: ok ✅
          type: Progressing
        deleted:
          state: ok ✅
        name: ramen-dr-cluster-operator
        namespace: ramen-system
        replicas:
          state: ok ✅
          value: 1
  - name: dr2
    ramen:
      configmap:
        deleted:
          state: ok ✅
        name: ramen-dr-cluster-operator-config
        namespace: ramen-system
        ramenControllerType:
          state: ok ✅
          value: dr-cluster
        s3StoreProfiles:
          state: ok ✅
          value:
          - bucket:
              state: ok ✅
              value: bucket
            caCertificate:
              state: ok ✅
            endpoint:
              state: ok ✅
              value: http://dr1-endpoint:30000
            profileName: minio-on-dr1
            region:
              state: ok ✅
              value: us-west-1
            secret:
              awsAccessKeyID:
                state: ok ✅
                value: 57:91:B6:22:67:43:C5:6F:44:F8:27:4C:C6:7B:8B:EB:97:51:B6:40:3E:69:72:43:AD:00:FD:37:AF:56:35:E3
              awsSecretAccessKey:
                state: ok ✅
                value: 17:9C:07:6A:C5:22:15:9D:BB:D5:0D:4F:1D:84:FA:F0:55:51:FE:5B:59:D7:E5:82:4A:80:0D:46:55:9F:B1:D1
              deleted:
                state: ok ✅
              name:
                state: ok ✅
                value: ramen-s3-secret-dr1
              namespace:
                state: ok ✅
                value: ramen-system
          - bucket:
              state: ok ✅
              value: bucket
            caCertificate:
              state: ok ✅
            endpoint:
              state: ok ✅
              value: http://dr2-endpoint:30000
            profileName: minio-on-dr2
            region:
              state: ok ✅
              value: us-east-1
            secret:
              awsAccessKeyID:
                state: ok ✅
                value: 57:91:B6:22:67:43:C5:6F:44:F8:27:4C:C6:7B:8B:EB:97:51:B6:40:3E:69:72:43:AD:00:FD:37:AF:56:35:E3
              awsSecretAccessKey:
                state: ok ✅
                value: 17:9C:07:6A:C5:22:15:9D:BB:D5:0D:4F:1D:84:FA:F0:55:51:FE:5B:59:D7:E5:82:4A:80:0D:46:55:9F:B1:D1
              deleted:
                state: ok ✅
              name:
                state: ok ✅
                value: ramen-s3-secret-dr2
              namespace:
                state: ok ✅
                value: ramen-system
      deployment:
        conditions:
        - state: ok ✅
          type: Progressing
        - state: ok ✅
          type: Available
        deleted:
          state: ok ✅
        name: ramen-dr-cluster-operator
        namespace: ramen-system
        replicas:
          state: ok ✅
          value: 1
  hub:
    drClusters:
      state: ok ✅
      value:
      - conditions:
        - state: ok ✅
          type: Fenced
        - state: ok ✅
          type: Clean
        - state: ok ✅
          type: Validated
        name: dr1
        phase: Available
      - conditions:
        - state: ok ✅
          type: Fenced
        - state: ok ✅
          type: Clean
        - state: ok ✅
          type: Validated
        name: dr2
        phase: Available
    drPolicies:
      state: ok ✅
      value:
      - conditions:
        - state: ok ✅
          type: Validated
        drClusters:
        - dr1
        - dr2
        name: dr-policy
        schedulingInterval: 1m
    ramen:
      configmap:
        deleted:
          state: ok ✅
        name: ramen-hub-operator-config
        namespace: ramen-system
        ramenControllerType:
          state: ok ✅
          value: dr-hub
        s3StoreProfiles:
          state: ok ✅
          value:
          - bucket:
              state: ok ✅
              value: bucket
            caCertificate:
              state: ok ✅
            endpoint:
              state: ok ✅
              value: http://dr1-endpoint:30000
            profileName: minio-on-dr1
            region:
              state: ok ✅
              value: us-west-1
            secret:
              awsAccessKeyID:
                state: ok ✅
                value: 57:91:B6:22:67:43:C5:6F:44:F8:27:4C:C6:7B:8B:EB:97:51:B6:40:3E:69:72:43:AD:00:FD:37:AF:56:35:E3
              awsSecretAccessKey:
                state: ok ✅
                value: 17:9C:07:6A:C5:22:15:9D:BB:D5:0D:4F:1D:84:FA:F0:55:51:FE:5B:59:D7:E5:82:4A:80:0D:46:55:9F:B1:D1
              deleted:
                state: ok ✅
              name:
                state: ok ✅
                value: ramen-s3-secret-dr1
              namespace:
                state: ok ✅
                value: ramen-system
          - bucket:
              state: ok ✅
              value: bucket
            caCertificate:
              state: ok ✅
            endpoint:
              state: ok ✅
              value: http://dr2-endpoint:30000
            profileName: minio-on-dr2
            region:
              state: ok ✅
              value: us-east-1
            secret:
              awsAccessKeyID:
                state: ok ✅
                value: 57:91:B6:22:67:43:C5:6F:44:F8:27:4C:C6:7B:8B:EB:97:51:B6:40:3E:69:72:43:AD:00:FD:37:AF:56:35:E3
              awsSecretAccessKey:
                state: ok ✅
                value: 17:9C:07:6A:C5:22:15:9D:BB:D5:0D:4F:1D:84:FA:F0:55:51:FE:5B:59:D7:E5:82:4A:80:0D:46:55:9F:B1:D1
              deleted:
                state: ok ✅
              name:
                state: ok ✅
                value: ramen-s3-secret-dr2
              namespace:
                state: ok ✅
                value: ramen-system
      deployment:
        conditions:
        - state: ok ✅
          type: Available
        - state: ok ✅
          type: Progressing
        deleted:
          state: ok ✅
        name: ramen-hub-operator
        namespace: ramen-system
        replicas:
          state: ok ✅
          value: 1
  s3:
    profiles:
      state: ok ✅
      value:
      - accessible:
          state: ok ✅
          value: true
        name: minio-on-dr2
      - accessible:
          state: ok ✅
          value: true
        name: minio-on-dr1
```

#### The validate-clusters.data directory

This directory contains all data gathered during validation. Use the gathered
data to investigate the problems reported in the `validate-clusters.yaml` report.

```console
$ tree -L3 out/validate-clusters.data
out/validate-clusters.data
├── dr1
│   ├── cluster
│   │   ├── apiextensions.k8s.io
│   │   ├── apiregistration.k8s.io
│   │   ├── cluster.open-cluster-management.io
│   │   ├── flowcontrol.apiserver.k8s.io
│   │   ├── namespaces
│   │   ├── networking.k8s.io
│   │   ├── nodes
│   │   ├── operator.open-cluster-management.io
│   │   ├── operators.coreos.com
│   │   ├── persistentvolumes
│   │   ├── ramendr.openshift.io
│   │   ├── rbac.authorization.k8s.io
│   │   ├── replication.storage.openshift.io
│   │   ├── scheduling.k8s.io
│   │   ├── snapshot.storage.k8s.io
│   │   ├── storage.k8s.io
│   │   ├── submariner.io
│   │   └── work.open-cluster-management.io
│   └── namespaces
│       └── ramen-system
├── dr2
│   ├── cluster
│   │   ├── apiextensions.k8s.io
│   │   ├── apiregistration.k8s.io
│   │   ├── cluster.open-cluster-management.io
│   │   ├── flowcontrol.apiserver.k8s.io
│   │   ├── namespaces
│   │   ├── networking.k8s.io
│   │   ├── nodes
│   │   ├── operator.open-cluster-management.io
│   │   ├── operators.coreos.com
│   │   ├── persistentvolumes
│   │   ├── ramendr.openshift.io
│   │   ├── rbac.authorization.k8s.io
│   │   ├── replication.storage.openshift.io
│   │   ├── scheduling.k8s.io
│   │   ├── snapshot.storage.k8s.io
│   │   ├── storage.k8s.io
│   │   ├── submariner.io
│   │   └── work.open-cluster-management.io
│   └── namespaces
│       └── ramen-system
└── hub
    ├── cluster
    │   ├── addon.open-cluster-management.io
    │   ├── admissionregistration.k8s.io
    │   ├── apiextensions.k8s.io
    │   ├── apiregistration.k8s.io
    │   ├── cluster.open-cluster-management.io
    │   ├── flowcontrol.apiserver.k8s.io
    │   ├── namespaces
    │   ├── networking.k8s.io
    │   ├── nodes
    │   ├── operator.open-cluster-management.io
    │   ├── operators.coreos.com
    │   ├── ramendr.openshift.io
    │   ├── rbac.authorization.k8s.io
    │   ├── scheduling.k8s.io
    │   └── storage.k8s.io
    └── namespaces
        └── ramen-system
```

#### The validate-clusters.log

This log includes detailed information that may help to troubleshoot the
validate clusters command. If the command failed, check the error details in
the log.

## gather

The gather command helps to troubleshoot disaster recovery issues by gathering
data about a protected application.

The command supports the following sub-commands:

* [application](#gather-application)

### gather application

The gather application command gathers data for a specific disaster recover
protected application. It gathers entire namespaces and S3 data related to the
protected application across the hub and the managed clusters.

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
   ✅ Inspected S3 profiles
   ✅ Gathered S3 profile "minio-on-dr1"
   ✅ Gathered S3 profile "minio-on-dr2"

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
├── hub
│   ├── cluster
│   │   └── namespaces
│   └── namespaces
│       ├── argocd
│       └── ramen-system
└── s3
    ├── minio-on-dr1
    │   └── test-appset-deploy-rbd
    └── minio-on-dr2
        └── test-appset-deploy-rbd
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
