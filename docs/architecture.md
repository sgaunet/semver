# Architecture

## System Overview

A dependency-free Go CLI (`semver`) organized as strictly layered packages: pure
domain logic at the bottom, thin CLI wrappers above it, and a minimal entrypoint
on top. Imports flow in one direction only — domain packages never know the CLI
exists.

## Components

- `internal/version`: leaf domain package — parsing, comparison, bump and
  pre-release operations on immutable `Version` values. Zero internal imports.
- `internal/constraint`: range-constraint parsing and evaluation
  (`satisfies`). Imports only `internal/version`.
- `internal/cli`: subcommand handlers, flag parsing (stdlib `flag.FlagSet`),
  output formatting, exit-code mapping. Imports both domain packages.
- `cmd/semver`: entrypoint — signal handling and a single call into
  `internal/cli`.

## Design Decisions

1. **One-directional layering**: domain logic stays testable and reusable;
   the CLI is a replaceable shell around it.
2. **Stdlib-only runtime**: no third-party dependencies (no go.sum) keeps the
   binary small, static, and auditable.
3. **Registry-driven dispatch**: `internal/cli/app.go` holds a static
   `[]command` table mapping subcommand names to handlers — adding a command
   is one table entry plus one handler file.
4. **Fixed exit-code contract**: `internal/cli/exit.go` centralizes codes
   (0 success, 1 semantic false/failure, 2 usage error, 10/11 domain errors)
   per `specs/001-semver-cli/contracts/cli-contract.md`.
5. **streams abstraction**: `internal/cli/output.go` separates stdout (machine
   data, text or JSON) from stderr (human diagnostics) so output stays
   script-safe.

## Integration Points

- None at runtime: no database, network, or external services. The CLI reads
  argv/stdin and writes stdout/stderr.
- Release tooling: GoReleaser (`.goreleaser.yaml`) publishes binaries, Docker
  images, and package-manager artifacts on tag push.

## Data Flow

argv → `cmd/semver/main.go` → command registry (`internal/cli/app.go`) →
handler parses flags → domain call (`internal/version` / `internal/constraint`)
→ result emitted via `streams` (stdout) or diagnostic (stderr) → exit code.
