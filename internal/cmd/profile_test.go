package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestResolveProfileName(t *testing.T) {
	for _, tc := range []struct {
		name                        string
		flag, env, binding, current string
		want                        string
	}{
		{"all empty → default", "", "", "", "", "default"},
		{"current only", "", "", "", "hosting", "hosting"},
		{"binding beats current", "", "", "dnsprof", "cloud", "dnsprof"},
		{"env beats binding", "", "envprof", "dnsprof", "cloud", "envprof"},
		{"flag wins", "flagprof", "envprof", "dnsprof", "cloud", "flagprof"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := resolveProfileName(tc.flag, tc.env, tc.binding, tc.current); got != tc.want {
				t.Errorf("resolveProfileName(%q,%q,%q,%q) = %q, want %q",
					tc.flag, tc.env, tc.binding, tc.current, got, tc.want)
			}
		})
	}
}

func TestTopLevelGroup(t *testing.T) {
	for _, tc := range []struct {
		cmd  *cobra.Command
		want string
	}{
		{dnsMxCmd, "dns"},      // dns mx → dns
		{dnsRecordsCmd, "dns"}, // dns records → dns
		{vpsListCmd, "vps"},    // vps list → vps
		{configureCmd, "configure"},
		{rootCmd, ""}, // root has no group
	} {
		if got := topLevelGroup(tc.cmd); got != tc.want {
			t.Errorf("topLevelGroup(%s) = %q, want %q", tc.cmd.Name(), got, tc.want)
		}
	}
}

func TestCredKeyNamespacing(t *testing.T) {
	// The default profile keeps legacy unprefixed keys; others are namespaced.
	if got := credKey("default", "login"); got != "login" {
		t.Errorf("credKey(default,login) = %q, want login (backward-compat)", got)
	}
	if got := credKey("", "token"); got != "token" {
		t.Errorf("credKey(\"\",token) = %q, want token", got)
	}
	if got := credKey("hosting", "password"); got != "hosting.password" {
		t.Errorf("credKey(hosting,password) = %q, want hosting.password", got)
	}
}

func TestConfigStoreRoundTrip(t *testing.T) {
	configDirOverride = t.TempDir()
	t.Cleanup(func() { configDirOverride = "" })

	if err := configSet("hosting", "current_profile"); err != nil {
		t.Fatalf("configSet current_profile: %v", err)
	}
	if err := configSet("hosting", "bindings", "dns"); err != nil {
		t.Fatalf("configSet binding: %v", err)
	}
	if err := configSet("https://staging.example", "profiles", "hosting", "endpoint"); err != nil {
		t.Fatalf("configSet endpoint: %v", err)
	}

	// Writes must not clobber earlier keys.
	if got := configGetString("current_profile"); got != "hosting" {
		t.Errorf("current_profile = %q, want hosting", got)
	}
	if got := configGetString("bindings", "dns"); got != "hosting" {
		t.Errorf("bindings.dns = %q, want hosting", got)
	}
	if got := configGetString("profiles", "hosting", "endpoint"); got != "https://staging.example" {
		t.Errorf("endpoint = %q, want https://staging.example", got)
	}

	// deleteProfileConfig clears the profile and its binding, resets current.
	if err := deleteProfileConfig("hosting"); err != nil {
		t.Fatalf("deleteProfileConfig: %v", err)
	}
	if got := configGetString("current_profile"); got != defaultProfile {
		t.Errorf("current after delete = %q, want %q", got, defaultProfile)
	}
	if got := configGetString("bindings", "dns"); got != "" {
		t.Errorf("bindings.dns after delete = %q, want empty", got)
	}
}
