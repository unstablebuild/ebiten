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

package graphicscommand

import (
	"sync/atomic"
)

var drawTrianglesCommandCount int64

// DrawTrianglesCommandCount returns the cumulative number of merged
// draw-triangles commands flushed to the graphics driver since program
// start.
//
// Draw operations merged by the command queue count as one command, so
// the counter reflects the number of driver DrawTriangles invocations
// rather than the number of application-level draw calls.
//
// DrawTrianglesCommandCount is concurrent-safe.
func DrawTrianglesCommandCount() int64 {
	return atomic.LoadInt64(&drawTrianglesCommandCount)
}