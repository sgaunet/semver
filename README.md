# semver

[![GitHub release](https://img.shields.io/github/release/sgaunet/semver.svg)](https://github.com/sgaunet/semver/releases/latest)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/semver/total)
[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/semver)](https://goreportcard.com/report/github.com/sgaunet/semver)
![Test Coverage](https://raw.githubusercontent.com/wiki/sgaunet/semver/coverage-badge.svg)
[![linter](https://github.com/sgaunet/semver/actions/workflows/linter.yml/badge.svg)](https://github.com/sgaunet/semver/actions/workflows/linter.yml)
[![coverage](https://github.com/sgaunet/semver/actions/workflows/coverage.yml/badge.svg)](https://github.com/sgaunet/semver/actions/workflows/coverage.yml)
[![Snapshot Build](https://github.com/sgaunet/semver/actions/workflows/snapshot.yml/badge.svg)](https://github.com/sgaunet/semver/actions/workflows/snapshot.yml)
[![Release Build](https://github.com/sgaunet/semver/actions/workflows/release.yml/badge.svg)](https://github.com/sgaunet/semver/actions/workflows/release.yml)
[![GoDoc](https://godoc.org/github.com/sgaunet/semver?status.svg)](https://godoc.org/github.com/sgaunet/semver)
[![License](https://img.shields.io/github/license/sgaunet/semver.svg)](LICENSE)

A single, statically linked CLI for manipulating [Semantic Versions](https://semver.org).
Bump versions, drive a pre-release lifecycle, compare, sort, validate, extract
components, and test versions against range constraints. No runtime dependencies.

- **stdout = data** (machine-parseable), **stderr = humans** (errors, progress)
- meaningful, documented **exit codes**
- `--output=text|json`, honors `NO_COLOR`, `--quiet`, `--verbose`
- composes in pipelines and shell conditionals

## Install

```sh
# From source (Go 1.25+)
go install github.com/sgaunet/semver/cmd/semver@latest

# Or build locally
git clone https://github.com/sgaunet/semver && cd semver
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o semver ./cmd/semver
```

Released binaries (cross-compiled with checksums and an SBOM) are attached to GitHub
releases.

## Usage

```text
semver <command> [flags] [args...]
```

### Bump a version

```sh
semver patch v1.0.0      # v1.0.1
semver minor v1.2.3      # v1.3.0
semver major v1.2.3      # v2.0.0
semver patch 1.0.0       # 1.0.1  (no v-prefix in -> none out)
```

A bump on a version that already carries a pre-release follows the node-semver
convention (e.g. `patch v1.2.0-rc.1` -> `v1.2.0`).

### Pre-release lifecycle

```sh
semver prerelease v1.0.0 --pre rc --bump minor   # v1.1.0-rc.1  (start)
semver prerelease v1.1.0-rc.1                     # v1.1.0-rc.2  (increment)
semver release    v1.1.0-rc.1                     # v1.1.0       (finalize)
```

### Compare

```sh
semver compare v1.0.0 v1.2.0     # prints "lower" (exit 10)
semver compare v2.0.0 v1.9.9     # prints "higher" (exit 11)
semver compare v1.0.0 v1.0.0     # prints "equal" (exit 0)
```

### Sort, validate, get, satisfies

```sh
semver sort v2.0.0 v1.0.0-rc.1 v1.0.0 v1.2.0     # ascending by precedence
printf 'v2.0.0\nv1.0.0\n' | semver sort           # reads stdin if no args
semver sort --desc v1.0.0 v2.0.0

semver validate v1.2.3-rc.1+build.7               # "valid" (exit 0)
semver validate 1.2 && echo ok || echo bad        # invalid (exit 10)

semver get major v2.5.7-rc.3                       # 2
semver get prerelease v2.5.7-rc.3                  # rc.3

semver satisfies v1.5.0 '^1.2.0'                   # true (exit 0)
semver satisfies v2.0.0 '>=1.2.0 <2.0.0'           # false (exit 10)
```

Constraint syntax supports comparators (`=`, `!=`, `>`, `>=`, `<`, `<=`), caret
(`^`), tilde (`~`), wildcards (`1.2.x`, `*`), hyphen ranges (`1.2.0 - 1.5.0`), AND
(space/comma), and OR (`||`), with standard pre-release-exclusion semantics — a
pre-release only satisfies a comparator that pins the same `major.minor.patch`.

### JSON output

```sh
semver patch v1.0.0 --output=json
# {"input":"v1.0.0","operation":"patch","result":"v1.0.1"}
```

## Exit codes

| Code | Meaning |
|------|---------|
| `0`  | success / `compare`: equal / `validate`: valid / `satisfies`: satisfied / `get`: present |
| `1`  | generic failure |
| `2`  | usage error (bad flags or args, malformed constraint, invalid version) |
| `10` | `compare`: lower · `satisfies`: no · `validate`: invalid · `get`: absent |
| `11` | `compare`: higher |

## Configuration precedence

`flags > environment variables (NO_COLOR) > defaults`. There is no config file in v1.

## Development

```sh
task check   # fmt, vet, lint, test
task test
task build
task release-snapshot
```

## License

See [LICENSE](LICENSE).
