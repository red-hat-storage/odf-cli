package set

import (
	"fmt"

	"github.com/red-hat-storage/odf-cli/cmd/odfcli"
	"github.com/red-hat-storage/odf-cli/pkg/set"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
)

var validRecoveryOptions = []string{"high_client_ops", "high_recovery_ops", "balanced"}

var recoveryProfile = &cobra.Command{
	Use:                "recovery-profile",
	Short:              "Set the recovery profile to favor new IO, recovery, or balanced mode with options high_client_ops, high_recovery_ops, or balanced",
	DisableFlagParsing: true,
	Args:               cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		recoveryInput := args[0]
		if isValidArg(recoveryInput) {
			clientsets := odfcli.GetClientsets(cmd.Context())
			set.SetProfile(cmd.Context(), clientsets, recoveryInput, odfcli.OperatorNamespace, odfcli.StorageClusterNamespace)
		} else {
			logging.Error(fmt.Errorf("Invalid arguments for recovery-profile command"))
		}
	},
}

func isValidArg(recoveryInput string) bool {
	for c, valid := range validRecoveryOptions {
		if recoveryInput == valid {
			break
		} else if c == len(validRecoveryOptions)-1 {
			logging.Warning("Invalid argument %s, valid arguments are: %v\n", recoveryInput, validRecoveryOptions)
			return false
		}
	}

	return true
}
