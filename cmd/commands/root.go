package command

import (
	"github.com/spf13/cobra"
)

var (
	KubeConfig              string
	KubeContext             string
	OperatorNamespace       string
	StorageClusterNamespace string
)

// RootCmd represents the odf command
var RootCmd = &cobra.Command{
	Use:              "odf",
	Short:            "Management and troubleshooting tools for ODF clusters.",
	TraverseChildren: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if StorageClusterNamespace != "" && OperatorNamespace == "" {
			OperatorNamespace = StorageClusterNamespace
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	// Hide autocompletion command
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	RootCmd.PersistentFlags().StringVar(&KubeConfig, "kubeconfig", "", "Openshift config path")
	RootCmd.PersistentFlags().StringVar(&OperatorNamespace, "operator-namespace", "", "Openshift namespace where the ODF operator is running")
	RootCmd.PersistentFlags().StringVarP(&StorageClusterNamespace, "namespace", "n", "openshift-storage", "Openshift namespace where the StorageCluster CR is created")
	RootCmd.PersistentFlags().StringVar(&KubeContext, "context", "", "Openshift context to use")
}
