package rook

import (
	"context"
	"fmt"

	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
)

func SetCephLogLevel(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNamespace, storageClusterNamespace string, component, subsystem, level string) {

	cephArgs := []string{"config", "set", component, fmt.Sprintf("debug_%s", subsystem), level}

	_, err := exec.RunCommandInOperatorPod(ctx, clientsets, "ceph", cephArgs, operatorNamespace, storageClusterNamespace, false)
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to run ceph command with args %v. %v", cephArgs, err))
	}
	logging.Info("successfully set the log levels for %q subsystem %q as %q", component, subsystem, level)
}
