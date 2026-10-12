package prompts

import (
	"context"
	"errors"
	"io"
	"time"
)

// Timer is an owned, stoppable clock event.
type Timer interface {
	// C returns the owned event channel; Stop need not close it.
	C() <-chan time.Time
	// Stop reports whether an active event was stopped; an already buffered value
	// can remain and the channel need not close.
	Stop() bool
}

// Ticker is an owned, stoppable sequence of clock events.
type Ticker interface {
	// C returns the owned event channel; Stop need not close it.
	C() <-chan time.Time
	// Stop prevents future tick generation; already buffered values can remain
	// and the channel need not close.
	Stop()
}

// Clock is an explicit caller-owned timing seam for sessions and presentation.
type Clock interface {
	// Now returns implementation time; callers choose its authority and concurrency.
	Now() time.Time
	// NewTimer creates an owned timer for caller stop/cleanup.
	NewTimer(time.Duration) Timer
	// NewTicker creates an owned ticker for caller stop/cleanup.
	NewTicker(time.Duration) Ticker
}

// Execution contains every ambient resource a prompt is allowed to use.
type Execution struct {
	// Input is for explicit adapter composition; core Run consumes Events, not this
	// reader.
	Input io.Reader
	// Output is a borrowed synchronous writer; context checks cannot interrupt a
	// blocking Write.
	Output io.Writer
	// Error is supplied diagnostic output for caller composition; core Run does not
	// write it.
	Error io.Writer
	// Capabilities supplies initial terminal evidence, not interaction permission.
	Capabilities Capabilities
	// Policy supplies caller authority for interaction, defaults and unlimited retries.
	Policy InteractionPolicy
	// Clock is an explicit timing seam for caller composition; core prompts add no
	// timers.
	Clock Clock
	// Dependencies is passed through to callbacks without copying or container
	// ownership.
	Dependencies any
	// Events is borrowed; implementations must honor cancellation without a hidden core
	// worker.
	Events EventSource
	// Terminal is acquired per interactive prompt; ordinary exits attempt restore/release.
	// Cleanup must return promptly rather than block or panic.
	Terminal TerminalController
	// Renderer is synchronous and borrowed; nil selects the capability-driven owned
	// renderer.
	Renderer Renderer
	// Theme is immutable default-renderer configuration; zero resolves to DefaultTheme.
	Theme Theme
	// Limits bounds interactive text/paste; nonpositive values select defaults, not
	// explicit Parse budgets.
	Limits InputLimits
	// Keys is the immutable local remap; zero uses default semantic bindings.
	Keys KeyMap
}

// Run executes a typed prompt without consulting process-wide streams or
// environment state. Nil/canceled contexts fail before execution. Headless
// defaults need PermitDefaults, but fallback and explicit cancel/EOF resolutions
// have their own configured authority. Validation failures can return a partial
// typed value with an error; errors govern success. Arbitrary callbacks/writers
// are synchronous and cannot be forcibly canceled; their retained values and
// resource concurrency remain caller-owned. Ordinary acquired-session cleanup
// attempts echo restoration and release; cleanup errors discard a successful
// value. External event/renderer/terminal panics are not globally isolated.
func Run[T any](ctx context.Context, prompt Prompt[T], execution Execution) (T, error) {
	var zero T
	if ctx == nil {
		return zero, &Error{
			Kind:      ErrorInvalidDefinition,
			Operation: "execute prompt",
			PromptID:  prompt.definition.id,
			Cause:     ErrInvalidDefinition,
		}
	}
	if err := ctx.Err(); err != nil {
		kind := ErrorCanceled
		if err == context.DeadlineExceeded {
			kind = ErrorDeadlineExceeded
		}

		return zero, &Error{
			Kind:      kind,
			Operation: "execute prompt",
			PromptID:  prompt.definition.id,
			Cause:     err,
		}
	}

	interactive, err := interactionAllowed(execution.Policy, execution.Capabilities)
	if err != nil {
		return zero, promptFailure(prompt.definition.id, err)
	}
	if interactive {
		return runInteractive(ctx, prompt, execution)
	}

	switch prompt.definition.headless {
	case HeadlessForbidden:
		return zero, promptFailure(prompt.definition.id, ErrInteractionNotPermitted)
	case HeadlessUseDefault:
		if execution.Policy.PermitDefaults {
			if value, ok := prompt.definition.defaultValue.Get(); ok {
				return applyPipeline(
					ctx,
					prompt.definition,
					value,
					execution.Dependencies,
					true,
				)
			}
		}

		return zero, promptFailure(prompt.definition.id, ErrInteractionNotPermitted)
	case HeadlessUseFallback:
		if value, ok := prompt.definition.fallbackValue.Get(); ok {
			return applyPipeline(
				ctx,
				prompt.definition,
				value,
				execution.Dependencies,
				true,
			)
		}

		return zero, promptFailure(prompt.definition.id, ErrInteractionNotPermitted)
	default:
		return zero, promptFailure(prompt.definition.id, ErrUnsupported)
	}
}

