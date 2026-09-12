# Prompt comparison benchmarks

This module is an engineering-only benchmark harness. It is not a public
package, supported application dependency, or independently releasable module.
It compares equivalent single-line prompt flows across `prompts`, Huh, Survey,
PromptUI, and Bubble Tea/Bubbles through controlled pseudo-terminals on Unix.

From this directory, run the latency and allocation comparison with:

```sh
BENCHTIME=1s COUNT=3 ./run-benchmarks.sh
```

Measure stripped minimum-import binaries with:

```sh
./measure-binaries.sh
```

The commands require the Go 1.27.0 toolchain. Results are observational across
machines; do not treat them as portable speed claims. See the root
[benchmark documentation](../../docs/benchmarks.md) for the compared boundary,
recorded results, and interpretation limits.
