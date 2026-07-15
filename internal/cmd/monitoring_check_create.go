package cmd

import (
	"fmt"

	"github.com/sanchpet/sweb-go-sdk/monitoring/checks"
	"github.com/spf13/cobra"
)

// addCheckSpecFlags registers the flags shared by `check create` and
// `check edit`, one per Spec field that makes sense on the command line. --type
// is create-only (edit is keyed by id), so it is registered by the callers.
func addCheckSpecFlags(cmd *cobra.Command) {
	cmd.Flags().String("target", "", "URL or IP to check")
	cmd.Flags().String("name", "", "display name")
	cmd.Flags().Int("interval", 0, "interval id (see 'monitoring check' intervals via getInfo)")
	cmd.Flags().IntSlice("contact", nil, "contact id to notify (repeatable)")
	cmd.Flags().Int("port", 0, "port to check (Port checks only)")
	cmd.Flags().Bool("ssl", false, "use SSL (HTTP checks)")
	cmd.Flags().StringSlice("keyword", nil, "keyword to match (HTTP checks; repeatable)")
	cmd.Flags().Int("keyword-mode", 0, "keyword mode id (HTTP checks)")
}

// checkSpecFromFlags builds a checks.Spec from the shared flags. --target and
// --name are required; their absence is reported.
func checkSpecFromFlags(cmd *cobra.Command) (checks.Spec, error) {
	target, _ := cmd.Flags().GetString("target")
	name, _ := cmd.Flags().GetString("name")
	if target == "" || name == "" {
		return checks.Spec{}, fmt.Errorf("--target and --name are required")
	}
	contacts, _ := cmd.Flags().GetIntSlice("contact")
	keywords, _ := cmd.Flags().GetStringSlice("keyword")
	ssl, _ := cmd.Flags().GetBool("ssl")
	return checks.Spec{
		Type:        flagInt(cmd, "type"),
		Target:      target,
		Name:        name,
		Interval:    flagInt(cmd, "interval"),
		ContactIDs:  contacts,
		Port:        flagInt(cmd, "port"),
		SSL:         ssl,
		Keywords:    keywords,
		KeywordMode: flagInt(cmd, "keyword-mode"),
	}, nil
}

var checkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a monitoring check — this BILLS against the check quota",
	Long: `Create a monitoring check (method "create"). This BILLS against the
check quota; you are asked to confirm unless --yes.

--type is the check type id (1 Ping, 2 HTTP, 3 Port); --port/--ssl/--keyword/
--keyword-mode are only meaningful for the matching type.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		spec, err := checkSpecFromFlags(cmd)
		if err != nil {
			return err
		}
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Create check %q? This bills against your check quota.", spec.Name), "Create") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.MonitoringChecks.Create(cmd.Context(), spec); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Created check:", spec.Name)
		return nil
	},
}

func init() {
	addCheckSpecFlags(checkCreateCmd)
	checkCreateCmd.Flags().Int("type", 0, "check type id (1 Ping, 2 HTTP, 3 Port)")
	checkCreateCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	checkCmd.AddCommand(checkCreateCmd)
}
