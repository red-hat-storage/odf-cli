package purgeosd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/red-hat-storage/odf-cli/pkg/rook/osd"
	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/rook/kubectl-rook-ceph/pkg/rook"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// osdDumpEntry represents a single OSD entry from `ceph osd dump --format json`.
type osdDumpEntry struct {
	OSD int `json:"osd"`
	Up  int `json:"up"` // 1 = up (running), 0 = down
	In  int `json:"in"` // 1 = in cluster, 0 = out (marked for removal)
}

type osdDumpOutput struct {
	OSDs []osdDumpEntry `json:"osds"`
}

// osdDetails holds all info about the OSD gathered from k8s before purging.
type osdDetails struct {
	id             int
	pvcName        string
	pvName         string
	nodeName       string
	devicePath     string
	encrypted      bool
	deploymentName string
}

var CephPurgeOsdCmd = &cobra.Command{
	Use:   "purge-osd",
	Short: "Permanently remove an OSD from the cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		clientsets := root.ClientSets
		operatorNS := root.OperatorNamespace
		storageNS := root.StorageClusterNamespace

		osdID := selectUnhealthyOSD(ctx, clientsets, operatorNS, storageNS)
		osdIDStr := strconv.Itoa(osdID)

		// LSO operator must be scaled down before purging
		// to prevent LSO from reprovisioning the disk we are about to wipe
		if !checkLSOScaledDown(ctx, clientsets) {
			fmt.Println("Purge cancelled.")
			return
		}

		// show cluster health before making any changes
		// a pre-existing HEALTH_WARN or HEALTH_ERR may indicate that safeToDestroy
		// will not return true, or that the cluster is already under stress.
		if !clusterHealthCheck(ctx, clientsets, operatorNS, storageNS) {
			fmt.Println("Purge cancelled.")
			return
		}

		// gather OSD details from k8s (deployment, PVC, PV) before any deletion
		details := gatherOSDDetails(ctx, clientsets, storageNS, osdID)

		if !stopOSDIfRunning(ctx, clientsets, storageNS, details) {
			fmt.Println("Purge cancelled.")
			return
		}

		if !markOSDOut(ctx, clientsets, operatorNS, storageNS, osdID, osdIDStr) {
			fmt.Println("Purge cancelled.")
			return
		}

		// returns true if force removal should be used (user chose not to wait)
		forceRemoval := dataMigrationPrompt(ctx, clientsets, operatorNS, storageNS, osdIDStr, osdID)

		printDryRun(details)

		if !promptYN("Proceed with purge?") {
			fmt.Println("Purge cancelled.")
			return
		}

		// execute the purge via the Rook operator.
		// on failure, the process exits. on success, it returns and we print the result below.
		forceFlagValue := "false"
		if forceRemoval {
			forceFlagValue = "true"
		}
		rook.PurgeOsds(ctx, clientsets, operatorNS, storageNS, osdIDStr, forceFlagValue)

		fmt.Printf("\nPurge complete.\n")
		if details.pvcName != "" {
			fmt.Printf("PVC %q is now Pending (awaiting disk).\n", details.pvcName)
		}
		fmt.Println("\nNext steps:")
		fmt.Println("  Replace the physical disk and scale the LSO operator back up.")
		fmt.Println("  If reusing the same disk: scale the LSO operator back up.")
	},
}

// selectUnhealthyOSD queries ceph for all down OSDs and prompts the user to pick one.
// Exits if the cluster cannot be reached or no unhealthy OSDs are found.
func selectUnhealthyOSD(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNS, storageNS string) int {
	logging.Info("Scanning cluster for unhealthy OSDs...")

	out, err := exec.RunCommandInOperatorPod(ctx, clientsets, "ceph",
		[]string{"osd", "dump", "--format", "json"}, operatorNS, storageNS, true)
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to get osd dump: %v", err))
	}

	var dump osdDumpOutput
	if err := json.Unmarshal([]byte(out), &dump); err != nil {
		logging.Fatal(fmt.Errorf("failed to parse osd dump: %v", err))
	}

	var unhealthy []osdDumpEntry
	for _, o := range dump.OSDs {
		if o.Up == 0 {
			unhealthy = append(unhealthy, o)
		}
	}

	if len(unhealthy) == 0 {
		logging.Fatal(fmt.Errorf("no unhealthy OSDs found"))
	}

	fmt.Printf("Found %d unhealthy OSD(s):\n", len(unhealthy))
	for i, o := range unhealthy {
		status := "down"
		if o.In == 0 {
			status = "out"
		}
		fmt.Printf("  [%d] osd.%d (%s)\n", i+1, o.OSD, status)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nWhich OSD would you like to purge? (enter number): ")
	input, _ := reader.ReadString('\n')
	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || choice < 1 || choice > len(unhealthy) {
		logging.Fatal(fmt.Errorf("invalid selection"))
	}

	return unhealthy[choice-1].OSD
}

