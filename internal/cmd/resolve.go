package cmd

import (
	"context"
	"fmt"
	"strings"

	sweb "github.com/sanchpet/sweb-go-sdk"
)

// resolveVPS maps a CLI argument — a VPS name (alias) OR a billing id — to its
// billing id, which is what the SDK/API operate on. An exact billing-id match
// wins; otherwise it matches by name. Names are not guaranteed unique, so an
// ambiguous name is an error (fall back to the billing id). This is a CLI-side
// convenience: the billing id (login_vps_N) is not easily visible in the panel,
// whereas the name is what the user knows.
func resolveVPS(ctx context.Context, c *sweb.Client, arg string) (string, error) {
	list, err := c.VPS.List(ctx)
	if err != nil {
		return "", fmt.Errorf("resolve %q: %w", arg, err)
	}
	var byName []string
	for _, v := range list {
		if v.BillingID == arg {
			return arg, nil // exact billing-id match — unambiguous
		}
		if v.Name == arg {
			byName = append(byName, v.BillingID)
		}
	}
	switch len(byName) {
	case 1:
		return byName[0], nil
	case 0:
		return "", fmt.Errorf("no VPS named or billed %q — see 'sweb vps list'", arg)
	default:
		return "", fmt.Errorf("name %q is ambiguous (%d VPS: %s) — use the billing id",
			arg, len(byName), strings.Join(byName, ", "))
	}
}
