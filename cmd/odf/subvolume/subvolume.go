package subvolume

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	subvolume "github.com/rook/kubectl-rook-ceph/pkg/filesystem"
	"github.com/spf13/cobra"
)

var SubvolumeCmd = &cobra.Command{
	Use:   "subvolume",
	Short: "Manages subvolumes",
	Args:  cobra.ExactArgs(1),
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "Print the list of stale subvolumes no longer in use.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		staleSubvol, _ := cmd.Flags().GetBool("stale")
		subvolume.List(ctx, root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, staleSubvol)
	},
}

var deleteCmd = &cobra.Command{
	Use:                "delete",
	Short:              "Deletes a stale subvolume",
	DisableFlagParsing: true,
	Args:               cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		subList := args[0]
		fs := args[1]
		svg := args[2]
		subvolume.Delete(ctx, root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, subList, fs, svg)
	},
}

func init() {
	SubvolumeCmd.AddCommand(listCmd)
	SubvolumeCmd.PersistentFlags().Bool("stale", false, "Only list stale subvolumes")
	SubvolumeCmd.AddCommand(deleteCmd)
}
