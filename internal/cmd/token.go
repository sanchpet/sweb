package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

const keyringService = "sweb"

// Keyring accounts / config-file fields. SpaceWeb tokens are short-lived and the
// API has no refresh-token flow, so we store login+password (not just the token)
// to let the SDK transparently re-authenticate.
const (
	keyLogin    = "login"
	keyPassword = "password"
	keyToken    = "token"
)

func tokenFilePath() string { return filepath.Join(configDir(), "config.yaml") }

// storeCredentials persists login, password and the initial token, preferring
// the OS keyring (macOS Keychain, Linux Secret Service, Windows Credential
// Manager). Falls back to a 0600 file when the keyring is unavailable or
// insecure is true. Returns where it landed and whether it fell back, so the
// caller can surface the fallback explicitly (never silently).
func storeCredentials(login, password, token string, insecure bool) (location string, fellBack bool, err error) {
	if !insecure {
		e1 := keyring.Set(keyringService, keyLogin, login)
		e2 := keyring.Set(keyringService, keyPassword, password)
		e3 := keyring.Set(keyringService, keyToken, token)
		if e1 == nil && e2 == nil && e3 == nil {
			return "OS keyring", false, nil
		}
		fellBack = true
		location = "keyring unavailable → "
	}

	dir := configDir()
	if err = os.MkdirAll(dir, 0o700); err != nil {
		return "", fellBack, err
	}
	if err = writeCredFile(login, password, token); err != nil {
		return "", fellBack, err
	}
	return location + tokenFilePath() + " (plaintext, 0600)", fellBack, nil
}

// saveToken updates only the cached token — used by the SDK's refresh callback
// after a transparent re-authentication.
func saveToken(token string) error {
	if err := keyring.Set(keyringService, keyToken, token); err == nil {
		return nil
	}
	login, password, _ := loadCredentials()
	if err := os.MkdirAll(configDir(), 0o700); err != nil {
		return err
	}
	return writeCredFile(login, password, token)
}

func writeCredFile(login, password, token string) error {
	content := fmt.Sprintf("login: %s\npassword: %s\ntoken: %s\n", login, password, token)
	return os.WriteFile(tokenFilePath(), []byte(content), 0o600)
}

// loadCredentials reads login/password/token from the keyring, falling back to
// the config file for any field the keyring does not hold.
func loadCredentials() (login, password, token string) {
	login, _ = keyring.Get(keyringService, keyLogin)
	password, _ = keyring.Get(keyringService, keyPassword)
	token, _ = keyring.Get(keyringService, keyToken)
	if login == "" {
		login = viper.GetString(keyLogin)
	}
	if password == "" {
		password = viper.GetString(keyPassword)
	}
	if token == "" {
		token = viper.GetString(keyToken)
	}
	return login, password, token
}
