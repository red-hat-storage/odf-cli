package dr

import (
	"github.com/ramendr/ramenctl/cmd/commands"
	"github.com/spf13/cobra"
)

// DrCmd is the dr sub command.
var DrCmd = commands.RootCmd

func init() {
	// Modify ramenctl RootCmd for odf-cli.
	DrCmd.Use = "dr"
	DrCmd.Short = "Troubleshoot ODF DR"
	DrCmd.Annotations = map[string]string{
		cobra.CommandDisplayNameAnnotation: "odf dr",
	}

	// Add a subset of ramenctl commands suitable for "odf dr".
	DrCmd.AddCommand(commands.InitCmd)
	DrCmd.AddCommand(commands.TestCmd)
}
