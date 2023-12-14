package purgeosd

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/red-hat-storage/odf-cli/cmd/odf/root"
	"github.com/red-hat-storage/odf-cli/pkg/rook/osd"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/rook/kubectl-rook-ceph/pkg/mons"
	"github.com/rook/kubectl-rook-ceph/pkg/rook"
	"github.com/spf13/cobra"
)

var forceFlag = "force"

var CephPurgeOsdCmd = &cobra.Command{
	Use:     "purge-osd",
	Short:   "Permanently remove an OSD from the cluster.",
	Args:    cobra.ExactArgs(1),
	Example: "odf purge-osd <ID>",
	PreRunE: validateOsdID,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		forceflagValue := "false"
		osdID := args[0]

		safeToDestroy, err := osd.SafeToDestroy(ctx, root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, osdID)
		if !safeToDestroy {
			var answer string
			yesForceDestroyOSD := "yes-force-destroy-osd"
			logging.Warning("Are you sure you want to purge osd.%s? The OSD is *not* safe to destroy. This may lead to data loss. If you are sure the OSD should be purged, enter '%s'", osdID, yesForceDestroyOSD)
			fmt.Scanf("%s", &answer)

			err := mons.PromptToContinueOrCancel("osd-purge", yesForceDestroyOSD, answer)
			if err != nil {
				logging.Fatal(err)
			}

			forceflagValue = "true"
		} else if err != nil {
			logging.Fatal(errors.Wrapf(err, "failed to check if osd.%s is safe to destroy", osdID))
		}

		rook.PurgeOsd(ctx, root.ClientSets, root.OperatorNamespace, root.StorageClusterNamespace, osdID, forceflagValue)
	},
}

func validateOsdID(cmd *cobra.Command, args []string) error {
	osdID := args[0]
	_, err := strconv.Atoi(osdID)
	if err != nil {
		return fmt.Errorf("Invalid ID %s, the OSD ID must be an integer. %v", osdID, err)
	}

	return nil
}
