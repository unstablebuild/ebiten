// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ebiten

import (
	"io/fs"
	"slices"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/internal/gamepad"
	"github.com/hajimehoshi/ebiten/v2/internal/gamepaddb"
	"github.com/hajimehoshi/ebiten/v2/internal/ui"
)

// KeyAction represents the action of a key event.
type KeyAction = ui.KeyAction

const (
	// KeyActionRelease indicates a key was released.
	KeyActionRelease KeyAction = ui.KeyRelease
	// KeyActionPress indicates a key was pressed.
	KeyActionPress KeyAction = ui.KeyPress
	// KeyActionRepeat indicates a key repeat from the OS.
	KeyActionRepeat KeyAction = ui.KeyRepeat
)

// KeyModifier is a bitmask of modifier keys held during a key event.
type KeyModifier = ui.KeyModifier

const (
	// KeyModShift indicates the Shift key was held.
	KeyModShift KeyModifier = ui.KeyModShift
	// KeyModControl indicates the Control key was held.
	KeyModControl KeyModifier = ui.KeyModControl
	// KeyModAlt indicates the Alt/Option key was held.
	KeyModAlt KeyModifier = ui.KeyModAlt
	// KeyModSuper indicates the Super/Meta/Command key was held.
	KeyModSuper KeyModifier = ui.KeyModSuper
)

// KeyEvent represents a discrete key event from the OS. Each KeyEvent
// carries a Key, an Action (Press/Release/Repeat), and a Mods bitmask.
type KeyEvent struct {
	Key    Key
	Action KeyAction
	Mods   KeyModifier
}

// InputEventKind distinguishes the kinds of input observation reported by
// the platform.
type InputEventKind = ui.InputEventKind

const (
	// InputEventKindKey is a key transition (press, release or OS repeat).
	InputEventKindKey InputEventKind = ui.InputEventKindKey
	// InputEventKindText is a committed Unicode code point.
	InputEventKindText InputEventKind = ui.InputEventKindText
)

// InputSource identifies the native key action that produced an input
// observation. Zero means the source is unknown.
type InputSource = ui.InputSource

// InputEvent is a single input observation: either a key transition or a
// committed Unicode code point.
type InputEvent struct {
	Kind InputEventKind

	// Key, Action and Mods describe an InputEventKindKey observation.
	Key    Key
	Action KeyAction
	Mods   KeyModifier

	// Rune is the committed code point of an InputEventKindText
	// observation; Mods carries the native modifier mask reported with it.
	Rune rune

	// NormalText is the platform's classification of the code point as
	// ordinary typed text rather than the side effect of a shortcut. It is
	// deliberately not derived from Mods, because text-producing layouts
	// such as AltGr report normal text while Ctrl+Alt is held.
	NormalText bool

	// Source ties a code point to the key action that produced it. All
	// observations of one physical action share a nonzero Source. Zero
	// means the platform could not establish causality, as for a standalone
	// IME commit.
	Source InputSource
}

// AppendInputChars appends "printable" runes, read from the keyboard at the time Update is called, to runes,
// and returns the extended buffer.
// Giving a slice that already has enough capacity works efficiently.
//
// AppendInputChars represents the environment's locale-dependent translation of keyboard
// input to Unicode characters. On the other hand, Key represents a physical key of US keyboard layout
//
// "Control" and modifier keys should be handled with IsKeyPressed.
//
// AppendInputChars is concurrent-safe.
//
// AppendInputChars is a projection of [AppendInputEvents] that keeps only
// the text observations and drops their ordering relative to key events,
// their modifiers and their source. Prefer [AppendInputEvents] when text
// must be attributed to the key action that produced it.
//
// On Android (ebitenmobile), EbitenView must be focusable to enable to handle keyboard keys.
func AppendInputChars(runes []rune) []rune {
	return theInputState.appendInputChars(runes)
}

// AppendKeyEvents appends discrete key events that occurred since the last Update
// to events, and returns the extended buffer.
// Giving a slice that already has enough capacity works efficiently.
//
// Each KeyEvent carries a Key, an Action (Press/Release/Repeat), and a Mods bitmask.
// Key repeat events come from the OS, not from polling — this provides precise
// press/release/repeat semantics without manual repeat tracking.
//
// AppendKeyEvents is concurrent-safe.
//
// AppendKeyEvents is a projection of [AppendInputEvents] that keeps only the
// key observations and drops their ordering relative to committed text.
// Prefer [AppendInputEvents] when a key action must be reconciled with the
// text it produced.
//
// On Android (ebitenmobile), EbitenView must be focusable to enable to handle keyboard keys.
func AppendKeyEvents(events []KeyEvent) []KeyEvent {
	return theInputState.appendKeyEvents(events)
}

