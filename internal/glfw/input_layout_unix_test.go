// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build darwin || freebsd || linux || netbsd || openbsd

package glfw

import "testing"

// TestIsPrintableKeyMatchesC verifies that the Go predicate the Windows key
// event path uses agrees with the C one the other platforms and
// glfwGetKeyName use, for every key.
func TestIsPrintableKeyMatchesC(t *testing.T) {
	for key := KeyUnknown; key <= KeyLast+1; key++ {
		if got, want := isPrintableKey(key), platformIsPrintableKey(key); got != want {
			t.Errorf("isPrintableKey(%d) = %v, C = %v", key, got, want)
		}
	}
}
