# odf-cli

The ODF CLI tool provides configuration and troubleshooting commands for OpenShift Data Foundation.

## Usage

`odf [root args] [command]`

### Root args

These are args currently supported:

1. `-h|--help`: this will print brief command help text.

    ```bash
    odf -h
    ```

2. `-n|--namespace='openshift-storage'`: the Openshift namespace in which the StorageCluster resides. (optional,  default: openshift-storage)

    ```bash
    odf -n test-cluster [commands]
    ```

3. `-o|--operator-namespace` : the Openshift namespace in which the rook operator resides, when the arg `-n` is passed but `-o` is not then `-o` will equal to the `-n`. (default: openshift-storage)

    ```bash
    odf --operator-namespace test-operator -n test-cluster [commands]
    ```

4. `--context`: the name of the Openshift context to be used (optional).

    ```bash
    odf --context=$(oc config current-context) [commands]
    ```

## Commands

- `set <args>`: Set the recovery profile to favor new IO, recovery, or balanced mode with options high_client_ops, high_recovery_ops, or balanced
- `help` : Output help text

## Documentation

Visit docs below for complete details about each command and their flags uses.

1. [Modify the mclock recovery profile](docs/set.md#recovery-profile)

## Examples

### Command to set mclock recovery profile

```code
odf -n openshift-storage set recovery-profile high_client_ops
```

To verify the mclock recovery profile settings

> ```text
> ceph config get osd
> WHO     MASK  LEVEL     OPTION                   VALUE            RO
> global        basic     log_to_file              false
> global        advanced  mon_allow_pool_delete    true
> global        advanced  mon_allow_pool_size_one  true
> global        advanced  mon_cluster_log_file
> osd           advanced  osd_mclock_profile       high_client_ops
>```
