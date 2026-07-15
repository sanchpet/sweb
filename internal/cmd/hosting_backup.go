package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/vh/backup"
	"github.com/spf13/cobra"
)

// backupCmd groups the shared-hosting account-backup operations (SDK /vh/backup):
// listing daily backups, browsing their file/MySQL contents, and the
// restore/receive/download lifecycle over the account's home directory and
// databases. It hangs off the hosting parent, so it inherits that group's
// profile binding.
//
// This is the hosting-account backup, DISTINCT from `sweb vps backup`, which
// snapshots a whole VPS disk.
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage shared-hosting account backups",
}

// bakFileRefs builds the SDK's [type, path] targets from the repeatable --file
// (regular files) and --folder (directories) flags. At least one entry is
// required, since restoreFiles/receiveFiles/downloadFile act on named targets.
func bakFileRefs(cmd *cobra.Command) ([]backup.FileRef, error) {
	files, _ := cmd.Flags().GetStringArray("file")
	folders, _ := cmd.Flags().GetStringArray("folder")
	refs := make([]backup.FileRef, 0, len(files)+len(folders))
	for _, p := range files {
		refs = append(refs, backup.FileRef{Dir: false, Path: p})
	}
	for _, p := range folders {
		refs = append(refs, backup.FileRef{Dir: true, Path: p})
	}
	if len(refs) == 0 {
		return nil, fmt.Errorf("at least one --file or --folder is required")
	}
	return refs, nil
}

var backupDatesCmd = &cobra.Command{
	Use:   "dates",
	Short: "List the account's daily backups",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		list, err := c.VHBackup.List(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "DATE\tFILES\tMYSQL\tFILE-BACKUP\tOVER-QUOTA")
			for _, d := range list {
				fmt.Fprintf(w, "%s\t%d\t%d\t%s\t%s\n",
					d.Date, int64(d.Files), int64(d.Mysql),
					yesNo(d.BackupFilesExists), yesNo(d.WarnQuota))
			}
		})
	},
}

var backupFilesCmd = &cobra.Command{
	Use:   "files <date>",
	Short: "Browse a day's file backup",
	Long: `List the contents of a directory inside a day's file backup (method
"getListFiles").

<date> is the backup folder in the server's strict format ("2023-02-27", NOT the
display format from 'sweb hosting backup dates'). --dir is the path within the
backup (default "/").`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		dir, _ := cmd.Flags().GetString("dir")
		list, err := c.VHBackup.ListFiles(cmd.Context(), args[0], dir)
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "NAME\tTYPE\tSIZE")
			for _, f := range list {
				kind := "file"
				if f.Dir {
					kind = "dir"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", f.Name, kind, f.Size)
			}
		})
	},
}

var backupMysqlCmd = &cobra.Command{
	Use:   "mysql <date>",
	Short: "Browse a day's MySQL backup",
	Long: `List the contents of a day's MySQL backup (method "getListMysql").

<date> is the backup folder in the server's strict format ("2023-02-27"). --dir
is the path within the backup (default "/").`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		dir, _ := cmd.Flags().GetString("dir")
		list, err := c.VHBackup.ListMysql(cmd.Context(), args[0], dir)
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "NAME\tTYPE")
			for _, d := range list {
				kind := "db"
				if d.Dir {
					kind = "dir"
				}
				fmt.Fprintf(w, "%s\t%s\n", d.Name, kind)
			}
		})
	},
}

var backupSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Take a fresh backup of the account's files and databases",
	Long: `Queue a fresh backup of all databases and the account's home directory
(method "makeAccountCopy").

You are asked to confirm unless --yes is given.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, "Take a fresh backup of the account's files and databases?", "Back up") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.VHBackup.MakeAccountCopy(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintln(w, "Backup started")
		})
	},
}

var backupRestoreFilesCmd = &cobra.Command{
	Use:   "restore-files <date>",
	Short: "Restore files from a day's backup in place — destructive",
	Long: `Restore the given files/folders from a day's backup in place (method
"restoreFiles"). DESTRUCTIVE: overwrites the live files.

<date> is the backup folder in the server's strict format ("2023-02-27").
Pass targets with repeatable --file (regular file) and --folder (directory);
at least one is required. You are asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		refs, err := bakFileRefs(cmd)
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Restore %d target(s) from backup %q? This overwrites the live files.", len(refs), args[0]), "Restore") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.VHBackup.RestoreFiles(cmd.Context(), args[0], refs)
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintf(w, "Restoring files from %s\n", args[0])
		})
	},
}

