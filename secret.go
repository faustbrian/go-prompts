package prompts

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
)

const redactedSecret = "[REDACTED]"

// SecretValue is a string-backed secret with redacted formatting. Reveal is
// explicit; Go strings cannot be reliably erased from memory.
type SecretValue struct {
	value string
}

// NewSecretValue copies a string reference into a redacting value wrapper.
func NewSecretValue(value string) SecretValue {
	return SecretValue{value: value}
}

// Reveal returns the underlying immutable Go string explicitly.
func (secret SecretValue) Reveal() string {
	return secret.value
}

// String returns a redacted representation without revealing payloads.
func (SecretValue) String() string {
	return redactedSecret
}

// GoString returns a redacted Go-syntax representation without exposing payloads.
func (SecretValue) GoString() string {
	return redactedSecret
}

// Format redacts every fmt formatting verb.
func (SecretValue) Format(state fmt.State, _ rune) {
	_, _ = state.Write([]byte(redactedSecret))
}

// MarshalText returns a redacted representation.
func (SecretValue) MarshalText() ([]byte, error) {
	return []byte(redactedSecret), nil
}

// MarshalJSON returns a redacted JSON string.
func (SecretValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(redactedSecret)
}

// LogValue returns a redacted structured logging value.
func (SecretValue) LogValue() slog.Value {
	return slog.StringValue(redactedSecret)
}

// SecretConfig defines a string-backed secret prompt.
type SecretConfig struct {
	// ID, Label, Description, Placeholder, Hint, Help identify and describe the prompt;
	// ID/Label are required, other text is optional and sanitized at owned presentation
	// boundaries.
	ID, Label, Description, Placeholder, Hint, Help string
	// Default, Fallback distinguish present typed values from omission; execution
	// follows the selected behavior and caller authority.
	Default, Fallback Optional[SecretValue]
	// Headless selects the noninteractive path; zero forbids headless execution.
	Headless HeadlessBehavior
	// Accessibility overrides localized label/description and adds textual navigation
	// guidance.
	Accessibility Accessibility
	// PreValidate, PostValidate are ordered callbacks before/after transformations;
	// slices are copied, callback state and context cooperation are not owned.
	PreValidate, PostValidate []Validator[SecretValue]
	// Transform contains ordered synchronous transformations; context cooperation and
	// retained values remain caller-owned.
	Transform []Transformer[SecretValue]
	// Retry bounds invalid interactive submissions; zero selects three attempts.
	Retry RetryPolicy
	// Cancel selects explicit cancel resolution; zero returns ErrCanceled.
	Cancel CancelBehavior
	// EndOfInput selects EOF/Ctrl-D resolution; zero returns ErrEndOfInput.
	EndOfInput EOFBehavior
	// Class must classify Password, Token or Other; secret constructors reject zero.
	Class SecretClass
}

// NewSecret creates a classified redacting string secret prompt.
func NewSecret(config SecretConfig) (Prompt[SecretValue], error) {
	if config.Class == SecretNone || config.Class > SecretOther {
		return Prompt[SecretValue]{}, invalidBehaviorDefinition(
			"define secret prompt",
			config.ID,
			fmt.Errorf("%w: secret classification is required", ErrInvalidDefinition),
		)
	}
	parser := func(input string) (SecretValue, error) {
		return NewSecretValue(input), nil
	}

	return newTypedPrompt(
		KindSecret,
		"define secret prompt",
		config.ID,
		config.Label,
		config.Description,
		config.Placeholder,
		config.Hint,
		config.Help,
		config.Default,
		config.Fallback,
		config.Headless,
		config.Accessibility,
		config.PreValidate,
		config.Transform,
		config.PostValidate,
		config.Retry,
		config.Cancel,
		config.EndOfInput,
		config.Class,
		parser,
	)
}

// SecretBytes synchronizes owned mutable bytes and is safe for concurrent
// Reveal/Destroy/Len calls; a zero wrapper is empty but not yet destroyed, and
// nil receivers are safe. Do not copy this mutex-bearing value. Destroy erases
// only its current owned slice, not prior Reveal/callback/runtime copies.
type SecretBytes struct {
	mu        sync.RWMutex
	value     []byte
	destroyed bool
}

// NewSecretBytes creates an owned copy of secret bytes.
func NewSecretBytes(value []byte) *SecretBytes {
	return &SecretBytes{value: append([]byte(nil), value...)}
}

// Reveal returns an owned copy of the current secret bytes.
func (secret *SecretBytes) Reveal() []byte {
	if secret == nil {
		return nil
	}
	secret.mu.RLock()
	defer secret.mu.RUnlock()

	return append([]byte(nil), secret.value...)
}

// Destroy overwrites and releases the wrapper's owned byte slice.
func (secret *SecretBytes) Destroy() {
	if secret == nil {
		return
	}
	secret.mu.Lock()
	defer secret.mu.Unlock()
	for index := range secret.value {
		secret.value[index] = 0
	}
	secret.value = nil
	secret.destroyed = true
}

// Len returns the current owned byte count.
func (secret *SecretBytes) Len() int {
	if secret == nil {
		return 0
	}
	secret.mu.RLock()
	defer secret.mu.RUnlock()

	return len(secret.value)
}

// Destroyed reports whether Destroy has been called.
func (secret *SecretBytes) Destroyed() bool {
	if secret == nil {
		return true
	}
	secret.mu.RLock()
	defer secret.mu.RUnlock()

	return secret.destroyed
}

