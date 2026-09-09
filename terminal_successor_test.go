//lint:file-ignore SA1019 Compatibility coverage intentionally exercises the deprecated facade.
package prompts_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	prompts "github.com/faustbrian/go-prompts"
	promptsterminal "github.com/faustbrian/go-prompts/adapters/terminal"
	legacyterminal "github.com/faustbrian/go-prompts/terminal" //nolint:staticcheck // Compatibility coverage requires the deprecated facade.
)

func TestTerminalCompatibilityFacadePreservesPublicContract(t *testing.T) {
	t.Parallel()

	legacyConfigType := reflect.TypeFor[legacyterminal.Config]() //nolint:staticcheck // Compatibility identity is the contract.
	successorConfigType := reflect.TypeFor[promptsterminal.Config]()
	if legacyConfigType.PkgPath() != "github.com/faustbrian/go-prompts/terminal" {
		t.Fatalf("legacy Config package = %q", legacyConfigType.PkgPath())
	}
	if successorConfigType.PkgPath() !=
		"github.com/faustbrian/go-prompts/adapters/terminal" {
		t.Fatalf("successor Config package = %q", successorConfigType.PkgPath())
	}
	legacyAdapterType := reflect.TypeFor[legacyterminal.Adapter]() //nolint:staticcheck // Compatibility identity is the contract.
	if legacyAdapterType.PkgPath() != "github.com/faustbrian/go-prompts/terminal" {
		t.Fatalf("legacy Adapter package = %q", legacyAdapterType.PkgPath())
	}
	if legacyAdapterType.Comparable() {
		t.Fatal("legacy Adapter became comparable")
	}
	successorAdapterType := reflect.TypeFor[promptsterminal.Adapter]()
	if successorAdapterType.PkgPath() !=
		"github.com/faustbrian/go-prompts/adapters/terminal" {
		t.Fatalf("successor Adapter package = %q", successorAdapterType.PkgPath())
	}

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error = %v", err)
	}
	defer reader.Close()
	defer writer.Close()

	legacyAdapter, err := legacyterminal.New(reader, writer, legacyterminal.Config{}) //nolint:staticcheck // Compatibility behavior is the contract.
	if err != nil {
		t.Fatalf("legacy New() error = %v", err)
	}
	successorAdapter, err := promptsterminal.New(reader, writer, promptsterminal.Config{})
	if err != nil {
		t.Fatalf("successor New() error = %v", err)
	}
	if legacyAdapter.Capabilities() != successorAdapter.Capabilities() { //nolint:staticcheck // Compatibility behavior is the contract.
		t.Fatal("legacy and successor capabilities differ")
	}

	if _, err := legacyterminal.New(nil, writer, legacyterminal.Config{}); !errors.Is(err, prompts.ErrInvalidDefinition) { //nolint:staticcheck // Compatibility behavior is the contract.
		t.Fatalf("legacy New(nil, writer) error = %v", err)
	}
	if _, err := promptsterminal.New(nil, writer, promptsterminal.Config{}); !errors.Is(err, prompts.ErrInvalidDefinition) {
		t.Fatalf("successor New(nil, writer) error = %v", err)
	}
}
