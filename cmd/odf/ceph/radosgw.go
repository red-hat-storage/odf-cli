package ceph

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

// RadosgwCmd represents the radosgw command
var RadosgwCmd = &cobra.Command{
	Use:                "radosgw-admin",
	Short:              "call a 'radosgw-admin' CLI command",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logging.Info("running 'radosgw-admin' command with args: %v", args)
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
