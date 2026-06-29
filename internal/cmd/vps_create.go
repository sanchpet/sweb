package cmd

import (
	"encoding/json"
	"fmt"

	sweb "github.com/sanchpet/sweb-go-sdk"
	"github.com/spf13/cobra"
)

var vpsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Provision a new VPS (mutates — bills your account)",
	Long: `Provision a new VPS. Pick the plan one of two ways:

  • a stock plan:        --plan <id>            (see 'sweb vps config')
  • the configurator:    --cpu N --ram N --disk N [--category id]
                         (builds a custom plan, like the panel's "Конфигуратор";
                          ram and disk are in GB; default category is NVMe)

--distributive, --datacenter, --alias and --ssh-key are always required.
This call MUTATES and BILLS your account. --dry-run prints the request without
creating (in configurator mode it still resolves the plan id — a read-only call).`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		f := cmd.Flags()
		plan, _ := f.GetInt("plan")
		cpu, _ := f.GetInt("cpu")
		ram, _ := f.GetInt("ram")
		disk, _ := f.GetInt("disk")
		category, _ := f.GetInt("category")
		distr, _ := f.GetInt("distributive")
		dc, _ := f.GetInt("datacenter")
		alias, _ := f.GetString("alias")
		sshKey, _ := f.GetString("ssh-key")
		ipCount, _ := f.GetInt("ip-count")
		dryRun, _ := f.GetBool("dry-run")

		if distr == 0 || dc == 0 || alias == "" || sshKey == "" {
			return fmt.Errorf("--distributive, --datacenter, --alias and --ssh-key are required")
		}

		var c *sweb.Client
		// Configurator mode: resolve a custom plan id from --cpu/--ram/--disk.
		if plan == 0 {
			if cpu == 0 || ram == 0 || disk == 0 {
				return fmt.Errorf("provide --plan, or --cpu/--ram/--disk to build a configurator plan")
			}
			if category == 0 {
				category = 1 // NVMe ("Быстрые") — see `sweb vps config`
			}
			var err error
			if c, err = client(); err != nil {
				return err
			}
			plan, err = c.VPS.GetConstructorPlanID(cmd.Context(), cpu, ram, disk, category)
			if err != nil {
				return fmt.Errorf("resolve configurator plan: %w", err)
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "configurator %dcpu/%dGB/%dGB (category %d) → plan %d\n",
				cpu, ram, disk, category, plan)
		}

		req := sweb.CreateVPSRequest{
			VPSPlanID:      plan,
			DistributiveID: distr,
			Datacenter:     dc,
			Alias:          alias,
			SSHKey:         sshKey,
			IPCount:        ipCount,
		}

		if dryRun {
			body, err := json.MarshalIndent(map[string]any{
				"jsonrpc": "2.0",
				"method":  "create",
				"params":  req,
			}, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "[dry-run] would POST the following JSON-RPC request to /vps (no VPS created):")
			fmt.Fprintln(cmd.OutOrStdout(), string(body))
			return nil
		}

		if c == nil {
			var err error
			if c, err = client(); err != nil {
				return err
			}
		}
		res, err := c.VPS.Create(cmd.Context(), req)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(res))
		return nil
	},
}

func init() {
	f := vpsCreateCmd.Flags()
	f.Int("plan", 0, "stock VPS plan ID (vpsPlanId) — see `sweb vps config`")
	f.Int("cpu", 0, "configurator: CPU cores")
	f.Int("ram", 0, "configurator: RAM in GB")
	f.Int("disk", 0, "configurator: disk in GB")
	f.Int("category", 0, "configurator: category id (default 1 = NVMe) — see `sweb vps config`")
	f.Int("distributive", 0, "OS distributive ID (distributiveId) — see `sweb vps config`")
	f.Int("datacenter", 0, "datacenter ID — see `sweb vps config`")
	f.String("alias", "", "human-readable name for the VPS")
	f.String("ssh-key", "", "SSH public key to install")
	f.Int("ip-count", 1, "number of IP addresses")
	f.Bool("dry-run", false, "print the request that would be sent, without creating")

	// Dynamic value completion from the live catalog (getAvailableConfig).
	_ = vpsCreateCmd.RegisterFlagCompletionFunc("plan", completePlans)
	_ = vpsCreateCmd.RegisterFlagCompletionFunc("datacenter", completeDatacenters)
	_ = vpsCreateCmd.RegisterFlagCompletionFunc("distributive", completeDistributives)
	_ = vpsCreateCmd.RegisterFlagCompletionFunc("category", completeCategories)

	vpsCmd.AddCommand(vpsCreateCmd)
}
