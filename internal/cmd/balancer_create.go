package cmd

import (
	"fmt"

	"github.com/sanchpet/sweb-go-sdk/balancer"
	"github.com/spf13/cobra"
)

var balancerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Order a new load balancer (mutates — bills your account)",
	Long: `Order a new load balancer. --datacenter, --type, --plan and at least one
--server and one --rule are required.

  --server ip[,weight[,vpsName]]     back-end server (repeatable, max 20)
                                     weight (1..5) applies to type roundrobin
  --rule   protoBal:portBal:protoSrv:portSrv   forwarding rule (repeatable)

See 'sweb balancer config' for the plan ids and protocols.
This call MUTATES and BILLS your account. You are asked to confirm unless --yes.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		f := cmd.Flags()
		dc, _ := f.GetInt("datacenter")
		typ, _ := f.GetString("type")
		plan, _ := f.GetInt("plan")
		serverVals, _ := f.GetStringArray("server")
		ruleVals, _ := f.GetStringArray("rule")
		alias, _ := f.GetString("alias")
		healthCheck, _ := f.GetBool("health-check")
		proxyProto, _ := f.GetBool("proxy-proto")
		keepalive, _ := f.GetBool("keepalive")
		saveSession, _ := f.GetBool("save-session")
		firstOrder, _ := f.GetBool("first-order")

		if dc == 0 || typ == "" || plan == 0 || len(serverVals) == 0 || len(ruleVals) == 0 {
			return fmt.Errorf("--datacenter, --type, --plan and at least one --server and --rule are required")
		}
		servers, err := parseServers(serverVals)
		if err != nil {
			return err
		}
		rules, err := parseRules(ruleVals)
		if err != nil {
			return err
		}

		if !confirmed(cmd, "Order a new load balancer? This bills your account.", "Order") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}

		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Balancer.Create(cmd.Context(), balancer.CreateOptions{
			Datacenter:   dc,
			Type:         typ,
			Servers:      servers,
			Rules:        rules,
			PlanID:       plan,
			HealthCheck:  healthCheck,
			ProxyProto:   proxyProto,
			Keepalive:    keepalive,
			SaveSession:  saveSession,
			Alias:        alias,
			IsFirstOrder: firstOrder,
		}); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Balancer order started")
		return nil
	},
}

func init() {
	f := balancerCreateCmd.Flags()
	f.Int("datacenter", 0, "datacenter ID — see `sweb vps config`")
	f.String("type", "", "balancing algorithm: roundrobin|leastconn")
	f.Int("plan", 0, "tariff plan ID — see `sweb balancer config`")
	f.StringArray("server", nil, "back-end server ip[,weight[,vpsName]] (repeatable, max 20)")
	f.StringArray("rule", nil, "forwarding rule protoBal:portBal:protoSrv:portSrv (repeatable)")
	f.String("alias", "", "human-readable name for the balancer (optional)")
	f.Bool("health-check", false, "enable back-end health checks")
	f.Bool("proxy-proto", false, "enable the PROXY protocol")
	f.Bool("keepalive", false, "enable keepalive")
	f.Bool("save-session", false, "enable session persistence")
	f.Bool("first-order", false, "mark this as a first order (isFirstOrder)")
	f.Bool("yes", false, "skip the confirmation prompt")

	balancerCmd.AddCommand(balancerCreateCmd)
}
