package filesystem

import (
	"context"
	"fmt"
	"slices"

	csiv1 "github.com/ceph/ceph-csi-operator/api/v1"
	ocsclientv1alpha1 "github.com/red-hat-storage/ocs-client-operator/api/v1alpha1"
	"github.com/rook/kubectl-rook-ceph/pkg/filesystem"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime/pkg/client"
)

type StorageClientConfig struct {
	SVG            string
	RadosNamespace string
	ExecConfig     *filesystem.CustomExecConfig
}

func ResolveStorageClientConfig(ctx context.Context, clientsets *k8sutil.Clientsets, ctrlClient ctrl.Client, storageClientName, namespace string) (*StorageClientConfig, error) {
	// 1. Get StorageClient (cluster-scoped) to find its UID
	var storageClient ocsclientv1alpha1.StorageClient
	if err := ctrlClient.Get(ctx, types.NamespacedName{Name: storageClientName}, &storageClient); err != nil {
		return nil, fmt.Errorf("failed to get StorageClient %q: %v", storageClientName, err)
	}

	// 2. Find ClientProfile owned by this StorageClient via ownerReferences
	var profileList csiv1.ClientProfileList
	if err := ctrlClient.List(ctx, &profileList, ctrl.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list ClientProfiles in namespace %q: %v", namespace, err)
	}

	idx := slices.IndexFunc(profileList.Items, func(p csiv1.ClientProfile) bool {
		return slices.ContainsFunc(p.GetOwnerReferences(), func(ref metav1.OwnerReference) bool {
			return ref.Kind == "StorageClient" && ref.UID == storageClient.UID
		})
	})
	if idx == -1 {
		return nil, fmt.Errorf("no ClientProfile found owned by StorageClient %q", storageClientName)
	}
	profile := &profileList.Items[idx]

	// 3. Extract CephFS config from ClientProfile
	if profile.Spec.CephFs == nil {
		return nil, fmt.Errorf("ClientProfile %q has no cephFs config", profile.Name)
	}
	svg := profile.Spec.CephFs.SubVolumeGroup
	if svg == "" {
		return nil, fmt.Errorf("ClientProfile %q has no cephFs.subVolumeGroup", profile.Name)
	}
	if profile.Spec.CephFs.RadosNamespace == nil || *profile.Spec.CephFs.RadosNamespace == "" {
		return nil, fmt.Errorf("ClientProfile %q has no cephFs.radosNamespace", profile.Name)
	}
	radosNS := *profile.Spec.CephFs.RadosNamespace

	connName := profile.Spec.CephConnectionRef.Name
	if connName == "" {
		return nil, fmt.Errorf("ClientProfile %q has no cephConnectionRef", profile.Name)
	}

	if profile.Spec.CephFs.CephCsiSecrets == nil || profile.Spec.CephFs.CephCsiSecrets.ControllerPublishSecret.Name == "" {
		return nil, fmt.Errorf("ClientProfile %q has no CephFS provisioner secret", profile.Name)
	}
	secretName := profile.Spec.CephFs.CephCsiSecrets.ControllerPublishSecret.Name

	// 4. Get CephConnection to find mon IPs
	var cephConn csiv1.CephConnection
	if err := ctrlClient.Get(ctx, types.NamespacedName{Name: connName, Namespace: namespace}, &cephConn); err != nil {
		return nil, fmt.Errorf("failed to get CephConnection %q: %v", connName, err)
	}
	if len(cephConn.Spec.Monitors) == 0 {
		return nil, fmt.Errorf("CephConnection %q has no monitors", connName)
	}

	// 5. Get Secret for auth credentials
	secret, err := clientsets.Kube.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %q: %v", secretName, err)
	}
	userID := string(secret.Data["userID"])
	userKey := string(secret.Data["userKey"])
	if userID == "" || userKey == "" {
		return nil, fmt.Errorf("secret %q missing userID or userKey", secretName)
	}

	// 6. Find CephFS CSI controller plugin pod
	label := fmt.Sprintf("app=%s.cephfs.csi.ceph.com-ctrlplugin", namespace)
	pods, err := clientsets.Kube.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, fmt.Errorf("failed to list cephfs ctrlplugin pods: %v", err)
	}
	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no cephfs ctrlplugin pod found with label %q", label)
	}

	return &StorageClientConfig{
		SVG:            svg,
		RadosNamespace: radosNS,
		ExecConfig: &filesystem.CustomExecConfig{
			PodName:      pods.Items[0].Name,
			PodNamespace: namespace,
			Container:    "csi-cephfsplugin",
			MonIP:        cephConn.Spec.Monitors[0],
			UserID:       userID,
			UserKey:      userKey,
		},
	}, nil
}
