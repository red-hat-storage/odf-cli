package get

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/health"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

var clusterHealth = &cobra.Command{
	Use:                "health",
	Short:              "check health of the cluster and common configuration issues",
	DisableFlagParsing: true,
	Args:               cobra.NoArgs,
	Example:            "odf get health",
	PreRun: func(cmd *cobra.Command, args []string) {
		// verify operator pod is running
		if _, err := k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator"); err != nil {
			logging.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, _ []string) {
		health.Health(cmd.Context(), root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace)
	},
}
