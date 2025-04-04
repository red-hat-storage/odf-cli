package benchmark

import (
	"github.com/red-hat-storage/odf-cli/pkg/benchmark"
	"github.com/spf13/cobra"
)

var (
	resourcesPath     string = "resources.json"
	daemonsetYamlPath string = "https://raw.githubusercontent.com/manishym/odf-benchmarker/refs/heads/main/daemonSet.yaml"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run benchmark tests on the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return benchmark.RunBenchmarkWorkflow(resourcesPath, daemonsetYamlPath)
	},
}

func init() {
	runCmd.Flags().StringVar(&resourcesPath, "resources", resourcesPath, "Path to combined benchmark.json")
	runCmd.Flags().StringVar(&daemonsetYamlPath, "daemonset", daemonsetYamlPath, "Path to daemonset YAML")
	BenchmarkCmd.AddCommand(runCmd)
}
