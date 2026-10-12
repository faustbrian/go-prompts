// Command go-prompts constructs a published v1.1.0 prompt for binary-size
// comparison, not the current root source or an interactive qualification.
package main

import prompts "github.com/faustbrian/go-prompts"

func main() {
	_, _ = prompts.NewText(prompts.TextConfig{ID: "name", Label: "Name"})
}
