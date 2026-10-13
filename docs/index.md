# prompts documentation map

[Repository README](../README.md) · [Current Golib portal](https://github.com/faustbrian/go-library-tools/blob/main/docs/ecosystem/README.md) · [Published ecosystem snapshot](https://github.com/faustbrian/go-library-tools/blob/v1.8.4/docs/ecosystem/README.md)

This map links package-owned adoption and operational guidance. Start with
the repository README for purpose, installation, and ownership boundaries.
Published compatibility sets are frozen multi-package baselines; they do not
certify every current main revision or a deployed application.

## Module boundaries

Import only the modules needed by your application. A deprecated facade,
preview adapter, or internal harness is not interchangeable with a stable
root API; follow its own README and migration guidance before composing it.

| Module | Status | Owner detail |
| --- | --- | --- |
| `github.com/faustbrian/go-prompts` | stable | [README](../README.md) |
| `github.com/faustbrian/go-prompts/benchmarks/comparison` | internal; non-releasable | [README](../benchmarks/comparison/README.md) |

## Choose and adopt

- [Documentation](README.md).
- [Accessibility review evidence](accessibility-review.md).
- [Accessibility](accessibility.md).
- [API](api.md).
- [Architecture](architecture.md).
- [Compatibility](compatibility.md).
- [Interactive engine evaluation](dependency-evaluation.md).
- [FAQ](faq.md).
- [Forms](forms.md).
- [Integrations](integrations.md).
- [Interactive input](interactive-input.md).
- [Migrations](migrations.md).
- [Mutation testing](mutation.md).
- [Progress and presentation](progress-and-presentation.md).
- [Typed prompt definitions](prompt-types.md).
- [Semantic rendering](rendering.md).
- [Secret handling](secrets.md).
- [Selection and search](selection.md).
- [Terminal adapter](terminal-adapter.md).
- [Validation and transformation](validation.md).

## Operate and evaluate boundaries

- [Security model](security.md).
- [Troubleshooting](troubleshooting.md).

## Maintainer, specification, and evidence references

- [Prompt comparison benchmarks](../benchmarks/comparison/README.md).
- [Benchmarks](benchmarks.md).
- [Hardening evidence](hardening-evidence.md).
- [Release process](release.md).
- [Source documentation audit](source-documentation-audit.md).

## Related next steps

[Contributor setup and checks](../CONTRIBUTING.md), [Compatibility policy](../COMPATIBILITY.md), [Release history](../CHANGELOG.md), [Private vulnerability reporting](../SECURITY.md), [Support](../SUPPORT.md).
