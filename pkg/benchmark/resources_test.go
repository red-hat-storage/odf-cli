// Unit tests
package benchmark

import "testing"

func TestParseNicLines(t *testing.T) {
	sampleOutput := []byte(
		`1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN mode DEFAULT group default qlen 1000
         2: eth0@if24: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 8901 qdisc noqueue state UP mode DEFAULT group default`,
	)
	expected := []string{"lo", "eth0@if24"} // Ensure @ interface names are preserved
	result := parseNicLines(sampleOutput)

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	}
}

func TestFilterPhysicalInterfaces(t *testing.T) {
	sampleNics := []string{"lo", "br-ex", "eth0", "ovs-system", "ens5", "ovn-k8s-mp0"}
	expected := []string{"eth0", "ens5"}
	result := filterPhysicalInterfaces(sampleNics)
	if len(result) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseMountedDisks(t *testing.T) {
	sampleOutput := []byte("sda /mnt/data\nsdb \nsdc /var/lib\nsdd \n")
	expected := []string{"sda", "sdc"}
	result := parseMountedDisks(sampleOutput)
	if len(result) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFilterUnmountedDisks(t *testing.T) {
	allDisks := []string{"sda", "sdb", "sdc", "sdd"}
	mountedDisks := []string{"sda", "sdc"}
	expected := []string{"sdb", "sdd"}
	result := filterUnmountedDisks(allDisks, mountedDisks)
	if len(result) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
