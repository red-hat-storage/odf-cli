package get

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/red-hat-storage/odf-cli/pkg/rook/osd"
	"github.com/spf13/cobra"
)

var getRecoveryProfile = &cobra.Command{
	Use:                "recovery-profile",
	Short:              "Get the recovery profile value currently set for the osd",
	Example:            "odf get recovery-profile",
	DisableFlagParsing: true,
	Args:               cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		osd.GetProfile(cmd.Context(), root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace)
	},
}
