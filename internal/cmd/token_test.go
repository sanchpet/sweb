package cmd

import (
	"testing"

	"github.com/zalando/go-keyring"
)

func TestStoreAndLoadCredentials(t *testing.T) {
	keyring.MockInit() // in-memory backend, cross-platform

	loc, fellBack, err := storeCredentials("user", "pass", "tok1", false)
	if err != nil {
		t.Fatalf("storeCredentials: %v", err)
	}
	if fellBack {
		t.Errorf("unexpected fallback to file: %s", loc)
	}

	login, password, token := loadCredentials()
	if login != "user" || password != "pass" || token != "tok1" {
		t.Errorf("loadCredentials = %q/%q/%q, want user/pass/tok1", login, password, token)
	}
}

func TestSaveTokenUpdatesCachedToken(t *testing.T) {
	keyring.MockInit()
	if _, _, err := storeCredentials("user", "pass", "tok1", false); err != nil {
		t.Fatalf("storeCredentials: %v", err)
	}

	if err := saveToken("tok2"); err != nil {
		t.Fatalf("saveToken: %v", err)
	}

	login, password, token := loadCredentials()
	if token != "tok2" {
		t.Errorf("token = %q, want tok2", token)
	}
	if login != "user" || password != "pass" {
		t.Errorf("credentials not preserved: %q/%q", login, password)
	}
}
