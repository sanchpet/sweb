package cmd

import (
	"fmt"

	"github.com/sanchpet/sweb-go-sdk/balancer"
	"github.com/spf13/cobra"
)

var balancerEditCmd = &cobra.Command{
	Use:   "edit <billing-id>",
	Short: "Reconfigure a load balancer in place (mutates)",
	Long: `Reconfigure a load balancer. <billing-id> is a Balancer.BillingID from
'sweb balancer list'. --type and at least one --server and one --rule are
required (edit replaces the full server/rule set).

  --server ip[,weight[,vpsName]]     back-end server (repeatable)
  --rule   protoBal:portBal:protoSrv:portSrv   forwarding rule (repeatable)

You are asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		typ, _ := f.GetString("type")
		serverVals, _ := f.GetStringArray("server")
		ruleVals, _ := f.GetStringArray("rule")
		alias, _ := f.GetString("alias")
		healthCheck, _ := f.GetBool("health-check")
		proxyProto, _ := f.GetBool("proxy-proto")
		keepalive, _ := f.GetBool("keepalive")
		saveSession, _ := f.GetBool("save-session")

		if typ == "" || len(serverVals) == 0 || len(ruleVals) == 0 {
			return fmt.Errorf("--type and at least one --server and --rule are required")
		}
		servers, err := parseServers(serverVals)
		if err != nil {
			return err
		}
		rules, err := parseRules(ruleVals)
		if err != nil {
			return err
		}

		if !confirmed(cmd, fmt.Sprintf("Reconfigure balancer %q?", args[0]), "Edit") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}

		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Balancer.Edit(cmd.Context(), balancer.EditOptions{
			BillingID:   args[0],
			Type:        typ,
			Servers:     servers,
			Rules:       rules,
			HealthCheck: healthCheck,
			ProxyProto:  proxyProto,
			Keepalive:   keepalive,
			SaveSession: saveSession,
			Alias:       alias,
		}); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Edited", args[0])
		return nil
	},
}

func init() {
	f := balancerEditCmd.Flags()
	f.String("type", "", "balancing algorithm: roundrobin|leastconn")
	f.StringArray("server", nil, "back-end server ip[,weight[,vpsName]] (repeatable)")
	f.StringArray("rule", nil, "forwarding rule protoBal:portBal:protoSrv:portSrv (repeatable)")
	f.String("alias", "", "human-readable name for the balancer (optional)")
	f.Bool("health-check", false, "enable back-end health checks")
	f.Bool("proxy-proto", false, "enable the PROXY protocol")
	f.Bool("keepalive", false, "enable keepalive")
	f.Bool("save-session", false, "enable session persistence")
	f.Bool("yes", false, "skip the confirmation prompt")

	balancerCmd.AddCommand(balancerEditCmd)
}
