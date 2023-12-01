package set

import (
	"fmt"
	"regexp"

	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/red-hat-storage/odf-cli/pkg/rook/common"
	"github.com/spf13/cobra"
)

var setLogLevelCmd = &cobra.Command{
	Use:                "log-level",
	Short:              "Set different log levels for ceph dameons like mon, osd and mds",
	Example:            "odf set log-level osd 20",
	DisableFlagParsing: true,
	PreRunE:            validateLogLevelArgs,
	Args:               cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		clientsets := root.GetClientsets(ctx)
		common.SetLogLevel(ctx, clientsets, root.OperatorNamespace, root.StorageClusterNamespace, args)
	},
}

func validateLogLevelArgs(_ *cobra.Command, args []string) error {
	isValidSubsytem := false
	for _, subsystem := range common.ValidCephSubSystems {
		if subsystem == args[0] {
			isValidSubsytem = true
			break
		}
	}

	if !isValidSubsytem {
		return fmt.Errorf("invalid ceph subsystem %q. Supported subsystems : %+v", args[0], common.ValidCephSubSystems)
	}

	// validate log level
	match, err := regexp.MatchString("^[1-9][0-9]?$", args[1])
	if err != nil {
		return err
	}

	if !match {
		return fmt.Errorf("invalid log level %q. Value must be in range [0, 99]", args[1])
	}

	return nil
}
