// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

// Support for the keyboard-layout end-to-end test. cgo is not allowed in
// _test.go files, so the bridge to the platform translation lives here.
// Run the test with `go test -tags e2e ./internal/glfw/`.

//go:build (linux || freebsd || openbsd || netbsd) && !android && e2e

package glfw

// #include <stdint.h>
//
// extern uint32_t _glfwPlatformGetScancodeCodepoint(int scancode, int shift);
import "C"

// scancodeCodepoint returns the code point the layout X currently reports for
// the given evdev keycode at its unshifted or shifted level. GLFW must be
// initialized first.
func scancodeCodepoint(scancode int, shift bool) rune {
	var s C.int
	if shift {
		s = 1
	}
	return rune(C._glfwPlatformGetScancodeCodepoint(C.int(scancode), s))
}