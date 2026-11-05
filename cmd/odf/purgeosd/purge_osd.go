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
	Long: `Safely remove an OSD from the cluster with interactive prompts.

This command walks through OSD removal step-by-step:
1. Selects an unhealthy OSD from the cluster
2. Verifies LSO operator is scaled down
3. Checks cluster health
4. Gathers OSD details (PVC, PV, device)
5. Stops the OSD process
6. Marks OSD out and waits for data migration
7. Shows dry-run of resources to be deleted
8. Executes purge via Rook backend

The operation is idempotent - you can re-run it if interrupted.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		osdID := ceph.SelectUnhealthyOSD(ctx, root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace)
		osdIDStr := strconv.Itoa(osdID)

		// LSO operator must be scaled down before purging
		// to prevent LSO from reprovisioning the disk we are about to wipe
		if !ceph.CheckLSOScaledDown(ctx, root.ClientSets) {
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
		if details.PVCName != "" {
			fmt.Printf("PVC %q is now Pending (awaiting disk).\n", details.PVCName)
		}
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("Replace the physical disk and scale the LSO operator back up.")
		fmt.Println("If reusing the same disk: scale the LSO operator back up.")
	},
}
