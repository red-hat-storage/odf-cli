# Get

The get command supports the following sub-commands:

* [recovery-profile](#recovery-profile)
* [health](#health)
* [dr-health](#dr-health)
* [dr-prereq](#dr-prereq)
* [mon-endpoints](#mon-endpoints)

## recovery-profile

This command will display the recovery-profile value set for the osd.

```bash
$ odf get recovery-profile
# high_recovery_ops
```

## health

`odf get health` prints the status of each StorageCluster and the CRs it manages. `odf get health --verbose` checks the health of the Ceph cluster and common configuration issues.

The command lists StorageClusters in all namespaces. Pass `-n` to limit the list to one namespace. StorageClusters in Phase Ignored are skipped. Each remaining StorageCluster gets a header and a table. A kind is listed only if at least one instance exists. Missing or empty phase and message fields print `-`.

```bash
$ odf get health

StorageCluster: ocs-storagecluster
Namespace:      openshift-storage
Phase:          Ready
Mode:           internal

COMPONENT     NAME                                              OBSERVED  STATUS MESSAGE
CephCluster   ocs-storagecluster-cephcluster                    Ready     Cluster created successfully
Block pool    ocs-storagecluster-cephblockpool                  Failure   Failed to create the block pool. See the operator log for details.
Filesystem    ocs-storagecluster-cephfilesystem                 Ready     Filesystem created successfully
Object store  ocs-storagecluster-cephobjectstore                Ready     Object store created successfully
NooBaa        noobaa                                            Ready     noobaa operator completed reconcile - system is ready
Consumer      internal                                          Ready     -
SVG           csi                                               Ready     SVG created successfully
RNS           ocs-storagecluster-cephblockpool-builtin-implicit Ready     RNS created successfully
StorageClient ocs-storagecluster                                Connected -
```

```bash
$ odf get health --verbose
```

The health command validates these configurations:

1. at least three mon pods should running on different nodes
2. mon quorum and ceph health details
3. at least three osd pods should running on different nodes
4. all pods 'Running' status
5. placement group status
6. at least one mgr pod is running


`--verbose` still prints the CR status table after the cluster health report, in a StorageCluster CRs section.

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

## dr-prereq

This command is used to get the status of all the pre-requisites for enabling Disaster Recovery on a pair of clusters. It takes a peer cluster as argument and uses it to compare current cluster configuration with peer cluster configuration. Based on the results, it will show status of all the pre-requisites.

```bash
$ odf get dr-prereq peer-cluster-1

Info: Submariner is installed.
Info: Globalnet is required.
Info: Globalnet is enabled.
```

## mon-endpoints

Prints the mon endpoints

```bash
$ odf get mon-endpoints

10.98.95.196:6789,10.106.118.240:6789,10.111.18.121:6789
```
