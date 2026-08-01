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

package ebiten

import (
	"github.com/hajimehoshi/ebiten/v2/internal/graphicscommand"
)

// DrawCommandCount returns the cumulative number of draw commands flushed
// to the graphics driver since program start.
//
// Consecutive draw operations that Ebitengine merges internally (same
// destination, sources, blend, shader, and uniforms) count as a single
// command, so the counter reflects the actual number of driver draw
// invocations rather than the number of DrawImage/DrawTriangles calls
// made by the application. Compute per-frame figures by subtracting two
// readings.
//
// This function is for measurement and/or debug, and your game logic
// should not rely on this value.
//
// DrawCommandCount is concurrent-safe.
func DrawCommandCount() int64 {
	return graphicscommand.DrawTrianglesCommandCount()
}