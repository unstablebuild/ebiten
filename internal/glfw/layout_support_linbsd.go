// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

// Support for the keyboard-layout end-to-end test. cgo is not allowed in
// _test.go files, so the bridge to the platform translation lives here.
// Run the test with `go test -tags e2e ./internal/glfw/`.

//go:build (linux || freebsd || openbsd || netbsd) && !android && e2e

package glfw

// // GLFW loads Xlib with dlopen; the test links it directly for the second
// // connection it uses to drive the server behind GLFW's back.
// #cgo LDFLAGS: -lX11
// #include <stdint.h>
// #include <X11/XKBlib.h>
//
// extern uint32_t _glfwPlatformGetScancodeCodepoint(int scancode, int shift);
//
// // Switches the server's keyboard group from a connection of its own, the
// // way a desktop layout switcher does, so GLFW only learns about it through
// // the XkbStateNotify it receives.
// static int lockGroup(int group) {
//     Display* display = XOpenDisplay(NULL);
//     if (!display) {
//         return 0;
//     }
//     int ok = XkbLockGroup(display, XkbUseCoreKbd, group);
//     XSync(display, False);
//     XCloseDisplay(display);
//     return ok;
// }
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

// lockKeyboardGroup selects the given XKB group on the server from a separate
// connection.
func lockKeyboardGroup(group int) bool {
	return C.lockGroup(C.int(group)) != 0
}
