package health

import (
	"bytes"
	"context"
	"testing"

	ocsv1 "github.com/red-hat-storage/ocs-operator/api/v4/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestNestedStringOrDash(t *testing.T) {
	phaseReady := map[string]interface{}{
		"status": map[string]interface{}{
			"phase":   "Ready",
			"message": "failed to create pool",
		},
	}
	assert.Equal(t, "Ready", nestedStringOrDash(phaseReady, "status", "phase"))
	assert.Equal(t, "failed to create pool", nestedStringOrDash(phaseReady, "status", "message"))

	emptyMessage := map[string]interface{}{
		"status": map[string]interface{}{
			"phase":   "Ready",
			"message": "",
		},
	}
	assert.Equal(t, missingValue, nestedStringOrDash(emptyMessage, "status", "message"))

	noStatus := map[string]interface{}{
		"metadata": map[string]interface{}{"name": "pool"},
	}
	assert.Equal(t, missingValue, nestedStringOrDash(noStatus, "status", "phase"))
	assert.Equal(t, missingValue, nestedStringOrDash(noStatus, "status", "message"))

	phaseOnly := map[string]interface{}{
		"status": map[string]interface{}{
			"phase": "Progressing",
		},
	}
	assert.Equal(t, "Progressing", nestedStringOrDash(phaseOnly, "status", "phase"))
	assert.Equal(t, missingValue, nestedStringOrDash(phaseOnly, "status", "message"))
}

func TestConditionMessage(t *testing.T) {
	withReady := map[string]interface{}{
		"status": map[string]interface{}{
			"conditions": []interface{}{
				map[string]interface{}{"type": "Progressing", "message": ""},
				map[string]interface{}{"type": "Available", "message": "NooBaa is ready"},
			},
		},
	}
	assert.Equal(t, "NooBaa is ready", conditionMessage(withReady))

	emptyConditions := map[string]interface{}{
		"status": map[string]interface{}{
			"conditions": []interface{}{
				map[string]interface{}{"type": "Available", "message": ""},
			},
		},
	}
	assert.Equal(t, missingValue, conditionMessage(emptyConditions))

	noConditions := map[string]interface{}{
		"status": map[string]interface{}{"phase": "Ready"},
	}
	assert.Equal(t, missingValue, conditionMessage(noConditions))
}

func TestRowFields(t *testing.T) {
	pool := unstructured.Unstructured{Object: map[string]interface{}{
		"metadata": map[string]interface{}{"name": "ocs-storagecluster-cephblockpool"},
		"status": map[string]interface{}{
			"phase":   "Failure",
			"message": "failed to create pool",
		},
	}}
	spec := kindSpec{component: "Block pool", observedPath: []string{"status", "phase"}}
	component, name, observed, message := rowFields(spec, pool)
	assert.Equal(t, "Block pool", component)
	assert.Equal(t, "ocs-storagecluster-cephblockpool", name)
	assert.Equal(t, "Failure", observed)
	assert.Equal(t, "failed to create pool", message)

	consumer := unstructured.Unstructured{Object: map[string]interface{}{
		"metadata": map[string]interface{}{"name": "internal"},
		"status":   map[string]interface{}{"state": "Ready"},
	}}
	consumerSpec := kindSpec{component: "Consumer", observedPath: []string{"status", "state"}}
	_, _, observed, message = rowFields(consumerSpec, consumer)
	assert.Equal(t, "Ready", observed)
	assert.Equal(t, missingValue, message)
}

func TestMatchingStorageClients(t *testing.T) {
	items := []unstructured.Unstructured{
		{Object: map[string]interface{}{"metadata": map[string]interface{}{"name": "ocs-storagecluster"}}},
		{Object: map[string]interface{}{"metadata": map[string]interface{}{"name": "other-client"}}},
	}
	got := matchingStorageClients(items, "ocs-storagecluster")
	require.Len(t, got, 1)
	assert.Equal(t, "ocs-storagecluster", got[0].GetName())

	assert.Empty(t, matchingStorageClients(items, "missing"))
}

func TestClusterMode(t *testing.T) {
	internal := ocsv1.StorageCluster{}
	assert.Equal(t, "internal", clusterMode(internal))

	external := ocsv1.StorageCluster{}
	external.Spec.ExternalStorage.Enable = true
	assert.Equal(t, "external", clusterMode(external))
}

func TestListStorageClustersSkipsIgnored(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, ocsv1.AddToScheme(scheme))

	active := &ocsv1.StorageCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "ocs-storagecluster", Namespace: "openshift-storage"},
		Status:     ocsv1.StorageClusterStatus{Phase: "Ready"},
	}
	ignored := &ocsv1.StorageCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "test-second-storagecluster", Namespace: "openshift-storage"},
		Status:     ocsv1.StorageClusterStatus{Phase: "Ignored"},
	}
	extended := &ocsv1.StorageCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "ocs-storagecluster", Namespace: "openshift-storage-extended"},
		Spec:       ocsv1.StorageClusterSpec{ExternalStorage: ocsv1.ExternalStorageClusterSpec{Enable: true}},
		Status:     ocsv1.StorageClusterStatus{Phase: "Ready"},
	}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(active, ignored, extended).Build()

	all, err := listStorageClusters(context.Background(), c, "", false)
	require.NoError(t, err)
	require.Len(t, all, 2)
	assert.Equal(t, "openshift-storage", all[0].Namespace)
	assert.Equal(t, "ocs-storagecluster", all[0].Name)
	assert.Equal(t, "openshift-storage-extended", all[1].Namespace)

	oneNS, err := listStorageClusters(context.Background(), c, "openshift-storage-extended", true)
	require.NoError(t, err)
	require.Len(t, oneNS, 1)
	assert.Equal(t, "openshift-storage-extended", oneNS[0].Namespace)
}

func TestPrintCRHealthHeader(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, ocsv1.AddToScheme(scheme))

	sc := &ocsv1.StorageCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "ocs-storagecluster", Namespace: "openshift-storage"},
		Status:     ocsv1.StorageClusterStatus{Phase: "Progressing"},
	}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(sc).Build()

	var buf bytes.Buffer
	err := PrintCRHealth(context.Background(), c, "", false, &buf)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "StorageCluster:")
	assert.Contains(t, out, "ocs-storagecluster")
	assert.Contains(t, out, "Namespace:")
	assert.Contains(t, out, "openshift-storage")
	assert.Contains(t, out, "Phase:")
	assert.Contains(t, out, "Progressing")
	assert.Contains(t, out, "Mode:")
	assert.Contains(t, out, "internal")
	assert.Contains(t, out, "COMPONENT")
	assert.Contains(t, out, "OBSERVED")
	assert.Contains(t, out, "STATUS MESSAGE")
}

func TestPrintCRHealthNoStorageCluster(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, ocsv1.AddToScheme(scheme))
	c := fake.NewClientBuilder().WithScheme(scheme).Build()

	err := PrintCRHealth(context.Background(), c, "", false, &bytes.Buffer{})
	require.EqualError(t, err, "no StorageCluster found")

	err = PrintCRHealth(context.Background(), c, "openshift-storage", true, &bytes.Buffer{})
	require.EqualError(t, err, `no StorageCluster found in namespace "openshift-storage"`)
}
