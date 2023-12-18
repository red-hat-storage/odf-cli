package set

import "github.com/spf13/cobra"

var cephCmd = &cobra.Command{
	Use:                "ceph",
	Short:              "Configure ceph components",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
}

func init() {
	cephCmd.AddCommand(setCephLogLevelCmd)
}
