package get

import (
	"github.com/red-hat-storage/odf-cli/pkg/odf/dr"
	"github.com/spf13/cobra"
)

var drPrereqCmd = &cobra.Command{
	Use:                "dr-prereq",
	Short:              "Print the status of pre-requisites for Disaster Recovery between peer clusters.",
	DisableFlagParsing: true,
	Args:               cobra.ExactArgs(1),
	Example:            "odf get dr-prereq <PeerManagedClusterName>",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		dr.GetDRPrerequisite(ctx, args[0])
	},
}
