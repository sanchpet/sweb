package cmd

import (
	"strings"
	"testing"

	"github.com/sanchpet/sweb-go-sdk/dns"
)

func TestRecordName(t *testing.T) {
	for _, tc := range []struct {
		name string
		rec  dns.Record
		want string
	}{
		{"named A", dns.Record{Name: "www"}, "www"},
		{"apex A", dns.Record{Name: ""}, "@"},
		{"apex TXT", dns.Record{Name: "", Domain: "@"}, "@"},
		{"subdomain TXT", dns.Record{Name: "", Domain: "_sweb-probe"}, "_sweb-probe"},
	} {
		if got := recordName(tc.rec); got != tc.want {
			t.Errorf("%s: recordName = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestTruncateCell(t *testing.T) {
	for _, tc := range []struct {
		in   string
		max  int
		want string
	}{
		{"short", 60, "short"},
		{"exactly-ten", 11, "exactly-ten"}, // len==max, unchanged
		{"abcdefghij", 5, "abcd…"},         // truncated to max-1 runes + ellipsis
	} {
		if got := truncateCell(tc.in, tc.max); got != tc.want {
			t.Errorf("truncateCell(%q, %d) = %q, want %q", tc.in, tc.max, got, tc.want)
		}
	}
	// A DKIM-length value is bounded so it cannot blow out the column.
	long := strings.Repeat("A", 300)
	if got := truncateCell(long, 60); len([]rune(got)) != 60 {
		t.Errorf("truncated length = %d runes, want 60", len([]rune(got)))
	}
}
