package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var vpsBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage a VPS's local backups and auto-backup schedule",
}

var vpsBackupListCmd = &cobra.Command{
	Use:               "list <vps>",
	Short:             "List a VPS's local backups",
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
		list, err := c.Backup.List(cmd.Context(), billingID)
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "NAME\tLABEL\tTYPE\tUPDATED")
			for _, b := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", b.Name, b.PrettyName, b.AttachType, b.UpdatedAt)
			}
		})
	},
}

var vpsBackupCreateCmd = &cobra.Command{
	Use:               "create <vps>",
	Short:             "Take a new local backup",
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
		if err := c.Backup.Create(cmd.Context(), billingID); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Backup of %s started\n", billingID)
		return nil
	},
}

var vpsBackupRestoreCmd = &cobra.Command{
	Use:               "restore <vps> <name>",
	Short:             "Restore a VPS from a local backup — destructive",
	Long:              "Restore a VPS from a local backup via the \"restore\" method. <name> is a backup name from 'sweb vps backup list'. DESTRUCTIVE: overwrites the current disk. Confirms unless --yes.",
	Args:              cobra.ExactArgs(2),
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
		if !confirmed(cmd, fmt.Sprintf("Restore %q from backup %q? This overwrites the current disk.", billingID, args[1]), "Restore") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Backup.Restore(cmd.Context(), billingID, args[1]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Restoring %s from %s\n", billingID, args[1])
		return nil
	},
}

var vpsBackupRemoveCmd = &cobra.Command{
	Use:               "remove <vps> <name>",
	Short:             "Delete a local backup",
	Args:              cobra.ExactArgs(2),
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
		if !confirmed(cmd, fmt.Sprintf("Delete backup %q of %q?", args[1], billingID), "Delete") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Backup.Remove(cmd.Context(), billingID, args[1]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Deleted backup %s of %s\n", args[1], billingID)
		return nil
	},
}

var vpsBackupAttachCmd = &cobra.Command{
	Use:               "attach <vps> <name>",
	Short:             "Mount a backup on the VPS as an extra disk",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeBillingIDs,
	RunE:              backupMountRunE(false),
}

var vpsBackupDetachCmd = &cobra.Command{
	Use:               "detach <vps> <name>",
	Short:             "Unmount a previously attached backup disk",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeBillingIDs,
	RunE:              backupMountRunE(true),
}

// backupMountRunE builds the attach/detach RunE (they differ only in the SDK call
// and the message).
func backupMountRunE(detach bool) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}
		verb := "Attached"
		if detach {
			err = c.Backup.Detach(cmd.Context(), billingID, args[1])
			verb = "Detached"
		} else {
			err = c.Backup.Attach(cmd.Context(), billingID, args[1])
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s backup %s on %s\n", verb, args[1], billingID)
		return nil
	}
}

var vpsBackupSettingsCmd = &cobra.Command{
	Use:   "settings <vps>",
	Short: "Show or set the auto-backup schedule",
	Long: `Show a VPS's auto-backup schedule, or set it by passing flags.

With no flags, prints the current settings. Pass --mode (and --frequency/--time
for auto) to change them via the "saveSettings" method.`,
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
		f := cmd.Flags()
		if f.Changed("mode") || f.Changed("frequency") || f.Changed("time") {
			mode, _ := f.GetString("mode")
			freq, _ := f.GetInt("frequency")
			t, _ := f.GetInt("time")
			if mode == "" {
				mode = "auto"
			}
			if err := c.Backup.SaveSettings(cmd.Context(), billingID, mode, freq, t); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Set auto-backup of %s: mode=%s frequency=%d time=%d\n", billingID, mode, freq, t)
			return nil
		}
		set, err := c.Backup.Settings(cmd.Context(), billingID)
		if err != nil {
			return err
		}
		return render(cmd, set, func(w io.Writer) {
			if set == nil {
				fmt.Fprintln(w, "no auto-backup settings")
				return
			}
			fmt.Fprintf(w, "MODE\t%s\n", set.Mode)
			fmt.Fprintf(w, "FREQUENCY\t%d\n", int64(set.Frequency))
			fmt.Fprintf(w, "TIME\t%d\n", int64(set.Time))
			fmt.Fprintf(w, "NEXT\t%s\n", set.NextDataBackup)
		})
	},
}

func init() {
	vpsBackupRestoreCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	vpsBackupRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	sf := vpsBackupSettingsCmd.Flags()
	sf.String("mode", "", "auto-backup mode: auto|manual")
	sf.Int("frequency", 0, "auto-backup frequency (days)")
	sf.Int("time", 0, "auto-backup hour of day")

	vpsBackupCmd.AddCommand(vpsBackupListCmd, vpsBackupCreateCmd, vpsBackupRestoreCmd,
		vpsBackupRemoveCmd, vpsBackupAttachCmd, vpsBackupDetachCmd, vpsBackupSettingsCmd)
	vpsCmd.AddCommand(vpsBackupCmd)
}
