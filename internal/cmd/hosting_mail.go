package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/vh/mail"
	"github.com/spf13/cobra"
)

// mailCmd groups the shared-hosting email service (endpoint /vh/mail): the
// account's mail domains and mailboxes, mailbox lifecycle, autoreply, forwarding
// and delivery lists, the mail collector, per-mailbox white/black lists, and the
// domain-level DKIM/SPF/sender-verify toggles. It hangs off `hosting`, so it
// inherits that group's profile binding.
var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Manage shared-hosting email (domains, mailboxes, forwarding, DKIM, …)",
}

// antispamLevels maps the API's filter-level ints to human labels (and back),
// so the mailbox table and the `mailbox antispam` command speak in words.
var antispamLevels = map[int]string{5: "hard", 8: "medium", 10: "soft", 0: "off"}

// antispamValue resolves a level label (hard|medium|soft|off) to its API int.
func antispamValue(label string) (int, error) {
	for v, l := range antispamLevels {
		if l == label {
			return v, nil
		}
	}
	return 0, fmt.Errorf("--level must be one of hard, medium, soft, off")
}

// antispamLabel renders a filter-level int as its label, falling back to the
// raw number for a value the API adds later.
func antispamLabel(v int) string {
	if l, ok := antispamLevels[v]; ok {
		return l
	}
	return fmt.Sprintf("%d", v)
}

// onOff renders a boolean as the "on"/"off" the mail panel uses.
func onOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

// parseOnOff reads the shared <on|off> positional argument.
func parseOnOff(arg string) (bool, error) {
	switch arg {
	case "on":
		return true, nil
	case "off":
		return false, nil
	default:
		return false, fmt.Errorf("expected on or off, got %q", arg)
	}
}

// emptyDash renders "" as "-" so an unset field reads cleanly in a table.
func emptyDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// kv writes a key/value row to a table writer.
func kv(w io.Writer, k, v string) { fmt.Fprintf(w, "%s\t%s\n", k, v) }

// mailClient aliases the concrete SDK mail service; the toggle helpers take a
// selector over it so the on/off commands stay one-liners.
type mailClient = *mail.Service

// toggleFn is the signature of the domain-level on/off SDK calls
// (ChangeDomainSpf/ChangeSenderVerify/ChangeAutoDiscover).
type toggleFn func(ctx context.Context, domain string, on bool) error

// domainToggle builds the RunE for a `mail domain <toggle> <domain> <on|off>`
// command: it resolves the client, parses on/off, calls the selected SDK method,
// and reports the new state. what is the human label for the result line.
func domainToggle(sel func(mailClient) toggleFn, what string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		on, err := parseOnOff(args[1])
		if err != nil {
			return err
		}
		if err := sel(c.Mail)(cmd.Context(), args[0], on); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s %s on %s\n", what, onOff(on), args[0])
		return nil
	}
}

func init() {
	hostingCmd.AddCommand(mailCmd)
}
