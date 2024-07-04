package maintenance

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
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
		k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator")
	},
}

func init() {
	MaintenanceCmd.AddCommand(startMaintenanceCmd)
	MaintenanceCmd.AddCommand(stopMaintenanceCmd)
}
