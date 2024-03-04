package mon

import (
	"github.com/spf13/cobra"
)

// MonsCmd represents the mon command
var MonsCmd = &cobra.Command{
	Use:                "mons",
	Short:              "Supports sub-commands for 'mons'",
	DisableFlagParsing: true,
	Hidden:             true,
	Args:               cobra.ExactArgs(1),
}

func init() {
	MonsCmd.AddCommand(RestoreQuorum)
}
