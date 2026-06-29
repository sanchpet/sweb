package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

const (
	keyringService = "sweb"
	keyringUser    = "default"
)

func tokenFilePath() string { return filepath.Join(configDir(), "config.yaml") }

// storeToken persists the token, preferring the OS keyring (macOS Keychain,
// Linux Secret Service, Windows Credential Manager). It falls back to a 0600
// config file when the keyring is unavailable (e.g. headless Linux) or when
// insecure is true. It returns a human-readable location and whether it fell
// back — the caller surfaces that explicitly, never a silent downgrade.
func storeToken(token string, insecure bool) (location string, fellBack bool, err error) {
	if !insecure {
		kerr := keyring.Set(keyringService, keyringUser, token)
		if kerr == nil {
			return "OS keyring", false, nil
		}
		fellBack = true
		location = fmt.Sprintf("keyring unavailable (%v) → ", kerr)
	}

	dir := configDir()
	if err = os.MkdirAll(dir, 0o700); err != nil {
		return "", fellBack, err
	}
	path := tokenFilePath()
	if err = os.WriteFile(path, []byte("token: "+token+"\n"), 0o600); err != nil {
		return "", fellBack, err
	}
	return location + path + " (plaintext, 0600)", fellBack, nil
}

// resolveToken returns the API token by precedence:
// --token flag → $SWEB_TOKEN → OS keyring → config file.
func resolveToken() string {
	if t, _ := rootCmd.PersistentFlags().GetString("token"); t != "" {
		return t
	}
	if t := os.Getenv("SWEB_TOKEN"); t != "" {
		return t
	}
	if t, err := keyring.Get(keyringService, keyringUser); err == nil && t != "" {
		return t
	}
	return viper.GetString("token") // from the config file, if any
}
