package get

import (
	"fmt"
	"os"
	"strings"
	"time"

	ocsv1 "github.com/red-hat-storage/ocs-operator/api/v4/v1"
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

const healthSeparatorWidth = 72

var clusterHealth = &cobra.Command{
	Use:     "health",
	Short:   "print cluster health, and include more details with --verbose",
	Long:    `Print the status of each StorageCluster and the CRs it manages, then the cluster health report (Ceph MON, MGR, and OSD checks, pod status, placement groups, and object storage).`,
	Args:    cobra.NoArgs,
	Example: "odf get health",
	PreRun: func(cmd *cobra.Command, args []string) {
		if _, err := k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, rookOperatorNamespace(cmd), "app=rook-ceph-operator"); err != nil {
			logging.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, _ []string) {
		namespace := storageClusterNamespace(cmd)

		clusters, err := odfhealth.ListStorageClusters(cmd.Context(), root.CtrlClient, namespace)
		if err != nil {
			if namespace != "" {
				logging.Fatal(fmt.Errorf("failed to list StorageClusters in namespace %q: %w", namespace, err))
			}
			logging.Fatal(fmt.Errorf("failed to list StorageClusters in all namespaces: %w", err))
		}
		if len(clusters) == 0 {
			if namespace != "" {
				logging.Fatal(fmt.Errorf("no StorageCluster found in namespace %q", namespace))
			}
			logging.Fatal(fmt.Errorf("no StorageCluster found"))
		}

		operatorNS := rookOperatorNamespace(cmd)
		// json/yaml reports are the rook check results only. The CR table is text.
		textOutput := healthOutput == "" || healthOutput == "text"

		namespaces := odfhealth.Namespaces(clusters)
		if textOutput {
			printClusterHealthHeader(namespaces)
		}
		for i, ns := range namespaces {
			nsClusters := odfhealth.ClustersInNamespace(clusters, ns)
			if textOutput {
				if i > 0 {
					logging.Plain("")
				}
				printCRHealthSection(cmd, nsClusters)
			}

			customChecks := []health.HealthChecker{
				&odfhealth.NooBaaHealthChecker{
					Ctx:              cmd.Context(),
					DynamicClient:    root.ClientSets.Dynamic,
					K8sClient:        root.ClientSets.Kube,
					ClusterNamespace: ns,
				},
			}
			health.Health(cmd.Context(), root.ClientSets, operatorNS, ns, healthVerbose, "", healthOutput, customChecks, false)
		}
	},
}

// storageClusterNamespace is the -n/--namespace value when the user passed it.
// An empty string means -n was not set, so StorageClusters are listed in all namespaces.
func storageClusterNamespace(cmd *cobra.Command) string {
	flagNamespace := cmd.Flag("namespace")
	if flagNamespace != nil && flagNamespace.Changed {
		return flagNamespace.Value.String()
	}
	return ""
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

// printClusterHealthHeader prints the CLUSTER HEALTH REPORT title once, before
// the StorageCluster tables. Health() is called with printHeader set to false
// so this block is not printed again.
func printClusterHealthHeader(namespaces []string) {
	logging.Plain("%s", strings.Repeat("=", healthSeparatorWidth))
	logging.Plain("CLUSTER HEALTH REPORT")
	logging.Plain("%s", strings.Repeat("=", healthSeparatorWidth))
	logging.Plain("Generated: %s", time.Now().UTC().Format("2006-01-02 15:04:05 UTC"))
	logging.Plain("Namespace: %s", strings.Join(namespaces, ", "))
	logging.Plain("")
}

func printCRHealthSection(cmd *cobra.Command, clusters []ocsv1.StorageCluster) {
	logging.Plain("%s", strings.Repeat("=", healthSeparatorWidth))
	logging.Plain("StorageCluster CRs")
	logging.Plain("%s", strings.Repeat("=", healthSeparatorWidth))
	logging.Plain("")
	if err := odfhealth.PrintCRHealthClusters(cmd.Context(), root.CtrlClient, clusters, os.Stdout); err != nil {
		logging.Fatal(err)
	}
	logging.Plain("")
	logging.Plain("%s", strings.Repeat("=", healthSeparatorWidth))
	logging.Plain("")
}

func init() {
	clusterHealth.Flags().BoolVar(&healthVerbose, "verbose", false, "include more details in the cluster health report")
	clusterHealth.Flags().StringVarP(&healthOutput, "output", "o", "text", "output format: text, json, yaml")
}
