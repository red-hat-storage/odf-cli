package ceph

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

	"github.com/red-hat-storage/odf-cli/pkg/rook/osd"
	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OSDDumpEntry represents a single OSD entry from `ceph osd dump --format json`.
type OSDDumpEntry struct {
	OSD int `json:"osd"`
	Up  int `json:"up"` // 1 = up (running), 0 = down
	In  int `json:"in"` // 1 = in cluster, 0 = out (marked for removal)
}

type osdDumpOutput struct {
	OSDs []OSDDumpEntry `json:"osds"`
}

// OSDDetails holds all info about the OSD gathered from k8s before purging.
type OSDDetails struct {
	ID             int
	PVCName        string
	PVName         string
	NodeName       string
	DevicePath     string
	Encrypted      bool
	DeploymentName string
}

type cephStatus struct {
	PGMap struct {
		PGsByState []struct {
			StateName string `json:"state_name"`
			Count     int    `json:"count"`
		} `json:"pgs_by_state"`
	} `json:"pgmap"`
}

// SelectUnhealthyOSD queries ceph for all down OSDs and prompts the user to pick one.
// Exits if the cluster cannot be reached or no unhealthy OSDs are found.
func SelectUnhealthyOSD(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNS, storageNS string) int {
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

	var unhealthy []OSDDumpEntry
	for _, o := range dump.OSDs {
		if o.Up == 0 {
			unhealthy = append(unhealthy, o)
		}
	}

	if len(unhealthy) == 0 {
		logging.Fatal(fmt.Errorf("no unhealthy OSDs found to be purged. " +
			"To force remove an OSD that is running, first scale down the OSD deployment: " +
			"kubectl scale deployment rook-ceph-osd-<ID> --replicas=0 -n <cluster-namespace>"))
	}

	fmt.Printf("Found %d unhealthy OSD(s):\n", len(unhealthy))
	for _, o := range unhealthy {
		status := "down"
		if o.In == 0 {
			status = "out"
		}
		fmt.Printf("  osd.%d (%s)\n", o.OSD, status)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter OSD ID to purge: ")
	input, _ := reader.ReadString('\n')
	osdID, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		logging.Fatal(fmt.Errorf("invalid OSD ID"))
	}

	// Validate the OSD is in the unhealthy list
	for _, o := range unhealthy {
		if o.OSD == osdID {
			return osdID
		}
	}
	logging.Fatal(fmt.Errorf("OSD %d is not unhealthy. Only down OSDs can be purged.", osdID))
	return -1
}

