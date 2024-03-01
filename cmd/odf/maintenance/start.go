package maintenance

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/maintenance"
	"github.com/spf13/cobra"
)

var startMaintenanceCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start a maintenance deployment with an optional alternative ceph container image",
	Args:    cobra.ExactArgs(1),
	Example: "odf maintenance start <deployment_name>",
	Run: func(cmd *cobra.Command, args []string) {
		alternateImage := cmd.Flag("alternate-image").Value.String()
		maintenance.StartMaintenance(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, args[0], alternateImage)
	},
}

func init() {
	startMaintenanceCmd.Flags().String("alternate-image", "", "To create deployment with alternate image")
}