// AppendInputEvents appends every input observation that occurred since the
// last Update to events, in the order the platform reported them, and
// returns the extended buffer.
// Giving a slice that already has enough capacity works efficiently.
//
// Key transitions and committed text are related observations, not disjoint
// streams. A single physical action reports one key transition and zero, one
// or many code points, and the platform may deliver those code points in a
// later update than the key transition. InputEvent.Source expresses the
// causality the platform actually knows: every observation of one key action
// shares a nonzero source, and zero means the source is unknown, as for a
// standalone IME commit.
//
// Consumers deciding whether a code point belongs to a key action they
// already handled must compare InputEvent.Source. Inferring the association
// from rune equality, from the modifier mask alone, or from membership in
// the same update is incorrect: it drops unrelated text, AltGr and
// international-layout input, dead-key and IME commits, and repeats.
//
// AppendInputEvents is concurrent-safe.
func AppendInputEvents(events []InputEvent) []InputEvent {
	return theInputState.appendInputEvents(events)
}

// InputChars return "printable" runes read from the keyboard at the time Update is called.
//
// Deprecated: as of v2.2. Use AppendInputChars instead.
func InputChars() []rune {
	return AppendInputChars(nil)
}

// IsKeyPressed returns a boolean indicating whether key is pressed.
//
// Deprecated: On GLFW desktop platforms, key state is no longer polled
// per frame. IsKeyPressed always returns false on those platforms.
// Use [AppendKeyEvents] instead, which provides discrete press/release/repeat
// events from the OS with correct ordering and modifier information.
//
// On other platforms (mobile, JS, Nintendo SDK), IsKeyPressed still works
// as before.
//
// IsKeyPressed is concurrent-safe.
func IsKeyPressed(key Key) bool {
	return theInputState.isKeyPressed(key)
}

// KeyName returns a key name for the current keyboard layout.
// For example, KeyName(KeyQ) returns 'q' for a QWERTY keyboard, and returns 'a' for an AZERTY keyboard.
//
// KeyName returns an empty string if 1) the key doesn't have a physical key name, 2) the platform doesn't support KeyName,
// or 3) the main loop doesn't start yet.
//
// KeyName is supported by desktops and browsers.
//
// KeyName is concurrent-safe.
func KeyName(key Key) string {
	return ui.Get().KeyName(ui.Key(key))
}

// CursorPosition returns a position of a mouse cursor relative to the game screen (window). The cursor position is
// 'logical' position and this considers the scale of the screen.
//
// CursorPosition returns (0, 0) before the main loop on desktops and browsers.
//
// CursorPosition always returns (0, 0) on mobiles.
//
// CursorPosition is concurrent-safe.
func CursorPosition() (x, y int) {
	cx, cy := theInputState.cursorPosition()
	return int(cx), int(cy)
}

// Wheel returns x and y offsets of the mouse wheel or touchpad scroll.
// It returns 0 if the wheel isn't being rolled.
//
// Wheel is concurrent-safe.
func Wheel() (xoff, yoff float64) {
	return theInputState.wheel()
}

// IsMouseButtonPressed returns a boolean indicating whether mouseButton is pressed.
//
// If you want to know whether the mouseButton started being pressed in the current tick,
// use inpututil.IsMouseButtonJustPressed
//
// IsMouseButtonPressed is concurrent-safe.
func IsMouseButtonPressed(mouseButton MouseButton) bool {
	return theInputState.isMouseButtonPressed(mouseButton)
}

// GamepadID represents a gamepad identifier.
type GamepadID = gamepad.ID

// GamepadSDLID returns a string with the GUID generated in the same way as SDL.
// To detect devices, see also the community project of gamepad devices database: https://github.com/gabomdq/SDL_GameControllerDB
//
// GamepadSDLID always returns an empty string on browsers and mobiles.
//
// GamepadSDLID is concurrent-safe.
func GamepadSDLID(id GamepadID) string {
	g := gamepad.Get(id)
	if g == nil {
		return ""
	}
	return g.SDLID()
}

