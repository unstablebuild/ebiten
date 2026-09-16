// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

// Support for the keyboard-layout tests. cgo is not allowed in _test.go
// files, so the bridge to the C predicate lives here.

//go:build darwin || freebsd || linux || netbsd || openbsd

package glfw

// extern int _glfwIsPrintableKey(int key);
import "C"

// platformIsPrintableKey is the C side of isPrintableKey.
func platformIsPrintableKey(key Key) bool {
	return C._glfwIsPrintableKey(C.int(key)) != 0
}
