package cmd

import (
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

// credKey namespaces a credential field by profile in the keyring. The default
// profile keeps the legacy unprefixed keys so a pre-profiles install is read
// unchanged; other profiles get a "<profile>." prefix.
func credKey(profile, field string) string {
	if profile == "" || profile == defaultProfile {
		return field
	}
	return profile + "." + field
}

// storeCredentials persists a profile's login, password and initial token,
// preferring the OS keyring (macOS Keychain, Linux Secret Service, Windows
// Credential Manager). Falls back to the 0600 config file (under
// credentials.<profile>) when the keyring is unavailable or insecure is true.
// Returns where it landed and whether it fell back, so the caller can surface
// the fallback explicitly (never silently).
func storeCredentials(profile, login, password, token string, insecure bool) (location string, fellBack bool, err error) {
	if !insecure {
		e1 := keyring.Set(keyringService, credKey(profile, keyLogin), login)
		e2 := keyring.Set(keyringService, credKey(profile, keyPassword), password)
		e3 := keyring.Set(keyringService, credKey(profile, keyToken), token)
		if e1 == nil && e2 == nil && e3 == nil {
			return "OS keyring", false, nil
		}
		fellBack = true
		location = "keyring unavailable → "
	}

	if err = configSet(map[string]any{
		keyLogin:    login,
		keyPassword: password,
		keyToken:    token,
	}, "credentials", profile); err != nil {
		return "", fellBack, err
	}
	return location + configPath() + " (plaintext, 0600)", fellBack, nil
}

// saveToken updates only a profile's cached token — used by the SDK's refresh
// callback after a transparent re-authentication.
func saveToken(profile, token string) error {
	if err := keyring.Set(keyringService, credKey(profile, keyToken), token); err == nil {
		return nil
	}
	return configSet(token, "credentials", profile, keyToken)
}

// loadCredentials reads a profile's login/password/token from the keyring,
// falling back to the config file for any field the keyring does not hold.
func loadCredentials(profile string) (login, password, token string) {
	login, _ = keyring.Get(keyringService, credKey(profile, keyLogin))
	password, _ = keyring.Get(keyringService, credKey(profile, keyPassword))
	token, _ = keyring.Get(keyringService, credKey(profile, keyToken))
	if login == "" {
		login = credFileField(profile, keyLogin)
	}
	if password == "" {
		password = credFileField(profile, keyPassword)
	}
	if token == "" {
		token = credFileField(profile, keyToken)
	}
	return login, password, token
}

// credFileField reads a fallback credential from the config file. The default
// profile also honours legacy top-level keys written before profiles existed.
func credFileField(profile, field string) string {
	if v := configGetString("credentials", profile, field); v != "" {
		return v
	}
	if profile == defaultProfile {
		return configGetString(field)
	}
	return ""
}

// removeCredentials deletes a profile's stored credentials from the keyring
// (config-file entries are cleared by deleteProfileConfig).
func removeCredentials(profile string) {
	_ = keyring.Delete(keyringService, credKey(profile, keyLogin))
	_ = keyring.Delete(keyringService, credKey(profile, keyPassword))
	_ = keyring.Delete(keyringService, credKey(profile, keyToken))
}