// CheckLSOScaledDown verifies that the LSO operator is scaled down before purging.
// LSO must be scaled down manually by the user before running this command.
// If LSO is running, it will reprovision the disk we are about to wipe.
// Returns false if LSO is running or its status cannot be verified.
func CheckLSOScaledDown(ctx context.Context, clientsets *k8sutil.Clientsets) bool {
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

// ClusterHealthCheck runs `ceph health detail` and prompts the user if the cluster
// is not healthy. Returns false if the user chooses to cancel.
func ClusterHealthCheck(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNS, storageNS string) bool {
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
	return PromptYN("Continue with removal?")
}

// GatherOSDDetails queries the k8s deployment, PVC, and PV for the given OSD.
// All resource names are captured here before any deletion begins.
// Exits if required k8s resources cannot be read.
func GatherOSDDetails(ctx context.Context, clientsets *k8sutil.Clientsets, storageNS string, osdID int) *OSDDetails {
	details := &OSDDetails{
		ID:             osdID,
		DeploymentName: fmt.Sprintf("rook-ceph-osd-%d", osdID),
	}

	deployment, err := clientsets.Kube.AppsV1().Deployments(storageNS).Get(ctx, details.DeploymentName, metav1.GetOptions{})
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to get deployment %q: %v", details.DeploymentName, err))
	}

	// PVC name and encryption state are stored as labels on the OSD deployment by Rook
	details.PVCName = deployment.GetLabels()["ceph.rook.io/pvc"]
	details.Encrypted = deployment.GetLabels()["encrypted"] == "true"

	if details.PVCName != "" {
		pvc, err := clientsets.Kube.CoreV1().PersistentVolumeClaims(storageNS).Get(ctx, details.PVCName, metav1.GetOptions{})
		if err != nil {
			logging.Fatal(fmt.Errorf("failed to get PVC %q: %v", details.PVCName, err))
		}
		details.PVName = pvc.Spec.VolumeName

		if details.PVName != "" {
			pv, err := clientsets.Kube.CoreV1().PersistentVolumes().Get(ctx, details.PVName, metav1.GetOptions{})
			if err != nil {
				logging.Fatal(fmt.Errorf("failed to get PV %q: %v", details.PVName, err))
			}
			// device path is the LSO symlink path stored in the local PV
			if pv.Spec.Local != nil {
				details.DevicePath = pv.Spec.Local.Path
			}
			// node name comes from the PV's node affinity, set by LSO at provision time
			if pv.Spec.NodeAffinity != nil && pv.Spec.NodeAffinity.Required != nil {
				for _, term := range pv.Spec.NodeAffinity.Required.NodeSelectorTerms {
					for _, expr := range term.MatchExpressions {
						if expr.Key == "kubernetes.io/hostname" && len(expr.Values) > 0 {
							details.NodeName = expr.Values[0]
						}
					}
				}
			}
		}
	}

	if details.DevicePath == "" {
		logging.Fatal(fmt.Errorf("could not determine device path for osd.%d: PV %q has no local path", osdID, details.PVName))
	}
	if details.NodeName == "" {
		logging.Fatal(fmt.Errorf("could not determine node for osd.%d: PV %q has no node affinity", osdID, details.PVName))
	}

	return details
}

// StopOSDIfRunning scales the OSD deployment to 0 replicas to stop the OSD process.
// The OSD process must be stopped before the PVC is deleted to avoid I/O errors.
// Skips the prompt if the deployment is already at 0 replicas.
// Returns false if the user chooses to cancel.
func StopOSDIfRunning(ctx context.Context, clientsets *k8sutil.Clientsets, storageNS string, details *OSDDetails) bool {
	deployment, err := clientsets.Kube.AppsV1().Deployments(storageNS).Get(ctx, details.DeploymentName, metav1.GetOptions{})
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to get deployment %q: %v", details.DeploymentName, err))
	}

	if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas == 0 {
		fmt.Printf("osd.%d process is already stopped.\n", details.ID)
		return true
	}

	if !PromptYN("The OSD process is currently running. Stop it before proceeding?") {
		return false
	}

	replicas := int32(0)
	deployment.Spec.Replicas = &replicas
	if _, err := clientsets.Kube.AppsV1().Deployments(storageNS).Update(ctx, deployment, metav1.UpdateOptions{}); err != nil {
		logging.Fatal(fmt.Errorf("failed to stop OSD deployment %q: %v", details.DeploymentName, err))
	}
	logging.Info("Stopped osd.%d deployment.", details.ID)

	// Wait for pod to actually terminate
	logging.Info("Waiting for OSD pod to terminate...")
	deadline := time.Now().Add(1 * time.Minute)
	for time.Now().Before(deadline) {
		dep, err := clientsets.Kube.AppsV1().Deployments(storageNS).Get(ctx, details.DeploymentName, metav1.GetOptions{})
		if err == nil && dep.Status.Replicas == 0 {
			logging.Info("OSD pod has terminated.")
			return true
		}
		time.Sleep(1 * time.Second)
	}
	logging.Warning("OSD pod did not terminate after 1 minute. Proceeding anyway...")
	return true
}

