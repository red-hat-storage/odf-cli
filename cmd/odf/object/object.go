package object

import (
	"fmt"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	noobaapkg "github.com/red-hat-storage/odf-cli/pkg/noobaa"
	"github.com/spf13/cobra"
)

const remoteOBCArg = "remote-obc"

// ObjectCmd groups object-storage cluster operations.
var ObjectCmd = &cobra.Command{
	Use:   "object",
	Short: "Object storage related cluster operations.",
	Annotations: map[string]string{
		cobra.CommandDisplayNameAnnotation: "odf object",
	},
}

var enableCmd = &cobra.Command{
	Use:       "enable",
	Short:     "Install the ObjectBucket and ObjectBucketClaim CRDs",
	Example:   fmt.Sprintf("odf object enable %s", remoteOBCArg),
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{remoteOBCArg},
	Run: func(cmd *cobra.Command, _ []string) {
		noobaapkg.InstallCRDs(cmd.Context(), root.APIExtensions)
	},
}

var disableCmd = &cobra.Command{
	Use:       "disable",
	Short:     "Remove the ObjectBucket and ObjectBucketClaim CRDs",
	Example:   fmt.Sprintf("odf object disable %s", remoteOBCArg),
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{remoteOBCArg},
	Run: func(cmd *cobra.Command, _ []string) {
		noobaapkg.UninstallCRDs(cmd.Context(), root.APIExtensions)
	},
}

func init() {
	ObjectCmd.AddCommand(enableCmd, disableCmd)
}
