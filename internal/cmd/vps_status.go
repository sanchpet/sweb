package cmd

import (
	"fmt"
	"io"

	sweb "github.com/sanchpet/sweb-go-sdk"
	"github.com/spf13/cobra"
)

var vpsStatusCmd = &cobra.Command{
	Use:   "status <vps>",
	Short: "Show a single VPS's state (power, current action, plan)",
	Long: `Show one VPS's current state via the "index" listing. <vps> is the VPS name
(alias) or its billing ID (login_vps_N), from 'sweb vps list'.

Reports the power state and the in-flight current_action (empty when idle),
plus plan, resources, addresses, OS and datacenter. Use -o json for the raw node.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeBillingIDs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}
		list, err := c.VPS.List(cmd.Context())
		if err != nil {
			return err
		}
		var node *sweb.VPS
		for i := range list {
			if list[i].BillingID == billingID {
				node = &list[i]
				break
			}
		}
		if node == nil {
			return fmt.Errorf("VPS %q not found — see 'sweb vps list'", billingID)
		}

		return render(cmd, node, func(w io.Writer) {
			action := node.CurrentAction
			if action == "" {
				action = "-"
			}
			localIP := node.LocalIP
			if localIP == "" {
				localIP = "-"
			}
			fmt.Fprintf(w, "NAME\t%s\n", node.Name)
			fmt.Fprintf(w, "BILLING_ID\t%s\n", node.BillingID)
			fmt.Fprintf(w, "RUNNING\t%s\n", runState(node))
			fmt.Fprintf(w, "CURRENT_ACTION\t%s\n", action)
			fmt.Fprintf(w, "PLAN\t%s\n", node.PlanName)
			fmt.Fprintf(w, "CPU\t%d\n", int64(node.CPU))
			fmt.Fprintf(w, "RAM(MB)\t%d\n", int64(node.RAM))
			fmt.Fprintf(w, "DISK\t%s\n", node.Disk)
			fmt.Fprintf(w, "IP\t%s\n", node.IP)
			fmt.Fprintf(w, "LOCAL_IP\t%s\n", localIP)
			fmt.Fprintf(w, "OS\t%s\n", node.OSDistribution)
			fmt.Fprintf(w, "DATACENTER\t%s\n", node.Datacenter)
		})
	},
}

func init() {
	vpsCmd.AddCommand(vpsStatusCmd)
}
