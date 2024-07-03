package benchdraw

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/atlas"
)

// BeginFrame updates the atlas package state to enable
// running benchmarks with a no-op graphicsdriver.
func BeginFrame(t testing.TB) {
	err := atlas.BeginFrame(NopGraphicsDriver{})
	if err != nil {
		t.Logf("begin frame: %v", err)
		t.FailNow()
	}
}

// EndFrame politely resets the atlas package state to
// indicate that no more benchmarks will be run.
func EndFrame(t testing.TB) {
	err := atlas.EndFrame()
	if err != nil {
		t.Logf("end frame: %v", err)
		t.FailNow()
	}
}
