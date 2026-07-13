package cmd

import (
	"fmt"

	"github.com/sanchpet/sweb-go-sdk/dns"
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage a domain's DNS zone",
}

// resolveDNSAction reads the shared --action flag (add|edit|remove, default
// edit) and validates it. For a remove it prompts for confirmation unless --yes,
// returning ok=false when the user cancels. Edit commands must register both the
// --action and --yes flags.
func resolveDNSAction(cmd *cobra.Command, what string) (dns.Action, bool, error) {
	raw, _ := cmd.Flags().GetString("action")
	action := dns.Action(raw)
	switch action {
	case dns.ActionAdd, dns.ActionEdit:
		return action, true, nil
	case dns.ActionRemove:
		if !confirmed(cmd, fmt.Sprintf("Remove %s?", what), "Remove") {
			return action, false, nil
		}
		return action, true, nil
	default:
		return "", false, fmt.Errorf("--action must be one of add, edit, remove")
	}
}

// flagInt reads an int flag, ignoring the lookup error (unregistered → 0).
func flagInt(cmd *cobra.Command, name string) int {
	v, _ := cmd.Flags().GetInt(name)
	return v
}

// pastTense renders a DNS action for a result message.
func pastTense(a dns.Action) string {
	switch a {
	case dns.ActionAdd:
		return "Added"
	case dns.ActionRemove:
		return "Removed"
	default:
		return "Edited"
	}
}

// addDNSEditFlags registers the flags every edit command shares.
func addDNSEditFlags(cmd *cobra.Command) {
	cmd.Flags().String("action", "edit", "operation on the record: add|edit|remove")
	cmd.Flags().Int("index", 0, "record index (from 'sweb dns records'); identifies the record to edit/remove")
	cmd.Flags().Bool("yes", false, "skip the confirmation prompt on remove")
	_ = cmd.RegisterFlagCompletionFunc("action", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"add", "edit", "remove"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}
