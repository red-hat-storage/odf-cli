package health

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"

	ocsv1 "github.com/red-hat-storage/ocs-operator/api/v4/v1"
	"golang.org/x/term"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const missingValue = "-"

// kindSpec is one StorageCluster-managed CR kind to list for the health table.
type kindSpec struct {
	component             string
	group                 string
	version               string
	kind                  string
	namespaced            bool
	observedPath          []string
	messageFromConditions bool
}

// managedKinds is the fixed list of CR types StorageCluster manages.
// A kind with no instances (or whose CRD is not installed) is omitted from the table.
var managedKinds = []kindSpec{
	{component: "CephCluster", group: "ceph.rook.io", version: "v1", kind: "CephCluster", namespaced: true, observedPath: []string{"status", "phase"}},
	{component: "Block pool", group: "ceph.rook.io", version: "v1", kind: "CephBlockPool", namespaced: true, observedPath: []string{"status", "phase"}},
	{component: "Filesystem", group: "ceph.rook.io", version: "v1", kind: "CephFilesystem", namespaced: true, observedPath: []string{"status", "phase"}},
	{component: "Object store", group: "ceph.rook.io", version: "v1", kind: "CephObjectStore", namespaced: true, observedPath: []string{"status", "phase"}},
	{component: "NooBaa", group: "noobaa.io", version: "v1alpha1", kind: "NooBaa", namespaced: true, observedPath: []string{"status", "phase"}, messageFromConditions: true},
	{component: "Consumer", group: "ocs.openshift.io", version: "v1alpha1", kind: "StorageConsumer", namespaced: true, observedPath: []string{"status", "state"}},
	{component: "SVG", group: "ceph.rook.io", version: "v1", kind: "CephFilesystemSubVolumeGroup", namespaced: true, observedPath: []string{"status", "phase"}},
	{component: "RNS", group: "ceph.rook.io", version: "v1", kind: "CephBlockPoolRadosNamespace", namespaced: true, observedPath: []string{"status", "phase"}},
	{component: "NFS", group: "ceph.rook.io", version: "v1", kind: "CephNFS", namespaced: true, observedPath: []string{"status", "phase"}},
	{component: "NVMe-oF", group: "ceph.rook.io", version: "v1", kind: "CephNVMeOFGateway", namespaced: true, observedPath: []string{"status", "phase"}},
	{component: "StorageClient", group: "ocs.openshift.io", version: "v1alpha1", kind: "StorageClient", namespaced: false, observedPath: []string{"status", "phase"}},
}

// PrintCRHealth prints one CR status section per StorageCluster.
// Each section has a header (name, namespace, phase, mode) and a table of managed CRs.
// If namespaceSet is true, only that namespace is listed; otherwise all namespaces are searched.
func PrintCRHealth(ctx context.Context, c client.Client, namespace string, namespaceSet bool, out io.Writer) error {
	clusters, err := listStorageClusters(ctx, c, namespace, namespaceSet)
	if err != nil {
		return err
	}
	if len(clusters) == 0 {
		if namespaceSet {
			return fmt.Errorf("no StorageCluster found in namespace %q", namespace)
		}
		return fmt.Errorf("no StorageCluster found")
	}

	for i := range clusters {
		// When multiple StorageClusters are found (multi-cluster scenario), print a
		// blank line between their sections so each header and table is visually
		// distinct. Skip it for the first cluster (i == 0) to avoid leading
		// whitespace at the top of the output.
		if i > 0 {
			fmt.Fprintln(out)
		}
		printClusterHeader(out, clusters[i])
		if err := printManagedTable(ctx, c, clusters[i], out); err != nil {
			return err
		}
	}
	return nil
}

// listStorageClusters lists all StorageClusters present on the cluster.
// If namespaceSet is true, the list is limited to that namespace.
// Phase Ignored clusters are skipped.
func listStorageClusters(ctx context.Context, c client.Client, namespace string, namespaceSet bool) ([]ocsv1.StorageCluster, error) {
	list := &ocsv1.StorageClusterList{}
	var opts []client.ListOption
	if namespaceSet {
		opts = append(opts, client.InNamespace(namespace))
	}
	if err := c.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list StorageClusters: %w", err)
	}

	var out []ocsv1.StorageCluster
	for i := range list.Items {
		if list.Items[i].Status.Phase == "Ignored" {
			continue
		}
		out = append(out, list.Items[i])
	}
	// Sort so the output is deterministic and consistent each time the command runs.
	sort.Slice(out, func(i, j int) bool {
		if out[i].Namespace != out[j].Namespace {
			return out[i].Namespace < out[j].Namespace
		}
		return out[i].Name < out[j].Name
	})
	return out, nil
}

// clusterMode returns "external" when spec.externalStorage.enable is true, otherwise "internal".
func clusterMode(sc ocsv1.StorageCluster) string {
	if sc.Spec.ExternalStorage.Enable {
		return "external"
	}
	return "internal"
}

// printClusterHeader prints the StorageCluster name, namespace, phase, and mode
// above that CR's table. An empty phase is printed as "-".
func printClusterHeader(out io.Writer, sc ocsv1.StorageCluster) {
	phase := sc.Status.Phase
	if phase == "" {
		phase = missingValue
	}
	fmt.Fprintf(out, "%-16s %s\n", "StorageCluster:", sc.Name)
	fmt.Fprintf(out, "%-16s %s\n", "Namespace:", sc.Namespace)
	fmt.Fprintf(out, "%-16s %s\n", "Phase:", phase)
	fmt.Fprintf(out, "%-16s %s\n", "Mode:", clusterMode(sc))
	fmt.Fprintln(out)
}

type tableRow struct {
	component, name, observed, message string
}

