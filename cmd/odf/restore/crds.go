package restore

import (
	"fmt"
	"slices"
	"strings"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	pkgrestore "github.com/rook/kubectl-rook-ceph/pkg/restore"
	"github.com/spf13/cobra"
)

// groupVersions defines the supported CRD groups and their corresponding API versions.
var groupVersions = map[string]string{
	"ocs.openshift.io":       "v1",
	"ceph.rook.io":           "v1",
	"csi.ceph.io":            "v1alpha1",
	"odf.openshift.io":       "v1alpha1",
	"noobaa.io":              "v1alpha1",
	"csiaddons.openshift.io": "v1alpha1",
}

// groupNameKeys returns the keys of a string map. It is used to print out supported group names.
func groupNameKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// parseFullyQualifiedCRD takes a fully qualified CRD type of the form "resource.group"
// (for example, "cephclusters.ceph.rook.io") and returns the resource name, group name, and
// the API version associated with that group. It returns an error if the format is invalid
// or the group is not recognized.
func parseFullyQualifiedCRD(fqcrd string) (resourceName, groupName, version string, err error) {
	parts := strings.SplitN(fqcrd, ".", 2)
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid CRD format %q; expected format <resource>.<group>", fqcrd)
	}
	resourceName = parts[0]
	groupName = parts[1]

	version, ok := groupVersions[groupName]
	if !ok {
		return "", "", "", fmt.Errorf("unsupported group %q; supported groups are: %v", groupName, groupNameKeys(groupVersions))
	}
	return resourceName, groupName, version, nil
}

// deletedCmd represents the deleted command
var deletedCmd = &cobra.Command{
	Use:                "deleted",
	Short:              "Restores a CR that was accidentally deleted and is still in terminating state.",
	DisableFlagParsing: true,
	Args:               cobra.RangeArgs(1, 2),
	Example:            "odf restore deleted <CRD> [CRNAME]",
	PreRun: func(cmd *cobra.Command, args []string) {
		k8sutil.WaitForPodToRun(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "app=rook-ceph-operator")

	},
	Run: func(cmd *cobra.Command, args []string) {
		k8sutil.SetDeploymentScale(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "ocs-operator", 0)
		// Parse the fully qualified CRD (e.g. "cephclusters.ceph.rook.io").
		resourceName, groupName, version, err := parseFullyQualifiedCRD(args[0])
		if err != nil {
			logging.Fatal(fmt.Errorf("Error parsing CRD type: %v\n", err))
		}
		// Construct a new args slice with the resource name as the first argument.
		newArgs := make([]string, len(args))
		newArgs[0] = resourceName
		if len(args) > 1 {
			newArgs[1] = args[1]
		}
		var customResources []pkgrestore.CustomResource
		if slices.Contains(newArgs, "storageclusters") {
			customResources = []pkgrestore.CustomResource{
				// ceph.rook.io/v1
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephblockpoolradosnamespaces"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephblockpools"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephbucketnotifications"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephbuckettopics"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephclients"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephclusters"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephcosidrivers"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephfilesystemmirrors"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephfilesystems"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephfilesystemsubvolumegroups"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephnfses"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephobjectrealms"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephobjectstores"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephobjectstoreusers"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephobjectzonegroups"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephobjectzones"},
				{Group: "ceph.rook.io", Version: "v1", Resource: "cephrbdmirrors"},

				// noobaa.io/v1alpha1
				{Group: "noobaa.io", Version: "v1alpha1", Resource: "backingstores"},
				{Group: "noobaa.io", Version: "v1alpha1", Resource: "bucketclasses"},
				{Group: "noobaa.io", Version: "v1alpha1", Resource: "namespacestores"},
				{Group: "noobaa.io", Version: "v1alpha1", Resource: "noobaaaccounts"},
				{Group: "noobaa.io", Version: "v1alpha1", Resource: "noobaas"},

				// ocs.openshift.io/v1alpha1
				{Group: "ocs.openshift.io", Version: "v1alpha1", Resource: "storageconsumers"},
			}
		} else {
			customResources = []pkgrestore.CustomResource{}
		}

		pkgrestore.RestoreCrd(cmd.Context(), root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, groupName, version, "ocs-operator", customResources, newArgs)
		k8sutil.SetDeploymentScale(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "ocs-operator", 1)
	},
}
