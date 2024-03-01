package get

import (
	"github.com/spf13/cobra"
)

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:                "get",
	Short:              "Get ODF configuration",
	DisableFlagParsing: true,
	Args:               cobra.ExactArgs(1),
}

func init() {
	GetCmd.AddCommand(getRecoveryProfile)
	GetCmd.AddCommand(rookCmd)
}
