package rook

import (
	"github.com/spf13/cobra"
)

// CephCmd represents the set command
var CephCmd = &cobra.Command{
	Use:                "ceph",
	Short:              "Commands for configuring Ceph",
	DisableFlagParsing: true,
}

func init() {
	CephCmd.AddCommand(cephPurgeOsdCmd)
}
