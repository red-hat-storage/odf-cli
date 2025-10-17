# Ceph Commands

Commands in this doc are an alternate to running the toolbox pod.

## ceph

Run the `ceph` cli command with with arbitrary args.

### Example

```bash
odf ceph status

#   cluster:
#     id:     b74c18dd-6ee3-44fe-90b5-ed12feac46a4
#     health: HEALTH_OK
#
#   services:
#     mon: 3 daemons, quorum a,b,c (age 62s)
#     mgr: a(active, since 23s)
#     osd: 1 osds: 1 up (since 12s), 1 in (since 30s)
#
#   data:
#     pools:   0 pools, 0 pg
#     objects: 0 objects, 0 B
#     usage:   0 B used, 0 B / 0 B avail
#     pgs:
```

## ceph daemon

Run the Ceph daemon command by connecting to its admin socket.

```bash
odf ceph daemon osd.0 dump_historic_ops

Info: running 'ceph' command with args: [daemon osd.0 dump_historic_ops]
{
    "size": 20,
    "duration": 600,
    "ops": []
}
```

## rados

Run any `rados` cli command with with arbitrary args.

### Example

```bash
odf rados df

# POOL_NAME     USED  OBJECTS  CLONES  COPIES  MISSING_ON_PRIMARY  UNFOUND  DEGRADED  RD_OPS      RD  WR_OPS       WR  USED COMPR  UNDER COMPR
# .mgr       452 KiB        2       0       2                   0        0         0      22  18 KiB      14  262 KiB         0 B          0 B

# total_objects    2
# total_used       27 MiB
# total_avail      10 GiB
# total_space      10 GiB
```

## rbd

Run any `rbd` cli command with with arbitrary args.

### Example

```bash
odf rbd ls replicapool

# csi-vol-427774b4-340b-11ed-8d66-0242ac110004
```

## radosgw-admin

Run any `radosgw-admin` cli command with arbitrary args.

### Example

```bash
odf radosgw-admin user create --display-name="my user" --uid=myuser

# Info: running 'radosgw-admin' command with args: [user create --display-name=my-user --uid=myuser]
# {
#     "user_id": "myuser",
#     "display_name": "my user",
#    ...
#    ...
#    ...
# }
```
