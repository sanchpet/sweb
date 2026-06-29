package cmd

import (
	"testing"

	"github.com/zalando/go-keyring"
)

func TestStoreAndResolveViaKeyring(t *testing.T) {
	keyring.MockInit() // in-memory backend, cross-platform
	t.Setenv("SWEB_TOKEN", "")

	loc, fellBack, err := storeToken("tok_keyring", false)
	if err != nil {
		t.Fatalf("storeToken: %v", err)
	}
	if fellBack {
		t.Errorf("unexpected fallback to file: %s", loc)
	}
	if got := resolveToken(); got != "tok_keyring" {
		t.Errorf("resolveToken = %q, want tok_keyring", got)
	}
}

func TestEnvOverridesKeyring(t *testing.T) {
	keyring.MockInit()
	if err := keyring.Set(keyringService, keyringUser, "tok_keyring"); err != nil {
		t.Fatalf("seed keyring: %v", err)
	}
	t.Setenv("SWEB_TOKEN", "tok_env")

	if got := resolveToken(); got != "tok_env" {
		t.Errorf("resolveToken = %q, want tok_env (env must beat keyring)", got)
	}
}
