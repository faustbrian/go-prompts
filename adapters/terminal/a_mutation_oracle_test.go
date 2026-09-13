package promptsterminal

import (
	"context"
	"os"
	"testing"
	"time"

	prompts "github.com/faustbrian/go-prompts"
)

func TestEarlyConfigurationAndReadOracles(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error = %v", err)
	}
	defer reader.Close()
	defer writer.Close()

	adapter, err := New(reader, writer, Config{})
	if err != nil ||
		adapter.readBuffer != defaultReadBuffer ||
		adapter.pollInterval != defaultPollInterval {
		t.Fatalf("default adapter = %#v, %v", adapter, err)
	}
	for count, want := range map[int]bool{-1: false, 0: false, 1: true} {
		if got := hasReadBytes(count); got != want {
			t.Fatalf("hasReadBytes(%d) = %v, want %v", count, got, want)
		}
	}

	reads := 0
	adapter.setDeadline = func(time.Time) error { return nil }
	adapter.read = func(buffer []byte) (int, error) {
		reads++
		if reads > 1 {
			t.Fatal("Next() read again instead of dequeuing the buffered event")
		}
		copy(buffer, "ab")

		return 2, nil
	}
	first, err := adapter.Next(context.Background())
	if err != nil || first != prompts.RuneEvent('a') {
		t.Fatalf("first Next() = %#v, %v", first, err)
	}
	second, err := adapter.Next(context.Background())
	if err != nil || second != prompts.RuneEvent('b') || reads != 1 {
		t.Fatalf("second Next() = %#v, %v after %d reads", second, err, reads)
	}
}
