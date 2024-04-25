package set

import (
	"fmt"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

var osdBackFillFullRatioCmd = &cobra.Command{
	Use:                "backfillfull",
	Short:              "Configure ceph osd 'backfillfull-ratio' setting",
	DisableFlagParsing: true,
	Args:               cobra.ExactArgs(1),
	Example:            "odf set backfillfull-ratio 0.8",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		cephArgs := []string{"osd", "set-backfillfull-ratio", args[0]}
		_, err := exec.RunCommandInOperatorPod(ctx, root.ClientSets, "ceph", cephArgs, root.OperatorNamespace, root.StorageClusterNamespace, true)
		if err != nil {
			logging.Fatal(fmt.Errorf("failed to run command ceph with args %v: %v", cephArgs, err))
		}
		logging.Info("Successfully set 'backfillfull-ratio' to '%s'", args[0])
	},
}
