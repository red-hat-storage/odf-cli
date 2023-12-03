package set

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/red-hat-storage/odf-cli/pkg/rook"
	"github.com/spf13/cobra"
)

var setCephLogLevelCmd = &cobra.Command{
	Use:                "ceph-log-level",
	Short:              "Set different log levels for ceph components like mon, osd and mds",
	Example:            "odf set ceph-log-level osd crush 10",
	DisableFlagParsing: true,
	Args:               cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		clientsets := root.GetClientsets(ctx)
		rook.SetCephLogLevel(ctx, clientsets, root.OperatorNamespace, root.StorageClusterNamespace, args[0], args[1], args[2])
	},
}
