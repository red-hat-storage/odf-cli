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
	Short:   "print StorageCluster-managed CR status, or detailed pod checks with --verbose",
	Args:    cobra.NoArgs,
	Example: "odf get health",
	PreRun: func(cmd *cobra.Command, args []string) {
		// verify operator pod is running
		if _, err := k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator"); err != nil {
			logging.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, _ []string) {
		namespace := ""
		namespaceSet := false
		if f := cmd.Flag("namespace"); f != nil && f.Changed {
			namespace = f.Value.String()
			namespaceSet = true
		}

		if healthVerbose {
			customChecks := []health.HealthChecker{
				&odfhealth.NooBaaHealthChecker{
					Ctx:              cmd.Context(),
					DynamicClient:    root.ClientSets.Dynamic,
					K8sClient:        root.ClientSets.Kube,
					ClusterNamespace: root.StorageClusterNamespace,
				},
			}
			health.Health(cmd.Context(), root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, healthVerbose, "", healthOutput, customChecks)
			// json/yaml reports are the rook check results only. The CR table is text.
			if healthOutput == "" || healthOutput == "text" {
				printCRHealthSection(cmd, namespace, namespaceSet)
			}
			return
		}

		if err := odfhealth.PrintCRHealth(cmd.Context(), root.CtrlClient, namespace, namespaceSet, os.Stdout); err != nil {
			logging.Fatal(err)
		}
	},
}

func printCRHealthSection(cmd *cobra.Command, namespace string, namespaceSet bool) {
	logging.Plain("")
	logging.Plain("%s", strings.Repeat("=", 72))
	logging.Plain("StorageCluster CRs")
	logging.Plain("%s", strings.Repeat("=", 72))
	logging.Plain("")
	if err := odfhealth.PrintCRHealth(cmd.Context(), root.CtrlClient, namespace, namespaceSet, os.Stdout); err != nil {
		logging.Fatal(err)
	}
}

func init() {
	clusterHealth.Flags().BoolVar(&healthVerbose, "verbose", false, "shows detailed check for pods")
	clusterHealth.Flags().StringVarP(&healthOutput, "output", "o", "text", "output format: text, json, yaml")
}
