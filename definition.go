package prompts

import "fmt"

// PromptKind identifies the semantic prompt type independently of rendering.
type PromptKind string

const (
	// KindText selects single-line string parsing; CR/LF input is rejected.
	KindText PromptKind = "text"
	// KindMultiline selects string input with interactive newline insertion.
	KindMultiline PromptKind = "multiline"
	// KindInteger selects signed base-10 int64 parsing.
	KindInteger PromptKind = "integer"
	// KindDecimal selects exact base-10 decimal parsing without binary floating point.
	KindDecimal PromptKind = "decimal"
	// KindDuration selects Go time.Duration parsing.
	KindDuration PromptKind = "duration"
	// KindDate selects calendar-date parsing without a timezone.
	KindDate PromptKind = "date"
	// KindTime selects wall-clock parsing without a date or timezone.
	KindTime PromptKind = "time"
	// KindPath selects nonempty NUL-free path text without filesystem access.
	KindPath PromptKind = "path"
	// KindConfirm selects localized normalized true/false spellings.
	KindConfirm PromptKind = "confirm"
	// KindSelect selects enabled stable option identity.
	KindSelect PromptKind = "select"
	// KindMultiSelect selects unique enabled identities returned in declaration order.
	KindMultiSelect PromptKind = "multi_select"
	// KindSearchSelect adds bounded normalized option ranking to single selection.
	KindSearchSelect PromptKind = "search_select"
	// KindSecret selects classified redacting string storage.
	KindSecret PromptKind = "secret"
	// KindSecretBytes selects classified owned byte storage with explicit cleanup.
	KindSecretBytes PromptKind = "secret_bytes"
)

// RetryPolicy bounds invalid submissions. Unlimited retry requires separate
// authority in the interactive execution policy.
type RetryPolicy struct {
	// MaxAttempts counts invalid interactive submissions; zero selects three unless
	// Unlimited is true.
	MaxAttempts uint
	// Unlimited requires zero MaxAttempts and explicit interactive
	// PermitUnlimitedRetries authority.
	Unlimited bool
}

// CancelBehavior defines the result of an explicit cancel event.
type CancelBehavior uint8

const (
	// CancelReturnError returns ErrCanceled for explicit cancellation.
	CancelReturnError CancelBehavior = iota
	// CancelUseDefault runs the present typed default through validation on explicit
	// cancel.
	CancelUseDefault
	// CancelUseFallback runs the present fallback through validation on explicit cancel.
	CancelUseFallback
)

// EOFBehavior defines the result of end-of-input or Ctrl-D.
type EOFBehavior uint8

const (
	// EOFReturnError returns ErrEndOfInput for EOF/Ctrl-D.
	EOFReturnError EOFBehavior = iota
	// EOFUseDefault runs the present typed default through validation on EOF.
	EOFUseDefault
	// EOFUseFallback runs the present fallback through validation on EOF.
	EOFUseFallback
)

// SecretClass controls redaction and terminal echo policy.
type SecretClass uint8

const (
	// SecretNone is public input; dedicated secret constructors reject it.
	SecretNone SecretClass = iota
	// SecretPassword classifies password input for owned display/error suppression.
	SecretPassword
	// SecretToken classifies token input for owned display/error suppression.
	SecretToken
	// SecretOther classifies other sensitive input without an erasure guarantee.
	SecretOther
)

// Descriptor is an immutable value snapshot used by renderers and tests.
type Descriptor struct {
	// Kind identifies the semantic parser/editor family.
	Kind PromptKind
	// ID is the stable nonempty identity, not a display label or authentication proof.
	ID string
	// Label is caller-localized display text; owned rendering sanitizes terminal
	// controls.
	Label string
	// Description is optional caller-localized explanatory text, with no universal byte
	// cap.
	Description string
	// Placeholder is display-only empty-input guidance, not a submitted default.
	Placeholder string
	// Hint is optional display guidance, not validation or execution authority.
	Hint string
	// Help is optional caller-localized help rendered as semantic text.
	Help string
	// Retry bounds invalid submissions; zero normalizes to three attempts.
	Retry RetryPolicy
	// Cancel selects explicit cancel resolution; zero returns a cancellation error.
	Cancel CancelBehavior
	// EndOfInput selects EOF/Ctrl-D resolution; zero returns an end-of-input error.
	EndOfInput EOFBehavior
	// Secret selects display/error redaction; classification does not erase string
	// memory.
	Secret SecretClass
	// Headless selects the noninteractive path; zero forbids headless execution.
	Headless HeadlessBehavior
	// Accessibility overrides localized label/description and adds textual navigation
	// guidance.
	Accessibility Accessibility
}

func normalizeBehavior(
	retry RetryPolicy,
	cancel CancelBehavior,
	eof EOFBehavior,
	secret SecretClass,
) (RetryPolicy, error) {
	if retry.Unlimited && retry.MaxAttempts != 0 {
		return RetryPolicy{}, fmt.Errorf(
			"%w: retry cannot be both bounded and unlimited",
			ErrInvalidDefinition,
		)
	}
	if !retry.Unlimited && retry.MaxAttempts == 0 {
		retry.MaxAttempts = 3
	}
	if cancel > CancelUseFallback || eof > EOFUseFallback || secret > SecretOther {
		return RetryPolicy{}, fmt.Errorf(
			"%w: invalid execution behavior",
			ErrInvalidDefinition,
		)
	}

	return retry, nil
}

func invalidBehaviorDefinition(operation string, promptID string, cause error) error {
	return &Error{
		Kind:      ErrorInvalidDefinition,
		Operation: operation,
		PromptID:  promptID,
		Cause:     cause,
	}
}
