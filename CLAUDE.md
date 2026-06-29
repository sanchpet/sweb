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
- **Config / token:** `token.go` resolves the token by precedence `--token` flag →
  `$SWEB_TOKEN` → OS keyring → config file. `configure` (Charm **huh** prompt) stores
  the token in the OS keyring (`zalando/go-keyring`), falling back to a 0600 file —
  the fallback is surfaced explicitly (never silent). Only the token is stored, never
  the password. `-o table|json` for output. Non-token config via Viper.
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
