package rook

import (
	"fmt"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:     "restart",
	Short:   "Restart rook-ceph-operator pod",
	Args:    cobra.NoArgs,
	Example: "odf operator rook restart",
	PreRun: func(cmd *cobra.Command, args []string) {
		// verify operator pod is running
		if _, err := k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator"); err != nil {
			logging.Fatal(fmt.Errorf("failed to wait for rook-ceph-operator pod: %v", err))
		}
	},
	Run: func(cmd *cobra.Command, _ []string) {
		k8sutil.RestartDeployment(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "rook-ceph-operator")
	},
}
