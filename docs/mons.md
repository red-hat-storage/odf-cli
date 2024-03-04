# Mon Commands

## Restore Quorum

Mon quorum is critical to the Ceph cluster. If majority of mons are not in quorum,
the cluster will be down. If the majority of mons are also lost permanently,
the quorum will need to be restore to a remaining good mon in order to bring
the Ceph cluster up again.

To restore the quorum in this disaster scenario:

1. Identify that mon quorum is lost. Some indications include:
   - The Rook operator log shows timeout errors and continuously fails to reconcile
   - All commands in the toolbox are unresponsive
   - Multiple mon pods are likely down
2. Identify which mon has good state.
   - Exec to a mon pod and run the following command
     - `ceph daemon mon.<name> mon_status`
     - For example, if connecting to mon.a, run: `ceph daemon mon.a mon_status`
   - If multiple mons respond, find the mon with the highest `election_epoch`
3. Start the toolbox pod if not already running
4. Run the command below to restore quorum to that good mon
5. Follow the prompts to confirm that you want to continue with each critical step of the restore
6. The final prompt will be to restart the operator, which will add new mons to restore the full quorum size

In this example, quorum is restored to mon **c**.

```bash
odf mons restore-quorum c
```

Before the restore proceeds, you will be prompted if you want to continue.
For example, here is the output in a test cluster up to the point of starting the restore.

```bash
$ odf mons restore-quorum a
Info: mon "a" state is "leader"
Info: mon=b, endpoints=172.30.157.248:3300

Info: mon=c, endpoints=172.30.220.5:3300

Info: mon=a, endpoints=172.30.124.246:3300

Info: printing fsid secret d9a71ab9-a497-4da2-af36-53dcd410a5e9

Info: Check for the running toolbox
Info: fetching the deployment rook-ceph-tools to be running

Info: deployment rook-ceph-tools exists

Info: Restoring mon quorum to mon a 172.30.124.246

Info: The mons to discard are: [b c]

Info: The cluster fsid is d9a71ab9-a497-4da2-af36-53dcd410a5e9

Warning: Are you sure you want to restore the quorum to mon a? If so, enter 'yes-really-restore'
```

After entering `yes-really-restore`, the restore continues with output as such:

```bash
Info: proceeding with resorting quorum
Info: Waiting for operator pod to stop
Info: rook-ceph-operator deployment scaled down
Info: Waiting for bad mon pod to stop
Info: deployment.apps/rook-ceph-mon-b scaled

Info: deployment.apps/rook-ceph-mon-c scaled

Info: fetching the deployment rook-ceph-mon-a to be running

Info: deployment rook-ceph-mon-a exists

Info: setting maintenance command to main container
Info: deployment rook-ceph-mon-a scaled down

Info: waiting for the deployment pod rook-ceph-mon-a-6849c8548f-p9lwb to be deleted

Info: ensure the maintenance deployment rook-ceph-mon-a is scaled up

Info: waiting for pod with label "ceph_daemon_type=mon,ceph_daemon_id=a" in namespace "openshift-storage" to be running
Info: waiting for pod with label "ceph_daemon_type=mon,ceph_daemon_id=a" in namespace "openshift-storage" to be running
Info: waiting for pod with label "ceph_daemon_type=mon,ceph_daemon_id=a" in namespace "openshift-storage" to be running
Info: pod rook-ceph-mon-a-maintenance-56bd9c6cfb-w9kjs is ready for maintenance operations
Info: fetching the deployment rook-ceph-mon-a-maintenance to be running

Info: deployment rook-ceph-mon-a-maintenance exists

Info: Started maintenance pod, restoring the mon quorum in the maintenance pod
Info: Extracting the monmap

# Lengthy rocksdb output removed

Info: Finished updating the monmap!
Info: Printing final monmap
monmaptool: monmap file /tmp/monmap
epoch 4
fsid d9a71ab9-a497-4da2-af36-53dcd410a5e9
last_changed 2024-03-01T07:24:34.567176+0000
created 2024-03-01T07:23:45.031715+0000
min_mon_release 17 (quincy)
election_strategy: 1
0: v2:172.30.124.246:3300/0 mon.a
Info: Restoring the mons in the rook-ceph-mon-endpoints configmap to the good mon
Info: Stopping the maintenance pod for mon a.

Info: fetching the deployment rook-ceph-mon-a-maintenance to be running

Info: deployment rook-ceph-mon-a-maintenance exists

Info: removing maintenance mode from deployment rook-ceph-mon-a-maintenance

Info: Successfully deleted maintenance deployment and restored deployment "rook-ceph-mon-a"
Info: Check that the restored mon is responding
Error: failed to get the status of ceph cluster. failed to run command. failed to run command. command terminated with exit code 1
Info: 1: waiting for ceph status to confirm single mon quorum.

Info: current ceph status output

Info: sleeping for 5 seconds

Info: finished waiting for ceph status   cluster:
    id:     d9a71ab9-a497-4da2-af36-53dcd410a5e9
    health: HEALTH_OK

  services:
    mon: 1 daemons, quorum a (age 18s)
    mgr: a(active, since 19m)
    mds: 1/1 daemons up, 1 hot standby
    osd: 3 osds: 3 up (since 18m), 3 in (since 18m)

  data:
    volumes: 1/1 healthy
    pools:   4 pools, 113 pgs
    objects: 31 objects, 582 KiB
    usage:   32 MiB used, 1.5 TiB / 1.5 TiB avail
    pgs:     113 active+clean

  io:
    client:   853 B/s rd, 1 op/s rd, 0 op/s wr



Info: Purging the bad mons []

Info: Mon quorum was successfully restored to mon a

Info: Only a single mon is currently running
Info: Enter 'continue' to start the operator and expand to full mon quorum again
```

After reviewing that the cluster is healthy with a single mon, press Enter to continue:

```bash
Info: continuing
Info: proceeding with resorting quorum
```
