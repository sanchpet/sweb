package cmd

import (
	"testing"

	"github.com/zalando/go-keyring"
)

func TestStoreAndLoadCredentials(t *testing.T) {
	keyring.MockInit() // in-memory backend, cross-platform

	loc, fellBack, err := storeCredentials(defaultProfile, "user", "pass", "tok1", false)
	if err != nil {
		t.Fatalf("storeCredentials: %v", err)
	}
	if fellBack {
		t.Errorf("unexpected fallback to file: %s", loc)
	}

	login, password, token := loadCredentials(defaultProfile)
	if login != "user" || password != "pass" || token != "tok1" {
		t.Errorf("loadCredentials = %q/%q/%q, want user/pass/tok1", login, password, token)
	}
}

func TestSaveTokenUpdatesCachedToken(t *testing.T) {
	keyring.MockInit()
	if _, _, err := storeCredentials(defaultProfile, "user", "pass", "tok1", false); err != nil {
		t.Fatalf("storeCredentials: %v", err)
	}

	if err := saveToken(defaultProfile, "tok2"); err != nil {
		t.Fatalf("saveToken: %v", err)
	}

	login, password, token := loadCredentials(defaultProfile)
	if token != "tok2" {
		t.Errorf("token = %q, want tok2", token)
	}
	if login != "user" || password != "pass" {
		t.Errorf("credentials not preserved: %q/%q", login, password)
	}
}

// TestProfilesAreIsolated verifies two profiles' credentials do not collide in
// the keyring — the core multi-account guarantee.
func TestProfilesAreIsolated(t *testing.T) {
	keyring.MockInit()
	if _, _, err := storeCredentials("cloud", "cloud-user", "cloud-pass", "cloud-tok", false); err != nil {
		t.Fatalf("store cloud: %v", err)
	}
	if _, _, err := storeCredentials("hosting", "host-user", "host-pass", "host-tok", false); err != nil {
		t.Fatalf("store hosting: %v", err)
	}

	if l, _, tok := loadCredentials("cloud"); l != "cloud-user" || tok != "cloud-tok" {
		t.Errorf("cloud = %q/%q, want cloud-user/cloud-tok", l, tok)
	}
	if l, _, tok := loadCredentials("hosting"); l != "host-user" || tok != "host-tok" {
		t.Errorf("hosting = %q/%q, want host-user/host-tok", l, tok)
	}

	// Removing one must not touch the other.
	removeCredentials("cloud")
	if l, _, _ := loadCredentials("cloud"); l != "" {
		t.Errorf("cloud login after remove = %q, want empty", l)
	}
	if l, _, _ := loadCredentials("hosting"); l != "host-user" {
		t.Errorf("hosting login after removing cloud = %q, want host-user", l)
	}
}