// checkLSOScaledDown verifies that the LSO operator is scaled down before purging.
// If LSO is running, it could reprovision the disk we are about to wipe.
// Returns false if LSO is running or its status cannot be verified.
func checkLSOScaledDown(ctx context.Context, clientsets *k8sutil.Clientsets) bool {
	lsoNS := "openshift-local-storage"
	deployment, err := clientsets.Kube.AppsV1().Deployments(lsoNS).Get(ctx, "local-storage-operator", metav1.GetOptions{})
	if err != nil {
		logging.Warning("could not verify LSO operator status: %v. Scale it down before running this command.", err)
		return false
	}

	if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas > 0 {
		logging.Warning("LSO operator is running. Scale it down before running this command.")
		return false
	}

	fmt.Println("✓ LSO operator is scaled down.")
	return true
}

// clusterHealthCheck runs `ceph health detail` and prompts the user if the cluster
// is not healthy. Returns false if the user chooses to cancel.
func clusterHealthCheck(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNS, storageNS string) bool {
	logging.Info("Checking cluster health...")

	out, err := exec.RunCommandInOperatorPod(ctx, clientsets, "ceph",
		[]string{"health", "detail"}, operatorNS, storageNS, true)
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to check cluster health: %v", err))
	}

	out = strings.TrimSpace(out)
	if strings.HasPrefix(out, "HEALTH_OK") {
		fmt.Println("✓ Cluster is healthy.")
		return true
	}

	logging.Warning("%s", out)
	return promptYN("Continue with removal?")
}

// gatherOSDDetails queries the k8s deployment, PVC, and PV for the given OSD
// and prints a summary. All resource names are captured here before any deletion begins.
// Exits if required k8s resources cannot be read.
func gatherOSDDetails(ctx context.Context, clientsets *k8sutil.Clientsets, storageNS string, osdID int) *osdDetails {
	details := &osdDetails{
		id:             osdID,
		deploymentName: fmt.Sprintf("rook-ceph-osd-%d", osdID),
	}

	deployment, err := clientsets.Kube.AppsV1().Deployments(storageNS).Get(ctx, details.deploymentName, metav1.GetOptions{})
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to get deployment %q: %v", details.deploymentName, err))
	}

	// PVC name and encryption state are stored as labels on the OSD deployment by Rook
	details.pvcName = deployment.GetLabels()["ceph.rook.io/pvc"]
	details.encrypted = deployment.GetLabels()["encrypted"] == "true"

	if details.pvcName != "" {
		pvc, err := clientsets.Kube.CoreV1().PersistentVolumeClaims(storageNS).Get(ctx, details.pvcName, metav1.GetOptions{})
		if err != nil {
			logging.Fatal(fmt.Errorf("failed to get PVC %q: %v", details.pvcName, err))
		}
		details.pvName = pvc.Spec.VolumeName

		if details.pvName != "" {
			pv, err := clientsets.Kube.CoreV1().PersistentVolumes().Get(ctx, details.pvName, metav1.GetOptions{})
			if err != nil {
				logging.Fatal(fmt.Errorf("failed to get PV %q: %v", details.pvName, err))
			}
			// device path is the LSO symlink path stored in the local PV
			if pv.Spec.Local != nil {
				details.devicePath = pv.Spec.Local.Path
			}
			// node name comes from the PV's node affinity, set by LSO at provision time
			if pv.Spec.NodeAffinity != nil && pv.Spec.NodeAffinity.Required != nil {
				for _, term := range pv.Spec.NodeAffinity.Required.NodeSelectorTerms {
					for _, expr := range term.MatchExpressions {
						if expr.Key == "kubernetes.io/hostname" && len(expr.Values) > 0 {
							details.nodeName = expr.Values[0]
						}
					}
				}
			}
		}
	}

	if details.devicePath == "" {
		logging.Fatal(fmt.Errorf("could not determine device path for osd.%d: PV %q has no local path", osdID, details.pvName))
	}
	if details.nodeName == "" {
		logging.Fatal(fmt.Errorf("could not determine node for osd.%d: PV %q has no node affinity", osdID, details.pvName))
	}

	fmt.Printf("\nOSD Details:\n")
	fmt.Printf("  OSD ID:    osd.%d\n", details.id)
	fmt.Printf("  Host:      %s\n", details.nodeName)
	fmt.Printf("  PVC:       %s\n", details.pvcName)
	fmt.Printf("  Device:    %s\n", details.devicePath)
	fmt.Printf("  Encrypted: %v\n\n", details.encrypted)

	return details
}

