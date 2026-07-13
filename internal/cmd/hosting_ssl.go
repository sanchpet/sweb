package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/sanchpet/sweb-go-sdk/vh/ssl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sslFlagString reads a string flag, ignoring the lookup error (unset → "").
func sslFlagString(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetString(name)
	return v
}

// sslIsJSON reports whether the active output format is JSON (-o json).
func sslIsJSON() bool {
	return viper.GetString("output") == "json"
}

// sslIsBase64Mimetype reports whether a CertFile.Mimetype declares base64
// content encoding (e.g. "application/zip;base64").
func sslIsBase64Mimetype(mimetype string) bool {
	return strings.Contains(strings.ToLower(mimetype), "base64")
}

// sslCmd groups the shared-hosting SSL-certificate operations (SDK /vh/ssl):
// the read side (list/orders/prolong-info/download) plus the autoprolong toggle,
// prolongation, Let's Encrypt install and certificate removal. It hangs off the
// hosting parent, so it inherits that group's profile binding.
var sslCmd = &cobra.Command{
	Use:   "ssl",
	Short: "Manage shared-hosting SSL certificates",
}

// sslResolveID resolves a certificate <domain> to its numeric id by scanning the
// account's certificate index. The SDK's per-certificate methods take an id, but
// the CLI addresses a certificate by the domain it covers; an unknown domain is
// an error and an ambiguous one names the collision.
func sslResolveID(c interface {
	CertList() ([]ssl.Certificate, error)
}, domain string) (int, error) {
	list, err := c.CertList()
	if err != nil {
		return 0, err
	}
	found := -1
	for _, cert := range list {
		if cert.Domain != domain {
			continue
		}
		if found != -1 {
			return 0, fmt.Errorf("domain %q matches more than one certificate; remove/select by inspecting `sweb hosting ssl list`", domain)
		}
		found = int(cert.ID)
	}
	if found == -1 {
		return 0, fmt.Errorf("no certificate found for domain %q (see `sweb hosting ssl list`)", domain)
	}
	return found, nil
}

// sslCertList adapts the SDK client to the CertList shape sslResolveID expects,
// so the resolver stays testable without a live client.
type sslCertLister struct {
	cmd *cobra.Command
}

func (l sslCertLister) CertList() ([]ssl.Certificate, error) {
	c, err := client()
	if err != nil {
		return nil, err
	}
	res, err := c.VHSSL.List(l.cmd.Context(), nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.List, nil
}

// sslResolve resolves a <domain> arg to its certificate id via the account index.
func sslResolve(cmd *cobra.Command, domain string) (int, error) {
	return sslResolveID(sslCertLister{cmd: cmd}, domain)
}

var sslListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's SSL certificates",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		res, err := c.VHSSL.List(cmd.Context(), &ssl.ListOptions{
			Page:        flagInt(cmd, "page"),
			PerPage:     flagInt(cmd, "per-page"),
			OrderField:  sslFlagString(cmd, "order-field"),
			OrderDirect: sslFlagString(cmd, "order-direct"),
		})
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tDOMAIN\tNAME\tSTATUS\tIP\tVALID TO\tAUTOPROLONG")
			if res == nil {
				return
			}
			for _, cert := range res.List {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
					int64(cert.ID), cert.Domain, cert.Name, cert.Status, cert.IP,
					cert.ValidTo, yesNo(cert.Autoprolong))
			}
		})
	},
}

var sslOrdersCmd = &cobra.Command{
	Use:   "orders",
	Short: "List the certificate products available for order",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		list, err := c.VHSSL.OrderList(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tNAME\tTYPE\tADVANTAGE")
			for _, o := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", o.ID, o.Name, o.Type, o.AdvantageText)
			}
		})
	},
}

