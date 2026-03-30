package ebiten

import "testing"

func TestToUIRunOptionsPropagatesSkipCocoaMenubar(t *testing.T) {
	op := toUIRunOptions(&RunGameOptions{
		InitUnfocused:    true,
		SkipCocoaMenubar: true,
	})

	if !op.InitUnfocused {
		t.Fatal("InitUnfocused = false, want true")
	}
	if !op.SkipCocoaMenubar {
		t.Fatal("SkipCocoaMenubar = false, want true")
	}
}
