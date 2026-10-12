package prompts

import (
	"bytes"
	"context"
	"io"
	"sync"
)

// VirtualTerminal is a parallel-safe deterministic event source, terminal
// controller, and output capture for tests. It never touches a real terminal.
type VirtualTerminal struct {
	mu                              sync.RWMutex
	events                          chan InputEvent
	output                          bytes.Buffer
	width, height                   int
	acquired, released, echoEnabled bool
	closed                          bool
	acquireErr, echoErr, releaseErr error
}

// NewVirtualTerminal creates a fixed-capability terminal with a bounded event
// queue of 4096 events. Use the constructor; its zero value is not usable.
// Closing input does not drain/destroy queued byte wrappers, which retain their
// transferred lifetime until consumed or explicitly released by the owner.
func NewVirtualTerminal(width, height int) *VirtualTerminal {
	return &VirtualTerminal{
		events:      make(chan InputEvent, 4096),
		width:       width,
		height:      height,
		echoEnabled: true,
	}
}

// Push queues events without copying byte wrappers. Capacity failure can follow
// a successfully queued prefix; ownership of enqueued events has transferred.
// Caller destroys unqueued byte payloads. Concurrent Push calls can interleave.
func (terminal *VirtualTerminal) Push(events ...InputEvent) error {
	terminal.mu.RLock()
	defer terminal.mu.RUnlock()
	if terminal.closed {
		return streamFailure(
			"virtual-terminal",
			ErrorEndOfInput,
			"queue virtual event",
			ErrEndOfInput,
		)
	}
	for _, event := range events {
		select {
		case terminal.events <- event:
		default:
			return invalidBehaviorDefinition(
				"queue virtual event",
				"virtual-terminal",
				ErrInvalidDefinition,
			)
		}
	}
	return nil
}

// CloseInput makes an exhausted event sequence return EOF.
func (terminal *VirtualTerminal) CloseInput() {
	terminal.mu.Lock()
	defer terminal.mu.Unlock()
	if !terminal.closed {
		close(terminal.events)
		terminal.closed = true
	}
}

// Next returns the next event, EOF, or the context cause.
func (terminal *VirtualTerminal) Next(ctx context.Context) (InputEvent, error) {
	select {
	case <-ctx.Done():
		return InputEvent{}, ctx.Err()
	case event, ok := <-terminal.events:
		if !ok {
			return InputEvent{Kind: EventEOF}, io.EOF
		}
		if event.Kind == EventResize {
			terminal.mu.Lock()
			terminal.width, terminal.height = event.Width, event.Height
			terminal.mu.Unlock()
		}
		return event, nil
	}
}

// Acquire marks the virtual terminal acquired.
func (terminal *VirtualTerminal) Acquire(context.Context) error {
	terminal.mu.Lock()
	defer terminal.mu.Unlock()
	if terminal.acquireErr != nil {
		return terminal.acquireErr
	}
	terminal.acquired = true
	terminal.released = false
	return nil
}

// SetEcho records virtual echo state.
func (terminal *VirtualTerminal) SetEcho(enabled bool) error {
	terminal.mu.Lock()
	defer terminal.mu.Unlock()
	if terminal.echoErr != nil {
		return terminal.echoErr
	}
	terminal.echoEnabled = enabled
	return nil
}

// Release marks the virtual terminal released.
func (terminal *VirtualTerminal) Release() error {
	terminal.mu.Lock()
	defer terminal.mu.Unlock()
	if terminal.releaseErr != nil {
		return terminal.releaseErr
	}
	terminal.acquired = false
	terminal.released = true
	return nil
}

// Write captures output.
func (terminal *VirtualTerminal) Write(value []byte) (int, error) {
	terminal.mu.Lock()
	defer terminal.mu.Unlock()
	return terminal.output.Write(value)
}

// Output returns a copy of captured output as a string.
func (terminal *VirtualTerminal) Output() string {
	terminal.mu.RLock()
	defer terminal.mu.RUnlock()
	return terminal.output.String()
}

// Width returns latest virtual output width under the fixture lock.
func (terminal *VirtualTerminal) Width() int {
	terminal.mu.RLock()
	defer terminal.mu.RUnlock()
	return terminal.width
}

// Height returns latest virtual output height under the fixture lock.
func (terminal *VirtualTerminal) Height() int {
	terminal.mu.RLock()
	defer terminal.mu.RUnlock()
	return terminal.height
}

// Acquired reports latest successful virtual acquisition state.
func (terminal *VirtualTerminal) Acquired() bool {
	terminal.mu.RLock()
	defer terminal.mu.RUnlock()
	return terminal.acquired
}

// Released reports latest successful virtual release state.
func (terminal *VirtualTerminal) Released() bool {
	terminal.mu.RLock()
	defer terminal.mu.RUnlock()
	return terminal.released
}

// EchoEnabled reports latest virtual echo state.
func (terminal *VirtualTerminal) EchoEnabled() bool {
	terminal.mu.RLock()
	defer terminal.mu.RUnlock()
	return terminal.echoEnabled
}

// FailAcquire injects an acquisition failure; nil clears it.
func (terminal *VirtualTerminal) FailAcquire(err error) {
	terminal.mu.Lock()
	defer terminal.mu.Unlock()
	terminal.acquireErr = err
}

// FailEcho injects an echo failure; nil clears it.
func (terminal *VirtualTerminal) FailEcho(err error) {
	terminal.mu.Lock()
	defer terminal.mu.Unlock()
	terminal.echoErr = err
}

// FailRelease injects a release failure; nil clears it.
func (terminal *VirtualTerminal) FailRelease(err error) {
	terminal.mu.Lock()
	defer terminal.mu.Unlock()
	terminal.releaseErr = err
}

var _ EventSource = (*VirtualTerminal)(nil)

var _ TerminalController = (*VirtualTerminal)(nil)

var _ io.Writer = (*VirtualTerminal)(nil)
