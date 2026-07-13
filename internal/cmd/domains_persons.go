package cmd

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sanchpet/sweb-go-sdk/domains/persons"
	"github.com/spf13/cobra"
)

var personsCmd = &cobra.Command{
	Use:   "persons",
	Short: "Manage domain registrant persons",
}

// personTypeLabel maps the API's "type" code to a readable label.
func personTypeLabel(t string) string {
	switch t {
	case persons.TypeIndividual:
		return "individual"
	case persons.TypeEntrepreneur:
		return "entrepreneur"
	case persons.TypeLegal:
		return "legal"
	default:
		return t
	}
}

var personsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's registrant contacts",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		people, propsFilled, err := c.Persons.List(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, people, func(w io.Writer) {
			fmt.Fprintf(w, "# props filled: %s\n", yesNo(propsFilled))
			fmt.Fprintln(w, "ID\tNAME\tHANDLE\tTYPE\tRESIDENT\tUSED\tVALID")
			for _, p := range people {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
					int64(p.ID), p.Name, p.SwebHandle, personTypeLabel(p.Type),
					yesNo(p.Resident == 1), yesNo(p.Used == 1), yesNo(p.Valid == 1))
			}
		})
	},
}

var personsInfoCmd = &cobra.Command{
	Use:   "info <id>",
	Short: "Show full details for a registrant contact",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid person id %q: %w", args[0], err)
		}
		c, err := client()
		if err != nil {
			return err
		}
		info, err := c.Persons.Info(cmd.Context(), id)
		if err != nil {
			return err
		}
		if info == nil {
			return fmt.Errorf("no registrant contact with id %d", id)
		}
		return render(cmd, info, func(w io.Writer) {
			row := func(k, v string) {
				if v != "" {
					fmt.Fprintf(w, "%s\t%s\n", k, v)
				}
			}
			row("NAME", info.Name)
			row("NAME (LATIN)", info.NameTrans)
			row("TYPE", personTypeLabel(info.Type))
			row("RESIDENT", yesNo(info.Resident))
			row("USED", yesNo(info.Used == 1))
			row("PHONES", strings.Join(info.Phones, ", "))
			row("EMAILS", strings.Join(info.Emails, ", "))
			row("INN", info.INN)
			row("POST INDEX", info.PostIndex)
			row("POST CITY", info.PostCity)
			row("POST ADDRESS", info.PostAddress)
			// Individual / sole proprietor.
			row("BIRTHDATE", info.Birthdate)
			row("PASSPORT SERIES", info.PassSeries)
			row("PASSPORT NUMBER", info.PassNum)
			row("PASSPORT DATE", info.PassDate)
			row("PASSPORT ORG", info.PassOrg)
			// Legal entity.
			row("FAXES", strings.Join(info.Faxes, ", "))
			row("JUR INDEX", info.JurIndex)
			row("JUR CITY", info.JurCity)
			row("JUR ADDRESS", info.JurAddress)
			row("KPP", info.KPP)
			row("REPRESENTATIVE", info.PersName)
			row("REPRESENTATIVE (LATIN)", info.PersNameTrans)
		})
	},
}

