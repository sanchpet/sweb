// Package cmd wires the sweb CLI commands over the sweb-go-sdk client.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/charmbracelet/fang"
	sweb "github.com/sanchpet/sweb-go-sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:          "sweb",
	Short:        "CLI for the SpaceWeb (sweb.ru) hosting API",
	SilenceUsage: true,
}

// Execute runs the root command with Fang (styled help/errors/version).
func Execute() {
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("config", "", "config file (default ~/.config/sweb/config.yaml)")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "output format: table|json")
	rootCmd.PersistentFlags().String("token", "", "API token (overrides config and $SWEB_TOKEN)")
	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "sweb")
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

// client builds an authenticated SDK client from the resolved token.
func client() (*sweb.Client, error) {
	token := viper.GetString("token")
	if token == "" {
		return nil, fmt.Errorf("no API token: run `sweb configure` or set SWEB_TOKEN")
	}
	return sweb.New(sweb.WithToken(token)), nil
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
