// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

package glfw

import (
	"unicode"
	"unicode/utf8"
)

// isPrintableKey reports whether key names a printable character under some
// layout, and so has a layout code point at all. The C key event path and
// glfwGetKeyName gate on the same predicate; TestIsPrintableKeyMatchesC keeps
// the two from drifting.
func isPrintableKey(key Key) bool {
	return key == KeyKPEqual ||
		(key >= KeyKP0 && key <= KeyKPAdd) ||
		(key >= KeyApostrophe && key <= KeyWorld2)
}

// layoutRune returns the code point a platform layout query produced for a
// key, or 0 if it is not a character a chord could be bound to. Cocoa reports
// control characters for editing and function keys, nothing upstream rejects
// a surrogate or a value past the Unicode range, and Windows decodes a lone
// surrogate to the replacement character.
func layoutRune(codepoint uint32) rune {
	if codepoint > unicode.MaxRune {
		return 0
	}
	r := rune(codepoint)
	if !utf8.ValidRune(r) || unicode.IsControl(r) || r == utf8.RuneError {
		return 0
	}
	return r
}

// firstRune returns the first code point of s, or 0 if s is empty.
func firstRune(s string) rune {
	for _, r := range s {
		return r
	}
	return 0
}