// stopOSDIfRunning scales the OSD deployment to 0 replicas to stop the OSD process.
// The OSD process must be stopped before the PVC is deleted to avoid I/O errors.
// Skips the prompt if the deployment is already at 0 replicas.
// Returns false if the user chooses to cancel.
func stopOSDIfRunning(ctx context.Context, clientsets *k8sutil.Clientsets, storageNS string, details *osdDetails) bool {
	deployment, err := clientsets.Kube.AppsV1().Deployments(storageNS).Get(ctx, details.deploymentName, metav1.GetOptions{})
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to get deployment %q: %v", details.deploymentName, err))
	}

	if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas == 0 {
		fmt.Printf("osd.%d process is already stopped.\n", details.id)
		return true
	}

	if !promptYN("The OSD process is currently running. Stop it before proceeding?") {
		return false
	}

	replicas := int32(0)
	deployment.Spec.Replicas = &replicas
	if _, err := clientsets.Kube.AppsV1().Deployments(storageNS).Update(ctx, deployment, metav1.UpdateOptions{}); err != nil {
		logging.Fatal(fmt.Errorf("failed to stop OSD deployment %q: %v", details.deploymentName, err))
	}
	fmt.Printf("Stopped osd.%d process.\n", details.id)
	return true
}

// markOSDOut marks the OSD out of the cluster, which triggers Ceph to start
// migrating data away from this OSD to the remaining OSDs.
// Skips the prompt if the OSD is already marked out. Earlier steps are idempotent
// so re-running the command picks up from here without repeating them.
// Returns false if the user chooses to cancel.
func markOSDOut(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNS, storageNS string, osdID int, osdIDStr string) bool {
	out, err := exec.RunCommandInOperatorPod(ctx, clientsets, "ceph",
		[]string{"osd", "dump", "--format", "json"}, operatorNS, storageNS, true)
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to get osd dump: %v", err))
	}

	var dump osdDumpOutput
	if err := json.Unmarshal([]byte(out), &dump); err != nil {
		logging.Fatal(fmt.Errorf("failed to parse osd dump: %v", err))
	}

	for _, o := range dump.OSDs {
		if o.OSD == osdID && o.In == 0 {
			fmt.Printf("osd.%d is already marked out, resuming purge\n", osdID)
			return true
		}
	}

	if !promptYN(fmt.Sprintf("Mark osd.%d out of the cluster?", osdID)) {
		return false
	}

	if _, err := exec.RunCommandInOperatorPod(ctx, clientsets, "ceph",
		[]string{"osd", "out", fmt.Sprintf("osd.%s", osdIDStr)}, operatorNS, storageNS, false); err != nil {
		logging.Fatal(fmt.Errorf("failed to mark osd.%d out: %v", osdID, err))
	}
	fmt.Printf("Marked osd.%d out.\n", osdID)
	return true
}

// dataMigrationPrompt asks whether to wait for data migration to complete.
// If yes: polls safeToDestroy every 15 seconds with Ctrl+C support for graceful exit.
// If no: returns true so the caller passes force=true to the purge.
func dataMigrationPrompt(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNS, storageNS, osdIDStr string, osdID int) bool {
	if !promptYN("Wait for data migration to complete?") {
		logging.Warning("Proceeding without waiting for data migration. Data redundancy may be reduced.")
		return true // force removal
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	logging.Info("Waiting for data migration...")
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigCh:
			fmt.Printf("\nosd.%d is marked out. Re-run this command to continue the purge.\n", osdID)
			os.Exit(0)
		case <-ticker.C:
			safe, err := osd.SafeToDestroy(ctx, clientsets, operatorNS, storageNS, osdIDStr)
			if err != nil {
				logging.Info("checking safe-to-destroy: %v", err)
				continue
			}
			if safe {
				fmt.Println("✓ Data migration complete.")
				fmt.Println("✓ OSD is safe to destroy.")
				return false // no force needed
			}
			logging.Info("  still migrating data...")
		}
	}
}

// printDryRun prints all k8s and Ceph resources that will be deleted during the purge.
func printDryRun(details *osdDetails) {
	fmt.Println("Resources to be deleted:")
	fmt.Printf("  Deployment: %s\n", details.deploymentName)
	if details.pvcName != "" {
		fmt.Printf("  PVC:        %s\n", details.pvcName)
	}
	if details.pvName != "" {
		fmt.Printf("  PV:         %s\n", details.pvName)
	}
	if details.encrypted {
		fmt.Println("  DM-crypt:   LUKS mapping will be closed")
	}
	if details.devicePath != "" {
		fmt.Printf("  Device:     %s (will be wiped)\n", details.devicePath)
	}
	fmt.Println()
}

// promptYN prints a question and reads a y/n answer from stdin.
func promptYN(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/n]: ", question)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}
