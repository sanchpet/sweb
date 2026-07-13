package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/sites"
	"github.com/spf13/cobra"
)

// sitesCmd groups the shared-hosting website operations (SDK /sites): the read
// side (list/info/backends) plus the add/edit/remove and set-domain/set-backend
// mutations. It hangs off the hosting parent, so it inherits that group's
// profile binding.
var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Manage shared-hosting websites",
}

var sitesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's websites",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		filter, _ := cmd.Flags().GetString("filter")
		list, err := c.Sites.List(cmd.Context(), &sites.ListOptions{
			Page:    flagInt(cmd, "page"),
			PerPage: flagInt(cmd, "per-page"),
			Filter:  filter,
		})
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tALIAS\tDOCROOT\tDOMAIN")
			for _, s := range list {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", int64(s.ID), s.Alias, s.DocRoot, s.DomainTech)
			}
		})
	},
}

var sitesInfoCmd = &cobra.Command{
	Use:   "info <docroot>",
	Short: "Show full information for a website",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		info, err := c.Sites.GetSiteInfo(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return render(cmd, info, func(w io.Writer) {
			row := func(k, v string) { fmt.Fprintf(w, "%s\t%s\n", k, v) }
			row("DOCROOT", args[0])
			row("BACKEND", info.BackEnd)
			row("BACKEND ID", fmt.Sprintf("%d", int64(info.BackEndID)))
			row("ENCODING", info.Encoding)
			row("VIEW FILES", yesNo(info.ViewFiles))
			row("RUN SCRIPTS", yesNo(info.RunScripts))
			row("REDIS ENABLED", yesNo(info.RedisEnabled))
			row("REDIS SESSION", yesNo(info.RedisSessionEnabled))
			for i, d := range info.Domains {
				if i == 0 {
					row("DOMAINS", d)
					continue
				}
				row("", d)
			}
		})
	},
}

var sitesBackendsCmd = &cobra.Command{
	Use:   "backends",
	Short: "List the web back-ends available to assign to a site",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		list, err := c.Sites.BackEndsList(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tNAME")
			for _, b := range list {
				fmt.Fprintf(w, "%d\t%s\n", int64(b.ID), b.Name)
			}
		})
	},
}

var sitesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a website",
	Long: `Create a website via the "add" method.

--alias, --docroot and --domain are required; --machine (subdomain) and
--redis-session are optional.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		alias, _ := cmd.Flags().GetString("alias")
		docRoot, _ := cmd.Flags().GetString("docroot")
		domain, _ := cmd.Flags().GetString("domain")
		machine, _ := cmd.Flags().GetString("machine")
		redisSession, _ := cmd.Flags().GetBool("redis-session")
		if err := c.Sites.Add(cmd.Context(), sites.AddOptions{
			Alias:              alias,
			DocRoot:            docRoot,
			Domain:             domain,
			Machine:            machine,
			EnableRedisSession: redisSession,
		}); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", alias)
		return nil
	},
}

var sitesEditCmd = &cobra.Command{
	Use:   "edit <docroot>",
	Short: "Rename a website and/or move its docroot",
	Long: `Rename a website and/or move its docroot via the "edit" method.

<docroot> is the current home directory; --alias is the new name (required);
--docroot-new moves the home directory (empty keeps the current one).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		alias, _ := cmd.Flags().GetString("alias")
		docRootNew, _ := cmd.Flags().GetString("docroot-new")
		if err := c.Sites.Edit(cmd.Context(), args[0], alias, docRootNew); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Edited %s\n", args[0])
		return nil
	},
}

var sitesRemoveCmd = &cobra.Command{
	Use:   "remove <docroot>",
	Short: "Delete a website — destructive",
	Long: `Delete a website via the "del" method.

This is DESTRUCTIVE. You are asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		docRoot := args[0]
		if !confirmed(cmd, fmt.Sprintf("Remove %s? This cannot be undone.", docRoot), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Sites.Del(cmd.Context(), docRoot); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Removed", docRoot)
		return nil
	},
}

var sitesSetDomainCmd = &cobra.Command{
	Use:   "set-domain <domain>",
	Short: "Repoint a domain at a different website",
	Long: `Repoint a domain at a website via the "changeDomainSite" method.

<domain> is the domain to move; --docroot is the target website (required);
--machine (subdomain) is optional.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		docRoot, _ := cmd.Flags().GetString("docroot")
		machine, _ := cmd.Flags().GetString("machine")
		if err := c.Sites.ChangeDomainSite(cmd.Context(), args[0], docRoot, machine); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Repointed %s to %s\n", args[0], docRoot)
		return nil
	},
}

var sitesSetBackendCmd = &cobra.Command{
	Use:   "set-backend <docroot>",
	Short: "Switch a website's web back-end",
	Long: `Switch a website's web back-end via the "changeBackEnd" method.

<docroot> is the target website; --backend is the back-end id from
'sweb hosting sites backends' (required).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Sites.ChangeBackEnd(cmd.Context(), args[0], flagInt(cmd, "backend")); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Changed back-end for %s\n", args[0])
		return nil
	},
}

func init() {
	sitesListCmd.Flags().Int("page", 0, "1-based page number")
	sitesListCmd.Flags().Int("per-page", 0, "records per page")
	sitesListCmd.Flags().String("filter", "", "filter by site name or domain")

	sitesAddCmd.Flags().String("alias", "", "site name")
	sitesAddCmd.Flags().String("docroot", "", "home directory")
	sitesAddCmd.Flags().String("domain", "", "domain")
	sitesAddCmd.Flags().String("machine", "", "subdomain")
	sitesAddCmd.Flags().Bool("redis-session", false, "store sessions in Redis")
	_ = sitesAddCmd.MarkFlagRequired("alias")
	_ = sitesAddCmd.MarkFlagRequired("docroot")
	_ = sitesAddCmd.MarkFlagRequired("domain")

	sitesEditCmd.Flags().String("alias", "", "new site name")
	sitesEditCmd.Flags().String("docroot-new", "", "new home directory (empty keeps the current one)")
	_ = sitesEditCmd.MarkFlagRequired("alias")

	sitesRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	sitesSetDomainCmd.Flags().String("docroot", "", "target website home directory")
	sitesSetDomainCmd.Flags().String("machine", "", "subdomain")
	_ = sitesSetDomainCmd.MarkFlagRequired("docroot")

	sitesSetBackendCmd.Flags().Int("backend", 0, "back-end id from 'sweb hosting sites backends'")
	_ = sitesSetBackendCmd.MarkFlagRequired("backend")

	sitesCmd.AddCommand(
		sitesListCmd,
		sitesInfoCmd,
		sitesBackendsCmd,
		sitesAddCmd,
		sitesEditCmd,
		sitesRemoveCmd,
		sitesSetDomainCmd,
		sitesSetBackendCmd,
	)
	hostingCmd.AddCommand(sitesCmd)
}