// GamepadName returns a string with the name.
// This function may vary in how it returns descriptions for the same device across platforms.
// for example the following drivers/platforms see an Xbox One controller as the following:
//
//   - Windows: "Xbox Controller"
//   - Chrome: "Xbox 360 Controller (XInput STANDARD GAMEPAD)"
//   - Firefox: "xinput"
//
// GamepadName is concurrent-safe.
func GamepadName(id GamepadID) string {
	g := gamepad.Get(id)
	if g == nil {
		return ""
	}
	return g.Name()
}

// AppendGamepadIDs appends available gamepad IDs to gamepadIDs, and returns the extended buffer.
// Giving a slice that already has enough capacity works efficiently.
//
// AppendGamepadIDs is concurrent-safe.
func AppendGamepadIDs(gamepadIDs []GamepadID) []GamepadID {
	return gamepad.AppendGamepadIDs(gamepadIDs)
}

// GamepadIDs returns a slice indicating available gamepad IDs.
//
// Deprecated: as of v2.2. Use AppendGamepadIDs instead.
func GamepadIDs() []GamepadID {
	return AppendGamepadIDs(nil)
}

// GamepadAxisCount returns the number of axes of the gamepad (id).
//
// GamepadAxisCount is concurrent-safe.
func GamepadAxisCount(id GamepadID) int {
	g := gamepad.Get(id)
	if g == nil {
		return 0
	}
	return g.AxisCount()
}

// GamepadAxisNum returns the number of axes of the gamepad (id).
//
// Deprecated: as of v2.4. Use GamepadAxisCount instead.
func GamepadAxisNum(id GamepadID) int {
	return GamepadAxisCount(id)
}

// GamepadAxisValue returns a float value [-1.0 - 1.0] of the given gamepad (id)'s axis (axis).
//
// GamepadAxisValue is concurrent-safe.
func GamepadAxisValue(id GamepadID, axis GamepadAxisType) float64 {
	g := gamepad.Get(id)
	if g == nil {
		return 0
	}
	return g.Axis(int(axis))
}

// GamepadAxis returns a float value [-1.0 - 1.0] of the given gamepad (id)'s axis (axis).
//
// Deprecated: as of v2.2. Use GamepadAxisValue instead.
func GamepadAxis(id GamepadID, axis GamepadAxisType) float64 {
	return GamepadAxisValue(id, axis)
}

// GamepadButtonCount returns the number of the buttons of the given gamepad (id).
//
// GamepadButtonCount is concurrent-safe.
func GamepadButtonCount(id GamepadID) int {
	g := gamepad.Get(id)
	if g == nil {
		return 0
	}

	// For backward compatibility, hats are treated as buttons in GLFW.
	return g.ButtonCount() + g.HatCount()*4
}

// GamepadButtonNum returns the number of the buttons of the given gamepad (id).
//
// Deprecated: as of v2.4. Use GamepadButtonCount instead.
func GamepadButtonNum(id GamepadID) int {
	return GamepadButtonCount(id)
}

// IsGamepadButtonPressed reports whether the given button of the gamepad (id) is pressed or not.
//
// If you want to know whether the given button of gamepad (id) started being pressed in the current tick,
// use inpututil.IsGamepadButtonJustPressed
//
// IsGamepadButtonPressed is concurrent-safe.
//
// The relationships between physical buttons and button IDs depend on environments.
// There can be differences even between Chrome and Firefox.
func IsGamepadButtonPressed(id GamepadID, button GamepadButton) bool {
	g := gamepad.Get(id)
	if g == nil {
		return false
	}

	nbuttons := g.ButtonCount()
	if int(button) < nbuttons {
		return g.Button(int(button))
	}

	// For backward compatibility, hats are treated as buttons in GLFW.
	if hat := (int(button) - nbuttons) / 4; hat < g.HatCount() {
		dir := (int(button) - nbuttons) % 4
		return g.Hat(hat)&(1<<dir) != 0
	}

	return false
}

// StandardGamepadAxisValue returns a float value [-1.0 - 1.0] of the given gamepad (id)'s standard axis (axis).
//
// StandardGamepadAxisValue returns 0 when the gamepad doesn't have a standard gamepad layout mapping.
//
// StandardGamepadAxisValue is concurrent safe.
func StandardGamepadAxisValue(id GamepadID, axis StandardGamepadAxis) float64 {
	g := gamepad.Get(id)
	if g == nil {
		return 0
	}
	return g.StandardAxisValue(axis)
}

