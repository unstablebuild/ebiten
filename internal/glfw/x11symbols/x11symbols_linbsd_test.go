// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build freebsd || linux || netbsd || openbsd

package x11symbols_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/glfw"
	"github.com/hajimehoshi/ebiten/v2/internal/glfw/x11symbols"
)

// The real check is that this binary linked at all; see the package comment.
func TestSymbolsAreExported(t *testing.T) {
	refs := x11symbols.Refs()
	if len(refs) == 0 {
		t.Fatal("_GLFW_X11_SYMBOLS is empty")
	}
	for i, ref := range refs {
		if ref == nil {
			t.Errorf("entry %d of _GLFW_X11_SYMBOLS resolved to NULL", i)
		}
	}
}

// Exercises the production loader without a display: the dlopen by SONAME
// and the dlsym of every listed entry point, then its idempotence.
func TestLoadX11Library(t *testing.T) {
	for _, call := range []string{"first", "repeated"} {
		if err := glfw.LoadX11LibraryForTesting(); err != nil {
			t.Fatalf("%s call: %v", call, err)
		}
	}
}
