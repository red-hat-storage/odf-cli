package get

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/rook"
	"github.com/spf13/cobra"
)

var rookCmd = &cobra.Command{
	Use:   "rook",
	Short: "Supports sub-commands for 'rook'",
	Args:  cobra.RangeArgs(1, 2),
}

var crStatusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Print the phase and/or conditions of Rook CRs",
	Example: "odf get rook status",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rook.PrintCustomResourceStatus(cmd.Context(), root.ClientSets, root.StorageClusterNamespace, args)
	},
}

func init() {
	rookCmd.AddCommand(crStatusCmd)
}
