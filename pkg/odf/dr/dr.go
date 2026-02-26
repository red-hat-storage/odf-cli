package dr

import (
	"context"
	"fmt"
	"reflect"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	submarinerv1alpha1 "github.com/submariner-io/submariner-operator/api/v1alpha1"
	submarinerv1 "github.com/submariner-io/submariner/pkg/apis/submariner.io/v1"
	"github.com/submariner-io/submariner/pkg/cidr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime/pkg/client"
)

const submarinerOperatorNamespace string = "submariner-operator"

func GetDRPrerequisite(ctx context.Context, peerClusterID string) {
	submarinerCR, err := isSubmarinerEnabled(ctx, root.CtrlClient)
	if err != nil {
		logging.Fatal(err)
	}

	if reflect.DeepEqual(submarinerCR, submarinerv1alpha1.Submariner{}) {
		logging.Info("Submariner is not installed.")
	} else {
		logging.Info("Submariner is installed.")
	}

	localClusterID := submarinerCR.Status.ClusterID
	globalnetRequired, err := isGlobalnetRequired(ctx, root.CtrlClient, localClusterID, peerClusterID)
	if err != nil {
		logging.Fatal(err)
	}
	if globalnetRequired {
		logging.Info("Globalnet is required.")
	} else {
		logging.Info("Globalnet is not required.")
	}

	globalnetEnabled := isGlobalnetEnabled(submarinerCR)
	if globalnetEnabled {
		logging.Info("Globalnet is enabled.")
	} else if !globalnetEnabled && globalnetRequired {
		logging.Warning("Globalnet is not enabled, but is required.")
	} else {
		logging.Info("Globalnet is not enabled.")
	}
}

func isSubmarinerEnabled(ctx context.Context, client ctrl.Client) (submarinerv1alpha1.Submariner, error) {
	submarinerCR := submarinerv1alpha1.Submariner{}
	err := client.Get(ctx, types.NamespacedName{Name: "submariner", Namespace: submarinerOperatorNamespace}, &submarinerCR)
	if err != nil {
		// These errors mean that Submariner is not installed.
		// IsNoMatchError -> Submariner CRD is not available.
		// IsNotFound -> Submariner CR is not created.
		if meta.IsNoMatchError(err) || errors.IsNotFound(err) {
			return submarinerCR, nil
		}
		return submarinerCR, err
	}
	return submarinerCR, nil
}

func isGlobalnetRequired(ctx context.Context, client ctrl.Client, clusterID, peerClusterID string) (bool, error) {
	if clusterID == peerClusterID {
		return false, fmt.Errorf("current ClusterID and peer ClusterID refer to the same cluster. Provide a different peer ClusterID")
	}

	clusterCR := &submarinerv1.Cluster{}
	err := client.Get(ctx, types.NamespacedName{Name: clusterID, Namespace: submarinerOperatorNamespace}, clusterCR)
	if err != nil {
		return false, err
	}

	peerClusterCR := &submarinerv1.Cluster{}
	err = client.Get(ctx, types.NamespacedName{Name: peerClusterID, Namespace: submarinerOperatorNamespace}, peerClusterCR)
	if err != nil {
		return false, err
	}

	err = cidr.OverlappingSubnets(clusterCR.Spec.ServiceCIDR, clusterCR.Spec.ClusterCIDR,
		append(peerClusterCR.Spec.ClusterCIDR, peerClusterCR.Spec.ServiceCIDR...))
	if err != nil {
		return true, nil
	}

	return false, nil
}

func isGlobalnetEnabled(submarinerCR submarinerv1alpha1.Submariner) bool {
	if submarinerCR.Status.GlobalnetDaemonSetStatus.Status == nil {
		return false
	}
	return submarinerCR.Status.GlobalnetDaemonSetStatus.Status.NumberAvailable > 0
}