// MarkOSDOut marks the OSD out of the cluster, which triggers Ceph to start
// migrating data away from this OSD to the remaining OSDs.
// Skips the prompt if the OSD is already marked out. Earlier steps are idempotent
// so re-running the command picks up from here without repeating them.
// Returns false if the user chooses to cancel.
func MarkOSDOut(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNS, storageNS string, osdID int, osdIDStr string) bool {
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

	if !PromptYN(fmt.Sprintf("Mark osd.%d out of the cluster?", osdID)) {
		return false
	}

	if _, err := exec.RunCommandInOperatorPod(ctx, clientsets, "ceph",
		[]string{"osd", "out", fmt.Sprintf("osd.%s", osdIDStr)}, operatorNS, storageNS, false); err != nil {
		logging.Fatal(fmt.Errorf("failed to mark osd.%d out: %v", osdID, err))
	}
	fmt.Printf("Marked osd.%d out.\n", osdID)
	return true
}

// DataMigrationPrompt checks if OSD is safe to destroy. If not, checks cluster status.
// If all PGs are active+clean, returns true (force removal internally, no prompt needed).
// Otherwise shows status and prompts user to wait for migration or proceed with force removal.
func DataMigrationPrompt(ctx context.Context, clientsets *k8sutil.Clientsets, operatorNS, storageNS, osdIDStr string, osdID int) bool {
	if safe, err := osd.SafeToDestroy(ctx, clientsets, operatorNS, storageNS, osdIDStr); err == nil && safe {
		return false
	}

	var status cephStatus
	raw, err := exec.RunCommandInOperatorPod(ctx, clientsets, "ceph",
		[]string{"status", "--format", "json"}, operatorNS, storageNS, true)

	if err == nil && json.Unmarshal([]byte(raw), &status) == nil {
		allClean := len(status.PGMap.PGsByState) > 0
		var pgStatusOutput []string
		for _, s := range status.PGMap.PGsByState {
			pgStatusOutput = append(pgStatusOutput, fmt.Sprintf("  %d PGs: %s", s.Count, s.StateName))
			if s.StateName != "active+clean" {
				allClean = false
			}
		}
		if allClean {
			logging.Info("All PGs are active+clean, proceeding with force.")
			return true
		}
		fmt.Printf("\nCeph is not fully clean right now:\n%s\n\n", strings.Join(pgStatusOutput, "\n"))
	}

	if !PromptYN("Wait for data migration to complete? (Ceph is moving data away from this OSD to other nodes)") {
		logging.Warning("Forcing removal while data is still migrating. This increases the risk of data loss.")
		return true
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
			fmt.Printf("\nosd.%d is marked out. Re-run this command to resume the purge.\n", osdID)
			os.Exit(0)
		case <-ticker.C:
			safe, err := osd.SafeToDestroy(ctx, clientsets, operatorNS, storageNS, osdIDStr)
			if err != nil {
				logging.Info("checking safe-to-destroy: %v", err)
				continue
			}
			if safe {
				fmt.Println("Data migration complete. OSD is safe to destroy.")
				return false
			}
			logging.Info("Data still migrating, checking again in 15s...")
		}
	}
}

// PrintDryRun prints all k8s and Ceph resources that will be deleted during the purge.
func PrintDryRun(details *OSDDetails) {
	fmt.Println("Resources to be deleted:")
	fmt.Printf("Deployment: %s\n", details.DeploymentName)
	if details.PVCName != "" {
		fmt.Printf("PVC:        %s\n", details.PVCName)
	}
	if details.PVName != "" {
		fmt.Printf("PV:         %s\n", details.PVName)
	}
	if details.Encrypted {
		fmt.Println("DM-crypt:   LUKS mapping will be closed")
	}
	if details.DevicePath != "" {
		fmt.Printf("Device:     %s (will be wiped)\n", details.DevicePath)
	}
	fmt.Println()
}

// PromptYN prints a question and reads a y/n answer from stdin.
func PromptYN(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/n]: ", question)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}
