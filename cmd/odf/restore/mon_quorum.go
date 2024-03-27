package restore

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/mons"
	"github.com/spf13/cobra"
)

// monQuorumCmd represents the mons command
var monQuorumCmd = &cobra.Command{
	Use:                "mon-quorum",
	Short:              "When quorum is lost, restore quorum to the remaining healthy mon",
	DisableFlagParsing: true,
	Args:               cobra.ExactArgs(1),
	Example:            "odf restore mon-quorum <mon_id>",
	Run: func(cmd *cobra.Command, args []string) {
		mons.RestoreQuorum(cmd.Context(), root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, args[0])
	},
}
