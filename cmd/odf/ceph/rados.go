package ceph

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

// RadosCmd represents the rados command
var RadosCmd = &cobra.Command{
	Use:                "rados",
	Short:              "call a 'rados' CLI command with arbitrary args",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logging.Info("running 'rados' command with args: %v", args)
		_, err := exec.RunCommandInOperatorPod(cmd.Context(), root.ClientSets, cmd.Use, args, root.OperatorNamespace, root.StorageClusterNamespace, false)
		if err != nil {
			logging.Fatal(err)
		}
	},
}