// String returns a redacted representation without revealing payloads.
func (*SecretBytes) String() string {
	return redactedSecret
}

// GoString returns a redacted Go-syntax representation without exposing payloads.
func (*SecretBytes) GoString() string {
	return redactedSecret
}

// Format redacts every fmt formatting verb.
func (*SecretBytes) Format(state fmt.State, _ rune) {
	_, _ = state.Write([]byte(redactedSecret))
}

// MarshalText returns a redacted representation.
func (*SecretBytes) MarshalText() ([]byte, error) {
	return []byte(redactedSecret), nil
}

// MarshalJSON returns a redacted JSON string.
func (*SecretBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(redactedSecret)
}

// LogValue returns a redacted structured logging value.
func (*SecretBytes) LogValue() slog.Value {
	return slog.StringValue(redactedSecret)
}

func (secret *SecretBytes) clone() *SecretBytes {
	return NewSecretBytes(secret.Reveal())
}

// SecretBytesConfig defines a cleanup-capable byte secret prompt.
type SecretBytesConfig struct {
	// ID, Label, Description, Placeholder, Hint, Help identify and describe the prompt;
	// ID/Label are required, other text is optional and sanitized at owned presentation
	// boundaries.
	ID, Label, Description, Placeholder, Hint, Help string
	// Default, Fallback distinguish present typed values from omission; execution
	// follows the selected behavior and caller authority.
	Default, Fallback Optional[*SecretBytes]
	// Headless selects the noninteractive path; zero forbids headless execution.
	Headless HeadlessBehavior
	// Accessibility overrides localized label/description and adds textual navigation
	// guidance.
	Accessibility Accessibility
	// PreValidate, PostValidate are ordered callbacks before/after transformations;
	// slices are copied, callback state and context cooperation are not owned.
	PreValidate, PostValidate []Validator[*SecretBytes]
	// Transform contains ordered synchronous transformations; context cooperation and
	// retained values remain caller-owned.
	Transform []Transformer[*SecretBytes]
	// Retry bounds invalid interactive submissions; zero selects three attempts.
	Retry RetryPolicy
	// Cancel selects explicit cancel resolution; zero returns ErrCanceled.
	Cancel CancelBehavior
	// EndOfInput selects EOF/Ctrl-D resolution; zero returns ErrEndOfInput.
	EndOfInput EOFBehavior
	// Class must classify Password, Token or Other; secret constructors reject zero.
	Class SecretClass
}

// NewSecretBytesPrompt creates a classified cleanup-capable secret prompt.
func NewSecretBytesPrompt(config SecretBytesConfig) (Prompt[*SecretBytes], error) {
	if config.Class == SecretNone || config.Class > SecretOther {
		return Prompt[*SecretBytes]{}, invalidBehaviorDefinition(
			"define secret bytes prompt",
			config.ID,
			fmt.Errorf("%w: secret classification is required", ErrInvalidDefinition),
		)
	}
	defaultValue := cloneOptionalSecret(config.Default)
	fallbackValue := cloneOptionalSecret(config.Fallback)
	parser := func(input string) (*SecretBytes, error) {
		return NewSecretBytes([]byte(input)), nil
	}
	prompt, err := newTypedPrompt(
		KindSecretBytes,
		"define secret bytes prompt",
		config.ID,
		config.Label,
		config.Description,
		config.Placeholder,
		config.Hint,
		config.Help,
		defaultValue,
		fallbackValue,
		config.Headless,
		config.Accessibility,
		config.PreValidate,
		config.Transform,
		config.PostValidate,
		config.Retry,
		config.Cancel,
		config.EndOfInput,
		config.Class,
		parser,
	)
	if err != nil {
		return Prompt[*SecretBytes]{}, err
	}
	prompt.definition.clone = func(secret *SecretBytes) *SecretBytes {
		return secret.clone()
	}
	prompt.definition.parseBytes = func(input []byte) (*SecretBytes, error) {
		return NewSecretBytes(input), nil
	}
	prompt.definition.destroy = func(secret *SecretBytes) {
		secret.Destroy()
	}

	return prompt, nil
}

func cloneOptionalSecret(optional Optional[*SecretBytes]) Optional[*SecretBytes] {
	secret, present := optional.Get()
	if !present {
		return Optional[*SecretBytes]{}
	}

	return Some(secret.clone())
}

// ParseBytes copies explicit caller-owned bytes and validates them without
// converting the input into an immutable Go string. It does not apply interactive
// InputLimits; callers bound explicit bytes. Failure can return a destroyed
// wrapper alongside an error; transformer-retained or replaced copies remain
// callback-owned, not a complete erasure guarantee.
func ParseBytes(
	ctx context.Context,
	prompt Prompt[*SecretBytes],
	input []byte,
	dependencies any,
) (*SecretBytes, error) {
	if ctx == nil {
		return nil, invalidBehaviorDefinition(
			"parse secret bytes",
			prompt.ID(),
			ErrInvalidDefinition,
		)
	}
	if err := ctx.Err(); err != nil {
		return nil, contextFailure(prompt.ID(), err)
	}
	if prompt.definition.kind != KindSecretBytes {
		return nil, promptFailure(prompt.ID(), ErrUnsupported)
	}

	value, err := prompt.definition.parseBytes(input)
	if err != nil {
		return nil, validationFailure(prompt.ID(), err, prompt.definition.secret)
	}
	value, err = applyPipeline(ctx, prompt.definition, value, dependencies, false)
	if err != nil {
		prompt.definition.destroy(value)
	}

	return value, err
}
