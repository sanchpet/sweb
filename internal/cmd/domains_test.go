package cmd

import (
	"errors"
	"testing"

	"github.com/sanchpet/sweb-go-sdk/apierr"
)

func TestPriceCells(t *testing.T) {
	if floatCell(nil) != "—" {
		t.Errorf("floatCell(nil) = %q, want —", floatCell(nil))
	}
	v := 189.0
	if floatCell(&v) != "189" {
		t.Errorf("floatCell(189) = %q, want 189", floatCell(&v))
	}
	if boolCell(nil) != "—" {
		t.Errorf("boolCell(nil) = %q, want —", boolCell(nil))
	}
	yes := true
	if boolCell(&yes) != "yes" {
		t.Errorf("boolCell(true) = %q, want yes", boolCell(&yes))
	}
}

func TestApiReason(t *testing.T) {
	// A SpaceWeb API error surfaces just its human message.
	if got := apiReason(&apierr.Error{Code: -32500, Message: "Домен занят"}); got != "Домен занят" {
		t.Errorf("apiReason(apierr.Error) = %q, want the message", got)
	}
	// A plain error falls back to its full string.
	if got := apiReason(errors.New("boom")); got != "boom" {
		t.Errorf("apiReason(plain) = %q, want boom", got)
	}
}

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
