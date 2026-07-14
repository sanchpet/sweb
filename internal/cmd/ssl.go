package cmd

import "github.com/spf13/cobra"

// cloudSSLFlagString reads a string flag, ignoring the lookup error (unset → "").
func cloudSSLFlagString(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetString(name)
	return v
}

// cloudSSLCmd groups the cloud/account SSL-certificate service (SDK /vps/ssl):
// list and order certificates, download an issued archive, inspect and toggle
// prolongation, and remove a certificate. The Go identifier avoids the `sslCmd`
// used by the shared-hosting group (`hosting ssl`, SDK /vh/ssl); the two are
// distinct SpaceWeb services.
var cloudSSLCmd = &cobra.Command{
	Use:   "ssl",
	Short: "Manage cloud/account SSL certificates — order, download, prolong; distinct from `hosting ssl`",
}

func init() {
	rootCmd.AddCommand(cloudSSLCmd)
}
