package cmd

import "testing"

func TestPayCommandTree(t *testing.T) {
	pay := findSub(rootCmd, "pay")
	if pay == nil {
		t.Fatal("pay command not registered")
	}

	psub := subNames(pay)
	for _, n := range []string{
		"balance", "index", "autopay", "recommendations",
		"recommendation-cost", "upcoming", "remains", "reserves", "deferment",
	} {
		if !psub[n] {
			t.Errorf("pay is missing subcommand %q", n)
		}
	}
}

// TestPayDefermentConfirmFlag guards that the one mutating pay command carries
// the --yes escape hatch its confirmed() call needs.
func TestPayDefermentConfirmFlag(t *testing.T) {
	pay := findSub(rootCmd, "pay")
	if pay == nil {
		t.Fatal("pay command not registered")
	}
	deferment := findSub(pay, "deferment")
	if deferment == nil {
		t.Fatal("pay deferment command not registered")
	}
	if deferment.Flags().Lookup("yes") == nil {
		t.Error("pay deferment is missing the --yes flag")
	}
}
