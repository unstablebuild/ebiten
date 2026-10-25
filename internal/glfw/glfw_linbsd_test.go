// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build freebsd || linux || netbsd || openbsd

package glfw

import (
	"errors"
	"os"
	"testing"
)

func TestInitWithoutDisplay(t *testing.T) {
	t.Setenv("DISPLAY", "")
	if err := os.Unsetenv("DISPLAY"); err != nil {
		t.Fatal(err)
	}

	err := Init()
	if err == nil {
		_ = Terminate()
		t.Fatal("Init() = nil, want an error")
	}
	if !errors.Is(err, PlatformError) {
		t.Fatalf("Init() = %v, want a PlatformError", err)
	}
}
