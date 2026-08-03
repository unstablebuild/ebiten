// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2002-2006 Marcus Geelnard
// SPDX-FileCopyrightText: 2006-2019 Camilla Löwy <elmindreda@glfw.org>
// SPDX-FileCopyrightText: 2022 The Ebitengine Authors

//go:build darwin || freebsd || linux || netbsd || openbsd || windows

package glfw

type VidMode struct {
	Width       int
	Height      int
	RedBits     int
	GreenBits   int
	BlueBits    int
	RefreshRate int
}

type Image struct {
	Width  int
	Height int
	Pixels []byte
}

// InputSource identifies a native key action. Every code point the platform
// translated from that action is reported with the same nonzero source, so a
// consumer that handled the key transition can recognize exactly the text it
// must not handle again. Zero means the platform could not establish
// causality, as for a standalone input method commit.
type InputSource uint64

// InputKeyCallback is the key callback that also reports the source of the
// key action.
type InputKeyCallback func(w *Window, key Key, scancode int, action Action, mods ModifierKey, source InputSource)

// InputCharCallback is the character callback that also reports the native
// modifier mask, the platform's classification of the code point as normal
// text, and the source of the key action that produced it.
type InputCharCallback func(w *Window, char rune, mods ModifierKey, normalText bool, source InputSource)
