# Multus Validation

The `multus validation` command verifies that the Multus and system configuration is compatible with Rook. Run it **before** installing Rook to ensure network compatibility.

See [Validating Multus configuration](https://rook.github.io/docs/rook/latest/CRDs/Cluster/network-providers/?h=valid#validating-multus-configuration) for more details.

## Subcommands

1. `config`: Generate a sample validation config file
2. `run`: Run Multus validation tests to verify network connectivity
3. `cleanup`: Clean up resources from a previous validation run

### Generate a config

The preferred way is to use a config file. Generate a sample (see `odf multus validation config --help` for presets), edit it for the deployment, then run validation with it:

```bash
odf multus validation config dedicated-storage-nodes > multus-validation-config.yaml
```

Edit `multus-validation-config.yaml` as needed, ex: set `publicNetwork` and `clusterNetwork` to the Multus network attachment names, adjust namespace, node placement etc. Then run:

```bash
odf multus validation run --config multus-validation-config.yaml
```

### Running with flags(Not Recommended)

To run without a config file, pass all options on the command line. Default namespace is `openshift-storage`:

```bash
odf multus validation run \
  --namespace openshift-storage \
  --public-network [namespace]/public-net \
  --cluster-network [namespace]/cluster-net
```

For all options:

```bash
odf multus validation run --help
```

## Cleanup

To remove leftover resources from a previous validation run:

```bash
odf multus validation cleanup --help
```