var sslDownloadCmd = &cobra.Command{
	Use:   "download <domain>",
	Short: "Download an issued certificate's archive files",
	Long: `Download an issued certificate's files via the "download" method.

<domain> selects the certificate; the account password is required and read from
--password or a masked prompt. The archive carries the private key — by default
the files are written to --dir (current directory); with -o json the base64
content is emitted instead. Handle the key material accordingly.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := sslResolve(cmd, args[0])
		if err != nil {
			return err
		}
		password, err := sslResolvePassword(cmd)
		if err != nil {
			return err
		}
		files, err := c.VHSSL.Download(cmd.Context(), id, password)
		if err != nil {
			return err
		}
		// -o json emits the archive verbatim (base64 content included).
		if sslIsJSON() {
			return render(cmd, files, nil)
		}
		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			dir = "."
		}
		for _, f := range files {
			body, err := sslDecodeContent(f)
			if err != nil {
				return fmt.Errorf("decode %s: %w", f.Name, err)
			}
			path := filepath.Join(dir, filepath.Base(f.Name))
			if err := os.WriteFile(path, body, 0o600); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Wrote", path)
		}
		return nil
	},
}

var sslProlongInfoCmd = &cobra.Command{
	Use:   "prolong-info <domain>",
	Short: "Show the prolongation options for a certificate",
	Long: `Show a certificate's prolongation options via the "getProlongInfo" method.

<domain> selects the certificate; the result lists the per-period prices and the
product ids to pass to 'sweb hosting ssl prolong --product-id'.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := sslResolve(cmd, args[0])
		if err != nil {
			return err
		}
		info, err := c.VHSSL.ProlongInfo(cmd.Context(), id)
		if err != nil {
			return err
		}
		return render(cmd, info, func(w io.Writer) {
			if info == nil {
				fmt.Fprintln(w, "no prolongation offered")
				return
			}
			fmt.Fprintf(w, "TITLE\t%s\n", info.Title)
			fmt.Fprintf(w, "CURRENT ID\t%d\n", int64(info.CurrentCertificateID))
			fmt.Fprintf(w, "FREE\t%s\n", yesNo(info.IsFreeCertificate))
			fmt.Fprintln(w, "PERIOD(MONTHS)\tPRICE\tPRODUCT ID")
			for period, price := range info.Prices {
				fmt.Fprintf(w, "%s\t%.2f\t%s\n", period, float64(price), info.IDs[period])
			}
		})
	},
}

var sslAutoprolongCmd = &cobra.Command{
	Use:   "autoprolong <domain>",
	Short: "Enable or disable a certificate's auto-prolongation",
	Long: `Toggle a certificate's auto-prolongation via the "editAutoprolong" method.

<domain> selects the certificate; --enable turns auto-prolongation on,
--enable=false turns it off.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := sslResolve(cmd, args[0])
		if err != nil {
			return err
		}
		enabled, _ := cmd.Flags().GetBool("enable")
		res, err := c.VHSSL.EditAutoprolong(cmd.Context(), id, enabled)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Set auto-prolongation for %s to %s (%s)\n", args[0], yesNo(enabled), res)
		return nil
	},
}

var sslProlongCmd = &cobra.Command{
	Use:   "prolong <domain>",
	Short: "Prolong a certificate — bills the account",
	Long: `Prolong a certificate via the "prolongCertificate" method.

<domain> selects the current certificate; --product-id is the prolongation
product id from 'sweb hosting ssl prolong-info' (required). This BILLS the
account — you are asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := sslResolve(cmd, args[0])
		if err != nil {
			return err
		}
		productID := flagInt(cmd, "product-id")
		if !confirmed(cmd, fmt.Sprintf("Prolong the certificate for %s? This bills the account.", args[0]), "Prolong") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.VHSSL.ProlongCertificate(cmd.Context(), id, productID)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Prolonged %s (%s)\n", args[0], res)
		return nil
	},
}

