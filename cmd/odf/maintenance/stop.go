package maintenance

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/maintenance"
	"github.com/spf13/cobra"
)

var stopMaintenanceCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stops the maintenance deployment",
	Args:    cobra.ExactArgs(1),
	Example: "odf debug stop <deployment_name>",
	Run: func(cmd *cobra.Command, args []string) {
		maintenance.StopMaintenance(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, args[0])
	},
}
