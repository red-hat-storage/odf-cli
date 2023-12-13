package root

import (
	"context"
	"fmt"

	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	rookclient "github.com/rook/rook/pkg/client/clientset/versioned"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	KubeConfig              string
	KubeContext             string
	OperatorNamespace       string
	StorageClusterNamespace string
	ClientSets              *k8sutil.Clientsets
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
		ClientSets = getClientsets(cmd.Context())
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

func getClientsets(ctx context.Context) *k8sutil.Clientsets {
	var err error

	clientsets := &k8sutil.Clientsets{}

	congfigOverride := &clientcmd.ConfigOverrides{}
	if KubeContext != "" {
		congfigOverride = &clientcmd.ConfigOverrides{CurrentContext: KubeContext}
	}

	// 1. Create Kubernetes Client
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		congfigOverride,
	)

	clientsets.KubeConfig, err = kubeconfig.ClientConfig()
	if err != nil {
		logging.Fatal(err)
	}

	clientsets.Rook, err = rookclient.NewForConfig(clientsets.KubeConfig)
	if err != nil {
		logging.Fatal(err)
	}

	clientsets.Kube, err = k8s.NewForConfig(clientsets.KubeConfig)
	if err != nil {
		logging.Fatal(err)
	}

	preValidationCheck(ctx, clientsets, OperatorNamespace, StorageClusterNamespace)

	return clientsets
}

func preValidationCheck(ctx context.Context, k8sclientset *k8sutil.Clientsets, operatorNamespace, storageClusterNamespace string) {
	_, err := k8sclientset.Kube.CoreV1().Namespaces().Get(ctx, operatorNamespace, v1.GetOptions{})
	if err != nil {
		logging.Fatal(fmt.Errorf("Operator namespace '%s' does not exist. %v", operatorNamespace, err))
	}
	_, err = k8sclientset.Kube.CoreV1().Namespaces().Get(ctx, storageClusterNamespace, v1.GetOptions{})
	if err != nil {
		logging.Fatal(fmt.Errorf("StorageCluster namespace '%s' does not exist. %v", storageClusterNamespace, err))
	}
}
