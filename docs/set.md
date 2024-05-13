# Set

The set command supports the following sub-commands:

- [Set](#set)
  - [recovery-profile](#recovery-profile)
  - [full-ratio](#full-ratio)
  - [backfill-ratio](#backfillfull-ratio)
  - [nearfull-ratio](#nearfull-ratio)
  - [ceph](#ceph)
    - [log-level](#log-level)

## recovery-profile

This command will set the recovery profile to favor either the new client IO, recovery IO, or a balanced mode. The default is `balanced` mode.
We can use the following built-in profile types:

* balanced
* high_client_ops
* high_recovery_ops

```bash
odf set recovery-profile high_client_ops
```

To verify the recovery profile setting run [odf get recovery-profile](get.md#recovery-profile).

## full

This command allows updating the ceph osd full ratio in case ceph prevents the I/O operation on OSDs that reached the capacity specified. The default is 0.85.

**Note**: If the value is set too close to 1.0, the cluster will be unrecoverable if the OSDs are full and there is nowhere to grow.

``` bash
odf set full 0.9
```

## backfillfull

This command allows updating the ceph osd backfillfull ratio in case ceph will deny backfilling to the OSD that reached the capacity specified. The default value is 0.80.

**Note**: If the value is set too close to 1.0, the OSDs will be full and the cluster will not able to backfill.

``` bash
odf set backfillfull 0.85
```

## nearfull

This command allows updating the ceph osd nearfull ratio in case Ceph returns the nearfull osds message when the cluster reaches the capacity specified. The default value is 0.75

**Note**: If the value is set too close to 1.0, the OSDs will be full and the cluster will not able to backfill.

``` bash
odf set nearfull 0.8
```

## ceph

The `ceph` command helps update the Ceph configuration.

### log-level

This command will set the log level for different ceph [subsystems](https://docs.ceph.com/en/latest/rados/troubleshooting/log-and-debug/#ceph-subsystems).
The `debug_` prefix will be automatically added to the subsystem when enabling the logging.

``` bash
odf set ceph log-level osd crush 20
```

Once the logging efforts are complete, restore the systems to their default or to a level suitable for normal operations.


