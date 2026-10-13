# Source documentation audit

[Documentation map](index.md) · [Repository README](../README.md)
[Next: API](api.md)

## Acceptance, ownership and named risks

This Tier A pass owns truthful source comments and directly contradicted current
prose at main 92f8caa27875e9e4ef2bf4b08476963315afbfaf. Root prompt definitions,
explicit/headless/interactive execution, parsing and validation, secrets,
selection/search/forms, semantic rendering, progress/status/task lifecycles and
canonical/legacy terminal adapters are direct consumer boundaries. The nested
comparison module and ignored build/review helpers remain separate operational
surfaces. Runtime AST/tokens/literals/directives, dependencies, fixtures, frozen
API/authority/history and executable documentation artifacts must remain exact.
No runtime, desired policy, deployment, accessibility or performance repair.

## Reproducible before-write inventory

Git and source-docs inventory identify 96 tracked Go files (43 source, 53 test),
two modules and ten source directories. Initial structural findings are 529
exported comment gaps and six package comment gaps. Common grouped configuration
members need named contracts without splitting declarations or changing implicit
semicolons. Operational packages include root, canonical and legacy terminal,
five comparison commands and two explicit ignored script entrypoints; the nested
comparison package itself is test-only. All tests, fifteen declared fuzz targets,
benchmarks, Examples, internal helpers and platform variants are inventoried.
Linux/BSD echo code and Linux/macOS poll fallbacks differ from unsupported forms.
No generated/vendor Go or generator/embed directive is present. Retain valid
legacy deprecations and existing analyzer/security directives.

## Missing versus technical investigation

Missing comments concentrate in public options, grouped members, semantic enums,
interfaces and formatting/accessor methods. Trace actual shallow generic copying,
secret wrapper ownership versus caller/runtime copies, synchronous callback and
writer cancellation, terminal release failures, queued-event precedence, decoder
admission and invalid-input cleanup, trusted clock calls under locks, generation
checks and form navigation. Presentation counts do not bound arbitrary text or
allocations. A successful virtual Example is not physical terminal or assistive
technology evidence; existing dated manual/hardening records remain historical.

## Native boundaries and smallest sufficient checks

.golib.yaml configures apidiff verification/api/api-v0.txt, not full GoDoc;
preserve both tracked API inputs. The nested comparison pins published Prompts
v1.1.0 rather than current source; no repin or peer benchmark campaign. Native
verification/package.mk docs runs scripts/check-docs.sh: required-file checks,
accessibility-review helper build/tests and existing root/canonical/legacy
Examples. Execute that bounded native gate only within an explicitly granted
Go 1.27.2 session with all three owned caches. Inspect native operational GoDoc
renders and ignored source comments, and prove every tracked Go file preserves
AST/raw tokens/literals/directives. Use structural source/local link/offline
specification checks and pinned full/changed/staged secrets with ignore annotations
disabled. Ordinary complete-diff truth/usefulness review precedes an explicit-path
Intent commit. Broad tests, PTY/manual terminal, concurrency/race/leak, active fuzz,
benchmarks, comparison composition, packaging and release campaigns are not prose
qualification. Immediately remove exact owned resources; root owns publication.

## Source-derived resolution

Public options/enums/interface methods and reachable formatting/accessor methods
now state named meanings and applicable defaults. Grouped declarations remain
intact; no field splitting or implicit-semicolon changes were used to satisfy
the structural checker. Generic values/defaults remain shallow except explicit
slice/SecretBytes clone policies. Typed defaults bypass string parsing and still
run validation. Headless default authority differs from fallback and explicit
cancel/EOF resolution. Callback panic containment belongs to the validation/
provider/form boundaries, not all external renderer/event/terminal code.

Execution and cleanup are synchronous: blocked caller code cannot be forcibly
canceled. Acquired cleanup attempts echo restoration before release; failure
invalidates a successful value. Terminal release clears ownership before restore,
so restoration failure cannot be repaired by another no-op Release. Queued events
and EOF precede cancellation checks. Files stay borrowed and input must remain
exclusive across poll/read. Linux/macOS support does not turn BSD echo variants
into a broader platform qualification.

SecretBytes synchronizes its owned slice, not caller/runtime/Reveal copies;
mutex-bearing values must not be copied. Callback replacement/retention needs
callback-owned cleanup. Explicit Parse/ParseBytes input is not capped by
interactive editor budgets. Decoder split-marker retention avoids false paste
admission, but neither decoder nor rendering counts bound arbitrary source IO,
text bytes or peak allocations. Wide graphemes stay intact across a wrapping
target. ANSI color and OSC hyperlink permissions are independent.

Dynamic providers run outside locks, while borrowed clocks run under session/
progress locks and must be prompt/non-reentrant. Scheduling discards superseded
results without canceling caller work or coalescing same-generation calls.
Form-local dependency fallback applies to cross-field validation only. Draft
invalidation and byte-secret cleanup protect back navigation; generic references
remain caller-owned. Progress and task presentation never schedules work;
independent task/entry/drop snapshots are not an atomic aggregate state view.
Comparison helpers bind published dependencies, and archive rewriting consumes
trusted producer input rather than validating arbitrary archive provenance.

Current maintained form/secret/terminal prose and the contributor tool-version
recipe were reconciled with source and .golib.yaml. Dated accessibility/security/
hardening records, dependency comparisons, frozen release source and baselines
remain historical authorities; no new physical, performance or conformance
claim is made. No runtime change or unresolved rationale decision is required.

## Executed evidence and remaining boundaries

Go 1.27.2 gofmt and exact AST/raw-token/literal/implicit-semicolon/directive
identity passed across all 96 tracked Go files, including tagged/ignored and
nested sources. All 789 current comment groups parsed/rendered. After final
ordinary cleanup-wording review, the three affected files received fresh exact
identity/HTML proof and root GoDoc refresh; unchanged native evidence was reused.
The root full GoDoc and representative ownership/lifecycle contracts rendered
successfully. A final one-line ProgressPending heading was narrowed to its
initial-state meaning; current structural parsing and manual inspection cover
that prose-only delta, with unchanged native execution evidence reused.
Canonical/legacy terminal and five nested comparison command GoDocs rendered;
both ignored commands received source-only stdlib GoDoc inspection.
Representative root ownership/callback/secret/form/progress contracts were read.

Native `make -f verification/package.mk docs` passed required-file validation,
accessibility-review helper build/tests and six root Examples. Four terminal
Examples without Output annotations compiled but did not execute; they are not
real-terminal evidence. Structural findings fell from 535 to zero. Offline
specification checking found zero applicable decision records, not a conformance
result. Both API inputs, all unowned tracked content and eighteen Markdown fenced
artifacts remain exact. Relative link checking excludes executable/inline code.
Pinned Gitleaks 8.30.1 full baseline/current scans passed with ignore annotations
disabled and complete redaction; final changed/staged bytes are checked at commit.

Two disposable verification assertions were corrected without source repair:
a literal GoDoc substring was split by normal wrapping, and a naive Markdown
link regex misread a generic Go call inside an executable fence. Wrap-insensitive
rendering and code-excluding link checks passed. Every successful/failed Go run
owned and removed all three caches/helpers/renders, including native helper temp
output. Shared caches and the dirty original checkout were not modified.
No broad tests, active fuzz, race/leak, PTY/manual accessibility, peer/composition,
benchmark, archive/release or online authority campaign ran. Parent owns remote
main and exact-source required CI; this evidence supports the local source pass.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.
