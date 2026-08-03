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
	"reflect"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/ui"
)

// TestInputProjections verifies that AppendInputEvents reports every
// observation in native order with its source, and that the older
// AppendInputChars and AppendKeyEvents remain the documented projections of
// that sequence.
func TestInputProjections(t *testing.T) {
	var s inputState
	s.state.InputEvents = []ui.InputEvent{
		{Kind: ui.InputEventKindKey, Key: ui.KeyL, Action: ui.KeyPress, Mods: ui.KeyModSuper, Source: 1},
		{Kind: ui.InputEventKindText, Mods: ui.KeyModSuper, Rune: 'l', NormalText: true, Source: 1},
		{Kind: ui.InputEventKindText, Rune: '漢', NormalText: true, Source: 0},
		{Kind: ui.InputEventKindKey, Key: ui.KeyL, Action: ui.KeyRelease, Mods: ui.KeyModSuper, Source: 2},
	}

	wantEvents := []InputEvent{
		{Kind: InputEventKindKey, Key: KeyL, Action: KeyActionPress, Mods: KeyModSuper, Source: 1},
		{Kind: InputEventKindText, Mods: KeyModSuper, Rune: 'l', NormalText: true, Source: 1},
		{Kind: InputEventKindText, Rune: '漢', NormalText: true},
		{Kind: InputEventKindKey, Key: KeyL, Action: KeyActionRelease, Mods: KeyModSuper, Source: 2},
	}
	if got := s.appendInputEvents(nil); !reflect.DeepEqual(got, wantEvents) {
		t.Errorf("appendInputEvents() = %+v, want %+v", got, wantEvents)
	}

	wantChars := []rune{'l', '漢'}
	if got := s.appendInputChars(nil); !reflect.DeepEqual(got, wantChars) {
		t.Errorf("appendInputChars() = %q, want %q", got, wantChars)
	}

	wantKeys := []KeyEvent{
		{Key: KeyL, Action: KeyActionPress, Mods: KeyModSuper},
		{Key: KeyL, Action: KeyActionRelease, Mods: KeyModSuper},
	}
	if got := s.appendKeyEvents(nil); !reflect.DeepEqual(got, wantKeys) {
		t.Errorf("appendKeyEvents() = %+v, want %+v", got, wantKeys)
	}
}