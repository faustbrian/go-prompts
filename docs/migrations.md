# Migrations

From Huh, keep Huh models out of application contracts. Define a `Prompt[T]`,
construct explicit execution resources, and move non-interactive values to
`Parse` or a declared fallback. Huh-specific translation belongs in an
optional adapter.

From Survey or PromptUI, replace label-based result matching with stable option
identity. Replace package-global stdio and templates with `Execution`, semantic
roles, and a local immutable theme.

From Laravel Prompts or Symfony Question Helper, preserve the expressive label,
hint, validation, selection, progress, and message concepts, but do not carry
over process-global IO or implicit interactivity. Go callers handle returned
errors and process exit themselves.

Migration is complete only when the same operation accepts explicit data in CI
or redirected execution without opening a terminal prompt.

## Terminal adapter import path

The successor path is available in v1.1.0 and later. Do not change imports
until that version is published and selected by the consuming module.

New integrations should replace:

```go
import "github.com/faustbrian/go-prompts/terminal"
```

with:

```go
import promptsterminal "github.com/faustbrian/go-prompts/adapters/terminal"
```

No constructor, configuration, method, error, lifecycle, or ownership change
is required. The original path remains a deprecated compatibility facade for
the longer of 180 days after successor publication and two later stable minor
releases containing both paths. Removal additionally requires owned-consumer
migration, external-consumer evidence, and an authorized `v2.0.0` or later
release.
