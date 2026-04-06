# Subvolume cleanup

The subvolume command is used to clean the stale subvolumes
which have no parent-pvc attached to them.
The command would list out all such subvolumes which needs to be removed.
This would consider all the cases where we can have stale subvolume
and delete them without impacting other resources and attached volumes.

The subvolume command supports the following sub commands:

* `ls` : [ls](#ls) lists all the subvolumes. Supported flags are listed below.
  * `--stale`: lists only stale subvolumes
  * `--svg <subvolumegroupname>`: lists subvolumes in a particular subvolume group (default is "csi")
  * `--rados-namespace <namespace>`: rados namespace used for OMAP lookups (default is "csi")
  * `--storage-client <name>`: StorageClient CR name to auto-resolve SVG, rados-namespace, and exec config
* `delete <filesystem> <subvolume> [subvolumegroup]`:
    [delete](#delete) a stale subvolume. Supported args and flags are listed below.
  * filesystem: the name of the CephFS filesystem.
  * subvolume: subvolume name.
  * subvolumegroup: subvolumegroup name to which the subvolume belongs (default is "csi")
  * `--rados-namespace <namespace>`: rados namespace used for OMAP lookups (default is "csi")
  * `--storage-client <name>`: StorageClient CR name to auto-resolve SVG, rados-namespace, and exec config

## ls

```bash
$ odf subvolume ls

Filesystem                         Subvolume                                     SubvolumeGroup  State
ocs-storagecluster-cephfilesystem  csi-vol-427774b4-340b-11ed-8d66-0242ac110004  csi             in-use
ocs-storagecluster-cephfilesystem  csi-vol-427774b4-340b-11ed-8d66-0242ac110005  csi             in-use
ocs-storagecluster-cephfilesystem  csi-vol-427774b4-340b-11ed-8d66-0242ac110006  csi             stale
ocs-storagecluster-cephfilesystem  csi-vol-427774b4-340b-11ed-8d66-0242ac110007  csi             stale-with-snapshot
```

```bash
$ odf subvolume ls --stale

Filesystem                         Subvolume                                     SubvolumeGroup  State
ocs-storagecluster-cephfilesystem  csi-vol-427774b4-340b-11ed-8d66-0242ac110006  csi             stale
ocs-storagecluster-cephfilesystem  csi-vol-427774b4-340b-11ed-8d66-0242ac110007  csi             stale-with-snapshot
```

```bash
$ odf subvolume ls --svg svg01

Filesystem                         Subvolume                                     SubvolumeGroup  State
ocs-storagecluster-cephfilesystem  csi-vol-427774b4-340b-11ed-8d66-0242ac110005  svg01           in-use
ocs-storagecluster-cephfilesystem  csi-vol-427774b4-340b-11ed-8d66-0242ac110007  svg01           stale
```

### StorageClient

When using a StorageClient CR, pass `--storage-client` to auto-resolve the SVG, rados-namespace, and exec config:

```bash
$ odf subvolume ls --storage-client my-storage-client
```

## delete

```bash
$ odf subvolume delete csi-vol-427774b4-340b-11ed-8d66-0242ac110005

Info: Deleting the omap object and key for subvolume "csi-vol-427774b4-340b-11ed-8d66-0242ac110005"
Info: omap object:"csi.volume.427774b4-340b-11ed-8d66-0242ac110005" deleted
Info: omap key:"csi.volume.pvc-78abf81c-5381-42ee-8d75-dc17cd0cf5de" deleted
Info: subvolume "csi-vol-427774b4-340b-11ed-8d66-0242ac110005" deleted
```

To delete using a custom rados namespace:

```bash
$ odf subvolume delete csi-vol-427774b4-340b-11ed-8d66-0242ac110005 svg01 --rados-namespace=svg01

Info: Deleting the omap object and key for subvolume "csi-vol-427774b4-340b-11ed-8d66-0242ac110005"
Info: omap object:"csi.volume.427774b4-340b-11ed-8d66-0242ac110005" deleted
Info: omap key:"csi.volume.pvc-78abf81c-5381-42ee-8d75-dc17cd0cf5de" deleted
Info: subvolume "csi-vol-427774b4-340b-11ed-8d66-0242ac110005" deleted
```
