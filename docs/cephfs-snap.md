# Snapshot cleanup

The cephfs-snap command is used to manage and clean up orphaned CephFS snapshots
which have no corresponding Kubernetes VolumeSnapshotContent resource.
The command lists all snapshots in a subvolumegroup and identifies their status as bound or orphaned.
Orphaned snapshots can then be safely deleted without impacting active resources.

The cephfs-snap command supports the following sub commands:

* `ls` : [ls](#ls) lists all the snapshots. Supported flags are listed below.
  * `--orphaned`: lists only orphaned snapshots
  * `--svg <subvolumegroupname>`: lists snapshots in a particular subvolume group (default is "csi")
  * `--filesystem <filesystemname>`: lists snapshots in a particular filesystem (default is "ocs-storagecluster-cephfilesystem")
  * `--rados-namespace <namespace>`: rados namespace used for OMAP lookups (default is "csi")
  * `--storage-client <name>`: StorageClient CR name to auto-resolve SVG, rados-namespace, and exec config
* `delete <subvol> <snapshot>`:
    [delete](#delete) an orphaned snapshot. Supported args and flags are listed below.
  * subvol: subvolume name to which the snapshot belongs.
  * snapshot: snapshot name.
  * `--svg <subvolumegroupname>`: subvolume group name (default is "csi")
  * `--filesystem <filesystemname>`: the name of the CephFS filesystem (default is "ocs-storagecluster-cephfilesystem")
  * `--rados-namespace <namespace>`: rados namespace used for OMAP lookups (default is "csi")
  * `--storage-client <name>`: StorageClient CR name to auto-resolve SVG, rados-namespace, and exec config

## ls

```bash
$ odf cephfs-snap ls

Filesystem                         Subvolume                                     SubvolumeGroup  Snapshot                                       State
ocs-storagecluster-cephfilesystem  csi-vol-aa0099b5-f7a0-49c2-bc97-a810005a9654  csi             csi-snap-3936435c-a14a-4a76-9d0f-71321ac084a9  bound
ocs-storagecluster-cephfilesystem  csi-vol-aa0099b5-f7a0-49c2-bc97-a810005a9654  csi             csi-snap-4047546d-b25b-5b87-da08-82432bd195ba  orphaned
```

```bash
$ odf cephfs-snap ls --filesystem fs01

Filesystem  Subvolume                                     SubvolumeGroup  Snapshot                                       State
fs01        csi-vol-aa0099b5-f7a0-49c2-bc97-a810005a9654  csi             csi-snap-3936435c-a14a-4a76-9d0f-71321ac084a9  bound
fs01        csi-vol-aa0099b5-f7a0-49c2-bc97-a810005a9654  csi             csi-snap-4047546d-b25b-5b87-da08-82432bd195ba  orphaned
```

```bash
$ odf cephfs-snap ls --orphaned

Filesystem                         Subvolume                                     SubvolumeGroup  Snapshot                                       State
ocs-storagecluster-cephfilesystem  csi-vol-aa0099b5-f7a0-49c2-bc97-a810005a9654  csi             csi-snap-4047546d-b25b-5b87-da08-82432bd195ba  orphaned
```

### StorageClient

When using a StorageClient CR, pass `--storage-client` to auto-resolve the SVG, rados-namespace, and exec config:

```bash
$ odf cephfs-snap ls --storage-client my-storage-client
```

## delete

```bash
$ odf cephfs-snap delete csi-vol-427774b4-340b-11ed-8d66-0242ac110005 csi-snap-b2c3d4e5-450e-11ed-8d66-0242ac110005

Info: Deleting the omap object and key for snapshot "csi-snap-b2c3d4e5-450e-11ed-8d66-0242ac110005"
Info: omap object:"csi.snap.b2c3d4e5-450e-11ed-8d66-0242ac110005" deleted
Info: omap key:"csi.snap.snapshot-a1b2c3d4-5678-9012-abcd-ef0123456789" deleted
snapshot csi-snap-b2c3d4e5-450e-11ed-8d66-0242ac110005 deleted successfully
```

```bash
$ odf cephfs-snap delete csi-vol-427774b4-340b-11ed-8d66-0242ac110005 csi-snap-a1b2c3d4-450e-11ed-8d66-0242ac110004

Error: snapshot "csi-snap-a1b2c3d4-450e-11ed-8d66-0242ac110004" is bound and cannot be deleted
```

To delete using a custom rados namespace:

```bash
$ odf cephfs-snap delete csi-vol-427774b4-340b-11ed-8d66-0242ac110005 csi-snap-b2c3d4e5-450e-11ed-8d66-0242ac110005 --rados-namespace=svg01

Info: Deleting the omap object and key for snapshot "csi-snap-b2c3d4e5-450e-11ed-8d66-0242ac110005"
Info: omap object:"csi.snap.b2c3d4e5-450e-11ed-8d66-0242ac110005" deleted
Info: omap key:"csi.snap.snapshot-a1b2c3d4-5678-9012-abcd-ef0123456789" deleted
snapshot "csi-snap-b2c3d4e5-450e-11ed-8d66-0242ac110005" deleted successfully
```
