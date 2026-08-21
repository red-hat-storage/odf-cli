package get

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	odfhealth "github.com/red-hat-storage/odf-cli/pkg/health"
	"github.com/rook/kubectl-rook-ceph/pkg/health"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

var (
	healthVerbose bool
	healthOutput  string
)

var clusterHealth = &cobra.Command{
	Use:     "health",
	Short:   "check health of the cluster and common configuration issues",
	Args:    cobra.NoArgs,
	Example: "odf get health",
	PreRun: func(cmd *cobra.Command, args []string) {
		// verify operator pod is running
		if _, err := k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator"); err != nil {
			logging.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, _ []string) {
		customChecks := []health.HealthChecker{
			&odfhealth.NooBaaHealthChecker{
				Ctx:              cmd.Context(),
				DynamicClient:    root.ClientSets.Dynamic,
				K8sClient:        root.ClientSets.Kube,
				ClusterNamespace: root.StorageClusterNamespace,
			},
		}
		health.Health(cmd.Context(), root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, healthVerbose, "", healthOutput, customChecks)
	},
}

func init() {
	clusterHealth.Flags().BoolVar(&healthVerbose, "verbose", false, "shows detailed check for pods")
	clusterHealth.Flags().StringVarP(&healthOutput, "output", "o", "text", "output format: text, json, yaml")
}
