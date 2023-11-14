# Set

The set command supports the following sub-commands:

* [recovery-profile](#recovery-profile)

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