var sslInstallLetsEncryptCmd = &cobra.Command{
	Use:   "install-letsencrypt <domain>",
	Short: "Install a free Let's Encrypt certificate",
	Long: `Install a free Let's Encrypt certificate via the "installLetsEncrypt" method.

<domain> is the base domain; --wildcard requests a wildcard certificate,
--subdomain covers a subdomain, --ip targets an IP ("sni" for SNI), and
--challenge selects the validation type ("acme" or "dns").`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		wildcard, _ := cmd.Flags().GetBool("wildcard")
		opts := &ssl.InstallLetsEncryptOptions{
			Virtdom:   sslFlagString(cmd, "subdomain"),
			IP:        sslFlagString(cmd, "ip"),
			Challenge: sslFlagString(cmd, "challenge"),
		}
		res, err := c.VHSSL.InstallLetsEncrypt(cmd.Context(), args[0], wildcard, opts)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Installed Let's Encrypt for %s (%s)\n", args[0], res)
		return nil
	},
}

var sslRemoveCmd = &cobra.Command{
	Use:   "remove <domain>",
	Short: "Delete a certificate — destructive",
	Long: `Delete a certificate via the "removeCertificate" method.

<domain> selects the certificate. This is DESTRUCTIVE. You are asked to confirm
unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := sslResolve(cmd, args[0])
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Remove the certificate for %s? This cannot be undone.", args[0]), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.VHSSL.RemoveCertificate(cmd.Context(), id)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Removed the certificate for %s (%s)\n", args[0], res)
		return nil
	},
}

// sslResolvePassword returns the account password for a download: the --password
// flag when set, otherwise a masked interactive prompt. An empty result
// (non-interactive with no flag) is an error, since the API requires one.
func sslResolvePassword(cmd *cobra.Command) (string, error) {
	if pw, _ := cmd.Flags().GetString("password"); pw != "" {
		return pw, nil
	}
	var pw string
	if err := huh.NewInput().
		Title("Account password").
		EchoMode(huh.EchoModePassword).
		Value(&pw).
		Run(); err != nil {
		return "", err
	}
	if pw == "" {
		return "", fmt.Errorf("a password is required: pass --password or enter one at the prompt")
	}
	return pw, nil
}

// sslDecodeContent decodes a downloaded certificate file body, base64-decoding it
// when the mimetype declares base64 encoding and returning it verbatim otherwise.
func sslDecodeContent(f ssl.CertFile) ([]byte, error) {
	if sslIsBase64Mimetype(f.Mimetype) {
		return base64.StdEncoding.DecodeString(f.Content)
	}
	return []byte(f.Content), nil
}

func init() {
	sslListCmd.Flags().Int("page", 0, "1-based page number")
	sslListCmd.Flags().Int("per-page", 0, "records per page")
	sslListCmd.Flags().String("order-field", "", "sort field: id|valid_to|fqdn|status|ip")
	sslListCmd.Flags().String("order-direct", "", "sort direction: asc|desc")

	sslDownloadCmd.Flags().String("password", "", "account password (prompted if unset)")
	sslDownloadCmd.Flags().String("dir", "", "directory to write the archive files (default current)")

	sslAutoprolongCmd.Flags().Bool("enable", true, "enable auto-prolongation (--enable=false to disable)")

	sslProlongCmd.Flags().Int("product-id", 0, "prolongation product id from 'sweb hosting ssl prolong-info'")
	_ = sslProlongCmd.MarkFlagRequired("product-id")
	sslProlongCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	sslInstallLetsEncryptCmd.Flags().Bool("wildcard", false, "request a wildcard certificate")
	sslInstallLetsEncryptCmd.Flags().String("subdomain", "", "subdomain to cover, e.g. sub.mysite.ru")
	sslInstallLetsEncryptCmd.Flags().String("ip", "", `target IP, or "sni" for SNI`)
	sslInstallLetsEncryptCmd.Flags().String("challenge", "", `validation type: "acme" or "dns"`)

	sslRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	sslCmd.AddCommand(
		sslListCmd,
		sslOrdersCmd,
		sslDownloadCmd,
		sslProlongInfoCmd,
		sslAutoprolongCmd,
		sslProlongCmd,
		sslInstallLetsEncryptCmd,
		sslRemoveCmd,
	)
	hostingCmd.AddCommand(sslCmd)
}
