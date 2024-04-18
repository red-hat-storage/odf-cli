package osd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
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

	out, err := exec.RunCommandInOperatorPod(ctx, clientsets, "ceph", cephArgs, operatorNamespace, storageClusterNamespace, true)
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to run command ceph with args %v: %v", cephArgs, err))
	}

	var config Config
	err = json.Unmarshal([]byte(out), &config)
	if err != nil {
		logging.Error(err)
		return
	}
	fmt.Println(config.OsdMclockProfile.Value)
}

func SetProfile(ctx context.Context, clientsets *k8sutil.Clientsets, recoveryOption, operatorNamespace, storageClusterNamespace string) {

	cephArgs := []string{"config", "set", "osd", "osd_mclock_profile", recoveryOption}

	_, err := exec.RunCommandInOperatorPod(ctx, clientsets, "ceph", cephArgs, operatorNamespace, storageClusterNamespace, false)
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to run ceph command with args %v. %v", cephArgs, err))
	}
}

type SafeToDestroyStatus struct {
	SafeToDestroy []int `json:"safe_to_destroy"`
}

func SafeToDestroy(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNamespace, storageClusterNamespace, osdID string) (bool, error) {
	args := []string{"osd", "safe-to-destroy", osdID}
	out, err := exec.RunCommandInOperatorPod(ctx, clientsets, "ceph", args, operatorNamespace, storageClusterNamespace, true)
	if err != nil {
		return false, errors.Wrapf(err, string(out))
	}

	var safeToDestroy SafeToDestroyStatus
	if err = json.Unmarshal([]byte(out), &safeToDestroy); err != nil {
		return false, errors.Wrapf(err, string(out))
	}
	if len(safeToDestroy.SafeToDestroy) != 0 && fmt.Sprint(safeToDestroy.SafeToDestroy[0]) == osdID {
		return true, nil
	}
	return false, nil
}
