# Get

The get command supports the following sub-commands:

* [recovery-profile](#recovery-profile)
* [health](#health)
* [dr-health](#dr-health)
* [mon-endpoints](#mon-endpoints)

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

## dr-health

The DR health command is used to get the connection status of one cluster from another cluster in mirroring-enabled clusters. The cephblockpool is queried with mirroring-enabled and If not found will exit with relevant logs. Optionally, args can be passed to be appended to ceph status, for example: --debug-ms 1.

```bash
$ odf get dr-health

Info: fetching the cephblockpools with mirroring enabled
Info: found "ocs-storagecluster-cephblockpool" cephblockpool with mirroring enabled
Info: running ceph status from peer cluster
Info:   cluster:
    id:     9a2e7e55-40e1-4a79-9bfa-c3e4750c6b0f
    health: HEALTH_OK

  services:
    mon:        3 daemons, quorum a,b,c (age 2w)
    mgr:        a(active, since 2w), standbys: b
    mds:        1/1 daemons up, 1 hot standby
    osd:        3 osds: 3 up (since 2w), 3 in (since 2w)
    rbd-mirror: 1 daemon active (1 hosts)
    rgw:        1 daemon active (1 hosts, 1 zones)

  data:
    volumes: 1/1 healthy
    pools:   12 pools, 185 pgs
    objects: 1.25k objects, 2.5 GiB
    usage:   9.9 GiB used, 290 GiB / 300 GiB avail
    pgs:     185 active+clean

  io:
    client:   18 KiB/s rd, 86 KiB/s wr, 22 op/s rd, 9 op/s wr


Info: running mirroring daemon health
health: WARNING
daemon health: OK
image health: WARNING
images: 4 total
    2 unknown
    2 replaying
```

## mon-endpoints

Prints the mon endpoints

```bash
$ odf get mon-endpoints

10.98.95.196:6789,10.106.118.240:6789,10.111.18.121:6789
```
