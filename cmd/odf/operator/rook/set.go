package rook

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:     "set",
	Short:   "Set the property in the rook-ceph-operator-config configmap.",
	Example: "odf operator rook set <KEY> <VALUE>",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		k8sutil.UpdateConfigMap(cmd.Context(), root.ClientSets.Kube, root.OperatorNamespace, "rook-ceph-operator-config", args[0], args[1])
	},
}
