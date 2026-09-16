// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build (linux || freebsd || openbsd || netbsd) && !android && e2e

package glfw

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

// X11 keycodes of the physical keys the cases below name. X adds 8 to the
// evdev code, and the names are the US positions, which is exactly the naming
// that layout translation has to override.
const (
	keycodeA            = 38 // US A
	keycodeZ            = 52 // US Z
	keycodeP            = 33 // US P
	keycodeQ            = 24 // US Q
	keycodeR            = 27 // US R
	keycodeSemicolon    = 47 // US ;
	keycodeMinus        = 20 // US -
	keycodeBracketRight = 35 // US ]
	keycodeKP7          = 79
	keycodeKPAdd        = 86
	keycodeReturn       = 36
	keycodeF1           = 67
	keycodeShiftLeft    = 50
)

// TestScancodeCodepoint verifies that a physical key is named by the active
// XKB layout rather than by its US position. Every case is a chord a non-US
// user reported as unreachable: the key labelled M on AZERTY sits on the US
// semicolon, Colemak P sits on the US R, and the Nordic +/? key sits on the
// US minus. The Cyrillic and Greek rows leave Latin-1, so the keysym table
// lookup is exercised rather than the identity mapping.
func TestScancodeCodepoint(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	startXvfb(t)

	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer Terminate()

	cases := []struct {
		layout    []string
		keycode   int
		want      rune
		wantShift rune
	}{
		{layout: []string{"us"}, keycode: keycodeSemicolon, want: ';', wantShift: ':'},
		{layout: []string{"us"}, keycode: keycodeA, want: 'a', wantShift: 'A'},
		{layout: []string{"us"}, keycode: keycodeKP7, want: '7', wantShift: '7'},
		{layout: []string{"us"}, keycode: keycodeKPAdd, want: '+', wantShift: '+'},
		{layout: []string{"fr"}, keycode: keycodeSemicolon, want: 'm', wantShift: 'M'},
		{layout: []string{"fr"}, keycode: keycodeQ, want: 'a', wantShift: 'A'},
		{layout: []string{"no"}, keycode: keycodeMinus, want: '+', wantShift: '?'},
		{layout: []string{"no"}, keycode: keycodeBracketRight, want: '¨', wantShift: '^'},
		{layout: []string{"us", "-variant", "colemak"}, keycode: keycodeR, want: 'p', wantShift: 'P'},
		{layout: []string{"us", "-variant", "colemak"}, keycode: keycodeP, want: ';', wantShift: ':'},
		{layout: []string{"us", "-variant", "dvorak"}, keycode: keycodeQ, want: '\'', wantShift: '"'},
		{layout: []string{"de"}, keycode: keycodeZ, want: 'y', wantShift: 'Y'},
		{layout: []string{"ru"}, keycode: keycodeA, want: 'ф', wantShift: 'Ф'},
		{layout: []string{"gr"}, keycode: keycodeA, want: 'α', wantShift: 'Α'},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("%v/%c", tc.layout, tc.want)
		t.Run(name, func(t *testing.T) {
			setLayout(t, tc.layout...)
			waitCodepoint(t, tc.keycode, false, tc.want)
			if got := scancodeCodepoint(tc.keycode, true); got != tc.wantShift {
				t.Errorf("shift+keycode %d = %q, want %q", tc.keycode, got, tc.wantShift)
			}
		})
	}
}

// TestScancodeCodepointNoCharacter verifies that keycodes outside the
// keyboard's range, keycodes GLFW has no key for, and keys that name no
// character all report 0 rather than an error or a control character. The
// key event path asks for every key it sees.
func TestScancodeCodepointNoCharacter(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	startXvfb(t)

	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer Terminate()

	setLayout(t, "us")
	waitCodepoint(t, keycodeA, false, 'a')

	for _, keycode := range []int{-1, 0, 7, 256, 1000, keycodeReturn, keycodeF1, keycodeShiftLeft} {
		for _, shift := range []bool{false, true} {
			if got := scancodeCodepoint(keycode, shift); got != 0 {
				t.Errorf("keycode %d shift=%v = %q, want 0", keycode, shift, got)
			}
		}
	}
}

// TestScancodeCodepointLayoutChangeAfterFirstLookup verifies that a layout
// change lands even when it happens right after the first lookup. Xlib only
// asks the server for map change notifications when that first lookup loads
// its copy of the map, and the request sits in its output buffer until the
// next flush; a change made in that window used to go unreported and the
// stale map was served for the rest of the session.
func TestScancodeCodepointLayoutChangeAfterFirstLookup(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	startXvfb(t)

	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer Terminate()

	// The first lookup of the session, with no Xlib traffic between it and
	// the layout change below.
	if got := scancodeCodepoint(keycodeSemicolon, false); got != ';' {
		t.Fatalf("keycode %d = %q, want ';'", keycodeSemicolon, got)
	}
	setLayout(t, "fr")
	waitCodepoint(t, keycodeSemicolon, false, 'm')
}

// TestScancodeCodepointFollowsGroupSwitch verifies that with several layouts
// loaded at once, the character follows the group the user switched to,
// which is how desktop layout switchers work.
func TestScancodeCodepointFollowsGroupSwitch(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	startXvfb(t)

	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer Terminate()

	setLayout(t, "-layout", "us,fr")
	if !lockKeyboardGroup(1) {
		t.Fatal("lock keyboard group 1")
	}
	waitCodepoint(t, keycodeA, false, 'q')
	if got := scancodeCodepoint(keycodeSemicolon, false); got != 'm' {
		t.Errorf("keycode %d in group 1 = %q, want 'm'", keycodeSemicolon, got)
	}
	// The keypad has a single group; it must still resolve from group 1.
	if got := scancodeCodepoint(keycodeKP7, false); got != '7' {
		t.Errorf("keycode %d in group 1 = %q, want '7'", keycodeKP7, got)
	}

	if !lockKeyboardGroup(0) {
		t.Fatal("lock keyboard group 0")
	}
	waitCodepoint(t, keycodeA, false, 'a')
}

// startXvfb runs a headless X server for the duration of the test and points
// DISPLAY at it, so the suite never touches the developer's own session.
func startXvfb(t *testing.T) {
	t.Helper()

	for _, bin := range []string{"Xvfb", "setxkbmap"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s is not installed", bin)
		}
	}

	const display = ":99"
	cmd := exec.Command("Xvfb", display, "-screen", "0", "640x480x24")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start Xvfb: %v", err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})
	t.Setenv("DISPLAY", display)

	// Xvfb has no readiness signal; retry until the socket accepts clients.
	deadline := time.Now().Add(10 * time.Second)
	for {
		if err := exec.Command("setxkbmap", "us").Run(); err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("Xvfb did not become ready")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// setLayout changes the server's keyboard map from a separate client. GLFW
// learns of it asynchronously; use waitCodepoint to observe the result.
func setLayout(t *testing.T, args ...string) {
	t.Helper()

	cmd := exec.Command("setxkbmap", args...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("setxkbmap %v: %v", args, err)
	}
}

// waitCodepoint pumps GLFW's event loop until the keycode resolves to want,
// failing if it has not within a generous deadline.
func waitCodepoint(t *testing.T, keycode int, shift bool, want rune) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for {
		if err := PollEvents(); err != nil {
			t.Fatalf("PollEvents: %v", err)
		}
		got := scancodeCodepoint(keycode, shift)
		if got == want {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("keycode %d shift=%v = %q, want %q", keycode, shift, got, want)
		}
		time.Sleep(5 * time.Millisecond)
	}
}
