package get

import (
	"os"
	"strings"

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
	Short:   "print ODF CR status, or detailed status with --verbose flag",
	Long:    `Print the status of each StorageCluster and the CRs it manages. --verbose prints that table first, then the full cluster health report, including Ceph MON, MGR, and OSD checks, pod status, placement groups, and object storage.`,
	Args:    cobra.NoArgs,
	Example: "odf get health",
	PreRun: func(cmd *cobra.Command, args []string) {
		if _, err := k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, rookOperatorNamespace(cmd), "app=rook-ceph-operator"); err != nil {
			logging.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, _ []string) {
		namespace := ""
		namespaceSet := false
		if flagNamespace := cmd.Flag("namespace"); flagNamespace != nil && flagNamespace.Changed {
			namespace = flagNamespace.Value.String()
			namespaceSet = true
		}

		if healthVerbose {
			// json/yaml reports are the rook check results only. The CR table is text.
			if healthOutput == "" || healthOutput == "text" {
				printCRHealthSection(cmd, namespace, namespaceSet)
			}
			clusterNS := root.StorageClusterNamespace
			if namespaceSet {
				clusterNS = namespace
			}
			operatorNS := rookOperatorNamespace(cmd)
			customChecks := []health.HealthChecker{
				&odfhealth.NooBaaHealthChecker{
					Ctx:              cmd.Context(),
					DynamicClient:    root.ClientSets.Dynamic,
					K8sClient:        root.ClientSets.Kube,
					ClusterNamespace: clusterNS,
				},
			}
			health.Health(cmd.Context(), root.ClientSets, operatorNS, clusterNS, healthVerbose, "", healthOutput, customChecks)
			return
		}

		if err := odfhealth.PrintCRHealth(cmd.Context(), root.CtrlClient, namespace, namespaceSet, os.Stdout); err != nil {
			logging.Fatal(err)
		}
	},
}

// rookOperatorNamespace is where rook-ceph-operator runs.
// -n only filters StorageClusters; it is not the operator namespace unless
// --operator-namespace is also set.
func rookOperatorNamespace(cmd *cobra.Command) string {
	if f := cmd.Flag("operator-namespace"); f != nil && f.Changed {
		return root.OperatorNamespace
	}
	if f := cmd.Flag("namespace"); f != nil && f.Changed {
		return "openshift-storage"
	}
	return root.OperatorNamespace
}

func printCRHealthSection(cmd *cobra.Command, namespace string, namespaceSet bool) {
	logging.Plain("%s", strings.Repeat("=", 72))
	logging.Plain("StorageCluster CRs")
	logging.Plain("%s", strings.Repeat("=", 72))
	logging.Plain("")
	if err := odfhealth.PrintCRHealth(cmd.Context(), root.CtrlClient, namespace, namespaceSet, os.Stdout); err != nil {
		logging.Fatal(err)
	}
	logging.Plain("")
	logging.Plain("%s", strings.Repeat("=", 72))
	logging.Plain("")
}

func init() {
	clusterHealth.Flags().BoolVar(&healthVerbose, "verbose", false, "shows detailed check for pods")
	clusterHealth.Flags().StringVarP(&healthOutput, "output", "o", "text", "output format: text, json, yaml")
}
