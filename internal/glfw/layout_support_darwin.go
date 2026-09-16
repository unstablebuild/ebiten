// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

// Support for the keyboard-layout tests. cgo is not allowed in _test.go
// files, so the bridge lives here, behind a build tag that keeps the extra
// Carbon linkage out of ordinary builds. Run the tests with
// `go test -tags glfwlayout ./internal/glfw/`.

//go:build darwin && glfwlayout

package glfw

// #cgo LDFLAGS: -framework Carbon
//
// #include <stdint.h>
// #include <stdlib.h>
// #include <Carbon/Carbon.h>
//
// // GLFWbool is an int, so the internal prototype is declared here rather
// // than pulling the platform headers into this translation unit.
// extern uint32_t _glfwTranslateScancodeNS(const UCKeyboardLayout* layout, UInt8 kbdType,
//                                          int scancode, int shift);
//
// static uint32_t translate(const UCKeyboardLayout* layout, int scancode, int shift) {
//     return _glfwTranslateScancodeNS(layout, LMGetKbdType(), scancode, shift);
// }
//
// // Returns the key layout Apple ships under the given input source ID,
// // whether or not the user has enabled it, or NULL if it is unavailable.
// // The layout data is retained for the lifetime of the process.
// static const UCKeyboardLayout* layoutForID(const char* identifier) {
//     CFStringRef sid = CFStringCreateWithCString(kCFAllocatorDefault, identifier,
//                                                 kCFStringEncodingUTF8);
//     if (!sid) {
//         return NULL;
//     }
//     const void* keys[] = { kTISPropertyInputSourceID };
//     const void* values[] = { sid };
//     CFDictionaryRef filter = CFDictionaryCreate(kCFAllocatorDefault, keys, values, 1,
//                                                 &kCFTypeDictionaryKeyCallBacks,
//                                                 &kCFTypeDictionaryValueCallBacks);
//     CFArrayRef list = TISCreateInputSourceList(filter, true);
//     CFRelease(filter);
//     CFRelease(sid);
//     if (!list) {
//         return NULL;
//     }
//     if (CFArrayGetCount(list) == 0) {
//         CFRelease(list);
//         return NULL;
//     }
//     TISInputSourceRef source = (TISInputSourceRef) CFArrayGetValueAtIndex(list, 0);
//     CFDataRef data = (CFDataRef) TISGetInputSourceProperty(source,
//                                                            kTISPropertyUnicodeKeyLayoutData);
//     if (!data) {
//         CFRelease(list);
//         return NULL;
//     }
//     CFRetain(data);
//     CFRelease(list);
//     return (const UCKeyboardLayout*) CFDataGetBytePtr(data);
// }
import "C"

import "unsafe"

// keyboardLayout is an opaque handle to a UCKeyboardLayout.
type keyboardLayout struct {
	ptr *C.UCKeyboardLayout
}

// loadKeyboardLayout returns the layout Apple ships under
// com.apple.keylayout.<name>, or ok == false when this machine does not have
// it. The user's own layout selection is left untouched.
func loadKeyboardLayout(name string) (layout keyboardLayout, ok bool) {
	identifier := C.CString("com.apple.keylayout." + name)
	defer C.free(unsafe.Pointer(identifier))

	ptr := C.layoutForID(identifier)
	if ptr == nil {
		return keyboardLayout{}, false
	}
	return keyboardLayout{ptr: ptr}, true
}

// translateScancode returns the code point layout produces for scancode at its
// unshifted or shifted level, or 0 if it produces none.
func translateScancode(layout keyboardLayout, scancode int, shift bool) rune {
	var s C.int
	if shift {
		s = 1
	}
	return rune(C.translate(layout.ptr, C.int(scancode), s))
}