var backupRestoreMysqlCmd = &cobra.Command{
	Use:   "restore-mysql <date>",
	Short: "Restore databases from a day's backup in place — destructive",
	Long: `Restore the named databases from a day's backup in place (method
"restoreMysql"). DESTRUCTIVE: overwrites the live databases.

<date> is the backup folder in the server's strict format ("2023-02-27").
Pass databases with repeatable --db; at least one is required. You are asked to
confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		dbs, _ := cmd.Flags().GetStringArray("db")
		if len(dbs) == 0 {
			return fmt.Errorf("at least one --db is required")
		}
		if !confirmed(cmd, fmt.Sprintf("Restore %d database(s) from backup %q? This overwrites the live databases.", len(dbs), args[0]), "Restore") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.VHBackup.RestoreMysql(cmd.Context(), args[0], dbs)
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintf(w, "Restoring databases from %s\n", args[0])
		})
	},
}

var backupReceiveFilesCmd = &cobra.Command{
	Use:   "receive-files <date>",
	Short: "Prepare files from a day's backup for download",
	Long: `Queue the given files/folders from a day's backup to be prepared for
download rather than restored in place (method "receiveFiles").

<date> is the backup folder in the server's strict format ("2023-02-27").
Pass targets with repeatable --file (regular file) and --folder (directory);
at least one is required.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		refs, err := bakFileRefs(cmd)
		if err != nil {
			return err
		}
		res, err := c.VHBackup.ReceiveFiles(cmd.Context(), args[0], refs)
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintf(w, "Preparing files from %s for download\n", args[0])
		})
	},
}

var backupReceiveMysqlCmd = &cobra.Command{
	Use:   "receive-mysql <date>",
	Short: "Prepare databases from a day's backup for download",
	Long: `Queue the named databases from a day's backup to be prepared for download
(method "receiveMysql").

<date> is the backup folder in the server's strict format ("2023-02-27").
Pass databases with repeatable --db; at least one is required.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		dbs, _ := cmd.Flags().GetStringArray("db")
		if len(dbs) == 0 {
			return fmt.Errorf("at least one --db is required")
		}
		res, err := c.VHBackup.ReceiveMysql(cmd.Context(), args[0], dbs)
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintf(w, "Preparing databases from %s for download\n", args[0])
		})
	},
}

var backupDownloadCmd = &cobra.Command{
	Use:   "download <date>",
	Short: "Download files from a day's backup",
	Long: `Download the given files from a day's backup (method "downloadFile").

<date> is the backup folder in the server's strict format ("2023-02-27").
Pass targets with repeatable --file (regular file) and --folder (directory);
at least one is required. The result carries base64-encoded content; prefer
-o json to capture it.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		refs, err := bakFileRefs(cmd)
		if err != nil {
			return err
		}
		res, err := c.VHBackup.DownloadFile(cmd.Context(), args[0], refs)
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintf(w, "Downloaded %d target(s) from %s (use -o json for the content)\n", len(refs), args[0])
		})
	},
}

func init() {
	backupFilesCmd.Flags().String("dir", "/", "path within the backup to list")
	backupMysqlCmd.Flags().String("dir", "/", "path within the backup to list")

	backupSnapshotCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	for _, c := range []*cobra.Command{backupRestoreFilesCmd, backupReceiveFilesCmd, backupDownloadCmd} {
		c.Flags().StringArray("file", nil, "a regular file to act on (repeatable)")
		c.Flags().StringArray("folder", nil, "a directory to act on (repeatable)")
	}
	for _, c := range []*cobra.Command{backupRestoreMysqlCmd, backupReceiveMysqlCmd} {
		c.Flags().StringArray("db", nil, "a database name to act on (repeatable)")
	}
	backupRestoreFilesCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	backupRestoreMysqlCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	backupCmd.AddCommand(
		backupDatesCmd,
		backupFilesCmd,
		backupMysqlCmd,
		backupSnapshotCmd,
		backupRestoreFilesCmd,
		backupRestoreMysqlCmd,
		backupReceiveFilesCmd,
		backupReceiveMysqlCmd,
		backupDownloadCmd,
	)
	hostingCmd.AddCommand(backupCmd)
}
