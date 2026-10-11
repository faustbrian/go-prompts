# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected package's docs.
2. Run the affected package's deterministic tests as described below.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected stable decision entries, and complete the Specification Decisions
section of the pull request template. An unresolved interpretation or stale
source pin is release-blocking; peer behavior cannot silently select policy.

Selected full/release mutation gates must finish with zero surviving viable
mutants; mutation is not an ordinary local or pull-request prerequisite.

Do not fork shared verification workflows, add permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use Go 1.27.2 for verification; the module/API minimum remains Go 1.27.0.
An isolated checkout needs no sibling repositories or coordination goal files.
Begin with deterministic tests and one active target when its terminal-input
boundary changes, using a disposable cache for each bounded session:

```bash
(
  set -eu
  GOCACHE="$(mktemp -d)"
  export GOCACHE
  trap 'rm -rf "${GOCACHE}"' EXIT
  GOWORK=off go test ./...
  GOWORK=off go test . -run '^$' -fuzz '^FuzzDecoderBoundsArbitraryBytes$' -fuzztime=2s -timeout=2m
)
```

The deterministic command executes fuzz seeds/corpus, not active fuzzing. It
retains secret suppression, event/editor clearing, terminal restoration,
accessibility metadata and allocation assertions.
Install the verified public v2.1.1 binary archive for your platform from
[go-library-tools releases](https://github.com/faustbrian/go-library-tools/releases/tag/v2.1.1),
checking `checksums.txt` against `.golib.yaml` and the archive against that
checksum set. Add the extracted `golib` binary to `PATH`; do not build Tools
from source. CI verifies the same public release while checking out its pinned
reusable workflow source.

```bash
make inventory
make check
make ci
make specification-check
golib mutation --module . --package .
```

Run adjacent Make commands in one disposable-cache session as shown above.
`make check` and `make ci` run the bounded local contract; no services or
test tags are declared for the root or benchmark comparison module. CI also
validates repository/cohesion structure. Specification checking is offline
unless explicitly requested online. The root race runs once with a two-minute
timeout and checks concurrent progress coalescing, stable non-ANSI output,
100 task increments without new workers, and cancellation rendering.
Base-free checks run seeds/corpus; changed-source CI selects causal active
fuzz smoke. Optional mutation is limited to one package, one worker and a
twenty-minute deadline.

Select full qualification and release rehearsal explicitly:

```bash
make check-full
make -f verification/package.mk docs
make -f verification/package.mk reproducible
make -f verification/package.mk benchmark
make specification-check-online
golib api check
golib release dry-run --module .
```

Full checks retain all fifteen declared fuzz targets at 1000x/20m each,
twenty-repeat broad concurrency, benchmarks and all original module gates.
The unchanged root reproducibility script compares two deterministic HEAD Git
archives with `gzip -n`; this is packaging proof, not a runtime integrity
contract, so it remains full/release-selected rather than an ordinary check.
The full documentation gate retains compilation and tests of the owned
accessibility-review tool plus examples. Online specification checking fetches
authority/errata sources when applicable. Release rehearsal can run locally.
Report unavailable or failing commands; bounded checks alone do not establish
release qualification.

## Selected Weekly Assurance

The owner-level [selected-assurance caller](.github/workflows/assurance.yml) is
configured for main only, on Sundays at 03:53 UTC. It selects one of the
15 current declared root-module fuzz targets using the shared UTC-week
rotation, with 30 seconds of active exploration. The selected risks are:
terminal-control sanitization, deterministic rendering, decoder bounds,
secret suppression, editor/selection/form state models and progress lifecycle.
The root targets use in-process models and VirtualTerminal fixtures, not a
physical terminal. Real terminal restoration, accessibility review and the
nested benchmark-comparison module retain their relevant full/manual routes;
this lane is not their qualification.

Base-free ordinary checks execute deterministic seeds, not active exploration;
committed-range CI selects affected active smoke. Manual assurance requires an
exact anchored selector from the root fuzz operation in `.golib.yaml`, for
example `^FuzzDecoderBoundsArbitraryBytes$`.
Empty, unknown and malformed selections refuse. The shared job bounds setup
and execution to ten minutes, active workers and CPU to two, owned memory to
2 GiB, native output to 1 MiB and private-cache lifetime. Failure, cancellation,
exhausted budgets and cleanup refusal remain attributable failures. This lane
does not run mutation or replace Required, original full fuzz budgets,
explicit `make check-full`, or relevant release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
