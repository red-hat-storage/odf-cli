package noobaa

import (
	"slices"

	"github.com/noobaa/noobaa-operator/v5/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	NoobaaCmd       = cli.Cmd()
	invalidCommands = []string{"install", "uninstall", "status", "bucket"}
)

func init() {
	// Modify noobaa RootCmd for odf-cli.
	NoobaaCmd.Use = "noobaa"
	NoobaaCmd.Short = "Run Noobaa cli commands with odf-cli."
	NoobaaCmd.Annotations = map[string]string{
		cobra.CommandDisplayNameAnnotation: "odf noobaa",
	}

	for _, i := range NoobaaCmd.Commands() {
		if slices.Contains(invalidCommands, i.Use) {
			i.Hidden = true
		}
	}
}
