package common

import (
	"context"
	"fmt"

	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
)

var ValidCephSubSystems = []string{"osd", "mds", "mon", "mgr", "auth"}

func SetLogLevel(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNamespace, storageClusterNamespace string, args []string) {

	cephArgs := []string{"config", "set", args[0], fmt.Sprintf("debug_%s", args[0]), args[1]}

	exec.RunCommandInOperatorPod(ctx, clientsets, "ceph", cephArgs, operatorNamespace, storageClusterNamespace, true, false)
}