// printManagedTable lists each kind in managedKinds and prints one table row per instance.
// Namespaced kinds are listed in the StorageCluster namespace. A kind with no instances,
// or whose CRD is not installed, is omitted. List errors skip that kind instead of failing the command.
// A long STATUS MESSAGE may not fit in the remaining terminal width. In that case we clip it
// instead of wrapping, so the table stays neat and consistent.
func printManagedTable(ctx context.Context, c client.Client, sc ocsv1.StorageCluster, out io.Writer) error {
	var rows []tableRow
	for _, spec := range managedKinds {
		items, err := listKind(ctx, c, spec, sc.Namespace)
		if err != nil {
			continue
		}
		if !spec.namespaced {
			items = matchingStorageClients(items, sc.Name)
		}
		for i := range items {
			component, name, observed, message := rowFields(spec, items[i])
			rows = append(rows, tableRow{component, name, observed, message})
		}
	}

	clipMessages(rows, outputWidth(out))

	w := tabwriter.NewWriter(out, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "COMPONENT\tNAME\tOBSERVED\tSTATUS MESSAGE")
	for _, row := range rows {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", row.component, row.name, row.observed, row.message)
	}
	return w.Flush()
}

// outputWidth returns how many columns the terminal has when out is a terminal.
// If out is not a terminal (tests, pipes, or a file), it returns 0 and the
// STATUS MESSAGE is left unclipped.
func outputWidth(out io.Writer) int {
	f, ok := out.(*os.File)
	if !ok {
		return 0
	}
	width, _, err := term.GetSize(int(f.Fd()))
	if err != nil || width <= 0 {
		return 0
	}
	return width
}

// clipMessages shortens each STATUS MESSAGE so the printed row fits in termWidth.
// tabwriter padding is 2 spaces between columns, so three gaps cost 6 columns.
// remaining is that leftover width for STATUS MESSAGE.
// This function ensures that long status messages do not wrap and the table stays consistent.
func clipMessages(rows []tableRow, termWidth int) {
	if termWidth <= 0 {
		return
	}
	compW, nameW, obsW := len("COMPONENT"), len("NAME"), len("OBSERVED")
	for _, row := range rows {
		if n := len(row.component); n > compW {
			compW = n
		}
		if n := len(row.name); n > nameW {
			nameW = n
		}
		if n := len(row.observed); n > obsW {
			obsW = n
		}
	}
	remaining := termWidth - (compW + nameW + obsW + 6)
	if remaining <= 0 {
		return
	}
	for i := range rows {
		rows[i].message = clipToWidth(rows[i].message, remaining)
	}
}

// clipToWidth shortens s so it fits in width characters.
// If s is longer than that, it is cut and an ellipsis is added at the end.
func clipToWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	return string(r[:width-1]) + "…"
}

// listKind lists all instances of one managed kind.
// Namespaced kinds are limited to the StorageCluster namespace.
// Cluster-scoped kinds (StorageClient) are listed without a namespace filter.
func listKind(ctx context.Context, c client.Client, spec kindSpec, namespace string) ([]unstructured.Unstructured, error) {
	list := &unstructured.UnstructuredList{}
	list.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   spec.group,
		Version: spec.version,
		Kind:    spec.kind + "List",
	})
	var opts []client.ListOption
	if spec.namespaced {
		opts = append(opts, client.InNamespace(namespace))
	}
	if err := c.List(ctx, list, opts...); err != nil {
		return nil, err
	}
	items := list.Items
	// Sort so the output is deterministic and consistent each time the command runs.
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetName() < items[j].GetName()
	})
	return items, nil
}

// matchingStorageClients keeps the cluster-scoped StorageClient whose name matches
// this StorageCluster. Remote clients on another cluster are not StorageClients here.
func matchingStorageClients(items []unstructured.Unstructured, storageClusterName string) []unstructured.Unstructured {
	var out []unstructured.Unstructured
	for i := range items {
		if items[i].GetName() == storageClusterName {
			out = append(out, items[i])
		}
	}
	return out
}

// rowFields returns the four table cells for one CR: COMPONENT from kindSpec,
// NAME from metadata.name, OBSERVED from the nested status field, and STATUS MESSAGE.
func rowFields(spec kindSpec, obj unstructured.Unstructured) (component, name, observed, message string) {
	return spec.component, obj.GetName(), nestedStringOrDash(obj.Object, spec.observedPath...), statusMessage(obj, spec.messageFromConditions)
}

// nestedStringOrDash reads a nested string field (for example status.phase).
// If the field is missing, is not a string, or is empty, it returns "-".
func nestedStringOrDash(obj map[string]interface{}, fields ...string) string {
	s, found, err := unstructured.NestedString(obj, fields...)
	if err != nil || !found || s == "" {
		return missingValue
	}
	return s
}

// statusMessage returns status.message for most kinds.
// NooBaa has no status.message, so fromConditions reads status.conditions[].message instead.
func statusMessage(obj unstructured.Unstructured, fromConditions bool) string {
	if fromConditions {
		return conditionMessage(obj.Object)
	}
	return nestedStringOrDash(obj.Object, "status", "message")
}

// conditionMessage returns the first non-empty status.conditions[].message.
// NooBaa writes the same text on the standard conditions, so any one of those is enough.
// Missing conditions or only empty messages return "-".
func conditionMessage(obj map[string]interface{}) string {
	conds, found, err := unstructured.NestedSlice(obj, "status", "conditions")
	if err != nil || !found {
		return missingValue
	}
	for _, raw := range conds {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		msg, found, err := unstructured.NestedString(m, "message")
		if err != nil || !found || msg == "" {
			continue
		}
		return msg
	}
	return missingValue
}
