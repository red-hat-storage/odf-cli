# Subvolume cleanup

The subvolume command is used to clean the stale subvolumes
which have no parent-pvc attached to them.
The command would list out all such subvolumes which needs to be removed.
This would consider all the cases where we can have stale subvolume
and delete them without impacting other resources and attached volumes.

The subvolume command supports the following sub commands:

* [ls](#ls)
* [delete](#delete)

## ls

This command will lists all the subvolumes. It also accepts the stale flag to check for stale subvolumes.

* `--stale`: lists only stale subvolumes
* `--svg <subvolumegroupname`: lists subvolumes in a particular subvolume(default is "csi")

```bash
odf subvolume ls

# Filesystem  Subvolume Subvolumegroup State
# ocs-storagecluster-cephfilesystem csi-vol-427774b4-340b-11ed-8d66-0242ac110004 csi
# ocs-storagecluster-cephfilesystem csi-vol-427774b4-340b-11ed-8d66-0242ac110005 csi
# ocs-storagecluster-cephfilesystem csi-vol-427774b4-340b-11ed-8d66-0242ac110006 csi
# ocs-storagecluster-cephfilesystem csi-vol-427774b4-340b-11ed-8d66-0242ac110007 stale
```

```bash
odf subvolume ls --stale

# Filesystem  Subvolume Subvolumegroup State
# ocs-storagecluster-cephfilesystem csi-vol-427774b4-340b-11ed-8d66-0242ac110004 csi stale
# ocs-storagecluster-cephfilesystem csi-vol-427774b4-340b-11ed-8d66-0242ac110005 csi stale
```

```bash
odf subvolume ls --svg svg01

# Filesystem  Subvolume Subvolumegroup State
# ocs-storagecluster-cephfilesystem csi-vol-427774b4-340b-11ed-8d66-0242ac110005 svg01 in-use
# ocs-storagecluster-cephfilesystem csi-vol-427774b4-340b-11ed-8d66-0242ac110007 svg01 stale
```

## delete

This command deletes stale subvolumes after user's confirmation.
`delete <subvolumes> <filesystem> <subvolumegroup>`:
It will delete only the stale subvolumes to prevent any loss of data.

* subvolumes: comma-separated list of subvolumes of same filesystem and subvolumegroup.

```bash
odf subvolume delete csi-vol-427774b4-340b-11ed-8d66-0242ac110004 ocs-storagecluster csi

# Info: subvolume csi-vol-427774b4-340b-11ed-8d66-0242ac110004 deleted
```

```bash
odf subvolume delete csi-vol-427774b4-340b-11ed-8d66-0242ac110004,csi-vol-427774b4-340b-11ed-8d66-0242ac110005 ocs-storagecluster csi

# Info: subvolume csi-vol-427774b4-340b-11ed-8d66-0242ac110004 deleted
# Info: subvolume csi-vol-427774b4-340b-11ed-8d66-0242ac110004 deleted
```
