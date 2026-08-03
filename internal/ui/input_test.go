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

package ui

import (
	"reflect"
	"testing"
)

// TestInputStateAppendRune verifies that the character callback keeps
// printable runes and the zero-width joiner / tag characters that
// compose emoji grapheme clusters, while still dropping ordinary
// control characters. The zero-width joiner in particular is delivered
// by the platform emoji picker for family and profession sequences; if
// it were filtered, such emoji would split into their component glyphs.
func TestInputStateAppendRune(t *testing.T) {
	cases := []struct {
		name string
		r    rune
		keep bool
	}{
		{"ascii", 'a', true},
		{"emoji base", '\U0001F468', true},
		{"skin tone modifier", '\U0001F3FC', true},
		{"variation selector 16", '\uFE0F', true},
		{"zero width joiner", '\u200d', true},
		{"flag tag letter", '\U000E0067', true},
		{"flag tag terminator", '\U000E007F', true},
		{"tab", '\t', false},
		{"newline", '\n', false},
		{"null", '\x00', false},
		{"bell", '\a', false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s InputState
			s.appendRune(tc.r)
			got := len(s.InputEvents) == 1 && s.InputEvents[0].Rune == tc.r
			if got != tc.keep {
				t.Fatalf("appendRune(%#U): kept=%v, want kept=%v (events=%v)",
					tc.r, got, tc.keep, s.InputEvents)
			}
		})
	}
}

// TestInputStateOrdering verifies that key transitions and committed text are
// recorded as one ordered sequence of observations. GLFW reports both from the
// same physical action, so keeping them in separate buffers would destroy the
// order and the association a consumer needs to tell an echo of a chord it
// already handled from unrelated text.
func TestInputStateOrdering(t *testing.T) {
	cases := []struct {
		name   string
		record func(*InputState)
		want   []InputEvent
	}{
		{
			name: "key press then its text",
			record: func(s *InputState) {
				s.appendKeyEvent(KeyL, KeyPress, KeyModSuper, 7)
				s.appendTextInput('l', KeyModSuper, true, 7)
			},
			want: []InputEvent{
				{Kind: InputEventKindKey, Key: KeyL, Action: KeyPress, Mods: KeyModSuper, Source: 7},
				{Kind: InputEventKindText, Mods: KeyModSuper, Rune: 'l', NormalText: true, Source: 7},
			},
		},
		{
			name: "text interleaved between two key actions keeps native order",
			record: func(s *InputState) {
				s.appendKeyEvent(KeyA, KeyPress, 0, 1)
				s.appendTextInput('a', 0, true, 1)
				s.appendKeyEvent(KeyA, KeyRelease, 0, 2)
				s.appendKeyEvent(KeyB, KeyPress, 0, 3)
				s.appendTextInput('b', 0, true, 3)
			},
			want: []InputEvent{
				{Kind: InputEventKindKey, Key: KeyA, Action: KeyPress, Source: 1},
				{Kind: InputEventKindText, Rune: 'a', NormalText: true, Source: 1},
				{Kind: InputEventKindKey, Key: KeyA, Action: KeyRelease, Source: 2},
				{Kind: InputEventKindKey, Key: KeyB, Action: KeyPress, Source: 3},
				{Kind: InputEventKindText, Rune: 'b', NormalText: true, Source: 3},
			},
		},
		{
			name: "one key action commits several code points",
			record: func(s *InputState) {
				s.appendKeyEvent(KeyEnter, KeyPress, 0, 11)
				s.appendTextInput('あ', 0, true, 11)
				s.appendTextInput('い', 0, true, 11)
				s.appendTextInput('う', 0, true, 11)
			},
			want: []InputEvent{
				{Kind: InputEventKindKey, Key: KeyEnter, Action: KeyPress, Source: 11},
				{Kind: InputEventKindText, Rune: 'あ', NormalText: true, Source: 11},
				{Kind: InputEventKindText, Rune: 'い', NormalText: true, Source: 11},
				{Kind: InputEventKindText, Rune: 'う', NormalText: true, Source: 11},
			},
		},
		{
			name: "each OS repeat is its own action",
			record: func(s *InputState) {
				s.appendKeyEvent(KeyA, KeyPress, 0, 4)
				s.appendTextInput('a', 0, true, 4)
				s.appendKeyEvent(KeyA, KeyRepeat, 0, 5)
				s.appendTextInput('a', 0, true, 5)
			},
			want: []InputEvent{
				{Kind: InputEventKindKey, Key: KeyA, Action: KeyPress, Source: 4},
				{Kind: InputEventKindText, Rune: 'a', NormalText: true, Source: 4},
				{Kind: InputEventKindKey, Key: KeyA, Action: KeyRepeat, Source: 5},
				{Kind: InputEventKindText, Rune: 'a', NormalText: true, Source: 5},
			},
		},
		{
			name: "AltGr text is normal text despite a Ctrl+Alt mask",
			record: func(s *InputState) {
				s.appendKeyEvent(KeyE, KeyPress, KeyModControl|KeyModAlt, 9)
				s.appendTextInput('€', KeyModControl|KeyModAlt, true, 9)
			},
			want: []InputEvent{
				{Kind: InputEventKindKey, Key: KeyE, Action: KeyPress, Mods: KeyModControl | KeyModAlt, Source: 9},
				{Kind: InputEventKindText, Mods: KeyModControl | KeyModAlt, Rune: '€', NormalText: true, Source: 9},
			},
		},
		{
			name: "standalone input method commit has no source",
			record: func(s *InputState) {
				s.appendTextInput('漢', 0, true, 0)
			},
			want: []InputEvent{
				{Kind: InputEventKindText, Rune: '漢', NormalText: true},
			},
		},
		{
			name: "emoji sequence survives verbatim",
			record: func(s *InputState) {
				for _, r := range []rune{'\U0001F468', '\u200d', '\U0001F469', '\U000E007F'} {
					s.appendTextInput(r, 0, true, 0)
				}
			},
			want: []InputEvent{
				{Kind: InputEventKindText, Rune: '\U0001F468', NormalText: true},
				{Kind: InputEventKindText, Rune: '\u200d', NormalText: true},
				{Kind: InputEventKindText, Rune: '\U0001F469', NormalText: true},
				{Kind: InputEventKindText, Rune: '\U000E007F', NormalText: true},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s InputState
			tc.record(&s)
			if !reflect.DeepEqual(s.InputEvents, tc.want) {
				t.Fatalf("InputEvents = %+v, want %+v", s.InputEvents, tc.want)
			}
		})
	}
}

