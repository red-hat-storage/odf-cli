package benchmark

import (
	"fmt"

	"github.com/rook/kubectl-rook-ceph/pkg/logging"

	"github.com/red-hat-storage/odf-cli/pkg/benchmark"
	"github.com/spf13/cobra"
)

var outputPath string

var ResourceCmd = &cobra.Command{
	Use:     "resources",
	Short:   "Collect cluster resources like disks and network interfaces",
	Example: "odf benchmark resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := benchmark.CollectResources(outputPath); err != nil {
			return fmt.Errorf("failed to collect resources: %v", err)
		}
		logging.Info("Resources written to %s\n", outputPath)
		return nil
	},
}

func init() {
	ResourceCmd.Flags().StringVarP(&outputPath, "output", "o", "resources.json", "Output file for resources JSON")
	BenchmarkCmd.AddCommand(ResourceCmd)
}
