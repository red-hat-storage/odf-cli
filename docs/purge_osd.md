# purge-osd

The purge-osd command permanently removes an OSD from the cluster.

If the OSD is not safe to remove, a prompt will confirm if you are sure the OSD should be removed.

```bash
odf purge-osd 0
# Error: command terminated with exit code 16
# Error EBUSY: OSD(s) 0 have 12 pgs currently mapped to them.
# Warning: Are you sure you want to purge osd.0? The OSD is *not* safe to destroy. This may lead to data loss. If you are sure the OSD should be purged, enter 'yes-force-destroy-osd'
# yes-force-destroy-osd
# Info: Running purge osd command
# 2023-12-08 16:29:23.381399 W | cephcmd: loaded admin secret from env var ROOK_CEPH_SECRET instead of from file
# 2023-12-08 16:29:23.381917 I | rookcmd: starting Rook v1.13.0-alpha.0.37.gb51577011 with arguments 'rook ceph osd remove --osd-ids=0 --force-osd-removal=true'
# 2023-12-08 16:29:23.381945 I | rookcmd: flag values: --force-osd-removal=true, --help=false, --log-level=INFO, --osd-ids=0, --preserve-pvc=false
# 2023-12-08 16:29:23.381955 I | ceph-spec: parsing mon endpoints: a=10.101.19.56:6789
# 2023-12-08 16:29:23.389943 I | cephclient: writing config file /var/lib/rook/rook-ceph/rook-ceph.config
# 2023-12-08 16:29:23.390412 I | cephclient: generated admin config in /var/lib/rook/rook-ceph
# 2023-12-08 16:29:23.853067 I | cephosd: validating status of osd.0
# 2023-12-08 16:29:23.853122 I | cephosd: osd.0 is marked 'DOWN'
# 2023-12-08 16:29:24.245057 I | cephosd: marking osd.0 out
# 2023-12-08 16:29:25.258200 I | cephosd: osd.0 is NOT ok to destroy but force removal is enabled so proceeding with removal
# 2023-12-08 16:29:25.262919 I | cephosd: removing the OSD deployment "rook-ceph-osd-0"
# 2023-12-08 16:29:25.263207 I | op-k8sutil: removing deployment rook-ceph-osd-0 if it exists
# 2023-12-08 16:29:25.280243 I | op-k8sutil: Removed deployment rook-ceph-osd-0
# 2023-12-08 16:29:25.301057 I | op-k8sutil: "rook-ceph-osd-0" still found. waiting...
# 2023-12-08 16:29:27.305646 I | op-k8sutil: confirmed rook-ceph-osd-0 does not exist
# 2023-12-08 16:29:27.305684 I | cephosd: did not find a pvc name to remove for osd "rook-ceph-osd-0"
# 2023-12-08 16:29:27.305689 I | cephosd: purging osd.0
# 2023-12-08 16:29:27.749319 I | cephosd: attempting to remove host "minikube-m02" from crush map if not in use
# 2023-12-08 16:29:28.935667 I | cephosd: removed CRUSH host "minikube-m02"
# 2023-12-08 16:29:29.322622 I | cephosd: no ceph crash to silence
# 2023-12-08 16:29:29.322741 I | cephosd: completed removal of OSD 0
```
