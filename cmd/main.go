package main

import (
	command "github.com/red-hat-storage/odf-cli/cmd/commands"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
)

func main() {
	addcommands()
	err := command.RootCmd.Execute()
	if err != nil {
		logging.Fatal(err)
	}
}

func addcommands() {
	command.RootCmd.AddCommand()
}
