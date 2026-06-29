package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	sweb "github.com/sanchpet/sweb-go-sdk"
	"github.com/spf13/cobra"
)

var configureInsecure bool

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Authenticate and store an API token",
	Long: `Prompt for your SpaceWeb login and password, exchange them for a token,
and store login + password + token in the OS keyring (macOS Keychain, Linux
Secret Service, Windows Credential Manager). Falls back to a 0600 config file
when no keyring is available.

SpaceWeb tokens are short-lived and the API has no refresh-token flow, so the
password is stored (in the keyring) to re-authenticate transparently — you are
not prompted again until your password changes.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		var login, password string
		form := huh.NewForm(huh.NewGroup(
			huh.NewInput().Title("SpaceWeb login").Value(&login),
			huh.NewInput().Title("SpaceWeb password").EchoMode(huh.EchoModePassword).Value(&password),
		))
		if err := form.Run(); err != nil {
			return err
		}

		token, err := sweb.New().CreateToken(cmd.Context(), login, password)
		if err != nil {
			return err
		}

		where, fellBack, err := storeCredentials(login, password, token, configureInsecure)
		if err != nil {
			return err
		}
		if fellBack {
			fmt.Fprintln(cmd.ErrOrStderr(), "warning: OS keyring unavailable — credentials written to a plaintext file")
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Credentials stored:", where)
		return nil
	},
}

func init() {
	configureCmd.Flags().BoolVar(&configureInsecure, "insecure-storage", false,
		"store the token in a plaintext 0600 file instead of the OS keyring")
	rootCmd.AddCommand(configureCmd)
}
