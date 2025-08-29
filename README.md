# odf-cli

The ODF CLI tool provides configuration and troubleshooting commands for OpenShift Data Foundation.

## Commands

- `odf set`:
  - `recovery-profile <profile>`: Set the recovery profile to favor new IO, recovery, or balanced mode with options `high_client_ops`, `high_recovery_ops`, or `balanced`. The default is `balanced`.
  - `full <ratio>`: Update the ceph osd full ratio after which Ceph automatically prevents any I/O operations on OSDs.
  - `nearfull <ratio>`: Update the ceph osd nearfull ratio in case Ceph returns the nearfull osds message when the cluster reaches the capacity specified.
  - `backfillfull <ratio>`:  Update the ceph osd backfillfull ratio in case ceph will deny backfilling to the OSD that reached the capacity specified.
  - `ceph`
    - `log-level <daemon> <subsystem> <log-level>`: Set the log level for Ceph daemons like OSD, mon, mds etc. More information about the ceph subsystems can be found [here](https://docs.ceph.com/en/latest/rados/troubleshooting/log-and-debug/#ceph-subsystems)
- `odf get`:
  - `recovery-profile`: Get the recovery profile value.
  - `health`: Check health of the cluster and common configuration issues.
  - `dr-health [ceph status args]`: Print the ceph status of a peer cluster in a mirroring-enabled environment thereby validating connectivity between ceph clusters. Ceph status args can be optionally passed, such as to change the log level: --debug-ms 1.
  - `dr-prereq <PeerManagedClusterName>` : Print the status of pre-requisites for Disaster Recovery between peer clusters.
  - `mon-endpoints`: Print mon endpoints.
- `odf purge-osd <ID>`: Permanently remove an OSD from the cluster.
- `odf maintenance`: [Perform maintenance operations](docs/maintenance.md) on mons or OSDs. The mon or OSD deployment will be scaled down and replaced temporarily by a maintenance deployment.
  - `start <deployment-name>`
    - `[--alternate-image <alternate-image>]` : (optional) Start a maintenance deployment with an optional alternative ceph container image
  - `stop <deployment-name>`: Stop the maintenance deployment and restore the mon or OSD deployment
- `odf subvolume`:
  - `ls`: Display all the subvolumes
  - `delete <subvolume> <filesystem> <subvolumegroup>`: Deletes the stale subvolumes
- `odf operator`:
  - `rook`:
    - `set`: Set the property in the rook-ceph-operator-config configmap.
    - `restart` : Restart the Rook-Ceph operator
- `odf restore`:
  - `mon-quorum`: Restore the mon quorum based on a single healthy mon since quorum was lost with the other mons
  - `deleted`: Restore the ceph resources which are stuck in deleting state due to underlying resources being present in the cluster
- `odf ceph <args>` : Run a Ceph CLI command. Supports any arguments the ceph command supports. See [Ceph](https://docs.ceph.com/en/pacific/start/intro/) docs for more.
- `odf rados <args>` : Run a Rados CLI command. Supports any arguments the rados command supports. See [Rados](https://docs.ceph.com/en/latest/man/8/rados/) docs for more.
- `odf radosgw-admin <args>` : Run an RGW CLI command. Supports any arguments the [radosgw-admin](https://docs.ceph.com/en/latest/man/8/radosgw-admin/) command supports. See the radosgw-admin docs for more.
- `odf rbd <args>` : Call a 'rbd' CLI command with arbitrary args
- `odf dr`:
  - `init`: Create configuration file for `odf dr` sub-commands.
  - `test`:
    - `run`: Run disaster recovery test with a tiny application
    - `clean`: Clean up after running tests
  - `validate`:
    - `application`: Detect problems in disaster recovery protected application
    - `clusters`: Detect problems in disaster recovery clusters
  - `gather`:
    - `application`: Collect diagnostic data from your clusters
- `odf noobaa:` Run a noobaa CLI command. Supports all the noobaa cli arguments available in odf.
- `odf help` : Display help text

## Documentation

Visit docs below for complete details about each command and their flags uses.

- [set](docs/set.md)
- [get](docs/get.md)
- [purge-osd](docs/purge_osd.md)
- [mon](docs/mons.md)
- [maintenance](docs/maintenance.md)
- [operator](docs/operator.md)
- [ceph](docs/ceph.md)
- [rbd](docs/ceph.md#rbd)
- [rados](docs/ceph.md#rados)
- [radosgw-admin](docs/ceph.md#radosgw-admin)
- [dr](docs/dr.md)

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

## Installation

### Build and run from source

1. Clone the repository

    ```bash
    git clone https://github.com/red-hat-storage/odf-cli.git
    ```

2. Change the directory and build the binary

    ```bash
    cd odf-cli/
    make
    ```

3. Use the binary present in the `bin/` directory to run the commands

    ```bash
    ./bin/odf -h
    ```

### Download nightly build

If you want to test the latest nightly build you can use the
download-odf-nightly.sh script:

```bash
scripts/download-odf-nightly.sh v4.19
```

The command downloads the latest 4.19 build image and extracts the odf
executable to the current directory.
