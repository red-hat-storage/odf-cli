package get

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/dr"
	"github.com/spf13/cobra"
)

var drHealthCmd = &cobra.Command{
	Use:                "dr-health",
	Short:              "Print the ceph status of a peer cluster in a mirroring-enabled cluster.",
	DisableFlagParsing: true,
	Args:               cobra.MaximumNArgs(2),
	Example:            "odf get dr-health [ceph status args]",
	Run: func(cmd *cobra.Command, args []string) {
		dr.Health(cmd.Context(), root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, args)
	},
}
