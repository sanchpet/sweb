package cmd

import (
	"fmt"
	"io"
	"sort"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage credential profiles (multi-account)",
	Long: `Manage credential profiles for multiple SpaceWeb accounts — e.g. a cloud
panel (VPS) on one account and a hosting panel (mail, domains) on another.

Create a profile with 'sweb configure --profile <name>'. Select one globally with
'sweb profile use', or bind a command group to a profile with 'sweb profile bind'
(e.g. bind dns to the hosting account) so those commands use it automatically.
Precedence: --profile flag > $SWEB_PROFILE > group binding > current profile.`,
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured profiles",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		current := configGetString("current_profile")
		if current == "" {
			current = defaultProfile
		}
		names := profileNames()
		bindings := configGetMap("bindings")
		// group bindings by profile for display
		groupsOf := map[string][]string{}
		for g, p := range bindings {
			if s, ok := p.(string); ok {
				groupsOf[s] = append(groupsOf[s], g)
			}
		}
		type row struct {
			Name     string   `json:"name"`
			Current  bool     `json:"current"`
			Endpoint string   `json:"endpoint,omitempty"`
			Groups   []string `json:"groups,omitempty"`
		}
		var rows []row
		for _, n := range names {
			g := groupsOf[n]
			sort.Strings(g)
			rows = append(rows, row{n, n == current, configGetString("profiles", n, "endpoint"), g})
		}
		return render(cmd, rows, func(w io.Writer) {
			fmt.Fprintln(w, "PROFILE\tCURRENT\tENDPOINT\tBOUND GROUPS")
			for _, r := range rows {
				star := ""
				if r.Current {
					star = "*"
				}
				ep := r.Endpoint
				if ep == "" {
					ep = "(default)"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Name, star, ep, joinOrDash(r.Groups))
			}
		})
	},
}

var profileUseCmd = &cobra.Command{
	Use:               "use <name>",
	Short:             "Set the current (global) profile",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeProfiles,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := configSet(name, "current_profile"); err != nil {
			return err
		}
		if err := registerProfile(name); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Current profile is now %q\n", name)
		if !hasCredentials(name) {
			fmt.Fprintf(cmd.ErrOrStderr(), "note: profile %q has no stored credentials — run `sweb configure --profile %s`\n", name, name)
		}
		return nil
	},
}

var profileCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Print the current profile",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		current := configGetString("current_profile")
		if current == "" {
			current = defaultProfile
		}
		fmt.Fprintln(cmd.OutOrStdout(), current)
		return nil
	},
}

var profileBindCmd = &cobra.Command{
	Use:   "bind <group> <profile>",
	Short: "Bind a command group to a profile (e.g. bind dns hosting)",
	Long: fmt.Sprintf(`Bind a command group to a profile so its commands use that account
automatically, without --profile. Valid groups: %v.`, apiCommandGroups),
	Args: cobra.ExactArgs(2),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return apiCommandGroups, cobra.ShellCompDirectiveNoFileComp
		}
		return profileNames(), cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		group, profile := args[0], args[1]
		if !isAPIGroup(group) {
			return fmt.Errorf("unknown group %q; valid groups: %v", group, apiCommandGroups)
		}
		if err := configSet(profile, "bindings", group); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Bound group %q to profile %q\n", group, profile)
		return nil
	},
}

var profileUnbindCmd = &cobra.Command{
	Use:   "unbind <group>",
	Short: "Remove a command group's profile binding",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return apiCommandGroups, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		m := loadConfigMap()
		if bindings, ok := m["bindings"].(map[string]any); ok {
			delete(bindings, args[0])
		}
		if err := saveConfigMap(m); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Unbound group %q\n", args[0])
		return nil
	},
}

var profileRemoveCmd = &cobra.Command{
	Use:               "remove <name>",
	Short:             "Delete a profile and its stored credentials",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeProfiles,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if !confirmed(cmd, fmt.Sprintf("Delete profile %q and its stored credentials?", name), "Delete") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		removeCredentials(name)
		if err := deleteProfileConfig(name); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Removed profile %q\n", name)
		return nil
	},
}

// profileNames returns the configured profile names (always including default),
// sorted.
func profileNames() []string {
	set := map[string]bool{defaultProfile: true}
	for n := range configGetMap("profiles") {
		set[n] = true
	}
	names := make([]string, 0, len(set))
	for n := range set {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func isAPIGroup(g string) bool {
	for _, v := range apiCommandGroups {
		if v == g {
			return true
		}
	}
	return false
}

func hasCredentials(profile string) bool {
	login, _, token := loadCredentials(profile)
	return login != "" || token != ""
}

func joinOrDash(s []string) string {
	if len(s) == 0 {
		return "-"
	}
	out := s[0]
	for _, v := range s[1:] {
		out += "," + v
	}
	return out
}

// completeProfiles completes profile-name arguments and the --profile flag.
func completeProfiles(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return profileNames(), cobra.ShellCompDirectiveNoFileComp
}

func init() {
	profileCmd.AddCommand(profileListCmd, profileUseCmd, profileCurrentCmd, profileBindCmd, profileUnbindCmd, profileRemoveCmd)
	profileRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	rootCmd.AddCommand(profileCmd)
}
