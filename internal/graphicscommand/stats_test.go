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

package graphicscommand_test

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicscommand"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
	"github.com/hajimehoshi/ebiten/v2/internal/ui"
)

func TestDrawTrianglesCommandCount(t *testing.T) {
	const w, h = 16, 16
	src := graphicscommand.NewImage(w, h, false)
	dst := graphicscommand.NewImage(w, h, false)
	vs := quadVertices(w, h)
	is := graphics.QuadIndices()
	dr := image.Rect(0, 0, w, h)
	srcs := [graphics.ShaderImageCount]*graphicscommand.Image{src}
	var srcRegions [graphics.ShaderImageCount]image.Rectangle

	// A read-pixels command needs a synchronous flush, so ReadPixels
	// both drains commands pending from UI initialization and works as
	// a deterministic settle point after each batch below.
	pix := make([]byte, 4*w*h)
	readPixels := func() {
		if err := dst.ReadPixels(ui.Get().GraphicsDriverForTesting(), []graphicsdriver.PixelsArgs{
			{Pixels: pix, Region: image.Rect(0, 0, w, h)},
		}); err != nil {
			t.Fatal(err)
		}
	}
	readPixels()

	// Two draws with identical state must merge into one command.
	before := graphicscommand.DrawTrianglesCommandCount()
	dst.DrawTriangles(srcs, vs, is, graphicsdriver.BlendSourceOver, dr, srcRegions, nearestFilterShader, nil, graphicsdriver.FillAll)
	dst.DrawTriangles(srcs, vs, is, graphicsdriver.BlendSourceOver, dr, srcRegions, nearestFilterShader, nil, graphicsdriver.FillAll)
	readPixels()
	if got := graphicscommand.DrawTrianglesCommandCount() - before; got != 1 {
		t.Errorf("merged draws: got %d commands, want 1", got)
	}

	// Two draws with different blend states cannot merge.
	before = graphicscommand.DrawTrianglesCommandCount()
	dst.DrawTriangles(srcs, vs, is, graphicsdriver.BlendSourceOver, dr, srcRegions, nearestFilterShader, nil, graphicsdriver.FillAll)
	dst.DrawTriangles(srcs, vs, is, graphicsdriver.BlendClear, dr, srcRegions, nearestFilterShader, nil, graphicsdriver.FillAll)
	readPixels()
	if got := graphicscommand.DrawTrianglesCommandCount() - before; got != 2 {
		t.Errorf("unmergeable draws: got %d commands, want 2", got)
	}
}