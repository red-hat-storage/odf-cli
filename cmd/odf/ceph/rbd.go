package ceph

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

// RbdCmd represents the rbd command
var RbdCmd = &cobra.Command{
	Use:                "rbd",
	Short:              "call a 'rbd' CLI command with arbitrary args",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logging.Info("running 'rbd' command with args: %v", args)
		// verify operator pod is running
		_, err := k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator")
		if err != nil {
			logging.Fatal(err)
		}

		_, err = exec.RunCommandInOperatorPod(cmd.Context(), root.ClientSets, cmd.Use, args, root.OperatorNamespace, root.StorageClusterNamespace, false)
		if err != nil {
			logging.Fatal(err)
		}
	},
}
