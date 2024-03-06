# Restore

## deleted

When a Rook CR is deleted, the Rook operator will respond to the deletion event to attempt to clean up the cluster resources. If any data is still present in the cluster, Rook will refuse to delete the CR to ensure data is not lost. The operator will refuse to remove the finalizer on the CR until the underlying data is deleted.

While the underlying Ceph data and daemons continue to be available, the CRs will be stuck indefinitely in a Deleting state in which the operator will not continue to ensure cluster health. Upgrades will be blocked, further updates to the CRs are prevented, and so on. Since Kubernetes does not allow undeleting resources, the command below will allow repairing the CRs without even necessarily suffering cluster downtime.

> [!NOTE]
> If there are multiple deleted resources in the cluster and no specific resource is mentioned, the first resource will be restored. To restore all deleted resources, re-run the command multiple times.

The `restore-deleted` command has one required and one optional parameter:

- `<CRD>`: The CRD type that is to be restored, such as CephClusters, CephFilesystems, CephBlockPools and so on.
- `[CRName]`: The name of the specific CR which you want to restore since there can be multiple instances under the same CRD. For example, if there are multiple CephFilesystems stuck in deleting state, a specific filesystem can be restored: `restore-deleted cephfilesystems filesystem-2`.

```bash
odf restore deleted <CRD> [CRName]
```

```bash
$ odf restore deleted cephclusters

Info: Detecting which resources to restore for crd "cephclusters"

Info: Restoring CR my-cluster
Warning: The resource my-cluster was found deleted. Do you want to restore it? yes | no
```

After entering `yes`, the restore continues with output as such:

```bash
yes
Info: Proceeding with restoring deleting CR
Info: Scaling down the operator
Info: Deleting validating webhook rook-ceph-webhook if present
Info: Removing ownerreferences from resources with matching uid 0dfd114c-a9bc-47b8-916b-08d7fd57f227
Info: Removing owner references for secret cluster-peer-token-ocs-storagecluster-cephcluster
Info: Removed ownerReference for Secret: cluster-peer-token-ocs-storagecluster-cephcluster

Info: Removing owner references for secret rook-ceph-admin-keyring
Info: Removed ownerReference for Secret: rook-ceph-admin-keyring

---
---
---

Info: Removing owner references for pvc rook-ceph-mon-e
Info: Removed ownerReference for pvc: rook-ceph-mon-e

Info: Removing finalizers from cephclusters/ocs-storagecluster-cephcluster
Info: Re-creating the CR cephclusters from dynamic resource
W0307 11:24:03.682209   20982 warnings.go:70] metadata.finalizers: "cephcluster.ceph.rook.io": prefer a domain-qualified finalizer name to avoid accidental conflicts with other finalizer writers
Info: Scaling up the operator
Info: CR is successfully restored. Please watch the operator logs and check the crd
```
