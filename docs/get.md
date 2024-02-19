# Get

The get command supports the following sub-commands:

* [recovery-profile](#recovery-profile)
* [health](#health)

## recovery-profile

This command will display the recovery-profile value set for the osd.

```bash
$ odf get recovery-profile
# high_recovery_ops
```

## health

The health command checks the health of the Ceph cluster and common configuration issues. The health command validates these configurations:

1. at least three mon pods should running on different nodes
2. mon quorum and ceph health details
3. at least three osd pods should running on different nodes
4. all pods 'Running' status
5. placement group status
6. at least one mgr pod is running

```bash
$ odf get health

Info: Checking if at least three mon pods are running on different nodes
rook-ceph-mon-a-7fb76597dc-98pxz        Running openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
rook-ceph-mon-b-885bdc59c-4vvcm Running openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
rook-ceph-mon-c-5f59bb5dbc-8vvlg        Running openshift-storage       ip-10-0-30-197.us-west-1.compute.internal

Info: Checking mon quorum and ceph health details
Info: HEALTH_OK

Info: Checking if at least three osd pods are running on different nodes
rook-ceph-osd-0-585bb4cbcf-g5clq        Running openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
rook-ceph-osd-1-5dd9c89487-22rvp        Running openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
rook-ceph-osd-2-7f7d4fd5f9-7lhkt        Running openshift-storage       ip-10-0-30-197.us-west-1.compute.internal

Info: Pods that are in 'Running' or `Succeeded` status
csi-addons-controller-manager-ccd58d558-rr6mr    Running         openshift-storage       ip-10-0-30-197.us-west-1.compute.internal
csi-cephfsplugin-gtpws   Running         openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
csi-cephfsplugin-provisioner-7764cd547f-5p4zg    Running         openshift-storage       ip-10-0-30-197.us-west-1.compute.internal
csi-cephfsplugin-provisioner-7764cd547f-pwlxd    Running         openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
csi-cephfsplugin-rpfgx   Running         openshift-storage       ip-10-0-30-197.us-west-1.compute.internal
csi-cephfsplugin-t7c9k   Running         openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
csi-rbdplugin-6hbwt      Running         openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
csi-rbdplugin-ddrfm      Running         openshift-storage       ip-10-0-30-197.us-west-1.compute.internal
csi-rbdplugin-provisioner-6df8c7664f-2scxq       Running         openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
csi-rbdplugin-provisioner-6df8c7664f-r7b2x       Running         openshift-storage       ip-10-0-30-197.us-west-1.compute.internal
csi-rbdplugin-vdwp2      Running         openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
noobaa-operator-6557d9459c-6d2zb         Running         openshift-storage       ip-10-0-30-197.us-west-1.compute.internal
ocs-metrics-exporter-8467fdcc4-wb95h     Running         openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
ocs-operator-6c9d95576-tvbcj     Running         openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
odf-console-84bd79c79d-25dqg     Running         openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
odf-operator-controller-manager-5dd94c64cf-jkfw8         Running         openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
rook-ceph-crashcollector-34dcd5cf09c4fd6e3de5b84f7f723728-b2nn7          Running         openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
rook-ceph-crashcollector-40d61025f06f85a221c7f61dd6d0d563-fkptw          Running         openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
rook-ceph-crashcollector-45c68baeb3d43ed324ef7a89e15ea97c-mqwl7          Running         openshift-storage       ip-10-0-30-197.us-west-1.compute.internal
rook-ceph-exporter-ip-10-0-30-197.us-west-1.compute.internj6w2t          Running         openshift-storage       ip-10-0-30-197.us-west-1.compute.internal
rook-ceph-exporter-ip-10-0-64-239.us-west-1.compute.internlg8jg          Running         openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
rook-ceph-exporter-ip-10-0-69-145.us-west-1.compute.internlskh4          Running         openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
rook-ceph-mds-ocs-storagecluster-cephfilesystem-a-84c86d75fgbz2          Running         openshift-storage       ip-10-0-30-197.us-west-1.compute.internal
rook-ceph-mds-ocs-storagecluster-cephfilesystem-b-56c86fdfjmm5c          Running         openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
rook-ceph-mgr-a-54f7dbddcd-h7j7m         Running         openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
rook-ceph-mon-a-7fb76597dc-98pxz         Running         openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
rook-ceph-mon-b-885bdc59c-4vvcm          Running         openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
rook-ceph-mon-c-5f59bb5dbc-8vvlg         Running         openshift-storage       ip-10-0-30-197.us-west-1.compute.internal
rook-ceph-operator-698b8bf74c-r2xfw      Running         openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
rook-ceph-osd-0-585bb4cbcf-g5clq         Running         openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
rook-ceph-osd-1-5dd9c89487-22rvp         Running         openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
rook-ceph-osd-2-7f7d4fd5f9-7lhkt         Running         openshift-storage       ip-10-0-30-197.us-west-1.compute.internal
rook-ceph-osd-prepare-6e3c8e8e8920b4b33586a7d07744ab01-vb6pg     Succeeded       openshift-storage       ip-10-0-69-145.us-west-1.compute.internal
rook-ceph-osd-prepare-b6143333237cc6da7383cc8f2e092154-ljksg     Succeeded       openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
rook-ceph-osd-prepare-fa07e0b1110d792c016bd8d38eb717a5-46rx4     Succeeded       openshift-storage       ip-10-0-30-197.us-west-1.compute.internal

Warning: Pods that are 'Not' in 'Running' status
noobaa-core-0    Pending         openshift-storage
noobaa-db-pg-0   Pending         openshift-storage

Info: Checking placement group status
Info:   PgState: active+clean, PgCount: 113

Info: Checking if at least one mgr pod is running
rook-ceph-mgr-a-54f7dbddcd-h7j7m        Running openshift-storage       ip-10-0-64-239.us-west-1.compute.internal
```
