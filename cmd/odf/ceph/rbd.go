package ceph

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/rook/kubectl-rook-ceph/pkg/rbd"
	"github.com/spf13/cobra"
)

// RbdCmd represents the rbd command
var RbdCmd = &cobra.Command{
	Use:                "rbd",
	Short:              "call a 'rbd' CLI command with arbitrary args",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
}

var listCmdRbd = &cobra.Command{
	Use:   "ls",
	Short: "Print the list of rbd images.",
	Run: func(cmd *cobra.Command, args []string) {
		// verify operator pod is running
		_, err := k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator")
		if err != nil {
			logging.Fatal(err)
		}
		rbd.ListImages(cmd.Context(), root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace)
	},
}

func init() {
	RbdCmd.AddCommand(listCmdRbd)
}
