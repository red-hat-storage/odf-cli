package purgeosd

import (
	"fmt"
	"strconv"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/red-hat-storage/odf-cli/pkg/ceph"
	"github.com/rook/kubectl-rook-ceph/pkg/rook"
	"github.com/spf13/cobra"
)

var CephPurgeOsdCmd = &cobra.Command{
	Use:   "purge-osd",
	Short: "Permanently remove an OSD from the cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		osdID := ceph.SelectUnhealthyOSD(ctx, root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace)
		osdIDStr := strconv.Itoa(osdID)

		// LSO operator must be scaled down before purging
		// to prevent LSO from reprovisioning the disk we are about to wipe
		stopPurge, isLSO := ceph.CheckLSOScaledDown(ctx, root.ClientSets)
		if stopPurge {
			fmt.Println("Purge cancelled.")
			return
		}

		// show cluster health before making any changes
		// a pre-existing HEALTH_WARN or HEALTH_ERR may indicate that safeToDestroy
		// will not return true, or that the cluster is already under stress.
		if !ceph.ClusterHealthCheck(ctx, root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace) {
			fmt.Println("Purge cancelled.")
			return
		}

		// gather OSD details from k8s (deployment, PVC, PV) before any deletion
		details := ceph.GatherOSDDetails(ctx, root.ClientSets, root.StorageClusterNamespace, osdID)

		if !ceph.StopOSDIfRunning(ctx, root.ClientSets, root.StorageClusterNamespace, details) {
			fmt.Println("Purge cancelled.")
			return
		}

		if !ceph.MarkOSDOut(ctx, root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, osdID, osdIDStr) {
			fmt.Println("Purge cancelled.")
			return
		}

		// returns true if force removal should be used (user chose not to wait)
		forceRemoval := ceph.DataMigrationPrompt(ctx, root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, osdIDStr, osdID)

		ceph.PrintDryRun(details)

		if !ceph.PromptYN("Proceed with purge?") {
			fmt.Println("Purge cancelled.")
			return
		}

		// execute the purge via the Rook operator.
		// on failure, the process exits. on success, it returns and we print the result below.
		forceFlagValue := "false"
		if forceRemoval {
			forceFlagValue = "true"
		}
		rook.PurgeOsds(ctx, root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, osdIDStr, forceFlagValue)

		fmt.Println()
		fmt.Println("Purge complete.")
		fmt.Println("A new PVC is pending (awaiting disk).")
		fmt.Println()

		if isLSO {
			fmt.Println("Next steps:")
			fmt.Println("Replace the physical disk and enable diskmaker-manager, then scale the LSO operator back up.")
			fmt.Println("If reusing the same disk: enable diskmaker-manager and scale the LSO operator back up.")
		}
	},
}
