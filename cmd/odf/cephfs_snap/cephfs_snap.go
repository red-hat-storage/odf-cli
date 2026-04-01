package cephfs_snap

import (
	"fmt"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	odffilesystem "github.com/red-hat-storage/odf-cli/pkg/filesystem"
	"github.com/rook/kubectl-rook-ceph/pkg/filesystem"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

var CephFSSnapshotCmd = &cobra.Command{
	Use:   "cephfs-snap",
	Short: "Manages CephFS snapshots",
}

var snapshotListCmd = &cobra.Command{
	Use:     "ls",
	Short:   "Print the list of CephFS snapshots.",
	Example: "odf cephfs-snap ls --filesystem=ocs-storagecluster-cephfilesystem",
	Run: func(cmd *cobra.Command, args []string) {
		orphanedOnly, _ := cmd.Flags().GetBool("orphaned")
		svg, _ := cmd.Flags().GetString("svg")
		fs, _ := cmd.Flags().GetString("filesystem")
		radosNamespace, _ := cmd.Flags().GetString("rados-namespace")
		var execCfg *filesystem.CustomExecConfig

		storageClient, _ := cmd.Flags().GetString("storage-client")
		if storageClient != "" {
			cfg, err := odffilesystem.ResolveStorageClientConfig(cmd.Context(), root.ClientSets, root.CtrlClient, storageClient, root.StorageClusterNamespace)
			if err != nil {
				logging.Fatal(fmt.Errorf("failed to resolve storage client config: %v", err))
			}
			svg = cfg.SVG
			radosNamespace = cfg.RadosNamespace
			execCfg = cfg.ExecConfig
		}

		f := &filesystem.CephFilesystem{
			Ctx:               cmd.Context(),
			Clientsets:        root.ClientSets,
			OperatorNamespace: root.OperatorNamespace,
			ClusterNamespace:  root.StorageClusterNamespace,
			RadosNamespace:    radosNamespace,
			CustomExecConfig:  execCfg,
		}
		f.SnapshotList(svg, fs, orphanedOnly)
	},
}

var snapshotDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Deletes a CephFS snapshot.",
	Args:    cobra.ExactArgs(2),
	Example: "odf cephfs-snap delete <subvol> <snapshot> --filesystem=ocs-storagecluster-cephfilesystem",
	Run: func(cmd *cobra.Command, args []string) {
		subvol := args[0]
		snap := args[1]
		fs, _ := cmd.Flags().GetString("filesystem")
		svg, _ := cmd.Flags().GetString("svg")
		radosNamespace, _ := cmd.Flags().GetString("rados-namespace")
		var execCfg *filesystem.CustomExecConfig

		storageClient, _ := cmd.Flags().GetString("storage-client")
		if storageClient != "" {
			cfg, err := odffilesystem.ResolveStorageClientConfig(cmd.Context(), root.ClientSets, root.CtrlClient, storageClient, root.StorageClusterNamespace)
			if err != nil {
				logging.Fatal(fmt.Errorf("failed to resolve storage client config: %v", err))
			}
			svg = cfg.SVG
			radosNamespace = cfg.RadosNamespace
			execCfg = cfg.ExecConfig
		}

		f := &filesystem.CephFilesystem{
			Ctx:               cmd.Context(),
			Clientsets:        root.ClientSets,
			OperatorNamespace: root.OperatorNamespace,
			ClusterNamespace:  root.StorageClusterNamespace,
			RadosNamespace:    radosNamespace,
			CustomExecConfig:  execCfg,
		}
		f.SnapshotDelete(fs, subvol, snap, svg)
	},
}

func init() {
	CephFSSnapshotCmd.AddCommand(snapshotListCmd)
	snapshotListCmd.Flags().Bool("orphaned", false, "List only orphaned snapshots")
	CephFSSnapshotCmd.PersistentFlags().String("svg", "csi", "The name of the subvolume group")
	CephFSSnapshotCmd.PersistentFlags().String("filesystem", "ocs-storagecluster-cephfilesystem", "The name of the CephFS filesystem")
	CephFSSnapshotCmd.PersistentFlags().String("rados-namespace", "csi", "The rados namespace for omap operations")
	CephFSSnapshotCmd.PersistentFlags().String("storage-client", "", "StorageClient CR name to auto-resolve SVG, rados-namespace, and exec config")
	CephFSSnapshotCmd.AddCommand(snapshotDeleteCmd)
}
