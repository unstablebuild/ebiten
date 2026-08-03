// Copyright 2022 The Ebitengine Authors
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

package ui

import (
	"io/fs"
	"unicode"
)

// KeyAction represents the action of a key event.
type KeyAction int

const (
	// KeyRelease indicates a key was released.
	KeyRelease KeyAction = 0
	// KeyPress indicates a key was pressed.
	KeyPress KeyAction = 1
	// KeyRepeat indicates a key repeat from the OS.
	KeyRepeat KeyAction = 2
)

// KeyModifier is a bitmask of modifier keys held during a key event.
type KeyModifier int

const (
	// KeyModShift indicates the Shift key was held.
	KeyModShift KeyModifier = 0x0001
	// KeyModControl indicates the Control key was held.
	KeyModControl KeyModifier = 0x0002
	// KeyModAlt indicates the Alt/Option key was held.
	KeyModAlt KeyModifier = 0x0004
	// KeyModSuper indicates the Super/Meta/Command key was held.
	KeyModSuper KeyModifier = 0x0008
)

// KeyEvent represents a discrete key event from the OS.
type KeyEvent struct {
	Key    Key
	Action KeyAction
	Mods   KeyModifier
}

type MouseButton int

const (
	MouseButton0   MouseButton = iota // The 'left' button
	MouseButton1                      // The 'right' button
	MouseButton2                      // The 'middle' button
	MouseButton3                      // The additional button (usually browser-back)
	MouseButton4                      // The additional button (usually browser-forward)
	MouseButtonMax = MouseButton4
)

type TouchID int

type Touch struct {
	ID TouchID
	X  int
	Y  int
}

type InputState struct {
	KeyPressed         [KeyMax + 1]bool
	MouseButtonPressed [MouseButtonMax + 1]bool
	CursorX            float64
	CursorY            float64
	WheelX             float64
	WheelY             float64
	Touches            []Touch
	Runes              []rune
	KeyEvents          []KeyEvent
	WindowBeingClosed  bool
	DroppedFiles       fs.FS
}

func (i *InputState) copyAndReset(dst *InputState) {
	dst.KeyPressed = i.KeyPressed
	dst.MouseButtonPressed = i.MouseButtonPressed
	dst.CursorX = i.CursorX
	dst.CursorY = i.CursorY
	dst.WheelX = i.WheelX
	dst.WheelY = i.WheelY
	dst.Touches = append(dst.Touches[:0], i.Touches...)
	dst.Runes = append(dst.Runes[:0], i.Runes...)
	dst.KeyEvents = append(dst.KeyEvents[:0], i.KeyEvents...)
	dst.WindowBeingClosed = i.WindowBeingClosed
	dst.DroppedFiles = i.DroppedFiles

	// Reset the members that are updated by deltas, rather than absolute values.
	i.WheelX = 0
	i.WheelY = 0
	i.Runes = i.Runes[:0]
	i.KeyEvents = i.KeyEvents[:0]

	// Reset the members that are never reset until they are explicitly done.
	i.WindowBeingClosed = false
	i.DroppedFiles = nil
}

func (i *InputState) appendRune(r rune) {
	if !unicode.IsPrint(r) && !isEmojiSequenceRune(r) {
		return
	}
	i.Runes = append(i.Runes, r)
}

// isEmojiSequenceRune reports whether r is a zero-width, non-printable
// code point that nonetheless composes an emoji grapheme cluster and so
// must be preserved as text input. IMEs and the platform emoji picker
// deliver these through the character callback: the zero-width joiner
// (U+200D) stitches family/profession sequences, and the tag characters
// (U+E0020..U+E007F) encode subdivision flags such as 🏴󠁧󠁢󠁳󠁣󠁴󠁿. Dropping them
// via the unicode.IsPrint filter would split a single emoji into its
// component glyphs. Emoji variation selectors and skin-tone modifiers are
// already printable, so they do not need to be listed here.
func isEmojiSequenceRune(r rune) bool {
	return r == '\u200d' || (r >= '\U000E0020' && r <= '\U000E007F')
}

func (i *InputState) appendKeyEvent(key Key, action KeyAction, mods KeyModifier) {
	i.KeyEvents = append(i.KeyEvents, KeyEvent{
		Key:    key,
		Action: action,
		Mods:   mods,
	})
}
