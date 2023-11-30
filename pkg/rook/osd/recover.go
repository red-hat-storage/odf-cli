package osd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
)

type Config struct {
	OsdMclockProfile OsdConfig `json:"osd_mclock_profile"`
}

type OsdConfig struct {
	Value string `json:"value"`
}

func GetProfile(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNamespace, storageClusterNamespace string) {
	cephArgs := []string{"config", "get", "osd", "--format", "json"}

	out := exec.RunCommandInOperatorPod(ctx, clientsets, "ceph", cephArgs, operatorNamespace, storageClusterNamespace, true, true)

	var config Config
	err := json.Unmarshal([]byte(out), &config)
	if err != nil {
		logging.Error(err)
		return
	}
	fmt.Println(config.OsdMclockProfile.Value)
}

func SetProfile(ctx context.Context, clientsets *k8sutil.Clientsets, recoveryOption, operatorNamespace, storageClusterNamespace string) {

	cephArgs := []string{"config", "set", "osd", "osd_mclock_profile", recoveryOption}

	exec.RunCommandInOperatorPod(ctx, clientsets, "ceph", cephArgs, operatorNamespace, storageClusterNamespace, false, true)
}
