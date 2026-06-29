# sweb

[![CI](https://github.com/sanchpet/sweb/actions/workflows/ci.yml/badge.svg)](https://github.com/sanchpet/sweb/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sanchpet/sweb)](https://goreportcard.com/report/github.com/sanchpet/sweb)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A command-line interface for the [SpaceWeb](https://sweb.ru) (sweb.ru) hosting
API, modeled after `kubectl` / `yc`. Built on
[`sweb-go-sdk`](https://github.com/sanchpet/sweb-go-sdk) (Cobra + Viper + Charm
Fang).

## Install

```sh
go install github.com/sanchpet/sweb@latest
```

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
