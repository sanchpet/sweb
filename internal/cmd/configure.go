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
	Long: `Prompt for your SpaceWeb login and password, exchange them for a personal
access token, and store the token in the OS keyring (macOS Keychain, Linux
Secret Service, Windows Credential Manager). Falls back to a 0600 config file
when no keyring is available. Only the token is stored — never your password.`,
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

		where, fellBack, err := storeToken(token, configureInsecure)
		if err != nil {
			return err
		}
		if fellBack {
			fmt.Fprintln(cmd.ErrOrStderr(), "warning: OS keyring unavailable — token written to a plaintext file")
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Token stored:", where)
		return nil
	},
}

func init() {
	configureCmd.Flags().BoolVar(&configureInsecure, "insecure-storage", false,
		"store the token in a plaintext 0600 file instead of the OS keyring")
	rootCmd.AddCommand(configureCmd)
}
