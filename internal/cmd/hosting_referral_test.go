package cmd

import "testing"

// TestHostingReferralTree asserts the referral subgroup hangs off hosting with
// the list/add/confirm/remove lifecycle. It reuses findSub (hosting_db_test.go).
func TestHostingReferralTree(t *testing.T) {
	hosting := findSub(rootCmd, "hosting")
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}

	referral := findSub(hosting, "referral")
	if referral == nil {
		t.Fatal("hosting referral command not registered")
	}

	sub := subNames(referral)
	for _, n := range []string{"list", "add", "confirm", "remove"} {
		if !sub[n] {
			t.Errorf("hosting referral is missing subcommand %q", n)
		}
	}
}
