package prompts

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
)

// EventKind identifies one bounded terminal input or capability event.
type EventKind uint8

const (
	// EventKey carries Key and optional insertion Rune.
	EventKey EventKind = iota
	// EventPaste carries Text or an owned byte wrapper.
	EventPaste
	// EventResize replaces interactive dimensions.
	EventResize
	// EventEOF resolves configured end-of-input behavior.
	EventEOF
	// EventDetached terminates with terminal-detachment classification.
	EventDetached
	// EventCapabilities replaces the entire active capability snapshot.
	EventCapabilities
)

// Key identifies a semantic key independently of terminal byte encoding.
type Key uint8

const (
	// KeyRune inserts Rune and cannot be remapped.
	KeyRune Key = iota
	// KeyEnter submits input unless rebound.
	KeyEnter
	// KeyEscape resolves cancellation by default.
	KeyEscape
	// KeyCtrlC maps to Escape in the default key map.
	KeyCtrlC
	// KeyCtrlD resolves EOF by default.
	KeyCtrlD
	// KeyTab advances selection focus or submits the current form field.
	KeyTab
	// KeyShiftTab moves focus backward or revisits the prior enabled form field.
	KeyShiftTab
	// KeyBackspace removes the preceding editor grapheme.
	KeyBackspace
	// KeyDelete removes the editor grapheme at the cursor.
	KeyDelete
	// KeyLeft moves the cursor one grapheme left.
	KeyLeft
	// KeyRight moves the cursor one grapheme right.
	KeyRight
	// KeyUp moves selection focus up; plain editors ignore it.
	KeyUp
	// KeyDown moves selection focus down; plain editors ignore it.
	KeyDown
	// KeyHome moves to the first editor position or enabled selection.
	KeyHome
	// KeyEnd moves to the last editor position or enabled selection.
	KeyEnd
	// KeyWordLeft moves across whitespace-delimited editor words.
	KeyWordLeft
	// KeyWordRight moves across whitespace-delimited editor words.
	KeyWordRight
	// KeyPageUp moves selection focus by a height-derived page.
	KeyPageUp
	// KeyPageDown moves selection focus by a height-derived page.
	KeyPageDown
	// KeyIgnored is an inert unbound shortcut.
	KeyIgnored
	// KeyNewline inserts LF only in multiline input unless rebound.
	KeyNewline
)

// InputEvent is a secret-redacting semantic input event.
type InputEvent struct {
	// Kind selects which event payload fields execution interprets.
	Kind EventKind
	// Key is the semantic key interpreted for EventKey, subject to execution remapping.
	Key Key
	// Rune carries KeyRune input; direct field access bypasses formatting redaction.
	Rune rune
	// Text is borrowed immutable paste text; clearing references cannot erase it.
	Text string
	// Bytes is a shared paste wrapper; transfer or explicit Destroy determines its
	// lifetime.
	Bytes *SecretBytes
	// Width, Height are caller resize dimensions; negative interactive dimensions are
	// rejected.
	Width, Height int
	// Capabilities is the full replacement snapshot for EventCapabilities.
	Capabilities Capabilities
}

// RuneEvent creates a text insertion event.
func RuneEvent(value rune) InputEvent {
	return InputEvent{Kind: EventKey, Key: KeyRune, Rune: value}
}

// KeyEvent creates a non-text key event.
func KeyEvent(key Key) InputEvent {
	return InputEvent{Kind: EventKey, Key: key}
}

// PasteEvent creates a bracketed or ordinary paste event.
func PasteEvent(value string) InputEvent {
	return InputEvent{Kind: EventPaste, Text: value}
}

// PasteBytesEvent creates an owned byte-oriented paste event. Call Destroy
// when the event is not transferred to prompt execution.
func PasteBytesEvent(value []byte) InputEvent {
	return InputEvent{Kind: EventPaste, Bytes: NewSecretBytes(value)}
}

// ResizeEvent creates a terminal dimension event.
func ResizeEvent(width, height int) InputEvent {
	return InputEvent{Kind: EventResize, Width: width, Height: height}
}

// CapabilityEvent reports a complete replacement capability snapshot.
func CapabilityEvent(capabilities Capabilities) InputEvent {
	return InputEvent{Kind: EventCapabilities, Capabilities: capabilities}
}

// Destroy overwrites an owned byte payload and releases text references.
func (event *InputEvent) Destroy() {
	if event == nil {
		return
	}
	event.Bytes.Destroy()
	event.Text = ""
	event.Rune = 0
	event.Capabilities = Capabilities{}
}

// String returns a redacted representation without revealing payloads.
func (InputEvent) String() string {
	return "[INPUT EVENT]"
}

// GoString returns a redacted Go-syntax representation without exposing payloads.
func (InputEvent) GoString() string {
	return "[INPUT EVENT]"
}

// Format prevents pasted or typed input from appearing through fmt.
func (InputEvent) Format(state fmt.State, _ rune) {
	_, _ = state.Write([]byte("[INPUT EVENT]"))
}

// MarshalText redacts event payloads.
func (InputEvent) MarshalText() ([]byte, error) {
	return []byte("[INPUT EVENT]"), nil
}

// MarshalJSON redacts event payloads.
func (InputEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal("[INPUT EVENT]")
}

// LogValue redacts event payloads in structured logging.
func (InputEvent) LogValue() slog.Value {
	return slog.StringValue("[INPUT EVENT]")
}

// EventSource supplies decoded, cancellable semantic terminal events. It must
// honor context cancellation; the core never starts a goroutine around it.
type EventSource interface {
	// Next returns one semantic event or context/read error; implementations must honor
	// cancellation and transfer byte-payload ownership.
	Next(context.Context) (InputEvent, error)
}

// TerminalController owns acquisition and restoration of explicit terminal state.
// Methods must return errors for failure rather than panic; the core does not
// isolate arbitrary terminal-control panics or bound blocking cleanup calls.
type TerminalController interface {
	// Acquire acquires explicit state; success requires a later Release.
	Acquire(context.Context) error
	// SetEcho configures owned terminal echo; implementations must not disclose secret
	// input.
	SetEcho(bool) error
	// Release attempts restoration; failures must remain visible.
	Release() error
}

// InputLimits bound retained input and individual paste events.
type InputLimits struct {
	// MaxPasteBytes caps an interactive paste; nonpositive values select 1 MiB.
	MaxPasteBytes int
	// MaxInputBytes caps retained interactive editor bytes; nonpositive values select 4
	// MiB.
	MaxInputBytes int
}
