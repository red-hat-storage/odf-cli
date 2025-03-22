package main

import (
	"github.com/red-hat-storage/odf-cli/cmd/odf/benchmark"
	"github.com/red-hat-storage/odf-cli/cmd/odf/ceph"
	"github.com/red-hat-storage/odf-cli/cmd/odf/dr"
	"github.com/red-hat-storage/odf-cli/cmd/odf/get"
	"github.com/red-hat-storage/odf-cli/cmd/odf/maintenance"
	"github.com/red-hat-storage/odf-cli/cmd/odf/operator"
	"github.com/red-hat-storage/odf-cli/cmd/odf/purgeosd"
	"github.com/red-hat-storage/odf-cli/cmd/odf/restore"
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/red-hat-storage/odf-cli/cmd/odf/set"
	"github.com/red-hat-storage/odf-cli/cmd/odf/subvolume"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
)

func main() {
	addcommands()
	err := root.RootCmd.Execute()
	if err != nil {
		logging.Fatal(err)
	}
}

func addcommands() {
	root.RootCmd.AddCommand(
		set.SetCmd,
		get.GetCmd,
		purgeosd.CephPurgeOsdCmd,
		subvolume.SubvolumeCmd,
		maintenance.MaintenanceCmd,
		operator.OperatorCmd,
		restore.RestoreCrd,
		ceph.CephCmd,
		ceph.RadosgwCmd,
		ceph.RbdCmd,
		ceph.RadosCmd,
		dr.DrCmd,
		benchmark.BenchmarkCmd,
	)
}
