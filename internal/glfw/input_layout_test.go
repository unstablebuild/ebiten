// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

package glfw

import (
	"testing"
	"unicode"
)

func TestIsPrintableKey(t *testing.T) {
	cases := []struct {
		key  Key
		want bool
	}{
		{KeyUnknown, false},
		{KeySpace, false},
		{KeyApostrophe, true},
		{KeyComma, true},
		{Key0, true},
		{Key9, true},
		{KeySemicolon, true},
		{KeyEqual, true},
		{KeyA, true},
		{KeyZ, true},
		{KeyLeftBracket, true},
		{KeyBackslash, true},
		{KeyGraveAccent, true},
		{KeyWorld1, true},
		{KeyWorld2, true},
		{KeyEscape, false},
		{KeyEnter, false},
		{KeyTab, false},
		{KeyBackspace, false},
		{KeyF1, false},
		{KeyF24, false},
		{KeyKP0, true},
		{KeyKP9, true},
		{KeyKPDecimal, true},
		{KeyKPDivide, true},
		{KeyKPMultiply, true},
		{KeyKPSubtract, true},
		{KeyKPAdd, true},
		{KeyKPEnter, false},
		{KeyKPEqual, true},
		{KeyLeftShift, false},
		{KeyRightSuper, false},
		{KeyMenu, false},
		{KeyLast, false},
	}
	for _, tc := range cases {
		if got := isPrintableKey(tc.key); got != tc.want {
			t.Errorf("isPrintableKey(%d) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestLayoutRune(t *testing.T) {
	cases := []struct {
		name      string
		codepoint uint32
		want      rune
	}{
		{"none", 0, 0},
		{"ascii", 'a', 'a'},
		{"space", ' ', ' '},
		{"latin-1", '¨', '¨'},
		{"cyrillic", 'ф', 'ф'},
		{"astral", '\U0001F600', '\U0001F600'},
		{"max rune", unicode.MaxRune, unicode.MaxRune},
		{"carriage return from Return on Cocoa", '\r', 0},
		{"function key marker from Cocoa", 0x10, 0},
		{"delete", 0x7f, 0},
		{"C1 control", 0x85, 0},
		{"high surrogate", 0xd800, 0},
		{"low surrogate", 0xdfff, 0},
		{"replacement character from a lone surrogate on Windows", 0xfffd, 0},
		{"past the Unicode range", 0x110000, 0},
		{"unbounded X11 keysym", 0xffffff, 0},
		{"sign bit set", 0xffffffff, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := layoutRune(tc.codepoint); got != tc.want {
				t.Errorf("layoutRune(%#x) = %#x, want %#x", tc.codepoint, got, tc.want)
			}
		})
	}
}

func TestFirstRune(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want rune
	}{
		{"empty", "", 0},
		{"single", "a", 'a'},
		{"dead key reported twice by ToUnicode", "^^", '^'},
		{"multi-byte", "ф", 'ф'},
		{"astral", "\U0001F600", '\U0001F600'},
		{"invalid UTF-8", "\xff", '\uFFFD'},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := firstRune(tc.s); got != tc.want {
				t.Errorf("firstRune(%q) = %q, want %q", tc.s, got, tc.want)
			}
		})
	}
}
