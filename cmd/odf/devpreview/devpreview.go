package devpreview

import (
	"strings"

	"github.com/spf13/cobra"
)

const (
	// Suffix for appending to the short command description.
	Suffix = " (developer preview)"

	// Note for appending to the long command description.
	Note = `
+------------------------------------------------------------------------------+
|                                                                              |
|    This command is a developer preview, unsupported and not fully tested.    |
|    Please see the following document for more info on developer preview:     |
|    https://access.redhat.com/support/offerings/devpreview                    |
|                                                                              |
+------------------------------------------------------------------------------+
`
)

// Configure a command as a developer preview command. Mark the command hidden
// and amend the command short and long descriptions. If the command has
// children commands their short and long descriptions are amended as well.
func Configure(cmd *cobra.Command) {
	cmd.Hidden = true
	amendDescriptions(cmd)
}

func amendDescriptions(cmd *cobra.Command) {
	cmd.Short += Suffix

	if cmd.Long != "" {
		cmd.Long = strings.TrimRight(cmd.Long, "\n") + "\n" + Note
	} else {
		// If the long description is set the short description is not shown in
		// the help text. Add the short description so we have more specific
		// long description.
		cmd.Long = strings.TrimRight(cmd.Short, "\n") + "\n" + Note
	}

	for _, child := range cmd.Commands() {
		amendDescriptions(child)
	}
}
