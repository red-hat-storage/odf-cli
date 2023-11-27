package set

import (
	"github.com/red-hat-storage/odf-cli/cmd/odfcli"
	"github.com/red-hat-storage/odf-cli/pkg/set"
	"github.com/spf13/cobra"
)

var logLevelCmd = &cobra.Command{
	Use:                "log-level",
	Short:              "Set different log levels for ceph dameons like mon, osd and mds",
	DisableFlagParsing: true,
	Args:               cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		clientsets := odfcli.GetClientsets(cmd.Context())
		set.SetLogLevel(cmd.Context(), clientsets, odfcli.OperatorNamespace, odfcli.StorageClusterNamespace, args)
	},
}
