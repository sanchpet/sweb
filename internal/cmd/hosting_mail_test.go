package cmd

import (
	"testing"
)

func TestMailCommandTree(t *testing.T) {
	hosting := findSub(rootCmd, "hosting")
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}
	mail := findSub(hosting, "mail")
	if mail == nil {
		t.Fatal("hosting mail command not registered")
	}

	// Top-level mail groups + the two flat reads.
	msub := subNames(mail)
	for _, n := range []string{
		"domains", "quota", "mailbox", "forwarding", "autoreply",
		"dkim", "whitelist", "blacklist", "delivery", "collector", "domain",
	} {
		if !msub[n] {
			t.Errorf("hosting mail is missing subcommand %q", n)
		}
	}

	// Each group's leaves.
	for group, want := range map[string][]string{
		"mailbox":    {"list", "create", "delete", "password", "comment", "antispam", "spf", "purge", "requisites"},
		"forwarding": {"list", "add", "remove", "delete-after"},
		"autoreply":  {"get", "set"},
		"dkim":       {"enable", "disable"},
		"whitelist":  {"list", "add", "remove"},
		"blacklist":  {"list", "add", "remove"},
		"delivery":   {"info", "list", "add", "remove"},
		"collector":  {"get", "set", "remove", "confirm"},
		"domain":     {"spf", "sender-verify", "autodiscover"},
	} {
		g := findSub(mail, group)
		if g == nil {
			t.Errorf("hosting mail %s group not registered", group)
			continue
		}
		gsub := subNames(g)
		for _, n := range want {
			if !gsub[n] {
				t.Errorf("hosting mail %s is missing subcommand %q", group, n)
			}
		}
	}
}

func TestAntispamValue(t *testing.T) {
	for label, want := range map[string]int{"hard": 5, "medium": 8, "soft": 10, "off": 0} {
		got, err := antispamValue(label)
		if err != nil || got != want {
			t.Errorf("antispamValue(%q) = %d, %v; want %d, nil", label, got, err, want)
		}
		if antispamLabel(want) != label {
			t.Errorf("antispamLabel(%d) = %q, want %q", want, antispamLabel(want), label)
		}
	}
	if _, err := antispamValue("bogus"); err == nil {
		t.Error("antispamValue(bogus) should error")
	}
	// An unknown level renders as its raw number rather than losing it.
	if got := antispamLabel(3); got != "3" {
		t.Errorf("antispamLabel(3) = %q, want \"3\"", got)
	}
}

func TestParseOnOff(t *testing.T) {
	for arg, want := range map[string]bool{"on": true, "off": false} {
		got, err := parseOnOff(arg)
		if err != nil || got != want {
			t.Errorf("parseOnOff(%q) = %v, %v; want %v, nil", arg, got, err, want)
		}
	}
	if _, err := parseOnOff("maybe"); err == nil {
		t.Error("parseOnOff(maybe) should error")
	}
}

func TestEmptyDashAndOnOff(t *testing.T) {
	if emptyDash("") != "-" || emptyDash("x") != "x" {
		t.Error("emptyDash mismatch")
	}
	if onOff(true) != "on" || onOff(false) != "off" {
		t.Error("onOff mismatch")
	}
}
