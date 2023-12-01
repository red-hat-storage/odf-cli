# Set

The set command supports the following sub-commands:

- [Set](#set)
  - [recovery-profile](#recovery-profile)
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

## log-level

This command will set the log level for different subsystems. Supported subsystems:

* osd
* mds
* mon
* mgr
* auth

``` bash
odf set log-level osd 20
```

Once the logging efforts are complete, restore the systems to their default or to a level suitable for normal operations.

