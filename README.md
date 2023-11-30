# odf-cli

The ODF CLI tool provides configuration and troubleshooting commands for OpenShift Data Foundation.

## Commands

- `odf set recovery-profile <profile>`: Set the recovery profile to favor new IO, recovery, or balanced mode with options `high_client_ops`, `high_recovery_ops`, or `balanced`. The default is `balanced`.
- `odf get recovery-profile`: Get the recovery profile value.
- `odf help` : Display help text

## Documentation

Visit docs below for complete details about each command and their flags uses.

- [set](docs/set.md)
- [get](docs/get.md)

### Root args

These are the arguments that apply to all commands:

- `-h|--help`: this will print brief command help text.

    ```bash
    odf -h
    ```

- `-n|--namespace='openshift-storage'`: the Openshift namespace in which the StorageCluster resides. (optional,  default: openshift-storage)

    ```bash
    odf -n test-cluster [commands]
    ```

- `-o|--operator-namespace` : the Openshift namespace in which the rook operator resides, when the arg `-n` is passed but `-o` is not then `-o` will equal to the `-n`. (default: openshift-storage)

    ```bash
    odf --operator-namespace test-operator -n test-cluster [commands]
    ```

- `--context`: the name of the Openshift context to be used (optional).

    ```bash
    odf --context=$(oc config current-context) [commands]
    ```
