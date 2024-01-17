package set

import (
	"github.com/spf13/cobra"
)

// SetCmd represents the set command
var SetCmd = &cobra.Command{
	Use:                "set",
	Short:              "Set ODF configuration",
	DisableFlagParsing: true,
}

func init() {
	SetCmd.AddCommand(setRecoveryProfile)
	SetCmd.AddCommand(cephCmd)
}
