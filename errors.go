package prompts

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// ErrorKind identifies a stable class of prompt failure.
type ErrorKind string

const (
	// ErrorInteractionNotPermitted classifies when caller authority refuses interaction
	// or the default path.
	ErrorInteractionNotPermitted ErrorKind = "interaction_not_permitted"
	// ErrorTerminalUnavailable classifies when required terminal resources or
	// capabilities are absent.
	ErrorTerminalUnavailable ErrorKind = "terminal_unavailable"
	// ErrorCanceled classifies when execution was canceled.
	ErrorCanceled ErrorKind = "canceled"
	// ErrorDeadlineExceeded classifies when the caller deadline ended.
	ErrorDeadlineExceeded ErrorKind = "deadline_exceeded"
	// ErrorEndOfInput classifies when input ended without configured resolution.
	ErrorEndOfInput ErrorKind = "end_of_input"
	// ErrorTerminalDetached classifies when terminal attachment was lost.
	ErrorTerminalDetached ErrorKind = "terminal_detached"
	// ErrorInvalidDefinition classifies when definition or configuration is invalid.
	ErrorInvalidDefinition ErrorKind = "invalid_definition"
	// ErrorUnsupported classifies when configuration or operation is unsupported.
	ErrorUnsupported ErrorKind = "unsupported_configuration"
	// ErrorValidationExhausted classifies when validation failed or retries ended.
	ErrorValidationExhausted ErrorKind = "validation_exhausted"
	// ErrorRenderer classifies when semantic rendering failed.
	ErrorRenderer ErrorKind = "renderer_failure"
	// ErrorReader classifies when input reading or decoding failed.
	ErrorReader ErrorKind = "reader_failure"
	// ErrorWriter classifies when output writing failed.
	ErrorWriter ErrorKind = "writer_failure"
	// ErrorTerminalControl classifies when configuration or restoration failed.
	ErrorTerminalControl ErrorKind = "terminal_control_failure"
	// ErrorAdapter classifies when an adapter or contained callback failed.
	ErrorAdapter ErrorKind = "adapter_failure"
)

type sentinelError ErrorKind

// Error returns the fixed machine-readable kind.
func (err sentinelError) Error() string {
	return string(err)
}

var (
	// ErrInteractionNotPermitted is the errors.Is sentinel for refused interaction or
	// default authority.
	ErrInteractionNotPermitted error = sentinelError(ErrorInteractionNotPermitted)
	// ErrTerminalUnavailable is the errors.Is sentinel for absent terminal
	// resources/capabilities.
	ErrTerminalUnavailable error = sentinelError(ErrorTerminalUnavailable)
	// ErrCanceled is the errors.Is sentinel for canceled execution.
	ErrCanceled error = sentinelError(ErrorCanceled)
	// ErrDeadlineExceeded is the errors.Is sentinel for an expired caller deadline.
	ErrDeadlineExceeded error = sentinelError(ErrorDeadlineExceeded)
	// ErrEndOfInput is the errors.Is sentinel for unresolved EOF.
	ErrEndOfInput error = sentinelError(ErrorEndOfInput)
	// ErrTerminalDetached is the errors.Is sentinel for lost terminal attachment.
	ErrTerminalDetached error = sentinelError(ErrorTerminalDetached)
	// ErrInvalidDefinition is the errors.Is sentinel for invalid
	// definitions/configuration.
	ErrInvalidDefinition error = sentinelError(ErrorInvalidDefinition)
	// ErrUnsupported is the errors.Is sentinel for unsupported operations/configuration.
	ErrUnsupported error = sentinelError(ErrorUnsupported)
	// ErrValidationExhausted is the errors.Is sentinel for rejected validation or
	// exhausted retries.
	ErrValidationExhausted error = sentinelError(ErrorValidationExhausted)
	// ErrRenderer is the errors.Is sentinel for renderer failures.
	ErrRenderer error = sentinelError(ErrorRenderer)
	// ErrReader is the errors.Is sentinel for input/event read failures.
	ErrReader error = sentinelError(ErrorReader)
	// ErrWriter is the errors.Is sentinel for output write failures.
	ErrWriter error = sentinelError(ErrorWriter)
	// ErrTerminalControl is the errors.Is sentinel for terminal
	// configuration/restoration failures.
	ErrTerminalControl error = sentinelError(ErrorTerminalControl)
	// ErrAdapter is the errors.Is sentinel for adapter or contained callback failures.
	ErrAdapter error = sentinelError(ErrorAdapter)
)

// Error describes a prompt failure without incorporating unsafe cause text.
type Error struct {
	// Kind is stable classification; manually supplied kinds must not carry secrets or
	// terminal controls.
	Kind ErrorKind
	// Operation is context sanitized during Error formatting, not a secret-redaction
	// boundary.
	Operation string
	// PromptID is stable prompt context sanitized by Error formatting.
	PromptID string
	// Cause is omitted from default Error text but remains accessible through Unwrap.
	Cause error
}

// Error formats kind and sanitized operation/identity, omitting Cause text.
func (err *Error) Error() string {
	if err == nil {
		return "prompt failure"
	}

	message := string(err.Kind)
	if err.Operation != "" {
		message = safeText(err.Operation) + ": " + message
	}
	if err.PromptID != "" {
		message += fmt.Sprintf(" (prompt %q)", safeText(err.PromptID))
	}

	return message
}

// Is matches stable error kind rather than cause text; nil receivers do not match.
func (err *Error) Is(target error) bool {
	if err == nil || target == nil {
		return false
	}

	if sentinel, ok := errors.AsType[sentinelError](target); ok {
		return err.Kind == ErrorKind(sentinel)
	}

	other, ok := errors.AsType[*Error](target)
	return ok && err.Kind == other.Kind
}

// Unwrap returns the original cause, which may carry sensitive caller/backend
// diagnostics.
func (err *Error) Unwrap() error {
	if err == nil {
		return nil
	}

	return err.Cause
}

func safeText(value string) string {
	return strings.Map(
		func(char rune) rune {
			if unicode.IsControl(char) || isBidiControl(char) {
				return '\ufffd'
			}

			return char
		},
		value,
	)
}
