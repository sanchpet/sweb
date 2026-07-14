package cmd

import (
	"fmt"
	"strings"

	"github.com/sanchpet/sweb-go-sdk/dbaas"
	"github.com/spf13/cobra"
)

var dbaasCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Provision a new managed-database cluster (mutates — bills your account)",
	Long: `Provision a new managed-database cluster. Pick the plan one of two ways:

  • a stock plan:      --plan <id>                       (see 'sweb dbaas config')
  • the configurator:  --cpu N --memory N --storage N [--replicas N]
                       (builds a custom plan; memory and storage are in GB;
                        replicas default to 0 = master only)

--engine, --version and at least one --user are always required. Each --user is
"name:password" and may be repeated.

This call MUTATES and BILLS your account. You are asked to confirm unless --yes.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		f := cmd.Flags()
		engine, _ := f.GetString("engine")
		version, _ := f.GetString("version")
		plan, _ := f.GetInt("plan")
		cpu, _ := f.GetInt("cpu")
		memory, _ := f.GetInt("memory")
		storage, _ := f.GetInt("storage")
		replicas, _ := f.GetInt("replicas")
		displayName, _ := f.GetString("name")
		userSpecs, _ := f.GetStringArray("user")

		if engine == "" || version == "" {
			return fmt.Errorf("--engine and --version are required")
		}
		if len(userSpecs) == 0 {
			return fmt.Errorf("at least one --user name:password is required")
		}
		users, err := parseDBaaSUsers(userSpecs)
		if err != nil {
			return err
		}

		c, err := client()
		if err != nil {
			return err
		}

		// Configurator mode: resolve a custom plan id from --cpu/--memory/--storage.
		if plan == 0 {
			if cpu == 0 || memory == 0 || storage == 0 {
				return fmt.Errorf("provide --plan, or --cpu/--memory/--storage to build a configurator plan")
			}
			id, err := c.DBaaS.ConstructorPlanID(cmd.Context(), cpu, memory, storage, replicas)
			if err != nil {
				return fmt.Errorf("resolve configurator plan: %w", err)
			}
			plan = int(id)
			fmt.Fprintf(cmd.ErrOrStderr(), "configurator %dcpu/%dGB/%dGB (%d replicas) → plan %d\n",
				cpu, memory, storage, replicas, plan)
		}

		if !confirmed(cmd, "Create this managed-database cluster? This mutates and bills your account.", "Create") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}

		req := dbaas.CreateInstanceRequest{
			EngineType:    engine,
			EngineVersion: version,
			Users:         users,
			PlanID:        plan,
			DisplayName:   displayName,
		}
		res, err := c.DBaaS.CreateInstance(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(res))
		return nil
	},
}

// parseDBaaSUsers parses "name:password" specs into UserCredentials. The
// password may contain ':' (split on the first colon only).
func parseDBaaSUsers(specs []string) ([]dbaas.UserCredentials, error) {
	users := make([]dbaas.UserCredentials, 0, len(specs))
	for _, s := range specs {
		name, password, ok := strings.Cut(s, ":")
		if !ok || name == "" || password == "" {
			return nil, fmt.Errorf("invalid --user %q: want name:password", s)
		}
		users = append(users, dbaas.UserCredentials{Name: name, Password: password})
	}
	return users, nil
}

func init() {
	f := dbaasCreateCmd.Flags()
	f.String("engine", "", "engine type, e.g. PostgreSQL or MySQL (see `sweb dbaas config`)")
	f.String("version", "", "engine version (see `sweb dbaas config`)")
	f.Int("plan", 0, "stock plan ID (planId) — see `sweb dbaas config`")
	f.Int("cpu", 0, "configurator: CPU cores")
	f.Int("memory", 0, "configurator: memory in GB")
	f.Int("storage", 0, "configurator: storage in GB")
	f.Int("replicas", 0, "configurator: replica count (0 = master only)")
	f.String("name", "", "display name for the cluster")
	f.StringArray("user", nil, "cluster user as name:password (repeatable)")
	f.Bool("yes", false, "skip the confirmation prompt")

	dbaasCmd.AddCommand(dbaasCreateCmd)
}
