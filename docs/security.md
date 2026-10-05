# Security model

**Model version:** 1.1

**Reviewed:** 2026-10-01

## Trust boundaries

The core consumes caller-supplied prompt definitions, input events, labels,
options, presentation text, and callback results. Callers own interaction
permission, contexts, streams, capabilities, callbacks, and revealed secrets.
The optional terminal adapter accesses only explicitly supplied files and owns
raw-mode acquisition, read polling, echo restoration, and terminal release.

## Defensive guarantees

The core has no ambient process access, package initialization, production
panic, or unowned goroutine. An executable architecture test enforces those
boundaries and rejects upstream prompt-engine dependencies in the core module.

Secret wrappers redact formatting, JSON, text marshaling, structured logging,
metadata, captured semantic output, and validation failures. Interactive
secret entry requests echo disablement and always attempts restoration and
terminal release. Cleanup failure is returned even when another operation also
failed.

Byte-mode terminal decoding clears consumed decoder and adapter buffers.
Byte-secret execution clears redacting event payloads, grapheme-editor cells,
validation intermediates, and failed result wrappers. It is opt-in because a
normal string paste event has already crossed an immutable string boundary.

Redaction is not access control. `Reveal` is explicit but returns sensitive
data to the caller. Go strings cannot be erased. `SecretBytes.Destroy` and
`FormResult.DestroySecrets` overwrite only owned byte slices; runtime copies,
caller copies, swap, crash dumps, terminals, and hardware remain outside the
guarantee.

Untrusted labels, options, messages, table cells, errors, and event formatting
neutralize terminal and bidi control characters. Interactive input and option
cardinality are bounded. Default error text omits wrapped causes and replaces
controls in operation and identity fields; explicit cause inspection retains
caller-owned diagnostics. Vulnerability and secret scans are local and CI
gates, but no scanner proves absence of a vulnerability.

Semantic hyperlinks accept only absolute HTTP, HTTPS, and mailto targets,
reject credentials and control or bidi characters, and emit owned OSC 8
controls only under explicit hyperlink capability. Plain fallback prints the
target textually.

## Accepted risks

| Risk | Owner | Rationale and mitigation | Review condition |
| --- | --- | --- | --- |
| Presentation text has no universal byte ceiling | Application owner | Callers own labels, cells, and messages; cap externally sourced text before constructing semantic frames. Control sanitization is not a size bound. | Revisit if the package takes ownership of serialized or remote presentation input. |
| Caller callbacks or event sources can ignore cancellation | Application owner | Arbitrary caller code cannot be forcibly stopped without changing execution ownership; use context-aware implementations and bounded downstream operations. | Revisit when callback or event-source execution ownership changes. |
| Secret strings and revealed copies cannot be reliably erased | Application owner | Go strings and caller copies are outside wrapper ownership; prefer byte secrets, destroy owned returned copies, and restrict dumps. | Revisit if runtime or API guarantees allow owned mutable storage throughout. |
| Terminal restoration can fail after device or process failure | Terminal adapter owner and operator | Cleanup failures are returned and release is attempted; supervise terminal use and reset affected terminals after failed cleanup. | Revisit when supported terminal-control mechanisms change. |

## Evidence scope

Deterministic error-formatting tests cover control replacement, benign Unicode,
default cause omission, and retained error classification. Existing secret,
decoder, interaction, and terminal lifecycle tests cover their owned behavior.
This model is not a claim that current release, scanner, consumer, or ecosystem
security verification is complete; those results require current execution
and delivery evidence.
