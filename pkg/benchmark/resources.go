package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rook/kubectl-rook-ceph/pkg/logging"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeResource struct {
	NodeName          string   `json:"node_name"`
	Disks             []string `json:"disks"`
	NetworkInterfaces []string `json:"network_interfaces"`
}

type ClusterResources struct {
	Nodes []NodeResource `json:"nodes"`
}

func CollectResources(outputPath string) error {
	nodes, err := getClusterNodes()
	if err != nil {
		return err
	}

	var resources ClusterResources

	for _, node := range nodes {
		logging.Info("Collecting data for node: %s\n", node)
		disks, err := getNodeDisks(node)
		if err != nil {
			return err
		}
		nics, err := getNodeNics(node)
		if err != nil {
			return err
		}

		resources.Nodes = append(resources.Nodes, NodeResource{
			NodeName:          node,
			Disks:             disks,
			NetworkInterfaces: nics,
		})
	}

	jsonBytes, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, jsonBytes, 0644)
}

func getClusterNodes() ([]string, error) {

	nodeList, err := root.ClientSets.Kube.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %v", err)
	}

	var nodeNames []string
	for _, node := range nodeList.Items {
		nodeNames = append(nodeNames, node.Name)
	}
	return nodeNames, nil
}

func getNodeDisks(nodeName string) ([]string, error) {
	cmd := exec.Command("oc", "debug", "node/"+nodeName, "--", "chroot", "/host", "lsblk", "-d", "-n", "-o", "NAME")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get disks for node %s: %v", nodeName, err)
	}
	disks := parseLines(output)

	mountedDisks, err := getMountedDisks(nodeName)
	if err != nil {
		return nil, err
	}

	return filterUnmountedDisks(disks, mountedDisks), nil
}

func getMountedDisks(nodeName string) ([]string, error) {
	cmd := exec.Command("oc", "debug", "node/"+nodeName, "--", "chroot", "/host", "lsblk", "-o", "NAME,MOUNTPOINT")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get mounted disks for node %s: %v", nodeName, err)
	}
	return parseMountedDisks(output), nil
}

func parseMountedDisks(output []byte) []string {
	lines := strings.Split(string(output), "\n")
	var mountedDisks []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 1 && fields[1] != "" { // Disk has a mountpoint
			mountedDisks = append(mountedDisks, fields[0])
		}
	}
	return mountedDisks
}

func filterUnmountedDisks(allDisks, mountedDisks []string) []string {
	mountedSet := make(map[string]bool)
	for _, disk := range mountedDisks {
		mountedSet[disk] = true
	}

	var unmountedDisks []string
	for _, disk := range allDisks {
		if !mountedSet[disk] {
			unmountedDisks = append(unmountedDisks, disk)
		}
	}
	return unmountedDisks
}

func parseLines(output []byte) []string {
	lines := strings.Split(string(output), "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func getNodeNics(nodeName string) ([]string, error) {
	cmd := exec.Command("oc", "debug", "node/"+nodeName, "--", "chroot", "/host", "ip", "-o", "link", "show")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get NICs for node %s: %v", nodeName, err)
	}
	return filterPhysicalInterfaces(parseNicLines(output)), nil
}

func parseNicLines(output []byte) []string {
	lines := strings.Split(string(output), "\n")
	var interfaces []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 1 {
			iface := strings.TrimSuffix(fields[1], ":") // Ensure we remove only the colon
			interfaces = append(interfaces, iface)
		}
	}
	return interfaces
}

func filterPhysicalInterfaces(nics []string) []string {
	var filtered []string
	for _, nic := range nics {
		if strings.HasPrefix(nic, "lo") || // Ignore loopback
			strings.HasPrefix(nic, "br-") || // Ignore OVS bridges
			strings.HasPrefix(nic, "ovs-") || // Ignore Open vSwitch
			strings.HasPrefix(nic, "ovn-") || // Ignore OVN tunnels
			strings.HasPrefix(nic, "genev_sys") || // Ignore geneve tunnels
			len(nic) > 15 { // Ignore hashed virtual interfaces
			continue
		}
		filtered = append(filtered, nic)
	}
	return filtered
}
