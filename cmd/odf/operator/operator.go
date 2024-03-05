package operator

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/operator/rook"
	"github.com/spf13/cobra"
)

// OperatorCmd represents the rook commands
var OperatorCmd = &cobra.Command{
	Use:    "operator",
	Short:  "Calls subcommands specific to various ODF operators",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
}

func init() {
	OperatorCmd.AddCommand(rook.RookCmd)
}
