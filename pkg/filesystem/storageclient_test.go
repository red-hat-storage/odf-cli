package filesystem

import (
	"context"
	"testing"

	csiv1 "github.com/ceph/ceph-csi-operator/api/v1"
	ocsclientv1alpha1 "github.com/red-hat-storage/ocs-client-operator/api/v1alpha1"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kubefake "k8s.io/client-go/kubernetes/fake"
	ctrl "sigs.k8s.io/controller-runtime/pkg/client"
	ctrlfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	_ = csiv1.AddToScheme(s)
	_ = ocsclientv1alpha1.AddToScheme(s)
	_ = corev1.AddToScheme(s)
	return s
}

func TestResolveStorageClientConfig(t *testing.T) {
	ctx := context.Background()
	storageClientName := "test-storage-client"
	namespace := "openshift-storage"
	storageClientUID := types.UID("test-uid-12345")
	secretName := "ceph-csi-secret" //nolint:gosec // test fixture, not a real credential
	connName := "ceph-connection"
	radosNS := "rados-ns-test"

	tests := []struct {
		name         string
		setupFn      func() (*k8sutil.Clientsets, ctrl.Client)
		expectErr    bool
		expectErrMsg string
		validateFn   func(*testing.T, *StorageClientConfig)
	}{
		{
			name: "success",
			setupFn: func() (*k8sutil.Clientsets, ctrl.Client) {
				sc := &ocsclientv1alpha1.StorageClient{
					ObjectMeta: metav1.ObjectMeta{Name: storageClientName, UID: storageClientUID},
				}
				profile := &csiv1.ClientProfile{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-profile", Namespace: namespace,
						OwnerReferences: []metav1.OwnerReference{
							{Kind: "StorageClient", UID: storageClientUID, APIVersion: "ocs.openshift.io/v1alpha1", Name: storageClientName},
						},
					},
					Spec: csiv1.ClientProfileSpec{
						CephConnectionRef: corev1.LocalObjectReference{Name: connName},
						CephFs: &csiv1.CephFsConfigSpec{
							SubVolumeGroup: "csi",
							RadosNamespace: &radosNS,
							CephCsiSecrets: &csiv1.CephCsiSecretsSpec{
								ControllerPublishSecret: corev1.SecretReference{Name: secretName},
							},
						},
					},
				}
				conn := &csiv1.CephConnection{
					ObjectMeta: metav1.ObjectMeta{Name: connName, Namespace: namespace},
					Spec:       csiv1.CephConnectionSpec{Monitors: []string{"10.0.0.1:6789"}},
				}
				fakeCtrl := ctrlfake.NewClientBuilder().WithScheme(newScheme()).WithObjects(sc, profile, conn).Build()

				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: namespace},
					Data:       map[string][]byte{"userID": []byte("client.admin"), "userKey": []byte("AQD1234567890==")},
				}
				pod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "csi-cephfsplugin-controller-0", Namespace: namespace,
						Labels: map[string]string{"app": namespace + ".cephfs.csi.ceph.com-ctrlplugin"},
					},
				}
				return &k8sutil.Clientsets{Kube: kubefake.NewSimpleClientset(secret, pod)}, fakeCtrl
			},
			validateFn: func(t *testing.T, cfg *StorageClientConfig) {
				assert.Equal(t, "csi", cfg.SVG)
				assert.Equal(t, "rados-ns-test", cfg.RadosNamespace)
				assert.NotNil(t, cfg.ExecConfig)
				assert.Equal(t, "csi-cephfsplugin-controller-0", cfg.ExecConfig.PodName)
				assert.Equal(t, namespace, cfg.ExecConfig.PodNamespace)
				assert.Equal(t, "csi-cephfsplugin", cfg.ExecConfig.Container)
				assert.Equal(t, "10.0.0.1:6789", cfg.ExecConfig.MonIP)
				assert.Equal(t, "client.admin", cfg.ExecConfig.UserID)
				assert.Equal(t, "AQD1234567890==", cfg.ExecConfig.UserKey)
			},
		},
		{
			name: "storage client not found",
			setupFn: func() (*k8sutil.Clientsets, ctrl.Client) {
				return &k8sutil.Clientsets{Kube: kubefake.NewSimpleClientset()},
					ctrlfake.NewClientBuilder().WithScheme(newScheme()).Build()
			},
			expectErr:    true,
			expectErrMsg: "failed to get StorageClient",
		},
		{
			name: "no matching client profile",
			setupFn: func() (*k8sutil.Clientsets, ctrl.Client) {
				sc := &ocsclientv1alpha1.StorageClient{
					ObjectMeta: metav1.ObjectMeta{Name: storageClientName, UID: storageClientUID},
				}
				profile := &csiv1.ClientProfile{
					ObjectMeta: metav1.ObjectMeta{
						Name: "other-profile", Namespace: namespace,
						OwnerReferences: []metav1.OwnerReference{
							{Kind: "StorageClient", UID: "different-uid", APIVersion: "ocs.openshift.io/v1alpha1", Name: "other"},
						},
					},
				}
				return &k8sutil.Clientsets{Kube: kubefake.NewSimpleClientset()},
					ctrlfake.NewClientBuilder().WithScheme(newScheme()).WithObjects(sc, profile).Build()
			},
			expectErr:    true,
			expectErrMsg: "no ClientProfile found",
		},
		{
			name: "missing SVG",
			setupFn: func() (*k8sutil.Clientsets, ctrl.Client) {
				sc := &ocsclientv1alpha1.StorageClient{
					ObjectMeta: metav1.ObjectMeta{Name: storageClientName, UID: storageClientUID},
				}
				profile := &csiv1.ClientProfile{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-profile", Namespace: namespace,
						OwnerReferences: []metav1.OwnerReference{
							{Kind: "StorageClient", UID: storageClientUID, APIVersion: "ocs.openshift.io/v1alpha1", Name: storageClientName},
						},
					},
					Spec: csiv1.ClientProfileSpec{
						CephConnectionRef: corev1.LocalObjectReference{Name: connName},
						CephFs:            &csiv1.CephFsConfigSpec{RadosNamespace: &radosNS},
					},
				}
				return &k8sutil.Clientsets{Kube: kubefake.NewSimpleClientset()},
					ctrlfake.NewClientBuilder().WithScheme(newScheme()).WithObjects(sc, profile).Build()
			},
			expectErr:    true,
			expectErrMsg: "has no cephFs.subVolumeGroup",
		},
		{
			name: "missing secret credentials",
			setupFn: func() (*k8sutil.Clientsets, ctrl.Client) {
				sc := &ocsclientv1alpha1.StorageClient{
					ObjectMeta: metav1.ObjectMeta{Name: storageClientName, UID: storageClientUID},
				}
				profile := &csiv1.ClientProfile{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-profile", Namespace: namespace,
						OwnerReferences: []metav1.OwnerReference{
							{Kind: "StorageClient", UID: storageClientUID, APIVersion: "ocs.openshift.io/v1alpha1", Name: storageClientName},
						},
					},
					Spec: csiv1.ClientProfileSpec{
						CephConnectionRef: corev1.LocalObjectReference{Name: connName},
						CephFs: &csiv1.CephFsConfigSpec{
							SubVolumeGroup: "csi",
							RadosNamespace: &radosNS,
							CephCsiSecrets: &csiv1.CephCsiSecretsSpec{
								ControllerPublishSecret: corev1.SecretReference{Name: secretName},
							},
						},
					},
				}
				conn := &csiv1.CephConnection{
					ObjectMeta: metav1.ObjectMeta{Name: connName, Namespace: namespace},
					Spec:       csiv1.CephConnectionSpec{Monitors: []string{"10.0.0.1:6789"}},
				}
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: namespace},
					Data:       map[string][]byte{"userKey": []byte("AQD1234567890==")},
				}
				pod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "csi-cephfsplugin-controller-0", Namespace: namespace,
						Labels: map[string]string{"app": namespace + ".cephfs.csi.ceph.com-ctrlplugin"},
					},
				}
				fakeCtrl := ctrlfake.NewClientBuilder().WithScheme(newScheme()).WithObjects(sc, profile, conn).Build()
				return &k8sutil.Clientsets{Kube: kubefake.NewSimpleClientset(secret, pod)}, fakeCtrl
			},
			expectErr:    true,
			expectErrMsg: "missing userID or userKey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientsets, ctrlClient := tt.setupFn()
			cfg, err := ResolveStorageClientConfig(ctx, clientsets, ctrlClient, storageClientName, namespace)
			if tt.expectErr {
				assert.Error(t, err)
				if tt.expectErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectErrMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.validateFn != nil {
					tt.validateFn(t, cfg)
				}
			}
		})
	}
}
