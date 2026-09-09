// Package terminal preserves the original terminal-adapter import path.
//
// Deprecated: use github.com/faustbrian/go-prompts/adapters/terminal. This
// package remains supported through the documented compatibility interval.
package terminal

import (
	"context"
	"os"
	"time"

	prompts "github.com/faustbrian/go-prompts"
	promptsterminal "github.com/faustbrian/go-prompts/adapters/terminal"
)

// Config bounds reads, cancellation polling, and byte decoding.
//
// Deprecated: use promptsterminal.Config.
type Config struct {
	Decoder      prompts.DecoderConfig
	ReadBuffer   int
	PollInterval time.Duration
}

// Adapter implements prompts.EventSource and prompts.TerminalController for
// explicit files. A single prompt execution owns an Adapter at a time.
//
// Deprecated: use promptsterminal.Adapter.
type Adapter struct {
	_        [0]func()
	delegate *promptsterminal.Adapter
}

// New constructs an inert adapter without reading or mutating either file.
//
// Deprecated: use promptsterminal.New.
func New(input, output *os.File, config Config) (*Adapter, error) {
	delegate, err := promptsterminal.New(input, output, promptsterminal.Config{
		Decoder:      config.Decoder,
		ReadBuffer:   config.ReadBuffer,
		PollInterval: config.PollInterval,
	})
	if err != nil {
		return nil, err
	}

	return &Adapter{delegate: delegate}, nil
}

func (adapter *Adapter) target() *promptsterminal.Adapter {
	if adapter.delegate != nil {
		return adapter.delegate
	}

	return &promptsterminal.Adapter{}
}

// Capabilities detects only the explicitly supplied files. It does not inspect
// environment variables or interaction policy.
//
// Deprecated: use promptsterminal.Adapter.Capabilities.
func (adapter *Adapter) Capabilities() prompts.Capabilities {
	return adapter.target().Capabilities()
}

// Acquire places the explicit input terminal into raw mode.
//
// Deprecated: use promptsterminal.Adapter.Acquire.
func (adapter *Adapter) Acquire(ctx context.Context) error {
	return adapter.target().Acquire(ctx)
}

// SetEcho changes only the echo flag of an acquired raw terminal.
//
// Deprecated: use promptsterminal.Adapter.SetEcho.
func (adapter *Adapter) SetEcho(enabled bool) error {
	return adapter.target().SetEcho(enabled)
}

// Release restores the state captured by Acquire. It is idempotent.
//
// Deprecated: use promptsterminal.Adapter.Release.
func (adapter *Adapter) Release() error {
	return adapter.target().Release()
}

// Next reads and decodes one semantic event. Short file deadlines keep context
// cancellation bounded without a hidden goroutine.
//
// Deprecated: use promptsterminal.Adapter.Next.
func (adapter *Adapter) Next(ctx context.Context) (prompts.InputEvent, error) {
	return adapter.target().Next(ctx)
}