// Catch callback panics without exposing their values. This is not a sandbox:
// callers can block or retain/reveal data, and a transformer replacing a mutable
// secret must manage the replaced wrapper/copies itself. Context checks occur
// after callbacks; no helper goroutine or callback deadline is introduced.
func applyPipeline[T any](
	ctx context.Context,
	definition definition[T],
	value T,
	dependencies any,
	clone bool,
) (result T, resultErr error) {
	defer func() {
		if recover() != nil {
			var zero T
			result = zero
			resultErr = &Error{
				Kind:      ErrorAdapter,
				Operation: "run prompt callback",
				PromptID:  definition.id,
				Cause:     ErrAdapter,
			}
		}
	}()
	if clone && definition.clone != nil {
		value = definition.clone(value)
	}

	validation := ValidationContext{Dependencies: dependencies}
	for _, validator := range definition.preValidate {
		if err := validator(ctx, value, validation); err != nil {
			return value, validationFailure(definition.id, err, definition.secret)
		}
		if err := ctx.Err(); err != nil {
			return value, contextFailure(definition.id, err)
		}
	}

	for _, transformer := range definition.transform {
		transformed, err := transformer(ctx, value, validation)
		if err != nil {
			return value, validationFailure(definition.id, err, definition.secret)
		}
		value = transformed
		if err := ctx.Err(); err != nil {
			return value, contextFailure(definition.id, err)
		}
	}

	for _, validator := range definition.postValidate {
		if err := validator(ctx, value, validation); err != nil {
			return value, validationFailure(definition.id, err, definition.secret)
		}
		if err := ctx.Err(); err != nil {
			return value, contextFailure(definition.id, err)
		}
	}

	return value, nil
}

func validationFailure(promptID string, cause error, secret SecretClass) error {
	if secret != SecretNone {
		cause = NewValidationIssue(
			"secret_validation",
			"Secret value was rejected",
			promptID,
		)
	}

	return &Error{
		Kind:      ErrorValidationExhausted,
		Operation: "validate prompt",
		PromptID:  promptID,
		Cause:     normalizeIssue(cause, promptID),
	}
}

func contextFailure(promptID string, cause error) error {
	kind := ErrorCanceled
	if errors.Is(cause, context.DeadlineExceeded) {
		kind = ErrorDeadlineExceeded
	}

	return &Error{Kind: kind, Operation: "execute prompt", PromptID: promptID, Cause: cause}
}

func interactionAllowed(policy InteractionPolicy, capabilities Capabilities) (bool, error) {
	switch policy.Mode {
	case InteractiveRequired:
		if !policy.PermitInteraction {
			return false, ErrInteractionNotPermitted
		}
		if !capabilities.InputTerminal || !capabilities.OutputTerminal {
			return false, ErrTerminalUnavailable
		}

		return true, nil
	case InteractivePreferred:
		return policy.PermitInteraction &&
			capabilities.InputTerminal &&
			capabilities.OutputTerminal, nil
	case NonInteractiveOnly:
		return false, nil
	case AutoDetect:
		if !policy.PermitInteraction {
			return false, nil
		}
		if policy.Auto.RequireInputTerminal && !capabilities.InputTerminal {
			return false, nil
		}
		if policy.Auto.RequireOutputTerminal && !capabilities.OutputTerminal {
			return false, nil
		}

		return true, nil
	default:
		return false, ErrUnsupported
	}
}

func promptFailure(promptID string, target error) error {
	kind := ErrorUnsupported
	switch {
	case errors.Is(target, ErrInteractionNotPermitted):
		kind = ErrorInteractionNotPermitted
	case errors.Is(target, ErrTerminalUnavailable):
		kind = ErrorTerminalUnavailable
	}

	return &Error{Kind: kind, Operation: "execute prompt", PromptID: promptID, Cause: target}
}
