# Development Workflows

## Feature Development

1. Create a feature branch from `main` (Spec-Kit features live under
   `specs/NNN-name/` with spec.md, plan.md, tasks.md).
2. Implement changes with black-box tests alongside.
3. Run the full gate locally: `task check-before-commit` (test + snapshot + lint).
4. Submit a PR for review.
5. Merge after approval and green CI.

## Code Review Process

- All PRs require green automated checks before merge.
- CI on every push: `linter.yml` (golangci-lint), `snapshot.yml` (GoReleaser
  snapshot build, cross-compilation check).
- `coverage.yml` regenerates the coverage badge on pushes to `main`.

## Testing Strategy

- Unit tests: `internal/**/*_test.go`, black-box (`package foo_test`),
  table-driven; run with race detection (`go test -count=2 -race ./...`).
- Property tests: `internal/version/property_test.go`.
- Integration tests: `test/integration/` builds the real binary
  (`harness_test.go`) and drives it via os/exec, asserting stdout/stderr and
  exit codes against the CLI contract.
- Pre-commit hooks: install once with `task dev:install-pre-commit`.

## Release Process

- Automated via GitHub Actions + GoReleaser (`release.yml`).
- Triggered by pushing a Git tag; publishes binaries, Docker images, and
  package distributions.
- Every push exercises the release path without publishing via
  `snapshot.yml` (`task snapshot` locally).
