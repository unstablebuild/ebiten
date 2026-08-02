package benchdraw

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/atlas"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicscommand"
	"github.com/hajimehoshi/ebiten/v2/internal/ui"
)

// BeginFrame updates the atlas package state to enable
// running benchmarks with a no-op graphicsdriver.
// This sould be called at the start of a benchmark loop.
func BeginFrame(t testing.TB) {
	// Route driver-backed operations that run during a game's Draw
	// (most importantly ReadPixels, used when rasterizing glyph masks
	// backed by ebiten images) through the no-op driver so headless
	// benchmarks do not dereference the nil windowed driver.
	ui.Get().SetGraphicsDriverForTesting(nopGraphicsDriver{})
	err := atlas.BeginFrame(nopGraphicsDriver{})
	if err != nil {
		t.Logf("begin frame: %v", err)
		t.FailNow()
	}
}

// EndFrame politely resets the atlas package and simulates
// a flush to the underlying graphics driver.
// This sould be called at the end of a benchmark loop.
func EndFrame(t testing.TB) {
	err := atlas.EndFrame()
	if err != nil {
		t.Logf("end frame: %v", err)
		t.FailNow()
	}

	err = graphicscommand.FlushCommands(nopGraphicsDriver{}, true)
	if err != nil {
		t.Logf("flush commands: %v", err)
		t.FailNow()
	}
}
