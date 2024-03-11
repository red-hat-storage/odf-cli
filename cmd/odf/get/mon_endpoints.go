package get

import (
	"fmt"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/mons"
	"github.com/spf13/cobra"
)

var monEndpoints = &cobra.Command{
	Use:     "mon-endpoints",
	Short:   "Print mon endpoints",
	Args:    cobra.NoArgs,
	Example: "odf get mon-endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(mons.GetMonEndpoint(cmd.Context(), root.ClientSets.Kube, root.StorageClusterNamespace))
	},
}
