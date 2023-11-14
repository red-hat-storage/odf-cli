package main

import (
	"github.com/red-hat-storage/odf-cli/cmd/odfcli"
	"github.com/red-hat-storage/odf-cli/cmd/odfcli/set"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
)

func main() {
	addcommands()
	err := odfcli.RootCmd.Execute()
	if err != nil {
		logging.Fatal(err)
	}
}

func addcommands() {
	odfcli.RootCmd.AddCommand(
		set.SetCmd,
	)
}
