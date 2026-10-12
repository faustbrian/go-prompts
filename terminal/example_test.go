package terminal_test

import (
	"os"

	prompts "github.com/faustbrian/go-prompts"
	"github.com/faustbrian/go-prompts/terminal"
)

// ExampleNew demonstrates inert construction on explicit files, without acquiring
// raw terminal mode or reading input.
func ExampleNew() {
	adapter, err := terminal.New(os.Stdin, os.Stdout, terminal.Config{})
	if err != nil {
		return
	}
	_ = adapter.Capabilities()
}

// ExampleNew_secretBytes demonstrates a dedicated byte-paste decoder; actual
// prompt execution and revealed-copy cleanup remain caller-owned.
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
