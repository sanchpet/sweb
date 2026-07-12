package cmd

import "testing"

func TestYesNo(t *testing.T) {
	if yesNo(true) != "yes" || yesNo(false) != "no" {
		t.Errorf("yesNo = %q/%q, want yes/no", yesNo(true), yesNo(false))
	}
}

func TestCompleteProlongModes(t *testing.T) {
	// Completes only the <mode> arg (after <domain> is already present).
	got, _ := completeProlongModes(nil, []string{"example.com"}, "")
	if len(got) != len(prolongModes) {
		t.Errorf("with one prior arg: got %v, want the prolong modes", got)
	}
	// No completion before the domain arg, or after both args are present.
	if got, _ := completeProlongModes(nil, nil, ""); got != nil {
		t.Errorf("with no prior arg: got %v, want nil", got)
	}
	if got, _ := completeProlongModes(nil, []string{"example.com", "manual"}, ""); got != nil {
		t.Errorf("with two prior args: got %v, want nil", got)
	}
}
