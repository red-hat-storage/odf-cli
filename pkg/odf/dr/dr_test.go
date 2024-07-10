package dr

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	submarinerv1alpha1 "github.com/submariner-io/submariner-operator/api/v1alpha1"
	submarinerv1 "github.com/submariner-io/submariner/pkg/apis/submariner.io/v1"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	ctrl "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

func Test_isSubmarinerEnabled(t *testing.T) {
	var scheme = runtime.NewScheme()
	err := submarinerv1alpha1.AddToScheme(scheme)
	assert.NoError(t, err)

	ctx := context.TODO()

	client := ctrl.NewClientBuilder().WithScheme(scheme).Build()
	submarinerCR, err := isSubmarinerEnabled(ctx, client)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(submarinerCR, submarinerv1alpha1.Submariner{}))

	err = client.Create(ctx, &submarinerv1alpha1.Submariner{ObjectMeta: metav1.ObjectMeta{Name: "submariner", Namespace: "submariner-operator"}})
	assert.NoError(t, err)

	submarinerCR, err = isSubmarinerEnabled(ctx, client)
	assert.NoError(t, err)
	assert.False(t, reflect.DeepEqual(submarinerCR, submarinerv1alpha1.Submariner{}))
}

func Test_isGlobalnetRequired(t *testing.T) {
	var scheme = runtime.NewScheme()
	err := submarinerv1.AddToScheme(scheme)
	assert.NoError(t, err)

	ctx := context.TODO()

	client := ctrl.NewClientBuilder().WithScheme(scheme).Build()
	res, err := isGlobalnetRequired(ctx, client, "cluster1", "cluster1")
	assert.Error(t, err)
	assert.False(t, res)

	clusterCR := &submarinerv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster1",
			Namespace: submarinerOperatorNamespace,
		},
		Spec: submarinerv1.ClusterSpec{
			ServiceCIDR: []string{"172.30.0.0/16"},
			ClusterCIDR: []string{"10.128.0.0/14"},
		},
	}
	peerClusterCR := &submarinerv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster2",
			Namespace: submarinerOperatorNamespace,
		},
		Spec: submarinerv1.ClusterSpec{
			ServiceCIDR: []string{"172.40.0.0/16"},
			ClusterCIDR: []string{"10.138.0.0/14"},
		},
	}
	assert.NoError(t, client.Create(ctx, clusterCR))
	assert.NoError(t, client.Create(ctx, peerClusterCR))

	res, err = isGlobalnetRequired(ctx, client, "cluster1", "cluster2")
	assert.NoError(t, err)
	assert.False(t, res)

	peerClusterWithOverlappingIP := &submarinerv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster3",
			Namespace: submarinerOperatorNamespace,
		},
		Spec: submarinerv1.ClusterSpec{
			ServiceCIDR: []string{"172.30.0.0/16"},
			ClusterCIDR: []string{"10.128.0.0/14"},
		},
	}
	assert.NoError(t, client.Create(ctx, peerClusterWithOverlappingIP))

	res, err = isGlobalnetRequired(ctx, client, "cluster1", "cluster3")
	assert.NoError(t, err)
	assert.True(t, res)

	res, err = isGlobalnetRequired(ctx, client, "cluster2", "cluster3")
	assert.NoError(t, err)
	assert.False(t, res)

	client = ctrl.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
		Get: func(ctx context.Context, client ctrlclient.WithWatch, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
			return errors.NewNotFound(submarinerv1.Resource("clusters"), "cluster1")
		},
	}).WithObjects(clusterCR, peerClusterCR, peerClusterWithOverlappingIP).WithScheme(scheme).Build()
	res, err = isGlobalnetRequired(ctx, client, "cluster1", "cluster3")
	assert.Error(t, err)
	assert.False(t, res)
}

func Test_isGlobalnetEnabled(t *testing.T) {
	submarinerCR := submarinerv1alpha1.Submariner{}
	assert.False(t, isGlobalnetEnabled(submarinerCR))

	submarinerCR.Status.GlobalnetDaemonSetStatus = submarinerv1alpha1.DaemonSetStatusWrapper{
		Status: &appsv1.DaemonSetStatus{
			NumberAvailable: 0,
		},
	}
	assert.False(t, isGlobalnetEnabled(submarinerCR))

	submarinerCR.Status.GlobalnetDaemonSetStatus.Status.NumberAvailable = 1
	assert.True(t, isGlobalnetEnabled(submarinerCR))
}
