// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build darwin && glfwlayout

package glfw

import "testing"

// Virtual key codes of the physical keys the cases below press. They are the
// US-positional codes AppKit reports, which is exactly the naming that layout
// translation has to override.
const (
	scancodeA            = 0x00 // US A
	scancodeZ            = 0x06 // US Z
	scancodeP            = 0x23 // US P
	scancodeQ            = 0x0C // US Q
	scancodeR            = 0x0F // US R
	scancodeSemicolon    = 0x29 // US ;
	scancodeMinus        = 0x1B // US -
	scancodeBracketRight = 0x1E // US ]
)

// TestTranslateScancode verifies that a physical key is named by the active
// keyboard layout rather than by its US position. Every case is a chord a
// non-US user reported as unreachable: the key labelled M on AZERTY sits on
// the US semicolon, Colemak P sits on the US R, and the Nordic +/? key sits
// on the US minus.
func TestTranslateScancode(t *testing.T) {
	cases := []struct {
		layout    string
		scancode  int
		want      rune
		wantShift rune
	}{
		{layout: "US", scancode: scancodeSemicolon, want: ';', wantShift: ':'},
		{layout: "US", scancode: scancodeA, want: 'a', wantShift: 'A'},
		{layout: "French", scancode: scancodeSemicolon, want: 'm', wantShift: 'M'},
		{layout: "French", scancode: scancodeQ, want: 'a', wantShift: 'A'},
		{layout: "Norwegian", scancode: scancodeMinus, want: '+', wantShift: '?'},
		{layout: "Norwegian", scancode: scancodeBracketRight, want: '¨', wantShift: '^'},
		{layout: "Colemak", scancode: scancodeR, want: 'p', wantShift: 'P'},
		{layout: "Colemak", scancode: scancodeP, want: ';', wantShift: ':'},
		{layout: "Dvorak", scancode: scancodeQ, want: '\'', wantShift: '"'},
		{layout: "German", scancode: scancodeZ, want: 'y', wantShift: 'Y'},
	}

	for _, tc := range cases {
		t.Run(tc.layout+"/"+string(tc.want), func(t *testing.T) {
			layout, ok := loadKeyboardLayout(tc.layout)
			if !ok {
				t.Skipf("keyboard layout %q is not available on this machine", tc.layout)
			}

			if got := translateScancode(layout, tc.scancode, false); got != tc.want {
				t.Errorf("scancode %#x = %q, want %q", tc.scancode, got, tc.want)
			}
			if got := translateScancode(layout, tc.scancode, true); got != tc.wantShift {
				t.Errorf("shift+scancode %#x = %q, want %q", tc.scancode, got, tc.wantShift)
			}
		})
	}
}

// TestTranslateScancodeNoLayout verifies that a missing layout yields no
// character instead of crashing. The key event path calls this for every
// printable key press, so it must tolerate the platform having no layout data.
func TestTranslateScancodeNoLayout(t *testing.T) {
	if got := translateScancode(keyboardLayout{}, scancodeA, false); got != 0 {
		t.Errorf("translateScancode(no layout) = %q, want 0", got)
	}
}