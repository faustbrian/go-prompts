# Documentation

`prompts` provides typed, bounded terminal interaction and deterministic
non-interactive parsing. It does not own command dispatch, configuration
loading, application lifecycle, or ambient terminal access.

## Start here

- [Root package example](../example_test.go) shows explicit input, headless
  fallback, virtual-terminal interaction, forms, progress, and tables.
- [Terminal adapter example](../adapters/terminal/example_test.go) shows construction
  and capability detection; the [terminal adapter guide](terminal-adapter.md)
  covers caller-owned acquisition and release.
- [API overview](api.md) links the public API and package reference.
- [Integrations](integrations.md) explains when to use the package and how it
  composes with caller-owned command and configuration layers.

## Package map

- [`github.com/faustbrian/go-prompts`](https://pkg.go.dev/github.com/faustbrian/go-prompts)
  owns prompt definitions, explicit execution, forms, rendering, secrets, and
  caller-driven presentation.
- [`github.com/faustbrian/go-prompts/adapters/terminal`](https://pkg.go.dev/github.com/faustbrian/go-prompts/adapters/terminal)
  is the optional Linux and macOS terminal adapter. Callers own its files,
  context, acquisition authority, and release.
- [`github.com/faustbrian/go-prompts/terminal`](https://pkg.go.dev/github.com/faustbrian/go-prompts/terminal)
  is the deprecated v1 compatibility facade for the target-oriented adapter.

The internal `benchmarks/comparison` module is an engineering harness, not a
consumer package or independently released module.

## Contracts and operations

- [Architecture](architecture.md)
- [Dependency evaluation](dependency-evaluation.md)
- [Prompt types](prompt-types.md)
- [Interactive input](interactive-input.md)
- [Forms](forms.md)
- [Selection and search](selection.md)
- [Validation](validation.md)
- [Rendering and themes](rendering.md)
- [Progress and presentation](progress-and-presentation.md)
- [Terminal adapter](terminal-adapter.md)
- [Security model](security.md) and [secret handling](secrets.md)
- [Compatibility](compatibility.md) and [migration guidance](migrations.md)
- [Performance and benchmarks](benchmarks.md)
- [Accessibility](accessibility.md) and
  [review evidence](accessibility-review.md)

## Help, testing, and maintenance

- [Executable examples](../example_test.go) and
  [terminal adapter examples](../adapters/terminal/example_test.go)
- [Testing and hardening evidence](hardening-evidence.md)
- [Mutation evidence](mutation.md)
- [FAQ](faq.md)
- [Troubleshooting](troubleshooting.md)
- [Support](../SUPPORT.md)
- [Security reporting](../SECURITY.md)
- [Compatibility policy](../COMPATIBILITY.md) and
  [deprecation policy](../DEPRECATION.md)
- [Changelog](../CHANGELOG.md)
- [Release process](release.md)
- [License](../LICENSE)
- [Contributing](../CONTRIBUTING.md), [governance](../GOVERNANCE.md), and
  [code of conduct](../CODE_OF_CONDUCT.md)

For ecosystem-wide selection and ownership guidance, see the versioned
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.5.5/docs/ecosystem/README.md)
and its
[Tooling family](https://github.com/faustbrian/go-library-tools/blob/v1.5.5/docs/ecosystem/design-language.md#package-families-and-selection).
