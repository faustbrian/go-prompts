// Package prompts provides explicit, typed terminal interactions without
// consulting process-wide streams, environment variables, or terminal state.
// Start with NewText or a typed constructor, Parse for explicit input, and Run
// for caller-authorized interaction. Definitions copy callback/option slices,
// not arbitrary generic values or callback state; callers own their concurrent
// reuse and dependencies. Interactive limits do not cap explicit Parse input or
// arbitrary presentation text. Callbacks, writers and terminal cleanup are
// synchronous and cooperative; no worker, timeout or physical terminal is owned.
// Secret redaction is not access control or complete memory erasure. The optional
// adapters/terminal package owns platform operations on borrowed files; benchmark
// comparison and manual review scripts are separate operational surfaces.
package prompts
