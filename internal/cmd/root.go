// Package cmd wires the sweb CLI commands over the sweb-go-sdk client.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"text/tabwriter"

	"github.com/charmbracelet/fang"
	sweb "github.com/sanchpet/sweb-go-sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// version is injected at release time via -ldflags by GoReleaser. When empty
// (go install / local build) Fang falls back to the module build info.
var version string

var rootCmd = &cobra.Command{
	Use:          "sweb",
	Short:        "CLI for the SpaceWeb (sweb.ru) hosting API",
	SilenceUsage: true,
}

// resolveActiveProfile resolves the credential profile for the running command.
// It is wired as the root PersistentPreRunE (in init, to avoid a rootCmd↔
// topLevelGroup initialization cycle) so it runs once — after cobra.OnInitialize
// loads config — before any command's RunE, letting client() stay parameterless.
func resolveActiveProfile(cmd *cobra.Command, _ []string) error {
	flag, _ := cmd.Flags().GetString("profile")
	var binding string
	if group := topLevelGroup(cmd); group != "" {
		binding = configGetString("bindings", group)
	}
	activeProfile = resolveProfileName(flag, os.Getenv("SWEB_PROFILE"), binding, configGetString("current_profile"))
	return nil
}

// Execute runs the root command with Fang (styled help/errors/version).
func Execute() {
	if err := fang.Execute(context.Background(), rootCmd, fang.WithVersion(versionString())); err != nil {
		os.Exit(1)
	}
}

// versionString resolves the reported version: the GoReleaser-injected value,
// else the module version from build info (go install), else "dev".
func versionString() string {
	if version != "" {
		return version
	}
	if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}
	return "dev"
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentPreRunE = resolveActiveProfile
	rootCmd.PersistentFlags().String("config", "", "config file (default ~/.config/sweb/config.yaml)")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "output format: table|json")
	rootCmd.PersistentFlags().String("token", "", "API token (overrides keyring/config and $SWEB_TOKEN)")
	rootCmd.PersistentFlags().String("profile", "", "credential profile to use (overrides $SWEB_PROFILE, group bindings, and the current profile)")
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	_ = rootCmd.RegisterFlagCompletionFunc("output", completeOutput)
	_ = rootCmd.RegisterFlagCompletionFunc("profile", completeProfiles)
}

func initConfig() {
	if cf, _ := rootCmd.PersistentFlags().GetString("config"); cf != "" {
		viper.SetConfigFile(cf)
	} else {
		viper.AddConfigPath(configDir())
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	viper.SetEnvPrefix("SWEB") // SWEB_TOKEN -> token
	viper.AutomaticEnv()
	_ = viper.ReadInConfig() // missing config is fine
}

// client builds an authenticated SDK client for the active profile.
func client() (*sweb.Client, error) {
	// Explicit token (flag or env): used as-is, no auto-refresh, no profile.
	if t, _ := rootCmd.PersistentFlags().GetString("token"); t != "" {
		return sweb.New(sweb.WithToken(t)), nil
	}
	if t := os.Getenv("SWEB_TOKEN"); t != "" {
		return sweb.New(sweb.WithToken(t)), nil
	}

	profile := activeProfile
	if profile == "" {
		profile = defaultProfile
	}
	login, password, token := loadCredentials(profile)

	var opts []sweb.Option
	if endpoint := configGetString("profiles", profile, "endpoint"); endpoint != "" {
		opts = append(opts, sweb.WithBaseURL(endpoint))
	}
	switch {
	case login != "" && password != "":
		// Cached token + credentials → the SDK refreshes transparently on
		// expiry and we persist the new token for this profile via the callback.
		opts = append(opts,
			sweb.WithToken(token),
			sweb.WithCredentials(login, password),
			sweb.WithOnTokenRefresh(func(t string) { _ = saveToken(profile, t) }),
		)
	case token != "":
		opts = append(opts, sweb.WithToken(token))
	default:
		hint := ""
		if profile != defaultProfile {
			hint = " --profile " + profile
		}
		return nil, fmt.Errorf("no credentials for profile %q: run `sweb configure%s` or set SWEB_TOKEN", profile, hint)
	}
	return sweb.New(opts...), nil
}

// render prints data as JSON (-o json) or via the supplied table writer.
func render(cmd *cobra.Command, data any, table func(io.Writer)) error {
	if viper.GetString("output") == "json" {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	}
	tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
	table(tw)
	return tw.Flush()
}
