package set

import (
	"context"

	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
)

func SetProfile(ctx context.Context, clientsets *k8sutil.Clientsets, recoveryOption, operatorNamespace, storageClusterNamespace string) {

	cephArgs := []string{"config", "set", "osd", "osd_mclock_profile", recoveryOption}

	exec.RunCommandInOperatorPod(ctx, clientsets, "ceph", cephArgs, operatorNamespace, storageClusterNamespace, true, false)
}
