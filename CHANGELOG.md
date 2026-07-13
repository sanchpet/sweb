# Changelog

## [0.17.0](https://github.com/sanchpet/sweb/compare/v0.16.3...v0.17.0) (2026-07-13)


### Features

* **cli:** add domains bonus commands ([#74](https://github.com/sanchpet/sweb/issues/74)) ([6a80c0e](https://github.com/sanchpet/sweb/commit/6a80c0e869c91ef97cc58d98e126ef8c0bb2e9b3))
* **cli:** add domains persons commands ([#73](https://github.com/sanchpet/sweb/issues/73)) ([7d6c7b4](https://github.com/sanchpet/sweb/commit/7d6c7b41e6fce9429580c7ad250cf187ebecb901))
* **cli:** add hosting backup commands ([#67](https://github.com/sanchpet/sweb/issues/67)) ([2139440](https://github.com/sanchpet/sweb/commit/21394408d2af4b885d8b66ff543254d35f9659a7))
* **cli:** add hosting cron commands ([#60](https://github.com/sanchpet/sweb/issues/60)) ([1d49641](https://github.com/sanchpet/sweb/commit/1d49641330b75cefe34f7de54dc2f1007878c4a6))
* **cli:** add hosting db commands ([#61](https://github.com/sanchpet/sweb/issues/61)) ([46b39b0](https://github.com/sanchpet/sweb/commit/46b39b0c8039f97764ffc625cb45d41041641f84))
* **cli:** add hosting ddg commands ([#65](https://github.com/sanchpet/sweb/issues/65)) ([938fa9d](https://github.com/sanchpet/sweb/commit/938fa9d2599592b7efc6b332aab83ee84a586d5d))
* **cli:** add hosting disk-usage commands ([#64](https://github.com/sanchpet/sweb/issues/64)) ([15f2aa6](https://github.com/sanchpet/sweb/commit/15f2aa6138f046b719697d1cafb87a90337dded4))
* **cli:** add hosting load commands ([#66](https://github.com/sanchpet/sweb/issues/66)) ([1ed3bb4](https://github.com/sanchpet/sweb/commit/1ed3bb421fe2f01fbd7b9faa785cf1642431cb9f))
* **cli:** add hosting mail commands ([#63](https://github.com/sanchpet/sweb/issues/63)) ([adf2893](https://github.com/sanchpet/sweb/commit/adf289321e893db0865338bc2e91244dd0c6865d))
* **cli:** add hosting partner commands ([#71](https://github.com/sanchpet/sweb/issues/71)) ([03dcbdf](https://github.com/sanchpet/sweb/commit/03dcbdf762e7b5498fd1808df87bd436ee01c161))
* **cli:** add hosting referral commands ([#69](https://github.com/sanchpet/sweb/issues/69)) ([08be049](https://github.com/sanchpet/sweb/commit/08be04941b0e32abf5539c2273c300d217ffdd84))
* **cli:** add hosting sites commands ([#59](https://github.com/sanchpet/sweb/issues/59)) ([54c6af3](https://github.com/sanchpet/sweb/commit/54c6af375f9acbe5b6907d5436b45e4d52a5f90d))
* **cli:** add hosting ssh commands ([#68](https://github.com/sanchpet/sweb/issues/68)) ([eca20f9](https://github.com/sanchpet/sweb/commit/eca20f9db8f3a4594eeeee62f14e63e951d28cde))
* **cli:** add hosting ssl commands ([#70](https://github.com/sanchpet/sweb/issues/70)) ([77404fa](https://github.com/sanchpet/sweb/commit/77404fa364175affb37c413f8a4deb14eabfdd10))
* **cli:** add pay commands ([#75](https://github.com/sanchpet/sweb/issues/75)) ([694504b](https://github.com/sanchpet/sweb/commit/694504b1bfcf3b68b1067f1ba6f1e8db3b8aa043))
* **cli:** add tariff commands ([#72](https://github.com/sanchpet/sweb/issues/72)) ([2e98e07](https://github.com/sanchpet/sweb/commit/2e98e07cd4795f1dbcee6300d2aa370d132ecf65))

## [0.16.3](https://github.com/sanchpet/sweb/compare/v0.16.2...v0.16.3) (2026-07-12)


### Bug Fixes

* make domains price resilient to per-check API refusals ([#54](https://github.com/sanchpet/sweb/issues/54)) ([9989196](https://github.com/sanchpet/sweb/commit/9989196e8fa6c50fc6f4d47a4a006adc3ac06301))

## [0.16.2](https://github.com/sanchpet/sweb/compare/v0.16.1...v0.16.2) (2026-07-12)


### Bug Fixes

* drop the misleading --prefix flag from dns record ([#52](https://github.com/sanchpet/sweb/issues/52)) ([54cd691](https://github.com/sanchpet/sweb/commit/54cd69158205e4176a3909c1d002e1cea488007f))

## [0.16.1](https://github.com/sanchpet/sweb/compare/v0.16.0...v0.16.1) (2026-07-12)


### Bug Fixes

* show the subdomain host for TXT records in the dns records table ([#50](https://github.com/sanchpet/sweb/issues/50)) ([21be1a0](https://github.com/sanchpet/sweb/commit/21be1a01b15ff88014c2ba4fe04b662f1ad98e35))

## [0.16.0](https://github.com/sanchpet/sweb/compare/v0.15.3...v0.16.0) (2026-07-12)


### Features

* add domains command group over the /domains lifecycle ([#48](https://github.com/sanchpet/sweb/issues/48)) ([cdabd5d](https://github.com/sanchpet/sweb/commit/cdabd5db3fc502d1b4f2b26d5a4ad8d509861c82))

## [0.15.3](https://github.com/sanchpet/sweb/compare/v0.15.2...v0.15.3) (2026-07-11)


### Bug Fixes

* cap the DNS records VALUE column width ([#46](https://github.com/sanchpet/sweb/issues/46)) ([00e77df](https://github.com/sanchpet/sweb/commit/00e77dff37e0327098d6a672f18f292a22b96986))

## [0.15.2](https://github.com/sanchpet/sweb/compare/v0.15.1...v0.15.2) (2026-07-11)


### Bug Fixes

* bound the DNS records VALUE column width ([#44](https://github.com/sanchpet/sweb/issues/44)) ([58ac2fd](https://github.com/sanchpet/sweb/commit/58ac2fd322db3353604c47b0a194f36802b107ca))

## [0.15.1](https://github.com/sanchpet/sweb/compare/v0.15.0...v0.15.1) (2026-07-11)


### Bug Fixes

* bump sweb-go-sdk to v0.13.1 for DNS info array decode ([#42](https://github.com/sanchpet/sweb/issues/42)) ([00be648](https://github.com/sanchpet/sweb/commit/00be648a592f055dda7b6d9714a118a2567bab0c))

## [0.15.0](https://github.com/sanchpet/sweb/compare/v0.14.0...v0.15.0) (2026-07-11)


### Features

* multi-account support via credential profiles ([#40](https://github.com/sanchpet/sweb/issues/40)) ([a1dfbf6](https://github.com/sanchpet/sweb/commit/a1dfbf6583dbd2d20d1502691952f8196fb29407))

## [0.14.0](https://github.com/sanchpet/sweb/compare/v0.13.0...v0.14.0) (2026-07-11)


### Features

* add dns command group ([#38](https://github.com/sanchpet/sweb/issues/38)) ([78a02ad](https://github.com/sanchpet/sweb/commit/78a02ad95c3d378a9e6743c10c302f022b6eca70))

## [0.13.0](https://github.com/sanchpet/sweb/compare/v0.12.0...v0.13.0) (2026-07-10)


### Features

* add 'vps backup' and 'vps cloud-backup' commands ([#36](https://github.com/sanchpet/sweb/issues/36)) ([acf1d4f](https://github.com/sanchpet/sweb/commit/acf1d4f5b929ca840125ae1e27af1eca0c4f9bf7))

## [0.12.0](https://github.com/sanchpet/sweb/compare/v0.11.0...v0.12.0) (2026-07-10)


### Features

* add 'vps ip' commands (public IPs + PTR) ([#34](https://github.com/sanchpet/sweb/issues/34)) ([b6e14e5](https://github.com/sanchpet/sweb/commit/b6e14e511043fc0a99eb8c376a29ee52d7f21587))

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
