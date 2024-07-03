package benchdraw

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/atlas"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicscommand"
)

// BeginFrame updates the atlas package state to enable
// running benchmarks with a no-op graphicsdriver.
// This sould be called at the start of a benchmark loop.
func BeginFrame(t testing.TB) {
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
