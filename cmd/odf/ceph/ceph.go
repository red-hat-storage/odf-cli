package ceph

import (
	"fmt"
	"strings"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

// CephCmd represents the ceph command
var CephCmd = &cobra.Command{
	Use:                "ceph",
	Short:              "call a 'ceph' CLI command with arbitrary args",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logging.Info("running 'ceph' command with args: %v", args)
		// verify operator pod is running
		_, err := k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator")
		if err != nil {
			logging.Fatal(err)
		}
		if args[0] == "daemon" {
			if len(args) < 2 {
				logging.Fatal(fmt.Errorf("too few arguments to run the 'ceph daemon' command"))
			}
			cephDaemonNameAndID := splitCephDaemonNameAndID(args[1])

			podLabel := fmt.Sprintf("%s=%s", cephDaemonNameAndID[0], cephDaemonNameAndID[1])
			_, err = exec.RunCommandInLabeledPod(cmd.Context(), root.ClientSets, podLabel, cephDaemonNameAndID[0], cmd.Use, args, root.StorageClusterNamespace, false)
			if err != nil {
				logging.Fatal(err)
			}

		} else {
			_, err = exec.RunCommandInOperatorPod(cmd.Context(), root.ClientSets, cmd.Use, args, root.OperatorNamespace, root.StorageClusterNamespace, false)
			if err != nil {
				logging.Fatal(err)
			}
		}
	},
}

// splitCephDaemonNameAndID splits the arg based on '.' so that we use it with label, example osd.0 to osd and 0
func splitCephDaemonNameAndID(args string) []string {
	if !strings.Contains(args, ".") {
		logging.Fatal(fmt.Errorf("invalid argument to 'ceph daemon' command: %v. The arg should be in the format of '<daemon.id>'", args))
	}
	return strings.Split(args, ".")
}