// TestInputStateCopyAndReset verifies that each update sees exactly the
// observations recorded since the previous one, in native order. Text a
// platform delivers late must keep the source of the key action that produced
// it even though that key action was reported in an earlier update.
func TestInputStateCopyAndReset(t *testing.T) {
	var src, dst InputState

	src.appendKeyEvent(KeyL, KeyPress, KeyModSuper, 21)
	src.copyAndReset(&dst)

	want := []InputEvent{
		{Kind: InputEventKindKey, Key: KeyL, Action: KeyPress, Mods: KeyModSuper, Source: 21},
	}
	if !reflect.DeepEqual(dst.InputEvents, want) {
		t.Fatalf("first update = %+v, want %+v", dst.InputEvents, want)
	}
	if len(src.InputEvents) != 0 {
		t.Fatalf("delta buffer not reset: %+v", src.InputEvents)
	}

	src.appendTextInput('l', KeyModSuper, true, 21)
	src.appendTextInput('x', 0, true, 0)
	src.copyAndReset(&dst)

	want = []InputEvent{
		{Kind: InputEventKindText, Mods: KeyModSuper, Rune: 'l', NormalText: true, Source: 21},
		{Kind: InputEventKindText, Rune: 'x', NormalText: true},
	}
	if !reflect.DeepEqual(dst.InputEvents, want) {
		t.Fatalf("second update = %+v, want %+v", dst.InputEvents, want)
	}

	src.copyAndReset(&dst)
	if len(dst.InputEvents) != 0 {
		t.Fatalf("third update = %+v, want no observations", dst.InputEvents)
	}
}