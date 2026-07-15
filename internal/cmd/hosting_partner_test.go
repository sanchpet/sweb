package cmd

import "testing"

// TestHostingPartnerCommandTree asserts the partner-program command tree hangs
// off `hosting partner` with its expected subgroups and leaves.
func TestHostingPartnerCommandTree(t *testing.T) {
	hosting := findSub(rootCmd, "hosting")
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}
	partner := findSub(hosting, "partner")
	if partner == nil {
		t.Fatal("hosting partner command not registered")
	}

	psub := subNames(partner)
	for _, n := range []string{
		"status", "join", "requisites", "os-config",
		"plans", "clients", "materials", "stats", "links",
		"withdraw", "order",
	} {
		if !psub[n] {
			t.Errorf("hosting partner is missing subcommand %q", n)
		}
	}

	for _, tc := range []struct {
		parent string
		subs   []string
	}{
		{"plans", []string{"standard", "vip"}},
		{"requisites", []string{"set"}},
		{"clients", []string{"list", "card", "comment", "log"}},
		{"materials", []string{"types", "list"}},
		{"withdraw", []string{"requisites", "send"}},
		{"order", []string{"vh", "vip", "vps"}},
	} {
		group := findSub(partner, tc.parent)
		if group == nil {
			t.Fatalf("hosting partner %s not registered", tc.parent)
		}
		gsub := subNames(group)
		for _, n := range tc.subs {
			if !gsub[n] {
				t.Errorf("hosting partner %s is missing subcommand %q", tc.parent, n)
			}
		}
	}

	clients := findSub(partner, "clients")
	log := findSub(clients, "log")
	if log == nil {
		t.Fatal("hosting partner clients log not registered")
	}
	lsub := subNames(log)
	for _, n := range []string{"events", "finance"} {
		if !lsub[n] {
			t.Errorf("hosting partner clients log is missing subcommand %q", n)
		}
	}
}

// TestHostingPartnerConfirmFlags checks that every mutating/billing partner
// command registers the --yes escape hatch that confirmed() reads.
func TestHostingPartnerConfirmFlags(t *testing.T) {
	for _, path := range [][]string{
		{"hosting", "partner", "join"},
		{"hosting", "partner", "withdraw", "send"},
		{"hosting", "partner", "order", "vh"},
		{"hosting", "partner", "order", "vip"},
		{"hosting", "partner", "order", "vps"},
	} {
		cmd, _, err := rootCmd.Find(path)
		if err != nil {
			t.Errorf("%v: not found: %v", path, err)
			continue
		}
		if cmd.Flags().Lookup("yes") == nil {
			t.Errorf("%v: missing --yes flag", path)
		}
	}
}
