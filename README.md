# prompts

[![CI](https://github.com/faustbrian/go-prompts/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-prompts/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-prompts/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-prompts.svg)](https://pkg.go.dev/github.com/faustbrian/go-prompts)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-prompts?sort=semver)](https://github.com/faustbrian/go-prompts/releases)
[![Go](https://img.shields.io/badge/go-1.27.0-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

`prompts` is an explicit, typed foundation for interactive terminal
prompts and deterministic non-interactive fallbacks. Callers own the context,
streams, terminal capabilities, interaction authority, and lifecycle.

The published v1.1.0 release provides a stable v1 API. Its automated prompt,
rendering, terminal, headless, security, test, and release contracts are
implemented and remain subject to compatibility review. The recorded manual
terminal and assistive-technology matrix is available in
[Accessibility review evidence](docs/accessibility-review.md).

## Installation

```sh
go get github.com/faustbrian/go-prompts@v1.1.0
```

The target-oriented terminal adapter described below is available in v1.1.0.
Existing consumers can continue to import the compatibility path at
`github.com/faustbrian/go-prompts/terminal`.

Explicit non-interactive input uses the same typed parser and validation
pipeline without authorizing a terminal read:

```go
count, err := prompts.NewInteger(prompts.IntegerConfig{
    ID: "count", Label: "Count",
})
if err != nil { /* invalid definition */ }

value, err := prompts.Parse(ctx, count, configuredCount, dependencies)
```

`Parse` supports text, multiline text, confirmation, signed integers, exact
decimals, Go durations, ISO dates, wall-clock times, and paths. Path parsing
does not inspect or mutate the filesystem.

Secret definitions require an explicit classification. Their values redact
formatting, serialization, and structured logging by default; access requires
an explicit `Reveal` call. See [Secret handling](docs/secrets.md) for the
limits of string and byte cleanup.

Interactive execution requires caller authority plus explicit `EventSource`,
`TerminalController`, and output resources. The core never decodes a process
stream or acquires global terminal state on its own. See
[Interactive input](docs/interactive-input.md).

The optional real-terminal implementation is
`github.com/faustbrian/go-prompts/adapters/terminal`. The original
`github.com/faustbrian/go-prompts/terminal` path remains as a deprecated,
behavior-identical compatibility facade.

Ordered heterogeneous prompts compose through `NewForm`, `AsField`, and
`When`; typed answers are recovered with `FormValue`. See
[Forms](docs/forms.md).

Progress, spinners, bounded status streams, messages, tables, and summaries
are caller-driven presentation values. See
[Progress and presentation](docs/progress-and-presentation.md).

## Documentation

- [Documentation index](docs/README.md)
- [API overview](docs/api.md)
- [Prompt types](docs/prompt-types.md)
- [Selection and search](docs/selection.md)
- [Validation](docs/validation.md)
- [Rendering and themes](docs/rendering.md)
- [Accessibility](docs/accessibility.md)
- [Accessibility review evidence](docs/accessibility-review.md)
- [Hardening evidence](docs/hardening-evidence.md)
- [Benchmarks](docs/benchmarks.md)
- [Security model](docs/security.md)
- [Compatibility](docs/compatibility.md)
- [Terminal adapter](docs/terminal-adapter.md)
- [Integrations](docs/integrations.md)
- [Migration guidance](docs/migrations.md)
- [Troubleshooting](docs/troubleshooting.md)
- [FAQ](docs/faq.md)
- [Release process](docs/release.md)

Executable examples in `example_test.go` cover explicit input, headless
fallback, virtual-terminal interaction, forms, progress, and tables.

### Package map

- [`github.com/faustbrian/go-prompts`](https://pkg.go.dev/github.com/faustbrian/go-prompts)
  owns typed prompt definitions, explicit execution, forms, semantic rendering,
  secrets, and caller-driven presentation.
- [`github.com/faustbrian/go-prompts/adapters/terminal`](https://pkg.go.dev/github.com/faustbrian/go-prompts/adapters/terminal)
  is the optional Linux and macOS adapter for caller-owned terminal files, raw
  mode, echo restoration, and cancellable input decoding.
- [`github.com/faustbrian/go-prompts/terminal`](https://pkg.go.dev/github.com/faustbrian/go-prompts/terminal)
  is the retained v1 compatibility facade. New integrations should use the
  target-oriented adapter path.

Use `prompts` when an application needs typed, bounded interaction or the same
validation pipeline for explicit non-interactive input. Do not use it as a
command framework, configuration loader, application lifecycle, or source of
ambient terminal authority.

## Development

Go 1.27.0 is the tested toolchain. All commands must run with
`GOWORK=off` so the module is verified independently of sibling checkouts.

```sh
make format
make format-check
make check
make fuzz
make benchmark
```

Formatting uses the exact Glippy development revision pinned in the module
Makefile. Do not run gofmt, goimports, gofumpt, or golines over this module;
their output is not compatible with the canonical Glippy layout.

## Ecosystem

For ecosystem-wide selection and ownership guidance, see the versioned
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.5.5/docs/ecosystem/README.md)
and its [Tooling family](https://github.com/faustbrian/go-library-tools/blob/v1.5.5/docs/ecosystem/design-language.md#package-families-and-selection).

Licensed under Apache-2.0.
