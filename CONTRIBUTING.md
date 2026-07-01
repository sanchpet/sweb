# Contributing

## Commits & PR titles (Conventional Commits)

This repo uses [Conventional Commits](https://www.conventionalcommits.org/) with
[release-please](https://github.com/googleapis/release-please) for versioning and
the changelog. PRs are **squash-merged**, so the **PR title** becomes the commit
on `main` and drives the release — it MUST be a valid Conventional Commit (CI
enforces this):

```
<type>[optional scope]: <description>
```

- `feat:` → minor bump · `fix:` → patch bump
- `feat!:` / `fix!:` or a `BREAKING CHANGE:` footer → major bump
- `docs:` `chore:` `refactor:` `test:` `ci:` `perf:` → no release on their own

Examples: `feat: add vps change-plan --wait`, `fix: render fractional price`.

## Branches

Branch from `main`; name it `<type>/<slug>` (e.g. `feat/change-plan-wait`).
Branches are deleted automatically on merge.

## Releases — automated, do not hand-edit

Do **not** tag manually or edit `CHANGELOG.md`. release-please opens a *release
PR* that maintains the version and changelog from merged commits; merging it tags
the release and GoReleaser builds/uploads the binaries. Refine wording in the
release PR if needed.

## Local checks

`mise run ci` (lint + test + build) must pass. Keep this a thin presentation layer
over `sweb-go-sdk` — no HTTP/JSON-RPC here.
