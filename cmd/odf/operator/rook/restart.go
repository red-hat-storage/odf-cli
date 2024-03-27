package rook

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:     "restart",
	Short:   "Restart rook-ceph-operator pod",
	Args:    cobra.NoArgs,
	Example: "odf operator rook restart",
	PreRun: func(cmd *cobra.Command, args []string) {
		// verify operator pod is running
		k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator")
	},
	Run: func(cmd *cobra.Command, _ []string) {
		k8sutil.RestartDeployment(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "rook-ceph-operator")
	},
}
