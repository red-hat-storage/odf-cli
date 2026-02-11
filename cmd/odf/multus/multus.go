package multus

import (
	multus "github.com/rook/rook/cmd/rook/userfacing/multus"
	multusdefault "github.com/rook/rook/pkg/daemon/multus"
)

func init() {
	multusdefault.DefaultValidationNamespace = "openshift-storage"
	multusdefault.DefaultStorageNodeLabelKey = "cluster.ocs.openshift.io/openshift-storage"
	multusdefault.DefaultStorageNodeLabelValue = ""
}

// MultusCmd is imported from rook/rook for Multus network validation
var MultusCmd = multus.Cmd
