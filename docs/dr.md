# DR Commands

The dr command supports the following sub-commands:

* [init](#init)
* [test](#test)

> [!IMPORTANT]
> This command is a developer preview, unsupported and not fully tested.
> Please see the following document for more info on developer preview:
> https://access.redhat.com/support/offerings/devpreview.

## init

The init command crates a configuration file required for all other dr commands.

```bash
$ odf dr init --allow-developer-preview

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
$ odf dr test run -o odf-dr-test --allow-developer-preview
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

✅ passed (1 passed, 0 failed, 0 skipped)
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

The `test-run.yaml` is a machine and human readable description of the the test run.

```yaml
host:
  arch: arm64
  cpus: 12
  os: darwin
name: test-run
status: passed
steps:
- name: validate
  status: passed
- name: setup
  status: passed
- items:
  - config:
      deployer: appset
      pvcSpec: rbd
      workload: deploy
    name: appset-deploy-rbd
    status: passed
  name: tests
  status: passed
summary:
  failed: 0
  passed: 1
  skipped: 0
```

You can query it with tools like `yq`:

```bash
$ yq .status < odf-dr-test/test-run.yaml
passed
```

To clean up after the test use the [clean](#test-clean) command.

### test clean

The clean command delete resources created by the [run](#test-run) command.

```bash
$ odf dr test clean -o odf-dr-test --allow-developer-preview
⭐ Using report "odf-dr-test"
⭐ Using config "config.yaml"

🔎 Validate config ...
   ✅ Config validated

🔎 Clean tests ...
   ✅ Application "appset-deploy-rbd" unprotected
   ✅ Application "appset-deploy-rbd" undeployed

🔎 Clean environment ...
   ✅ Environment cleaned

✅ passed (1 passed, 0 failed, 0 skipped)
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
