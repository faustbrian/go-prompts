# Mutation testing

Run the canonical content-addressed campaign with:

```sh
golib mutation --module .
```

This full/release command delegates to the canonical runner; mutation is not
an ordinary local or pull-request prerequisite. Optional
`golib mutation --module . --package .` selects only the root package with one
worker and a twenty-minute deadline. Use the disposable-cache session in
[contributor verification](../CONTRIBUTING.md#verification).
The full campaign discovers
production packages from the repository manifest, uses the pinned toolchain,
and requires exact 100.00% efficacy and mutant coverage. Any lived, timed-out,
uncovered, malformed, missing, or unclassified mutant remains a release
blocker; package-local exclusions and threshold overrides are forbidden.
