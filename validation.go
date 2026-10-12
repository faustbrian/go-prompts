package prompts

import (
	"context"
	"errors"
	"strings"
)

// ValidationContext supplies caller-owned dependencies to validation and
// transformation callbacks.
type ValidationContext struct {
	// Dependencies is the exact borrowed callback dependency value; callers own
	// synchronization.
	Dependencies any
}

// Validator checks a typed value synchronously. Callbacks must cooperate with
// context and avoid mutating retained caller data; pipeline panics become ErrAdapter.
// Secret callbacks must not independently log or retain revealed input.
type Validator[T any] func(context.Context, T, ValidationContext) error

// Transformer returns the next pipeline value synchronously. Returned aliases,
// replaced mutable secret wrappers and callback-owned copies remain its cleanup
// responsibility; pipeline context checks cannot stop arbitrary blocked code.
type Transformer[T any] func(context.Context, T, ValidationContext) (T, error)

// ValidationIssue is a terminal-safe field-addressable validation failure.
// It neither bounds message bytes nor identifies/redacts arbitrary secrets;
// callers supply non-sensitive localized text. Dedicated secret prompts apply
// their own generic failure replacement during execution.
type ValidationIssue struct {
	code    string
	message string
	fields  []string
}

// NewValidationIssue creates a validation failure with defensively copied,
// terminal-safe content.
func NewValidationIssue(code string, message string, fields ...string) *ValidationIssue {
	safeFields := make([]string, len(fields))
	for index, field := range fields {
		safeFields[index] = safeText(field)
	}

	return &ValidationIssue{
		code:    safeText(code),
		message: safeText(message),
		fields:  safeFields,
	}
}

// Error returns message, then code, then generic text; nil is safe.
func (issue *ValidationIssue) Error() string {
	if issue == nil {
		return "validation failed"
	}
	if issue.message != "" {
		return issue.message
	}
	if issue.code != "" {
		return issue.code
	}

	return "validation failed"
}

// Code returns the stable machine-readable issue code.
func (issue *ValidationIssue) Code() string {
	if issue == nil {
		return ""
	}

	return issue.code
}

// Message returns safe caller-facing text.
func (issue *ValidationIssue) Message() string {
	if issue == nil {
		return ""
	}

	return issue.message
}

// Fields returns a defensive copy of relevant prompt identities.
func (issue *ValidationIssue) Fields() []string {
	if issue == nil {
		return nil
	}

	return append([]string(nil), issue.fields...)
}

func normalizeIssue(err error, promptID string) error {
	if issue, ok := errors.AsType[*ValidationIssue](err); ok {
		return issue
	}

	message := strings.TrimSpace(safeText(err.Error()))
	if message == "" {
		message = "validation failed"
	}

	return NewValidationIssue("validation", message, promptID)
}
