# Operator Commands

Operator is parent command which requires sub-command. Currently, operator supports these sub-commands:

- `set <property> <value>` : [set](#set) the property in the rook-ceph-operator-config configmap.
- `restart`: [restart](#restart) the Rook-Ceph operator

## set

Set the property in the rook-ceph-operator-config configmap

```bash
$ odf operator rook set ROOK_LOG_LEVEL DEBUG

configmap/rook-ceph-operator-config patched
```

## restart

Restart the Rook-Ceph operator.

``` bash
$ odf operator rook restart

deployment.apps/rook-ceph-operator restarted
```
