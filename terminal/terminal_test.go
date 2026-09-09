package terminal_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	prompts "github.com/faustbrian/go-prompts"
	"github.com/faustbrian/go-prompts/terminal"
)

func TestNewDelegatesToTerminalAdapter(t *testing.T) {
	t.Parallel()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error = %v", err)
	}
	defer reader.Close()
	defer writer.Close()

	adapter, err := terminal.New(
		reader,
		writer,
		terminal.Config{PollInterval: time.Millisecond},
	)
	if err != nil || adapter == nil {
		t.Fatalf("New() = %#v, %v", adapter, err)
	}
	if capabilities := adapter.Capabilities(); capabilities.Unicode != true {
		t.Fatalf("Capabilities() = %#v", capabilities)
	}
	var nilContext context.Context
	if err := adapter.Acquire(nilContext); !errors.Is(err, prompts.ErrInvalidDefinition) {
		t.Fatalf("Acquire(nil) error = %v", err)
	}
	if err := adapter.SetEcho(false); !errors.Is(err, prompts.ErrAdapter) {
		t.Fatalf("SetEcho(false) error = %v", err)
	}
	if err := adapter.Release(); err != nil {
		t.Fatalf("Release() error = %v", err)
	}
	if _, err := adapter.Next(nilContext); !errors.Is(err, prompts.ErrInvalidDefinition) {
		t.Fatalf("Next(nil) error = %v", err)
	}
	if _, err := writer.Write([]byte("x")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	event, err := adapter.Next(context.Background())
	if err != nil || event != prompts.RuneEvent('x') {
		t.Fatalf("Next() = %#v, %v", event, err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if _, err := adapter.Next(context.Background()); !errors.Is(err, io.EOF) {
		t.Fatalf("Next() EOF error = %v", err)
	}
	if _, err := terminal.New(nil, writer, terminal.Config{}); !errors.Is(err, prompts.ErrInvalidDefinition) {
		t.Fatalf("New(nil, writer) error = %v", err)
	}
}

func TestZeroValueAdapterPreservesLegacyTerminalErrors(t *testing.T) {
	t.Parallel()

	var adapter terminal.Adapter
	var nilContext context.Context
	if err := adapter.Release(); err != nil {
		t.Fatalf("Release() error = %v", err)
	}
	if err := adapter.SetEcho(false); !errors.Is(err, prompts.ErrAdapter) {
		t.Fatalf("SetEcho(false) error = %v", err)
	}
	if _, err := adapter.Next(nilContext); !errors.Is(err, prompts.ErrInvalidDefinition) {
		t.Fatalf("Next(nil) error = %v", err)
	}
}
