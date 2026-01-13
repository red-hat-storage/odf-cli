package maintenance

import (
	"fmt"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

// MaintenanceCmd represents the maintenance commands
var MaintenanceCmd = &cobra.Command{
	Use:                "maintenance",
	Short:              "Perform maintenance operation on mons and OSDs deployment by scaling it down and creating a maintenance deployment.",
	DisableFlagParsing: true,
	Args:               cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// verify operator pod is running
		if _, err := k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator"); err != nil {
			logging.Fatal(fmt.Errorf("failed to wait for rook-ceph-operator pod: %v", err))
		}
	},
}

func init() {
	MaintenanceCmd.AddCommand(startMaintenanceCmd)
	MaintenanceCmd.AddCommand(stopMaintenanceCmd)
}
