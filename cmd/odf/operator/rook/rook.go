package rook

import (
	"github.com/spf13/cobra"
)

// RookCmd represents the rook commands
var RookCmd = &cobra.Command{
	Use:   "rook",
	Short: "Calls subcommands for the rook operator",
	Args:  cobra.ExactArgs(1),
}

func init() {
	RookCmd.AddCommand(setCmd)
	RookCmd.AddCommand(restartCmd)
}
