package cmd

import "testing"

func TestHostingSSHCommandTree(t *testing.T) {
	hosting := findSub(rootCmd, "hosting")
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}
	ssh := findSub(hosting, "ssh")
	if ssh == nil {
		t.Fatal("hosting ssh command not registered")
	}

	for _, n := range []string{"on", "off"} {
		if findSub(ssh, n) == nil {
			t.Errorf("hosting ssh is missing subcommand %q", n)
		}
	}
}

// TestHostingSSHPeriodValidation guards the on-command's period whitelist: the
// documented 3/8/24 hours pass and anything else is rejected before the API call.
func TestHostingSSHPeriodValidation(t *testing.T) {
	for _, h := range []int{3, 8, 24} {
		if !sshValidPeriod(h) {
			t.Errorf("period %d should be valid", h)
		}
	}
	for _, h := range []int{0, 1, 12, 48} {
		if sshValidPeriod(h) {
			t.Errorf("period %d should be rejected", h)
		}
	}
}
