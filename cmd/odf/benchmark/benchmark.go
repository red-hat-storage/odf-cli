package benchmark

import (
	"github.com/spf13/cobra"
)

var BenchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark commands for ODF CLI",
}
