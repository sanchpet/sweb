# sweb

[![CI](https://github.com/sanchpet/sweb/actions/workflows/ci.yml/badge.svg)](https://github.com/sanchpet/sweb/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sanchpet/sweb)](https://goreportcard.com/report/github.com/sanchpet/sweb)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A command-line interface for the [SpaceWeb](https://sweb.ru) (sweb.ru) hosting
API, modeled after `kubectl` / `yc`. Built on
[`sweb-go-sdk`](https://github.com/sanchpet/sweb-go-sdk) (Cobra + Viper + Charm
Fang).

## Install

**[mise](https://mise.jdx.dev)** (installs the released binary via the `github` backend,
with SLSA provenance verification):

```sh
mise use -g github:sanchpet/sweb        # latest release
mise use -g github:sanchpet/sweb@0.1.0  # a specific version
```

**Go** (builds from source; needs a Go toolchain):

```sh
go install github.com/sanchpet/sweb@latest
```

**Pre-built binary** — download the archive for your OS/arch from the
[Releases](https://github.com/sanchpet/sweb/releases) page, verify it against
`checksums.txt`, extract, and put `sweb` on your `PATH`.

**Debian/Ubuntu** and **RHEL/Fedora** — grab the `.deb` / `.rpm` from the
[Releases](https://github.com/sanchpet/sweb/releases) page:

```sh
sudo dpkg -i sweb_*_linux_amd64.deb   # Debian/Ubuntu
sudo rpm -i  sweb_*_linux_amd64.rpm   # RHEL/Fedora
```

> Homebrew (`brew install`) is planned — it needs a dedicated tap repository.

Releases are produced by [GoReleaser](https://goreleaser.com) on every `v*` tag
(cross-platform binaries + checksums + deb/rpm).

## Usage

```sh
# 1. Authenticate (prompts for login/password, stores both + the token in the OS
#    keyring — falls back to a 0600 file when none is available). SpaceWeb tokens
#    are short-lived with no refresh flow, so the password is kept (in the
#    keyring) to re-auth transparently — you are not prompted again until it
#    changes. Or set SWEB_TOKEN in the environment for a one-off token.
sweb configure
sweb configure --insecure-storage   # force the plaintext-file backend

# 2. Inspect the catalog (plan / OS / datacenter IDs).
sweb vps config

# 3. List your VPS instances (-o json for machine output).
sweb vps list
sweb vps list -o json

# 4. Provision a VPS (mutates — bills your account).
#    --dry-run prints the request without calling the API.
sweb vps create --plan 4 --distributive 32 --datacenter 1 \
  --alias hub --ssh-key "$(cat ~/.ssh/id_ed25519.pub)" --dry-run
sweb vps create --plan 4 --distributive 32 --datacenter 1 \
  --alias hub --ssh-key "$(cat ~/.ssh/id_ed25519.pub)"

#    Or build a custom spec with the configurator (ram/disk in GB; NVMe default).
sweb vps create --cpu 2 --ram 6 --disk 15 --distributive 113 --datacenter 1 \
  --alias infra-hub --ssh-key "$(cat ~/.ssh/id_ed25519.pub)" --dry-run

# 5. Rename a VPS (change its alias) in place — no reprovision, no bill.
sweb vps rename login_vps_6 infra-01

# 6. Delete a VPS by its BILLING_ID (from `vps list`); asks to confirm.
sweb vps delete login_vps_6
sweb vps delete login_vps_6 --yes   # skip the confirmation prompt

# 7. Mint a fresh token to stdout, for env-based tooling (like `yc iam create-token`).
#    Feeds the Terraform provider / any $SWEB_TOKEN consumer without exporting creds.
export SWEB_TOKEN=$(sweb token)
```

### Feeding env-based tools (Terraform, mise)

`sweb token` prints only the token, so tools that read `$SWEB_TOKEN` need no login/password
in the environment. Wire it into [mise](https://mise.jdx.dev) so a token is injected
automatically per directory (analogous to `YC_TOKEN` via direnv):

```toml
# mise.toml — inject a fresh token when entering this dir
[env]
SWEB_TOKEN = "{{ exec(command='sweb token') }}"
```

Configure once (`sweb configure`); the token is minted fresh on each mise activation.

Auth precedence: `--token` flag → `$SWEB_TOKEN` (both one-off, no refresh) →
stored credentials (auto-refreshing). The token is cached and silently renewed
from the stored login+password when it expires.

### Exporting a DNS zone to Terraform

`dns records -o json` is machine-readable so it can seed declarative config. The
[`terraform-provider-sweb`](https://github.com/sanchpet/terraform-provider-sweb)
ships a `tf-dns-import` helper that turns a zone dump into Terraform `import {}`
blocks; `terraform plan -generate-config-out` then writes the HCL:

```sh
# domains live on the hosting account — mint that account's token
export TF_VAR_sweb_token="$(sweb token --profile hosting)"

sweb dns records example.com -o json | tf-dns-import example.com > imports.tf
terraform plan -generate-config-out=generated.tf
```

See the provider's [Importing an existing DNS zone](https://github.com/sanchpet/terraform-provider-sweb/blob/main/docs/guides/import-existing-zone.md)
guide for the full walkthrough. Use `--profile` to target the account that owns
the domain: with the wrong account's credentials every read returns
`-32500 Нет доступа к домену`.

## Shell completion

The `.deb` / `.rpm` packages install completions automatically (bash, zsh, fish).
For `go install` or the prebuilt binary, generate and install them yourself:

```sh
# zsh
sweb completion zsh > ~/.zsh/completions/_sweb   # a dir on your $fpath, then: autoload -Uz compinit && compinit

# bash (needs bash-completion)
sweb completion bash | sudo tee /etc/bash_completion.d/sweb >/dev/null

# quick try in the current shell (any shell)
source <(sweb completion zsh)   # or bash / fish
```

`vps create` completes `--plan`, `--datacenter` and `--distributive` with live
values from your account catalog (requires a configured token).

## Develop

```sh
mise install        # Go + golangci-lint + pre-commit
mise run ci         # lint + test + build
```

## License

MIT — see [LICENSE](LICENSE).
