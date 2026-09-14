// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build freebsd || linux || netbsd || openbsd

// Package x11symbols exists for its test. It is the one place that links
// against libX11, which glfw itself no longer does, so that the entries of
// _GLFW_X11_SYMBOLS are checked against the library at link time rather
// than at the first glfwInit. Nothing imports it.
package x11symbols

// The list is included on its own so the table references the real Xlib
// declarations, not the dlsym'd pointers glfw renames them to. A name that
// libX11 does not export is an undefined reference when the test binary
// links.

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
// #include <X11/Xutil.h>
// #include <X11/Xresource.h>
// #include <X11/XKBlib.h>
// #include "../x11_symbols_linbsd.h"
//
// #define _GLFW_X11_REF(name) (void*) name,
// void* const x11SymbolRefs[] = { _GLFW_X11_SYMBOLS(_GLFW_X11_REF) };
// #undef _GLFW_X11_REF
import "C"

import "unsafe"

// Refs returns the address libX11 exports for every entry of
// _GLFW_X11_SYMBOLS, in list order.
func Refs() []unsafe.Pointer {
	return unsafe.Slice(&C.x11SymbolRefs[0], len(C.x11SymbolRefs))
}
