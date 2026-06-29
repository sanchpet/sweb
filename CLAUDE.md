# sweb — instructions for Claude Code

A CLI for the SpaceWeb (sweb.ru) hosting API, modeled after `kubectl` / `yc`.
It is a **thin presentation layer** over [`sweb-go-sdk`](https://github.com/sanchpet/sweb-go-sdk):
all API/domain logic lives in the SDK; this repo handles commands, flags,
config, and output. Keep it that way — no HTTP/JSON-RPC here.

## Architecture

- `main.go` → `internal/cmd.Execute()` (Cobra root wrapped by Charm **Fang**).
- `internal/cmd/` — one file per command: `configure`, `vps`, `vps_list`,
  `vps_create`, `vps_config`. `root.go` holds config wiring + the `client()` and
  `render()` helpers.
- **Config / auth:** SpaceWeb tokens are short-lived with no refresh flow, so
  `configure` (Charm **huh** prompt) stores **login + password + token** in the OS
  keyring (`zalando/go-keyring`), falling back to a 0600 file — the fallback is
  surfaced explicitly (never silent). `client()` (root.go) passes the cached token
  plus credentials to the SDK (`WithCredentials` + `WithOnTokenRefresh`), which
  re-auths transparently on expiry; the callback persists the new token via
  `saveToken`. Precedence: `--token` flag / `$SWEB_TOKEN` (one-off, no refresh) →
  stored credentials. `-o table|json` for output; non-auth config via Viper.
- **Output:** `-o table|json`. Tables use stdlib `text/tabwriter` (no extra dep).

## Build & test (mise-first)

```sh
mise install
mise run ci    # lint + test + build
pre-commit install && pre-commit run -a
```

## Conventions

- **English** for all repo artifacts (code, comments, docs, commits, PRs).
- Small focused commits; `--signoff` on every commit + a `Co-Authored-By: Claude`
  trailer (personal repo). Branch + PR; do not self-merge.

## Security / opsec (BLOCKING)

- **Never commit tokens or credentials.** `configure` writes the token to the OS
  keyring (or the user's `~/.config/sweb/`), never into the repo. `gitleaks` runs in
  pre-commit.
- `sweb vps create` mutates and bills — never invoke it in tests/CI.
