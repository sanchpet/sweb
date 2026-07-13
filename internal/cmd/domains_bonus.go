package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/domains/bonus"
	"github.com/spf13/cobra"
)

var bonusCmd = &cobra.Command{
	Use:   "bonus",
	Short: "Domain bonuses",
}

var bonusListCmd = &cobra.Command{
	Use:   "list",
	Short: "List purchasable domain-bonus packages",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		packages, err := c.Bonus.GetList(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, packages, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tTITLE\tPRICE\tPRICE OLD\tDOMAINS\tPER DOMAIN")
			for _, p := range packages {
				fmt.Fprintf(w, "%d\t%s\t%d\t%d\t%d\t%s\n",
					p.ID, p.Title, p.Price, p.PriceOld, p.Domains, p.PriceForDomain)
			}
		})
	},
}

var bonusStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the account's domain bonuses and counts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		page, _ := cmd.Flags().GetInt("page")
		res, err := c.Bonus.Index(cmd.Context(), bonus.IndexOptions{Page: page})
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintf(w, "TOTAL\t%d\n", res.Count)
			fmt.Fprintf(w, "UNUSED\t%d\n", res.UnusedCount)
			fmt.Fprintln(w, "ID\tTLD\tDOMAIN\tTYPE\tUSED\tVALID TILL")
			for _, b := range res.Bonuses {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
					b.ID, b.TLD, b.Domain, b.TypeTitle, b.Used, b.ValidTill)
			}
		})
	},
}

var bonusBuyCmd = &cobra.Command{
	Use:   "buy <id>",
	Short: "Buy a domain-bonus package — mutating, bills the account",
	Long: `Buy a domain-bonus package by its id via the "buy" method.

This BILLS the account. You are asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		var id int
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
			return fmt.Errorf("invalid bonus id %q: %w", args[0], err)
		}
		if !confirmed(cmd, fmt.Sprintf("Buy bonus package %d? This bills your account.", id), "Buy") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Bonus.Buy(cmd.Context(), id); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Bought bonus package", id)
		return nil
	},
}

func init() {
	bonusStatusCmd.Flags().Int("page", 0, "page number (0-based)")
	bonusBuyCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	bonusCmd.AddCommand(bonusListCmd)
	bonusCmd.AddCommand(bonusStatusCmd)
	bonusCmd.AddCommand(bonusBuyCmd)
	domainsCmd.AddCommand(bonusCmd)
}
