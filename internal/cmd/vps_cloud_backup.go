package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/spf13/cobra"
)

var vpsCloudBackupCmd = &cobra.Command{
	Use:     "cloud-backup",
	Aliases: []string{"cbackup"},
	Short:   "Manage off-node cloud backups",
}

var vpsCloudBackupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cloud backups on the account",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		list, err := c.RemoteBackup.List(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tVPS\tNAME\tSTATUS\tSIZE\tCREATED\tCOMMENT")
			for _, b := range list {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%s\t%s\n",
					int64(b.ID), b.BillingID, b.Name, b.Status, int64(b.Size), b.TSCreate, b.Comment)
			}
		})
	},
}

var vpsCloudBackupCreateCmd = &cobra.Command{
	Use:               "create <vps>",
	Short:             "Take a new cloud backup of a VPS",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeBillingIDs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name is required")
		}
		comment, _ := cmd.Flags().GetString("comment")
		if _, err := c.RemoteBackup.Create(cmd.Context(), billingID, name, comment); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Cloud backup %q of %s started\n", name, billingID)
		return nil
	},
}

var vpsCloudBackupRestoreCmd = &cobra.Command{
	Use:   "restore <id>",
	Short: "Restore a cloud backup — destructive",
	Long: `Restore a cloud backup via "restore" (into its source VPS) or "restoreInto"
(--into a different VPS). <id> is a cloud-backup id from 'sweb vps cloud-backup
list'. DESTRUCTIVE: overwrites the target disk. Confirms unless --yes.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("id must be an integer: %w", err)
		}
		into, _ := cmd.Flags().GetString("into")
		target := "its source VPS"
		var billingID string
		if into != "" {
			if billingID, err = resolveVPS(cmd.Context(), c, into); err != nil {
				return err
			}
			target = billingID
		}
		if !confirmed(cmd, fmt.Sprintf("Restore cloud backup %d into %s? This overwrites its disk.", id, target), "Restore") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if into != "" {
			err = c.RemoteBackup.RestoreInto(cmd.Context(), id, billingID)
		} else {
			err = c.RemoteBackup.Restore(cmd.Context(), id)
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Restoring cloud backup %d into %s\n", id, target)
		return nil
	},
}

var vpsCloudBackupRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Delete a cloud backup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("id must be an integer: %w", err)
		}
		if !confirmed(cmd, fmt.Sprintf("Delete cloud backup %d?", id), "Delete") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.RemoteBackup.Remove(cmd.Context(), id); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Deleted cloud backup %d\n", id)
		return nil
	},
}

var vpsCloudBackupCommentCmd = &cobra.Command{
	Use:   "comment <id> <comment>",
	Short: "Edit a cloud backup's comment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("id must be an integer: %w", err)
		}
		if err := c.RemoteBackup.EditComment(cmd.Context(), id, args[1]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Updated comment of cloud backup %d\n", id)
		return nil
	},
}

func init() {
	vpsCloudBackupCreateCmd.Flags().String("name", "", "name for the cloud backup (required)")
	vpsCloudBackupCreateCmd.Flags().String("comment", "", "optional comment")
	vpsCloudBackupRestoreCmd.Flags().String("into", "", "restore into a different VPS (name or billing id) instead of the source")
	vpsCloudBackupRestoreCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	_ = vpsCloudBackupRestoreCmd.RegisterFlagCompletionFunc("into", completeBillingIDs)
	vpsCloudBackupRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	vpsCloudBackupCmd.AddCommand(vpsCloudBackupListCmd, vpsCloudBackupCreateCmd,
		vpsCloudBackupRestoreCmd, vpsCloudBackupRemoveCmd, vpsCloudBackupCommentCmd)
	vpsCmd.AddCommand(vpsCloudBackupCmd)
}
