package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// defaultProfile is the profile used when none is selected. For backward
// compatibility its keyring keys are unprefixed, so a pre-profiles install keeps
// working as the default profile with no migration.
const defaultProfile = "default"

// apiCommandGroups are the top-level command groups a profile can be bound to
// (SpaceWeb serves both panels from one api.sweb.ru, so a binding only changes
// which credentials are used, never the endpoint).
var apiCommandGroups = []string{"vps", "dns", "domains"}

// activeProfile is the profile resolved for the running command. It is set once
// by the root PersistentPreRunE before any RunE executes, so client() need not
// know the command it serves.
var activeProfile string

// configDirOverride redirects the config directory in tests.
var configDirOverride string

func configDir() string {
	if configDirOverride != "" {
		return configDirOverride
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "sweb")
}

func configPath() string { return filepath.Join(configDir(), "config.yaml") }

// resolveProfileName applies the profile precedence: explicit --profile flag,
// then $SWEB_PROFILE, then the running command group's binding, then the
// configured current profile, then "default".
func resolveProfileName(flag, env, binding, current string) string {
	for _, v := range []string{flag, env, binding, current} {
		if v != "" {
			return v
		}
	}
	return defaultProfile
}

// topLevelGroup returns the name of cmd's top-level ancestor (the child of root),
// e.g. "dns" for `dns mx`. Empty for the root command itself. This is the key a
// profile binding is looked up by.
func topLevelGroup(cmd *cobra.Command) string {
	if cmd == nil {
		return ""
	}
	for cmd.Parent() != nil && cmd.Parent() != rootCmd {
		cmd = cmd.Parent()
	}
	if cmd.Parent() == nil { // cmd is root
		return ""
	}
	return cmd.Name()
}

// --- config store (case-preserving YAML, unlike viper which lowercases keys) ---

// loadConfigMap reads config.yaml into a nested map, or an empty map when it is
// missing or unreadable.
func loadConfigMap() map[string]any {
	b, err := os.ReadFile(configPath())
	if err != nil {
		return map[string]any{}
	}
	var m map[string]any
	if err := yaml.Unmarshal(b, &m); err != nil || m == nil {
		return map[string]any{}
	}
	return m
}

// saveConfigMap writes the config back as 0600 YAML, creating the dir 0700.
func saveConfigMap(m map[string]any) error {
	if err := os.MkdirAll(configDir(), 0o700); err != nil {
		return err
	}
	b, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), b, 0o600)
}

// configGetString reads a nested string value (empty if any level is absent).
func configGetString(keys ...string) string {
	var cur any = loadConfigMap()
	for _, k := range keys {
		m, ok := cur.(map[string]any)
		if !ok {
			return ""
		}
		cur = m[k]
	}
	s, _ := cur.(string)
	return s
}

// configGetMap reads a nested map (nil if absent).
func configGetMap(keys ...string) map[string]any {
	var cur any = loadConfigMap()
	for _, k := range keys {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[k]
	}
	m, _ := cur.(map[string]any)
	return m
}

// configSet writes a nested value, creating intermediate maps, preserving all
// other keys in the file.
func configSet(value any, keys ...string) error {
	m := loadConfigMap()
	cur := m
	for _, k := range keys[:len(keys)-1] {
		next, ok := cur[k].(map[string]any)
		if !ok {
			next = map[string]any{}
			cur[k] = next
		}
		cur = next
	}
	cur[keys[len(keys)-1]] = value
	return saveConfigMap(m)
}

// registerProfile ensures a profile exists in config and, if no current profile
// is set, makes it current. Never overwrites an existing profile's settings.
func registerProfile(name string) error {
	m := loadConfigMap()
	profiles, _ := m["profiles"].(map[string]any)
	if profiles == nil {
		profiles = map[string]any{}
		m["profiles"] = profiles
	}
	if _, ok := profiles[name]; !ok {
		profiles[name] = map[string]any{}
	}
	if s, _ := m["current_profile"].(string); s == "" {
		m["current_profile"] = name
	}
	return saveConfigMap(m)
}

// deleteProfileConfig removes a profile and any references to it (bindings,
// fallback credentials); if it was current, resets current to default.
func deleteProfileConfig(name string) error {
	m := loadConfigMap()
	if profiles, ok := m["profiles"].(map[string]any); ok {
		delete(profiles, name)
	}
	if creds, ok := m["credentials"].(map[string]any); ok {
		delete(creds, name)
	}
	if bindings, ok := m["bindings"].(map[string]any); ok {
		for g, p := range bindings {
			if p == name {
				delete(bindings, g)
			}
		}
	}
	if s, _ := m["current_profile"].(string); s == name {
		m["current_profile"] = defaultProfile
	}
	return saveConfigMap(m)
}
