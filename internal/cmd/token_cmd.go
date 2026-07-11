package cmd

import (
	"fmt"
	"os"

	sweb "github.com/sanchpet/sweb-go-sdk"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Print a fresh API token to stdout (for env-based tooling)",
	Long: `Mint a fresh SpaceWeb API token from the stored credentials and print it
to stdout — nothing else — so it can feed env-based tooling, e.g.

  export SWEB_TOKEN=$(sweb token)

or a mise/direnv env, analogous to 'yc iam create-token'. Credentials come from
'sweb configure' (OS keyring / config file) or $SWEB_LOGIN + $SWEB_PASSWORD.

SpaceWeb tokens are short-lived, so mint on demand rather than caching in a
dotfile. The freshly minted token is also written back to the credential store.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		profile := activeProfile
		if profile == "" {
			profile = defaultProfile
		}
		login, password := os.Getenv("SWEB_LOGIN"), os.Getenv("SWEB_PASSWORD")
		var cached string
		if login == "" || password == "" {
			login, password, cached = loadCredentials(profile)
		}

		if login != "" && password != "" {
			token, err := sweb.New().CreateToken(cmd.Context(), login, password)
			if err != nil {
				return err
			}
			_ = saveToken(profile, token) // keep the store's cached token fresh; best-effort
			fmt.Fprintln(cmd.OutOrStdout(), token)
			return nil
		}

		if cached != "" { // no credentials to mint with, but a token was stored
			fmt.Fprintln(cmd.OutOrStdout(), cached)
			return nil
		}
		return fmt.Errorf("no credentials: run `sweb configure` or set SWEB_LOGIN and SWEB_PASSWORD")
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
