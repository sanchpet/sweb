package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/monitoring/contacts"
	"github.com/spf13/cobra"
)

var contactListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's monitoring contacts",
	Long: `List the account's monitoring contacts (method "index"). Pass --all to
use "getAllContacts", which also returns admin contacts.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if all, _ := cmd.Flags().GetBool("all"); all {
			list, err := c.MonitoringContacts.GetAllContacts(cmd.Context())
			if err != nil {
				return err
			}
			return render(cmd, list, func(w io.Writer) {
				writeContactRows(w, list)
			})
		}
		res, err := c.MonitoringContacts.Index(cmd.Context(), &contacts.ListOptions{
			Page:    flagInt(cmd, "page"),
			PerPage: flagInt(cmd, "per-page"),
		})
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			writeContactRows(w, res.List)
		})
	},
}

// writeContactRows prints the shared contact table for both list paths.
func writeContactRows(w io.Writer, list []contacts.Contact) {
	fmt.Fprintln(w, "ID\tTYPE\tNAME\tVALUE\tVERIFIED\tADMIN")
	for _, ct := range list {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%t\t%t\n",
			int64(ct.ID), ct.Type, ct.Name, ct.Value, ct.Verified, ct.Admin)
	}
}

func init() {
	contactListCmd.Flags().Bool("all", false, "list every contact including admin contacts (getAllContacts)")
	contactListCmd.Flags().Int("page", 0, "page of results (1-based; 0 lets the API default)")
	contactListCmd.Flags().Int("per-page", 0, "rows per page (0 lets the API default)")
	contactCmd.AddCommand(contactListCmd)
}