// StandardGamepadButtonValue returns a float value [0.0 - 1.0] of the given gamepad (id)'s standard button (button).
//
// StandardGamepadButtonValue returns 0 when the gamepad doesn't have a standard gamepad layout mapping.
//
// StandardGamepadButtonValue is concurrent safe.
func StandardGamepadButtonValue(id GamepadID, button StandardGamepadButton) float64 {
	g := gamepad.Get(id)
	if g == nil {
		return 0
	}
	return g.StandardButtonValue(button)
}

// IsStandardGamepadButtonPressed reports whether the given gamepad (id)'s standard gamepad button (button) is pressed.
//
// IsStandardGamepadButtonPressed returns false when the gamepad doesn't have a standard gamepad layout mapping.
//
// IsStandardGamepadButtonPressed is concurrent safe.
func IsStandardGamepadButtonPressed(id GamepadID, button StandardGamepadButton) bool {
	g := gamepad.Get(id)
	if g == nil {
		return false
	}
	return g.IsStandardButtonPressed(button)
}

// IsStandardGamepadLayoutAvailable reports whether the gamepad (id) has a standard gamepad layout mapping.
//
// IsStandardGamepadLayoutAvailable is concurrent-safe.
func IsStandardGamepadLayoutAvailable(id GamepadID) bool {
	g := gamepad.Get(id)
	if g == nil {
		return false
	}
	return g.IsStandardLayoutAvailable()
}

// IsStandardGamepadAxisAvailable reports whether the standard gamepad axis is available on the gamepad (id).
//
// IsStandardGamepadAxisAvailable is concurrent-safe.
func IsStandardGamepadAxisAvailable(id GamepadID, axis StandardGamepadAxis) bool {
	g := gamepad.Get(id)
	if g == nil {
		return false
	}
	return g.IsStandardAxisAvailable(axis)
}

// IsStandardGamepadButtonAvailable reports whether the standard gamepad button is available on the gamepad (id).
//
// IsStandardGamepadButtonAvailable is concurrent-safe.
func IsStandardGamepadButtonAvailable(id GamepadID, button StandardGamepadButton) bool {
	g := gamepad.Get(id)
	if g == nil {
		return false
	}
	return g.IsStandardButtonAvailable(button)
}

// UpdateStandardGamepadLayoutMappings parses the specified string mappings in SDL_GameControllerDB format and
// updates the gamepad layout definitions.
//
// UpdateStandardGamepadLayoutMappings reports whether the mappings were applied,
// and returns an error in case any occurred while parsing the mappings.
//
// One or more input definitions can be provided separated by newlines.
// In particular, it is valid to pass an entire gamecontrollerdb.txt file.
// Note though that Ebitengine already includes its own copy of this file,
// so this call should only be necessary to add mappings for hardware not supported yet;
// ideally games using the StandardGamepad* functions should allow the user to provide mappings and
// then call this function if provided.
// When using this facility to support new hardware, please also send a pull request to
// https://github.com/gabomdq/SDL_GameControllerDB to make your mapping available to everyone else.
//
// A platform field in a line corresponds with a GOOS like the following:
//
//	"Windows":  GOOS=windows
//	"Mac OS X": GOOS=darwin (not ios)
//	"Linux":    GOOS=linux (not android)
//	"Android":  GOOS=android
//	"iOS":      GOOS=ios
//	"":         Any GOOS
//
// On platforms where gamepad mappings are not managed by Ebitengine, this always returns false and nil.
//
// UpdateStandardGamepadLayoutMappings is concurrent-safe.
//
// UpdateStandardGamepadLayoutMappings mappings take effect immediately even for already connected gamepads.
//
// UpdateStandardGamepadLayoutMappings works atomically. If an error happens, nothing is updated.
func UpdateStandardGamepadLayoutMappings(mappings string) (bool, error) {
	if err := gamepaddb.Update([]byte(mappings)); err != nil {
		return false, err
	}
	return true, nil
}

// TouchID represents a touch's identifier.
type TouchID = ui.TouchID

// AppendTouchIDs appends the current touch states to touches, and returns the extended buffer.
// Giving a slice that already has enough capacity works efficiently.
//
// If you want to know whether a touch started being pressed in the current tick,
// use inpututil.JustPressedTouchIDs
//
// AppendTouchIDs doesn't append anything when there are no touches.
// AppendTouchIDs always does nothing on desktops.
//
// AppendTouchIDs is concurrent-safe.
func AppendTouchIDs(touches []TouchID) []TouchID {
	return theInputState.appendTouchIDs(touches)
}

// TouchIDs returns the current touch states.
//
// Deprecated: as of v2.2. Use AppendTouchIDs instead.
func TouchIDs() []TouchID {
	return AppendTouchIDs(nil)
}

