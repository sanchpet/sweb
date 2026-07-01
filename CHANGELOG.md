# Changelog

From v0.8.0 on, this file is maintained automatically by
[release-please](https://github.com/googleapis/release-please) from
[Conventional Commit](https://www.conventionalcommits.org/) messages — see
[CONTRIBUTING.md](CONTRIBUTING.md).

## Releases up to 0.7.0 (pre-automation)

- **0.7.0** — `vps change-plan` waits for the resize by default (`--async` to skip).
- **0.6.0** — `vps change-plan --wait`: poll until the resize settles, printing each phase.
- **0.5.1** — sweb-go-sdk v0.6.1 (`changePlan` `planId` fix).
- **0.5.0** — `vps change-plan` accepts configurator flags (`--cpu`/`--ram`/`--disk`).
- **0.4.1** — sweb-go-sdk v0.6.0 (configurator sold-out guard).
- **0.4.0** — `vps change-plan` (in-place resize).
- **0.3.1** — fractional price rendering; sweb-go-sdk v0.3.0 (reconciled index, Flex types).
- **0.3.0** — `sweb token` command.
- **0.2.0** — `vps rename` command.
- **0.1.4** — configurator mode (`create --cpu`/`--ram`/`--disk`).
- **0.1.3** — transparent token refresh (store login+password).
- **0.1.2** — `vps delete`, `create --dry-run`, `BILLING_ID` in list.
- **0.1.1** — dynamic shell completion + completions bundled in releases.
- **0.1.0** — initial: `configure` (keyring) + `vps list`/`create`/`config`.
