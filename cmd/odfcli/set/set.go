package set

import (
	"github.com/spf13/cobra"
)

// SetCmd represents the set command
var SetCmd = &cobra.Command{
	Use:                "set",
	Short:              "Call subcommands like 'recovery-profile'",
	DisableFlagParsing: true,
	Args:               cobra.ExactArgs(1),
}

func init() {
	SetCmd.AddCommand(recoveryProfile)
}
