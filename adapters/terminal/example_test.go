package promptsterminal_test

import (
	"os"

	prompts "github.com/faustbrian/go-prompts"
	terminal "github.com/faustbrian/go-prompts/adapters/terminal"
)

func ExampleNew() {
	adapter, err := terminal.New(os.Stdin, os.Stdout, terminal.Config{})
	if err != nil {
		return
	}
	_ = adapter.Capabilities()
}

func ExampleNew_secretBytes() {
	adapter, err := terminal.New(
		os.Stdin,
		os.Stdout,
		terminal.Config{Decoder: prompts.DecoderConfig{ByteInput: true}},
	)
	if err != nil {
		return
	}
	_ = adapter.Capabilities()
}
