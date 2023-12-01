package set

import (
	root "github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/red-hat-storage/odf-cli/pkg/rook/osd"
	"github.com/spf13/cobra"
)

var setRecoveryProfile = &cobra.Command{
	Use:                "recovery-profile",
	Short:              "Set the recovery profile to favor new IO, recovery, or balanced mode with options high_client_ops, high_recovery_ops, or balanced",
	Example:            "odf set recovery-profile <option>",
	DisableFlagParsing: true,
	ValidArgs:          []string{"high_client_ops", "high_recovery_ops", "balanced"},
	Args:               cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		clientsets := root.GetClientsets(cmd.Context())
		osd.SetProfile(cmd.Context(), clientsets, args[0], root.OperatorNamespace, root.StorageClusterNamespace)
	},
}
