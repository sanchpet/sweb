package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	sweb "github.com/sanchpet/sweb-go-sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Authenticate and store an API token",
	Long: `Prompt for your SpaceWeb login and password, exchange them for a personal
access token, and store the token in ~/.config/sweb/config.yaml (mode 0600).
Only the token is stored — not your password.`,
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

		dir := configDir()
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
		path := filepath.Join(dir, "config.yaml")
		viper.Set("token", token)
		if err := viper.WriteConfigAs(path); err != nil {
			return err
		}
		if err := os.Chmod(path, 0o600); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Token stored in", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
