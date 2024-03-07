package restore

import (
	"github.com/spf13/cobra"
)

// restoreCrd represents the mons command
var RestoreCrd = &cobra.Command{
	Use:                "restore",
	DisableFlagParsing: true,
	Hidden:             true,
}

func init() {
	RestoreCrd.AddCommand(deletedCmd)
	RestoreCrd.AddCommand(monQuorumCmd)
}
