package main

import prompts "github.com/faustbrian/go-prompts"

func main() {
	_, _ = prompts.NewText(prompts.TextConfig{ID: "name", Label: "Name"})
}
