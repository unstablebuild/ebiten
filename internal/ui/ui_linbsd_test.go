// Copyright 2026 The Ebitengine Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build (freebsd || (linux && !android) || netbsd || openbsd) && !nintendosdk && !playstation5

package ui

import (
	"errors"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/glfw"
)

func TestInitWithoutDisplay(t *testing.T) {
	t.Setenv("DISPLAY", "")
	if err := os.Unsetenv("DISPLAY"); err != nil {
		t.Fatal(err)
	}

	err := Init()
	if !errors.Is(err, glfw.PlatformError) {
		t.Fatalf("Init() = %v, want a glfw.PlatformError", err)
	}
	if again := Init(); again != err {
		t.Errorf("second Init() = %v, want the first error %v", again, err)
	}

	defer func() {
		if r := recover(); r != err {
			t.Errorf("Get() panicked with %v, want %v", r, err)
		}
	}()
	Get()
}
