package benchmark

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rook/kubectl-rook-ceph/pkg/logging"
)

func RunBenchmarkWorkflow(resourcesJsonPath, daemonsetYamlPath string) error {
	ctx := context.Background()

	logging.Info("🚀 Creating ConfigMap from resources.json...")
	//nolint:gosec // resourcesJsonPath is a file path provided by the CLI tool, not user input
	createCmd := exec.CommandContext(ctx, "oc", "create", "configmap", "benchmark-metrics", "--from-file=benchmark.json="+resourcesJsonPath, "--dry-run=client", "-o", "yaml")
	applyCmd := exec.CommandContext(ctx, "oc", "apply", "-f", "-")

	pipeReader, pipeWriter := io.Pipe()
	createCmd.Stdout = pipeWriter
	applyCmd.Stdin = pipeReader

	if err := createCmd.Start(); err != nil {
		return fmt.Errorf("failed to start configmap creation: %v", err)
	}
	if err := applyCmd.Start(); err != nil {
		return fmt.Errorf("failed to start configmap apply: %v", err)
	}

	if err := createCmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait for configmap creation: %v", err)
	}
	pipeWriter.Close()
	if err := applyCmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait for configmap apply: %v", err)
	}
	pipeReader.Close()

	logging.Info("✅ ConfigMap applied. Applying DaemonSet...")
	if out, err := exec.CommandContext(ctx, "oc", "apply", "-f", daemonsetYamlPath).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to apply daemonset: %v\n%s", err, string(out))
	}

	logging.Info("⏳ Waiting for benchmark pods to complete...")
	for {
		out, err := exec.CommandContext(ctx, "oc", "get", "pods", "-l", "app=odf-preinstall-benchmark", "--no-headers").Output()
		if err != nil {
			return fmt.Errorf("failed to get pods: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		complete := true
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) > 2 && fields[2] != "Completed" && fields[2] != "Succeeded" {
				complete = false
				break
			}
		}
		if complete {
			logging.Info("✅ All benchmark pods completed.")
			break
		}
		logging.Info("... still waiting")
		time.Sleep(10 * time.Second)
	}

	logging.Info("📥 Collecting logs from all pods...")
	out, err := exec.CommandContext(ctx, "oc", "get", "pods", "-l", "app=odf-preinstall-benchmark", "-o", "name").Output()
	if err != nil {
		return fmt.Errorf("failed to list benchmark pods: %v", err)
	}

	if err := os.MkdirAll("benchmark-logs", 0o755); err != nil {
		return fmt.Errorf("failed to create benchmark-logs directory: %v", err)
	}

	for _, pod := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		name := strings.Split(pod, "/")[1]
		logOut, err := exec.CommandContext(ctx, "oc", "logs", pod).Output()
		if err != nil {
			fmt.Printf("❌ Failed to get logs from %s: %v\n", name, err)
			continue
		}
		if err := os.WriteFile("benchmark-logs/"+name+".log", logOut, 0o600); err != nil {
			fmt.Printf("❌ Failed to write logs for %s: %v\n", name, err)
			continue
		}
		fmt.Printf("✅ Logs saved for pod %s\n", name)
	}

	logging.Info("🧹 Cleaning up DaemonSet and pods...")
	if out, err := exec.CommandContext(ctx, "oc", "delete", "daemonset", "odf-preinstall-benchmark").CombinedOutput(); err != nil {
		fmt.Printf("⚠️ Failed to delete DaemonSet: %v\n%s", err, string(out))
	} else {
		logging.Info("✅ DaemonSet deleted")
	}

	return nil
}
