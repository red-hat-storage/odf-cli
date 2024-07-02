package root

import (
	"context"
	"fmt"

	ocsv1 "github.com/red-hat-storage/ocs-operator/api/v4/v1"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	rookclient "github.com/rook/rook/pkg/client/clientset/versioned"
	"github.com/spf13/cobra"
	submarinerv1alpha1 "github.com/submariner-io/submariner-operator/api/v1alpha1"
	submarinerv1 "github.com/submariner-io/submariner/pkg/apis/submariner.io/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	KubeConfig              string
	KubeContext             string
	OperatorNamespace       string
	StorageClusterNamespace string
	ClientSets              *k8sutil.Clientsets
	CtrlClient              ctrl.Client
	scheme                  = runtime.NewScheme()
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
		CtrlClient = getControllerRuntimeClient()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	if err := ocsv1.AddToScheme(scheme); err != nil {
		logging.Fatal(err)
	}
	if err := submarinerv1.AddToScheme(scheme); err != nil {
		logging.Fatal(err)
	}
	if err := submarinerv1alpha1.AddToScheme(scheme); err != nil {
		logging.Fatal(err)
	}

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

	clientsets.Dynamic, err = dynamic.NewForConfig(clientsets.KubeConfig)
	if err != nil {
		logging.Fatal(err)
	}

	preValidationCheck(ctx, clientsets, OperatorNamespace, StorageClusterNamespace)

	return clientsets
}

func getControllerRuntimeClient() ctrl.Client {
	client, err := ctrl.New(ClientSets.KubeConfig, ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		logging.Fatal(err)
	}
	return client
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
