//go:build !android && !ios && !js && !nintendosdk && !playstation5

package ui

import "testing"

func TestShouldDisableCocoaMenubarByEnv(t *testing.T) {
	t.Setenv("EBITENGINE_COCOA_MENUBAR", "")
	if shouldDisableCocoaMenubarByEnv() {
		t.Fatal("shouldDisableCocoaMenubarByEnv() = true, want false")
	}

	t.Setenv("EBITENGINE_COCOA_MENUBAR", "1")
	if shouldDisableCocoaMenubarByEnv() {
		t.Fatal("shouldDisableCocoaMenubarByEnv() = true, want false")
	}

	t.Setenv("EBITENGINE_COCOA_MENUBAR", "0")
	if !shouldDisableCocoaMenubarByEnv() {
		t.Fatal("shouldDisableCocoaMenubarByEnv() = false, want true")
	}
}

func TestShouldHideDockAfterFocus(t *testing.T) {
	t.Setenv("EBITENGINE_COCOA_HIDE_DOCK", "")
	if shouldHideDockAfterFocus() {
		t.Fatal("shouldHideDockAfterFocus() = true, want false")
	}

	t.Setenv("EBITENGINE_COCOA_HIDE_DOCK", "0")
	if shouldHideDockAfterFocus() {
		t.Fatal("shouldHideDockAfterFocus() = true, want false")
	}

	t.Setenv("EBITENGINE_COCOA_HIDE_DOCK", "1")
	if !shouldHideDockAfterFocus() {
		t.Fatal("shouldHideDockAfterFocus() = false, want true")
	}
}