// TouchPosition returns the position for the touch of the specified ID.
//
// If the touch of the specified ID is not present, TouchPosition returns (0, 0).
//
// TouchPosition is concurrent-safe.
func TouchPosition(id TouchID) (int, int) {
	return theInputState.touchPosition(id)
}

var theInputState inputState

type inputState struct {
	state ui.InputState
	m     sync.Mutex
}

func (i *inputState) update(fn func(*ui.InputState)) {
	i.m.Lock()
	defer i.m.Unlock()
	fn(&i.state)
}

func (i *inputState) appendInputChars(runes []rune) []rune {
	i.m.Lock()
	defer i.m.Unlock()
	for _, ie := range i.state.InputEvents {
		if ie.Kind != ui.InputEventKindText {
			continue
		}
		runes = append(runes, ie.Rune)
	}
	return runes
}

func (i *inputState) appendKeyEvents(events []KeyEvent) []KeyEvent {
	i.m.Lock()
	defer i.m.Unlock()
	for _, ie := range i.state.InputEvents {
		if ie.Kind != ui.InputEventKindKey {
			continue
		}
		events = append(events, KeyEvent{
			Key:    Key(ie.Key),
			Action: ie.Action,
			Mods:   ie.Mods,
		})
	}
	return events
}

func (i *inputState) appendInputEvents(events []InputEvent) []InputEvent {
	i.m.Lock()
	defer i.m.Unlock()
	for _, ie := range i.state.InputEvents {
		events = append(events, InputEvent{
			Kind:       ie.Kind,
			Key:        Key(ie.Key),
			Action:     ie.Action,
			Mods:       ie.Mods,
			Rune:       ie.Rune,
			NormalText: ie.NormalText,
			Source:     ie.Source,
		})
	}
	return events
}

func (i *inputState) isKeyPressed(key Key) bool {
	if !key.isValid() {
		return false
	}

	i.m.Lock()
	defer i.m.Unlock()

	switch key {
	case KeyAlt:
		return i.state.KeyPressed[ui.KeyAltLeft] || i.state.KeyPressed[ui.KeyAltRight]
	case KeyControl:
		return i.state.KeyPressed[ui.KeyControlLeft] || i.state.KeyPressed[ui.KeyControlRight]
	case KeyShift:
		return i.state.KeyPressed[ui.KeyShiftLeft] || i.state.KeyPressed[ui.KeyShiftRight]
	case KeyMeta:
		return i.state.KeyPressed[ui.KeyMetaLeft] || i.state.KeyPressed[ui.KeyMetaRight]
	default:
		return i.state.KeyPressed[key]
	}
}

func (i *inputState) cursorPosition() (float64, float64) {
	i.m.Lock()
	defer i.m.Unlock()
	return i.state.CursorX, i.state.CursorY
}

func (i *inputState) wheel() (float64, float64) {
	i.m.Lock()
	defer i.m.Unlock()
	return i.state.WheelX, i.state.WheelY
}

func (i *inputState) isMouseButtonPressed(mouseButton MouseButton) bool {
	i.m.Lock()
	defer i.m.Unlock()
	return i.state.MouseButtonPressed[mouseButton]
}

func (i *inputState) appendTouchIDs(touches []TouchID) []TouchID {
	i.m.Lock()
	defer i.m.Unlock()

	for _, t := range i.state.Touches {
		touches = append(touches, t.ID)
	}
	return touches
}

func (i *inputState) touchPosition(id TouchID) (int, int) {
	i.m.Lock()
	defer i.m.Unlock()

	for _, t := range i.state.Touches {
		if id != t.ID {
			continue
		}
		return t.X, t.Y
	}
	return 0, 0
}

func (i *inputState) windowBeingClosed() bool {
	i.m.Lock()
	defer i.m.Unlock()
	return i.state.WindowBeingClosed
}

func (i *inputState) droppedFiles() fs.FS {
	i.m.Lock()
	defer i.m.Unlock()
	return i.state.DroppedFiles
}

func (i *inputState) droppedFilePaths() []string {
	i.m.Lock()
	defer i.m.Unlock()
	if len(i.state.DroppedFilePaths) == 0 {
		return nil
	}
	return slices.Clone(i.state.DroppedFilePaths)
}

func (i *inputState) dragging() (float64, float64, bool) {
	i.m.Lock()
	defer i.m.Unlock()
	return i.state.DragX, i.state.DragY, i.state.Dragging
}
