# Changelog

## [0.11.0](https://github.com/sanchpet/sweb/compare/v0.10.0...v0.11.0) (2026-07-10)


### Features

* add 'vps reinstall|clone|logs' commands ([#32](https://github.com/sanchpet/sweb/issues/32)) ([9bef775](https://github.com/sanchpet/sweb/commit/9bef77598948cbdf06a979d6cd661e4026c36662))

## [0.10.0](https://github.com/sanchpet/sweb/compare/v0.9.2...v0.10.0) (2026-07-10)


### Features

* add 'vps start|stop|reboot|status' power/state commands ([#29](https://github.com/sanchpet/sweb/issues/29)) ([a97ec28](https://github.com/sanchpet/sweb/commit/a97ec28beb1b6ca9519a5c882d19e10d3fdf78d6))

## [0.9.2](https://github.com/sanchpet/sweb/compare/v0.9.1...v0.9.2) (2026-07-02)


### Bug Fixes

* bump sweb-go-sdk to v0.8.2 (local_ip object-shape decode) ([#24](https://github.com/sanchpet/sweb/issues/24)) ([6673c0c](https://github.com/sanchpet/sweb/commit/6673c0c1d1e0336641590cf30bfbdcea98ef34ee))

## [0.9.1](https://github.com/sanchpet/sweb/compare/v0.9.0...v0.9.1) (2026-07-02)


### Bug Fixes

* bump sweb-go-sdk to v0.8.1 (IP price FlexFloat) ([#22](https://github.com/sanchpet/sweb/issues/22)) ([0d43c94](https://github.com/sanchpet/sweb/commit/0d43c9415e918cbe88aedc8691c1549b9e56adf2))

## [0.9.0](https://github.com/sanchpet/sweb/compare/v0.8.0...v0.9.0) (2026-07-02)


### Features

* accept VPS name (alias), not just billing id, in vps commands ([#20](https://github.com/sanchpet/sweb/issues/20)) ([90da204](https://github.com/sanchpet/sweb/commit/90da2042be89ed04881d7441cfbea013a2a792ec))

## [0.8.0](https://github.com/sanchpet/sweb/compare/v0.7.0...v0.8.0) (2026-07-02)


### Features

* add 'vps local-ip' (private network attach/detach/show) ([#18](https://github.com/sanchpet/sweb/issues/18)) ([a30e2ee](https://github.com/sanchpet/sweb/commit/a30e2eed719aa642065b1a838bc2d55a09ea2c8a))

## Changelog

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
