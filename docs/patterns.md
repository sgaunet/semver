# Code Patterns & Best Practices

## Error Handling

Each package declares private sentinel errors in an `errors.go` var block
(satisfies the err113 "no dynamic errors" lint rule), then wraps them with
`fmt.Errorf("...: %w", ...)` to add positional context while staying
`errors.Is`-matchable.

```go
// From internal/version/errors.go
var (
    errEmptyVersion = errors.New("empty string")
    errFormat       = errors.New("expected major.minor.patch")
    // ...
)
```

Exit-code mapping happens once, centrally: handlers return errors and
`internal/cli` maps them to the fixed contract in `internal/cli/exit.go`.

## Testing Patterns

- Test file naming: `*_test.go` next to the code under test.
- Black-box only: external test packages (`version_test`, `cli_test`,
  `integration`) exercising exported APIs.
- Table-driven cases as the default shape; property tests in
  `internal/version/property_test.go`.
- Integration suite builds the actual binary (`test/integration/harness_test.go`)
  and asserts stdout, stderr, and exit codes.

## Go-Specific Patterns

- Immutable value semantics: `Version` operations (`bump`, parse) return new
  values via a private `clone()`; receivers are never mutated.
- Presentation fidelity via private state: `Version` keeps a private `prefix`
  field so `v1.2.3` round-trips with its prefix intact.
- Stdlib `flag.FlagSet` per subcommand — no third-party CLI framework.
- stdout is for data, stderr is for humans: always emit through the `streams`
  abstraction in `internal/cli/output.go`, never `fmt.Println` directly.

## Common Utilities

- `internal/cli/output.go`: `streams` — centralized text/JSON emission and
  diagnostics.
- `internal/cli/exit.go`: exit-code constants (0/1/2/10/11) — the single
  source of truth for process exit status.
