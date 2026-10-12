package prompts

// InteractionMode defines whether execution may acquire interactive input.
type InteractionMode uint8

const (
	// InteractiveRequired requires permission and both input/output terminal
	// capabilities.
	InteractiveRequired InteractionMode = iota
	// InteractivePreferred interacts only with permission and both terminal
	// capabilities.
	InteractivePreferred
	// NonInteractiveOnly never acquires terminal input.
	NonInteractiveOnly
	// AutoDetect checks explicit AutoRules while still requiring interaction permission.
	AutoDetect
)

// AutoRules are caller-selected conditions for AutoDetect. Terminal detection
// is never sufficient without PermitInteraction.
type AutoRules struct {
	// RequireInputTerminal requires input-terminal evidence only in AutoDetect mode.
	RequireInputTerminal bool
	// RequireOutputTerminal requires output-terminal evidence only in AutoDetect mode.
	RequireOutputTerminal bool
}

// InteractionPolicy records the caller's explicit authority to interact.
type InteractionPolicy struct {
	// Mode selects interaction policy; zero requires both terminal capabilities and
	// permission.
	Mode InteractionMode
	// PermitInteraction is explicit caller authorization; detection alone never grants
	// it.
	PermitInteraction bool
	// PermitDefaults permits HeadlessUseDefault; explicit cancel/EOF defaults use their
	// own configured behavior.
	PermitDefaults bool
	// PermitUnlimitedRetries permits explicitly unlimited interactive retry definitions.
	PermitUnlimitedRetries bool
	// Auto selects capability prerequisites only for AutoDetect.
	Auto AutoRules
}

// HeadlessBehavior defines how a prompt behaves without interactive input.
type HeadlessBehavior uint8

const (
	// HeadlessForbidden refuses noninteractive execution.
	HeadlessForbidden HeadlessBehavior = iota
	// HeadlessUseDefault uses a present default only when PermitDefaults is true.
	HeadlessUseDefault
	// HeadlessUseFallback uses a present fallback without PermitDefaults.
	HeadlessUseFallback
)

// ColorProfile describes the color output explicitly supported by a terminal.
type ColorProfile uint8

const (
	// ColorNone disables ANSI styling.
	ColorNone ColorProfile = iota
	// ColorANSI16 selects sixteen-color output and downsampling.
	ColorANSI16
	// ColorANSI256 selects indexed 256-color output.
	ColorANSI256
	// ColorTrueColor permits explicit RGB ANSI colors.
	ColorTrueColor
)

// Capabilities are supplied by the caller or an explicit application adapter.
type Capabilities struct {
	// InputTerminal reports explicit input-terminal evidence, independent of authority.
	InputTerminal bool
	// OutputTerminal reports explicit output-terminal evidence, independent of
	// authority.
	OutputTerminal bool
	// Width is initial output width; zero disables owned-renderer wrapping.
	Width int
	// Height controls selection pagination, not rendered-output bytes.
	Height int
	// Color selects output profile; zero is plain.
	Color ColorProfile
	// CursorMovement is metadata; owned rendering is linear, not cursor-updating.
	CursorMovement bool
	// Hyperlinks permits owned OSC 8 output with ANSIRenderer.
	Hyperlinks bool
	// Animation permits caller-selected spinner frames, not a timer or scheduler.
	Animation bool
	// Unicode permits native Unicode; false uses visible ASCII escapes.
	Unicode bool
}