var personsCreateIndividualCmd = &cobra.Command{
	Use:   "create-individual",
	Short: "Create an individual / sole-proprietor registrant contact",
	Long: `Create (or, with --id, edit) an individual or sole-proprietor registrant
contact ("createFizIp"). Name, phones, emails, the postal address and birthdate
are required by the API; the passport fields and INN are optional.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		if !confirmed(cmd, fmt.Sprintf("Create individual registrant %q?", name), "Create") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		resident, _ := cmd.Flags().GetBool("resident")
		opts := persons.FizIPOptions{
			Name:        name,
			Resident:    resident,
			Phones:      flagString(cmd, "phones"),
			Emails:      flagString(cmd, "emails"),
			PostIndex:   flagString(cmd, "post-index"),
			PostCity:    flagString(cmd, "post-city"),
			PostAddress: flagString(cmd, "post-address"),
			Birthdate:   flagString(cmd, "birthdate"),
			PassSeries:  flagString(cmd, "pass-series"),
			PassNum:     flagString(cmd, "pass-num"),
			PassDate:    flagString(cmd, "pass-date"),
			PassOrg:     flagString(cmd, "pass-org"),
			INN:         flagString(cmd, "inn"),
			ID:          flagString(cmd, "id"),
		}
		if err := c.Persons.CreateFizIP(cmd.Context(), opts); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Saved individual registrant %q\n", name)
		return nil
	},
}

var personsCreateCompanyCmd = &cobra.Command{
	Use:   "create-company",
	Short: "Create a legal-entity registrant contact",
	Long: `Create a legal-entity ("juridical") registrant contact ("createJur").
All fields are required by the API; --phones1 is the notification phone.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		if !confirmed(cmd, fmt.Sprintf("Create legal-entity registrant %q?", name), "Create") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		resident, _ := cmd.Flags().GetBool("resident")
		opts := persons.JurOptions{
			Name:        name,
			NameTrans:   flagString(cmd, "name-trans"),
			Resident:    resident,
			Phones1:     flagString(cmd, "phones1"),
			Phones2:     flagString(cmd, "phones2"),
			Faxes:       flagString(cmd, "faxes"),
			Emails:      flagString(cmd, "emails"),
			PostIndex:   flagString(cmd, "post-index"),
			PostCity:    flagString(cmd, "post-city"),
			PostAddress: flagString(cmd, "post-address"),
			JurIndex:    flagString(cmd, "jur-index"),
			JurCity:     flagString(cmd, "jur-city"),
			JurAddress:  flagString(cmd, "jur-address"),
			INN:         flagString(cmd, "inn"),
			KPP:         flagString(cmd, "kpp"),
			PersName:    flagString(cmd, "pers-name"),
		}
		if err := c.Persons.CreateJur(cmd.Context(), opts); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Saved legal-entity registrant %q\n", name)
		return nil
	},
}

// flagString returns a string flag's value, ignoring the lookup error the way the
// existing commands do for their string flags.
func flagString(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetString(name)
	return v
}

func init() {
	// Fields shared by both registrant kinds.
	for _, c := range []*cobra.Command{personsCreateIndividualCmd, personsCreateCompanyCmd} {
		c.Flags().String("name", "", "display name / organization name")
		c.Flags().Bool("resident", true, "resident of the RF")
		c.Flags().String("emails", "", "email address(es)")
		c.Flags().String("post-index", "", "postal address: index")
		c.Flags().String("post-city", "", "postal address: city")
		c.Flags().String("post-address", "", "postal address: street")
		c.Flags().String("inn", "", "taxpayer number (INN)")
		c.Flags().Bool("yes", false, "skip the confirmation prompt")
	}

	// Individual / sole proprietor.
	personsCreateIndividualCmd.Flags().String("phones", "", "phone number(s)")
	personsCreateIndividualCmd.Flags().String("birthdate", "", "date of birth (YYYY-MM-DD)")
	personsCreateIndividualCmd.Flags().String("pass-series", "", "passport series")
	personsCreateIndividualCmd.Flags().String("pass-num", "", "passport number")
	personsCreateIndividualCmd.Flags().String("pass-date", "", "passport issue date (YYYY-MM-DD)")
	personsCreateIndividualCmd.Flags().String("pass-org", "", "passport issuing authority")
	personsCreateIndividualCmd.Flags().String("id", "", "id of the person to edit (empty creates a new one)")

	// Legal entity.
	personsCreateCompanyCmd.Flags().String("name-trans", "", "organization name, Latin transliteration")
	personsCreateCompanyCmd.Flags().String("phones1", "", "notification phone")
	personsCreateCompanyCmd.Flags().String("phones2", "", "secondary phone")
	personsCreateCompanyCmd.Flags().String("faxes", "", "fax number(s)")
	personsCreateCompanyCmd.Flags().String("jur-index", "", "legal address: index")
	personsCreateCompanyCmd.Flags().String("jur-city", "", "legal address: city")
	personsCreateCompanyCmd.Flags().String("jur-address", "", "legal address: street")
	personsCreateCompanyCmd.Flags().String("kpp", "", "tax registration reason code (KPP)")
	personsCreateCompanyCmd.Flags().String("pers-name", "", "contact representative")

	personsCmd.AddCommand(personsListCmd)
	personsCmd.AddCommand(personsInfoCmd)
	personsCmd.AddCommand(personsCreateIndividualCmd)
	personsCmd.AddCommand(personsCreateCompanyCmd)
	domainsCmd.AddCommand(personsCmd)
}
