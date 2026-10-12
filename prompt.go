package prompts

// Optional distinguishes an explicit zero value from an omitted value.
type Optional[T any] struct {
	value T
	set   bool
}

// Some marks a present value by ordinary assignment, not deep copying references.
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, set: true}
}

// Get returns the stored value by assignment and presence; referenced data can alias.
func (optional Optional[T]) Get() (T, bool) {
	return optional.value, optional.set
}

// Accessibility provides caller-localized linear rendering metadata.
type Accessibility struct {
	// Label is caller-localized raw display text, sanitized by owned rendering.
	Label string
	// Description is optional raw explanatory text with no universal byte ceiling.
	Description string
	// TextualHint adds localized linear navigation guidance independent of color.
	TextualHint string
}

// TextConfig defines an immutable single-line text prompt.
type TextConfig struct {
	// ID is the stable nonempty identity, independent of display text.
	ID string
	// Label is caller-localized raw display text, sanitized by owned rendering.
	Label string
	// Description is optional raw explanatory text with no universal byte ceiling.
	Description string
	// Placeholder is empty-input display guidance, not a submitted default.
	Placeholder string
	// Hint is optional semantic guidance, not validation authority.
	Hint string
	// Help is optional localized extended guidance.
	Help string
	// Default distinguishes a present typed default from omission; Some preserves zero
	// values.
	Default Optional[string]
	// Fallback distinguishes a present typed fallback from omission.
	Fallback Optional[string]
	// Headless selects the noninteractive path; zero forbids headless execution.
	Headless HeadlessBehavior
	// Accessibility overrides localized label/description and adds textual navigation
	// guidance.
	Accessibility Accessibility
	// PreValidate contains ordered validators for the parsed value; constructors copy
	// callback slices but not callback state.
	PreValidate []Validator[string]
	// Transform contains ordered synchronous transformations; context cooperation and
	// retained values remain caller-owned.
	Transform []Transformer[string]
	// PostValidate contains ordered validators for the transformed value; constructors
	// copy callback slices but not callback state.
	PostValidate []Validator[string]
	// Retry bounds invalid interactive submissions; zero selects three attempts.
	Retry RetryPolicy
	// Cancel selects explicit cancel resolution; zero returns ErrCanceled.
	Cancel CancelBehavior
	// EndOfInput selects EOF/Ctrl-D resolution; zero returns ErrEndOfInput.
	EndOfInput EOFBehavior
	// Secret selects display/error redaction, not string-memory erasure.
	Secret SecretClass
}

type definition[T any] struct {
	kind          PromptKind
	id            string
	label         string
	description   string
	placeholder   string
	hint          string
	help          string
	defaultValue  Optional[T]
	fallbackValue Optional[T]
	headless      HeadlessBehavior
	accessibility Accessibility
	preValidate   []Validator[T]
	transform     []Transformer[T]
	postValidate  []Validator[T]
	retry         RetryPolicy
	cancel        CancelBehavior
	endOfInput    EOFBehavior
	secret        SecretClass
	parse         func(string) (T, error)
	parseBytes    func([]byte) (T, error)
	clone         func(T) T
	destroy       func(T)
	selection     *selectionDetails
}

// Prompt is a reusable definition created by a constructor; its zero value
// is not executable. Callback slices are copied, but callback state and generic
// referenced values are borrowed. Concurrent use requires safe dependencies and
// separate execution resources. Typed defaults/fallbacks bypass string parsing
// and still run the validation pipeline; constructors do not validate their
// business meaning. Secret-byte definitions retain copied defaults until unreachable.
type Prompt[T any] struct {
	definition definition[T]
}

// NewText creates a single-line text prompt.
func NewText(config TextConfig) (Prompt[string], error) {
	return newStringPrompt(
		KindText,
		"define text prompt",
		stringDefinitionConfig{
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
			config.Secret,
		},
		parseSingleLine,
	)
}

// ID returns the stable prompt identity.
func (prompt Prompt[T]) ID() string {
	return prompt.definition.id
}

// Describe returns a value snapshot of the prompt's execution contract.
func (prompt Prompt[T]) Describe() Descriptor {
	definition := prompt.definition

	return Descriptor{
		Kind:          definition.kind,
		ID:            definition.id,
		Label:         definition.label,
		Description:   definition.description,
		Placeholder:   definition.placeholder,
		Hint:          definition.hint,
		Help:          definition.help,
		Retry:         definition.retry,
		Cancel:        definition.cancel,
		EndOfInput:    definition.endOfInput,
		Secret:        definition.secret,
		Headless:      definition.headless,
		Accessibility: definition.accessibility,
	}
}
