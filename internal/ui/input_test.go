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

import "testing"

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
			got := len(s.Runes) == 1 && s.Runes[0] == tc.r
			if got != tc.keep {
				t.Fatalf("appendRune(%#U): kept=%v, want kept=%v (runes=%v)",
					tc.r, got, tc.keep, s.Runes)
			}
		})
	}
}