# sweb — instructions for Claude Code

A CLI for the SpaceWeb (sweb.ru) hosting API, modeled after `kubectl` / `yc`.
It is a **thin presentation layer** over [`sweb-go-sdk`](https://github.com/sanchpet/sweb-go-sdk):
all API/domain logic lives in the SDK; this repo handles commands, flags,
config, and output. Keep it that way — no HTTP/JSON-RPC here.

## Architecture

- `main.go` → `internal/cmd.Execute()` (Cobra root wrapped by Charm **Fang**).
- `internal/cmd/` — **one file per command**, `<group>_<verb>.go` (e.g.
  `vps_create.go`, `vps_change_plan.go`, `vps_power.go` = start/stop/reboot,
  `vps_status.go`). Top level: `configure`, `token`, `vps`. `root.go` holds config
  wiring + the shared `client()` and `render()` helpers; `resolve.go` and
  `completion.go` hold cross-command helpers. A new command is a new file with an
  `init()` that hangs it off its parent — no central registry to edit.
- **Config / auth:** SpaceWeb tokens are short-lived with no refresh flow, so
  `configure` (Charm **huh** prompt) stores **login + password + token** in the OS
  keyring (`zalando/go-keyring`), falling back to a 0600 file — the fallback is
  surfaced explicitly (never silent). `client()` (root.go) passes the cached token
  plus credentials to the SDK (`WithCredentials` + `WithOnTokenRefresh`), which
  re-auths transparently on expiry; the callback persists the new token via
  `saveToken`. Precedence: `--token` flag / `$SWEB_TOKEN` (one-off, no refresh) →
  stored credentials. `-o table|json` for output; non-auth config via Viper.
- **Output:** `-o table|json`. Every command renders through `render(cmd, data,
  tableFn)` — `-o json` marshals `data`, otherwise `tableFn` writes a `tabwriter`
  table. Keep both paths in sync.

## Command patterns (match these)

- **Address a VPS by name or id.** Commands that take a `<vps>` accept the name
  (alias) *or* the billing id and pass the arg through `resolveVPS` (`resolve.go`);
  the SDK only ever sees the billing id. An ambiguous name is an error.
- **Async ops wait explicitly.** Long/async operations (resize, power) return as
  soon as the API accepts, and take a flag to block on `WaitForIdle` — printing
  each `current_action` phase. Follow the existing polarity per command
  (`change-plan` waits by default with `--async`; `start/stop/reboot` are
  fire-and-report with `--wait`).
- **Destructive ops confirm.** Mutating/irreversible commands prompt via Charm
  `huh` unless `--yes` (see `vps delete`).
- **Shell completion.** Register `ValidArgsFunction` on new commands — reuse
  `completeBillingIDs` / `completeCategories` (`completion.go`) so `<vps>` and id
  args complete.
- **Testing.** `cmd_test.go` asserts the command tree (add new subcommands there).
  Never call mutating/billing API in tests — there is no live-API test path.

## Build & test (mise-first)

```sh
mise install
mise run ci    # lint + test + build
pre-commit install && pre-commit run -a
```

## Conventions

- **English** for all repo artifacts (code, comments, docs, commits, PRs).
- Small focused commits; `--signoff` on every commit + an
  `Assisted-By: Claude <noreply@anthropic.com>` trailer (personal repo — the owner
  authors, Claude assists; not `Co-Authored-By`). Branch + PR; do not self-merge.
- **Conventional Commits + release-please (BLOCKING):** commit / PR-title format is
  `<type>[scope]: <desc>` (`feat`→minor, `fix`→patch, `!` or `BREAKING CHANGE`→major).
  PRs are squash-merged, so the **PR title is the release commit** — CI enforces its
  format (`pr-title` workflow). Versioning and `CHANGELOG.md` are automated by
  **release-please** (merging its release PR tags + runs GoReleaser) — never `git tag`
  or edit the changelog by hand. See `CONTRIBUTING.md`.

## Security / opsec (BLOCKING)

- **Never commit tokens or credentials.** `configure` writes the token to the OS
  keyring (or the user's `~/.config/sweb/`), never into the repo. `gitleaks` runs in
  pre-commit.
- `sweb vps create` mutates and bills — never invoke it in tests/CI.
