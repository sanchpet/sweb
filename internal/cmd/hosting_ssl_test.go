package cmd

import (
	"testing"

	"github.com/sanchpet/sweb-go-sdk/vh/ssl"
)

// sslStubLister feeds sslResolveID a fixed certificate set for the resolver tests.
type sslStubLister struct {
	certs []ssl.Certificate
	err   error
}

func (s sslStubLister) CertList() ([]ssl.Certificate, error) { return s.certs, s.err }

func TestSSLResolveID(t *testing.T) {
	certs := []ssl.Certificate{
		{ID: 10, Domain: "a.example"},
		{ID: 20, Domain: "b.example"},
		{ID: 30, Domain: "dup.example"},
		{ID: 31, Domain: "dup.example"},
	}

	id, err := sslResolveID(sslStubLister{certs: certs}, "b.example")
	if err != nil {
		t.Fatalf("resolve b.example: %v", err)
	}
	if id != 20 {
		t.Errorf("resolve b.example = %d, want 20", id)
	}

	if _, err := sslResolveID(sslStubLister{certs: certs}, "missing.example"); err == nil {
		t.Error("expected an error for an unknown domain")
	}

	if _, err := sslResolveID(sslStubLister{certs: certs}, "dup.example"); err == nil {
		t.Error("expected an error for an ambiguous domain")
	}
}

func TestHostingSSLCommandTree(t *testing.T) {
	hosting := findSub(rootCmd, "hosting")
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}
	sslGroup := findSub(hosting, "ssl")
	if sslGroup == nil {
		t.Fatal("hosting ssl command not registered")
	}
	for _, n := range []string{
		"list", "orders", "download", "prolong-info",
		"autoprolong", "prolong", "install-letsencrypt", "remove",
	} {
		if findSub(sslGroup, n) == nil {
			t.Errorf("hosting ssl is missing subcommand %q", n)
		}
	}
}
