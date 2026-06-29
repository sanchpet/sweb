# sweb

[![CI](https://github.com/sanchpet/sweb/actions/workflows/ci.yml/badge.svg)](https://github.com/sanchpet/sweb/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sanchpet/sweb)](https://goreportcard.com/report/github.com/sanchpet/sweb)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A command-line interface for the [SpaceWeb](https://sweb.ru) (sweb.ru) hosting
API, modeled after `kubectl` / `yc`. Built on
[`sweb-go-sdk`](https://github.com/sanchpet/sweb-go-sdk) (Cobra + Viper + Charm
Fang).

## Install

**[mise](https://mise.jdx.dev)** (installs the released binary via the `ubi` backend):

```sh
mise use -g ubi:sanchpet/sweb        # latest release
mise use -g ubi:sanchpet/sweb@0.1.0  # a specific version
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
# 1. Authenticate (prompts for login/password, stores the token in the OS
#    keyring — falls back to a 0600 file when none is available). Or set
#    SWEB_TOKEN in the environment.
sweb configure
sweb configure --insecure-storage   # force the plaintext-file backend

# 2. Inspect the catalog (plan / OS / datacenter IDs).
sweb vps config

# 3. List your VPS instances (-o json for machine output).
sweb vps list
sweb vps list -o json

# 4. Provision a VPS (mutates — bills your account).
sweb vps create --plan 4 --distributive 32 --datacenter 1 \
  --alias hub --ssh-key "$(cat ~/.ssh/id_ed25519.pub)"
```

Token resolution precedence: `--token` flag → `$SWEB_TOKEN` → OS keyring → config file.

## Develop

```sh
mise install        # Go + golangci-lint + pre-commit
mise run ci         # lint + test + build
```

## License

MIT — see [LICENSE](LICENSE).
