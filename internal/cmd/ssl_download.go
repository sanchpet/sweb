package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/sanchpet/sweb-go-sdk/ssl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cloudSSLDownloadCmd = &cobra.Command{
	Use:   "download <id>",
	Short: "Download an issued certificate's archive files",
	Long: `Download an issued certificate's files via the "download" method.

<id> is a certificate id (see 'sweb ssl list'); the account password is required
and read from --password or a masked prompt. The archive carries the private key
— by default the files are written to --dir (current directory); with -o json the
base64 content is emitted instead. Handle the key material accordingly.`,
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
		password, err := cloudSSLResolvePassword(cmd)
		if err != nil {
			return err
		}
		files, err := c.SSL.Download(cmd.Context(), id, password)
		if err != nil {
			return err
		}
		// -o json emits the archive verbatim (base64 content included).
		if viper.GetString("output") == "json" {
			return render(cmd, files, nil)
		}
		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			dir = "."
		}
		for _, f := range files {
			body, err := cloudSSLDecodeContent(f)
			if err != nil {
				return fmt.Errorf("decode %s: %w", f.Name, err)
			}
			path := filepath.Join(dir, filepath.Base(f.Name))
			if err := os.WriteFile(path, body, 0o600); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Wrote", path)
		}
		return nil
	},
}

// cloudSSLResolvePassword returns the account password for a download: the
// --password flag when set, otherwise a masked interactive prompt. An empty
// result (non-interactive with no flag) is an error, since the API requires one.
func cloudSSLResolvePassword(cmd *cobra.Command) (string, error) {
	if pw, _ := cmd.Flags().GetString("password"); pw != "" {
		return pw, nil
	}
	var pw string
	if err := huh.NewInput().
		Title("Account password").
		EchoMode(huh.EchoModePassword).
		Value(&pw).
		Run(); err != nil {
		return "", err
	}
	if pw == "" {
		return "", fmt.Errorf("a password is required: pass --password or enter one at the prompt")
	}
	return pw, nil
}

// cloudSSLDecodeContent decodes a downloaded certificate file body,
// base64-decoding it when the mimetype declares base64 encoding and returning it
// verbatim otherwise.
func cloudSSLDecodeContent(f ssl.CertFile) ([]byte, error) {
	if strings.Contains(strings.ToLower(f.Mimetype), "base64") {
		return base64.StdEncoding.DecodeString(f.Content)
	}
	return []byte(f.Content), nil
}

func init() {
	cloudSSLDownloadCmd.Flags().String("password", "", "account password (prompted if unset)")
	cloudSSLDownloadCmd.Flags().String("dir", "", "directory to write the archive files (default current)")

	cloudSSLCmd.AddCommand(cloudSSLDownloadCmd)
}
