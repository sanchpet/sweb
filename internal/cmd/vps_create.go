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
	Long: `Provision a new VPS. Use 'sweb vps config' to find the numeric IDs for
--plan, --distributive and --datacenter.

This call MUTATES and BILLS your SpaceWeb account. Use --dry-run to print the
request that would be sent without calling the API.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		f := cmd.Flags()
		plan, _ := f.GetInt("plan")
		distr, _ := f.GetInt("distributive")
		dc, _ := f.GetInt("datacenter")
		alias, _ := f.GetString("alias")
		sshKey, _ := f.GetString("ssh-key")
		ipCount, _ := f.GetInt("ip-count")
		dryRun, _ := f.GetBool("dry-run")

		if plan == 0 || distr == 0 || dc == 0 || alias == "" || sshKey == "" {
			return fmt.Errorf("--plan, --distributive, --datacenter, --alias and --ssh-key are required")
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
			fmt.Fprintln(cmd.OutOrStdout(), "[dry-run] would POST the following JSON-RPC request to /vps (no API call made):")
			fmt.Fprintln(cmd.OutOrStdout(), string(body))
			return nil
		}

		c, err := client()
		if err != nil {
			return err
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
	f.Int("plan", 0, "VPS plan ID (vpsPlanId) — see `sweb vps config`")
	f.Int("distributive", 0, "OS distributive ID (distributiveId) — see `sweb vps config`")
	f.Int("datacenter", 0, "datacenter ID — see `sweb vps config`")
	f.String("alias", "", "human-readable name for the VPS")
	f.String("ssh-key", "", "SSH public key to install")
	f.Int("ip-count", 1, "number of IP addresses")
	f.Bool("dry-run", false, "print the request that would be sent, without calling the API")

	// Dynamic value completion from the live catalog (getAvailableConfig).
	_ = vpsCreateCmd.RegisterFlagCompletionFunc("plan", completePlans)
	_ = vpsCreateCmd.RegisterFlagCompletionFunc("datacenter", completeDatacenters)
	_ = vpsCreateCmd.RegisterFlagCompletionFunc("distributive", completeDistributives)

	vpsCmd.AddCommand(vpsCreateCmd)
}
