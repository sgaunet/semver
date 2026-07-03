# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Operating Guidelines

**Read `docs/operating-guidelines.md` at the start of every session.** It
defines how to plan, verify, and iterate in this repository: plan mode,
subagent strategy, verification gates, self-improvement loop, and the
communication contract. Treat it as load-bearing context.

## Repository Overview

A single, statically linked Go CLI for manipulating Semantic Versions: bump,
pre-release lifecycle, compare, sort, validate, extract components, and test
against range constraints. Stdlib-only — no third-party runtime dependencies
(there is no go.sum). Go 1.25 per `go.mod`; the local toolchain (go,
golangci-lint, goreleaser, task) is pinned in `mise.toml`.

## Architecture

- Strict one-directional layering: `internal/version` (leaf domain package) ← `internal/constraint` ← `internal/cli` ← `cmd/semver`. Domain packages never import CLI code.
- Registry-driven subcommand dispatch (`internal/cli/app.go`) with a fixed exit-code contract in `internal/cli/exit.go` (0/1/2/10/11).
- stdout = machine data, stderr = human diagnostics — enforced by the `streams` abstraction in `internal/cli/output.go` (text and JSON emission).
- Sentinel errors per package (`internal/version/errors.go`, `internal/constraint/errors.go`) wrapped with `fmt.Errorf("...: %w", ...)`; keep them `errors.Is`-matchable (err113 rule).
- Immutable value-type domain model: transformations return new `Version` values via `clone()`, never mutate the receiver.

See docs/architecture.md for detailed design decisions.

## Development Commands

```bash
task build                # CGO_ENABLED=0 go build -o semver ./cmd/semver
task test                 # go test -count=2 -race ./...
task lint                 # golangci-lint run
task check-before-commit  # test + snapshot + lint
task snapshot             # goreleaser snapshot (local release test)
go run ./cmd/semver       # run without building
```

## Code Quality Standards

**Linters configured** (do not duplicate rules):
- golangci-lint: see `.golangci.yml` for the complete rule set (all linters enabled minus 14 explicitly disabled)
- pre-commit: see `.pre-commit-config.yaml` (install hooks with `task dev:install-pre-commit`)

**Key conventions:**
- Black-box tests only: every `_test.go` file uses an external `_test` package and table-driven cases.
- No dynamic errors: declare sentinels in the package's `errors.go`, wrap with `%w`.

## File Locations

- **Source**: `cmd/semver/` (entrypoint), `internal/cli/`, `internal/version/`, `internal/constraint/`
- **Tests**: `internal/**/*_test.go` (black-box unit), `test/integration/` (builds the real binary, drives it via os/exec)
- **Docs**: `docs/`, `specs/001-semver-cli/` (Spec-Kit artifacts)
- **Config**: `Taskfile.yml`, `mise.toml`, `.goreleaser.yaml`, `.github/workflows/`

## Documentation

- docs/architecture.md: System design and component overview
- docs/workflows.md: Development processes and git workflow
- docs/patterns.md: Code patterns and best practices

<!-- SPECKIT START -->
For additional context about technologies to be used, project structure,
shell commands, and other important information, read the current plan:
specs/001-semver-cli/plan.md

Active feature: 001-semver-cli — a single static Go 1.25 CLI for semver
manipulation (bump, pre-release lifecycle, compare, sort, validate, get,
satisfies). Layering: domain logic in internal/version and internal/constraint
(no CLI imports); thin CLI wrappers in internal/cli (stdlib flag, no third-party
runtime deps); entrypoint cmd/semver/main.go. stdout = data, stderr = humans;
exit codes 0/1/2/≥10 per specs/001-semver-cli/contracts/cli-contract.md.
<!-- SPECKIT END -->
