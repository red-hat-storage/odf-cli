package restore

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	pkgrestore "github.com/rook/kubectl-rook-ceph/pkg/restore"
	"github.com/spf13/cobra"
)

// deletedCmd represents the deleted command
var deletedCmd = &cobra.Command{
	Use:                "deleted",
	Short:              "Restores a CR that was accidentally deleted and is still in terminating state.",
	DisableFlagParsing: true,
	Args:               cobra.RangeArgs(1, 2),
	Example:            "odf restore deleted <CRD> [CRNAME]",
	PreRun: func(cmd *cobra.Command, args []string) {
		k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator")

	},
	Run: func(cmd *cobra.Command, args []string) {
		k8sutil.SetDeploymentScale(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "ocs-operator", 0)
		pkgrestore.RestoreCrd(cmd.Context(), root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, args)
		k8sutil.SetDeploymentScale(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "ocs-operator", 1)

	},
}
