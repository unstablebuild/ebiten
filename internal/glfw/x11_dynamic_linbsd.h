// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build freebsd || linux || netbsd || openbsd

// On Linux, Xlib is loaded with dlopen instead of being linked against, so
// that a build that never opens a window (a terminal UI, a headless server)
// starts on a machine that has no graphical libraries installed at all. The
// BSDs still link libX11 through pkg-config in build_bsd.go, which is where
// their -I flags come from; the loader below runs there too but only finds
// the already-mapped library.
//
// The extension libraries in x11_platform_linbsd.h keep hand-written PFN_
// typedefs in _glfw.x11. Xlib is handled differently on purpose: with close
// to a hundred entry points, __typeof__ on the real declaration is the only
// way to keep the pointer types from drifting from the ABI, and the pointers
// live outside _glfw so that the memset in glfwTerminate does not force a
// reload on the next glfwInit.

#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/Xresource.h>
#include <X11/XKBlib.h>

#include "x11_symbols_linbsd.h"

// The typedefs must come before the renames below: __typeof__ has to see the
// Xlib prototype, not the pointer it is about to become.
#define _GLFW_X11_DECLARE_SYMBOL(name)       \
    typedef __typeof__(name) _glfw_pfn_t_##name; \
    extern _glfw_pfn_t_##name* _glfw_dyn_##name;
_GLFW_X11_SYMBOLS(_GLFW_X11_DECLARE_SYMBOL)
#undef _GLFW_X11_DECLARE_SYMBOL

// Load libX11 and resolve every symbol listed above. Idempotent, so that
// glfwInit after glfwTerminate does not reopen the library.
GLFWbool _glfwLoadX11Library(void);

// One rename per entry of _GLFW_X11_SYMBOLS. The toolchain keeps the two in
// sync: a listed symbol used without its rename is an undefined reference,
// a rename without its entry is an undeclared identifier.
#define XAllocClassHint            _glfw_dyn_XAllocClassHint
#define XAllocSizeHints            _glfw_dyn_XAllocSizeHints
#define XAllocWMHints              _glfw_dyn_XAllocWMHints
#define XChangeProperty            _glfw_dyn_XChangeProperty
#define XChangeWindowAttributes    _glfw_dyn_XChangeWindowAttributes
#define XCheckIfEvent              _glfw_dyn_XCheckIfEvent
#define XCheckTypedWindowEvent     _glfw_dyn_XCheckTypedWindowEvent
#define XCloseDisplay              _glfw_dyn_XCloseDisplay
#define XCloseIM                   _glfw_dyn_XCloseIM
#define XConvertSelection          _glfw_dyn_XConvertSelection
#define XCreateColormap            _glfw_dyn_XCreateColormap
#define XCreateFontCursor          _glfw_dyn_XCreateFontCursor
#define XCreateIC                  _glfw_dyn_XCreateIC
#define XCreateRegion              _glfw_dyn_XCreateRegion
#define XCreateWindow              _glfw_dyn_XCreateWindow
#define XDefineCursor              _glfw_dyn_XDefineCursor
#define XDeleteContext             _glfw_dyn_XDeleteContext
#define XDeleteProperty            _glfw_dyn_XDeleteProperty
#define XDestroyIC                 _glfw_dyn_XDestroyIC
#define XDestroyRegion             _glfw_dyn_XDestroyRegion
#define XDestroyWindow             _glfw_dyn_XDestroyWindow
#define XDisplayKeycodes           _glfw_dyn_XDisplayKeycodes
#define XEventsQueued              _glfw_dyn_XEventsQueued
#define XFilterEvent               _glfw_dyn_XFilterEvent
#define XFindContext               _glfw_dyn_XFindContext
#define XFlush                     _glfw_dyn_XFlush
#define XFree                      _glfw_dyn_XFree
#define XFreeColormap              _glfw_dyn_XFreeColormap
#define XFreeCursor                _glfw_dyn_XFreeCursor
#define XFreeEventData             _glfw_dyn_XFreeEventData
#define XGetErrorText              _glfw_dyn_XGetErrorText
#define XGetEventData              _glfw_dyn_XGetEventData
#define XGetICValues               _glfw_dyn_XGetICValues
#define XGetIMValues               _glfw_dyn_XGetIMValues
#define XGetInputFocus             _glfw_dyn_XGetInputFocus
#define XGetKeyboardMapping        _glfw_dyn_XGetKeyboardMapping
#define XGetScreenSaver            _glfw_dyn_XGetScreenSaver
#define XGetSelectionOwner         _glfw_dyn_XGetSelectionOwner
#define XGetVisualInfo             _glfw_dyn_XGetVisualInfo
#define XGetWMNormalHints          _glfw_dyn_XGetWMNormalHints
#define XGetWindowAttributes       _glfw_dyn_XGetWindowAttributes
#define XGetWindowProperty         _glfw_dyn_XGetWindowProperty
#define XGrabPointer               _glfw_dyn_XGrabPointer
#define XIconifyWindow             _glfw_dyn_XIconifyWindow
#define XInitThreads               _glfw_dyn_XInitThreads
#define XInternAtom                _glfw_dyn_XInternAtom
#define XLookupString              _glfw_dyn_XLookupString
#define XMapRaised                 _glfw_dyn_XMapRaised
#define XMapWindow                 _glfw_dyn_XMapWindow
#define XMoveResizeWindow          _glfw_dyn_XMoveResizeWindow
#define XMoveWindow                _glfw_dyn_XMoveWindow
#define XNextEvent                 _glfw_dyn_XNextEvent
#define XOpenDisplay               _glfw_dyn_XOpenDisplay
#define XOpenIM                    _glfw_dyn_XOpenIM
#define XPeekEvent                 _glfw_dyn_XPeekEvent
#define XPending                   _glfw_dyn_XPending
#define XQLength                   _glfw_dyn_XQLength
#define XQueryExtension            _glfw_dyn_XQueryExtension
#define XQueryPointer              _glfw_dyn_XQueryPointer
#define XRaiseWindow               _glfw_dyn_XRaiseWindow
#define XResizeWindow              _glfw_dyn_XResizeWindow
#define XResourceManagerString     _glfw_dyn_XResourceManagerString
#define XSaveContext               _glfw_dyn_XSaveContext
#define XSelectInput               _glfw_dyn_XSelectInput
#define XSendEvent                 _glfw_dyn_XSendEvent
#define XSetClassHint              _glfw_dyn_XSetClassHint
#define XSetErrorHandler           _glfw_dyn_XSetErrorHandler
#define XSetICFocus                _glfw_dyn_XSetICFocus
#define XSetInputFocus             _glfw_dyn_XSetInputFocus
#define XSetLocaleModifiers        _glfw_dyn_XSetLocaleModifiers
#define XSetScreenSaver            _glfw_dyn_XSetScreenSaver
#define XSetSelectionOwner         _glfw_dyn_XSetSelectionOwner
#define XSetWMHints                _glfw_dyn_XSetWMHints
#define XSetWMNormalHints          _glfw_dyn_XSetWMNormalHints
#define XSetWMProtocols            _glfw_dyn_XSetWMProtocols
#define XSupportsLocale            _glfw_dyn_XSupportsLocale
#define XSync                      _glfw_dyn_XSync
#define XTranslateCoordinates      _glfw_dyn_XTranslateCoordinates
#define XUndefineCursor            _glfw_dyn_XUndefineCursor
#define XUngrabPointer             _glfw_dyn_XUngrabPointer
#define XUnmapWindow               _glfw_dyn_XUnmapWindow
#define XUnsetICFocus              _glfw_dyn_XUnsetICFocus
#define XWarpPointer               _glfw_dyn_XWarpPointer
#define XkbFreeKeyboard            _glfw_dyn_XkbFreeKeyboard
#define XkbFreeNames               _glfw_dyn_XkbFreeNames
#define XkbGetMap                  _glfw_dyn_XkbGetMap
#define XkbGetNames                _glfw_dyn_XkbGetNames
#define XkbGetState                _glfw_dyn_XkbGetState
#define XkbLookupKeySym            _glfw_dyn_XkbLookupKeySym
#define XkbQueryExtension          _glfw_dyn_XkbQueryExtension
#define XkbSelectEventDetails      _glfw_dyn_XkbSelectEventDetails
#define XkbSetDetectableAutoRepeat _glfw_dyn_XkbSetDetectableAutoRepeat
#define XmbSetWMProperties         _glfw_dyn_XmbSetWMProperties
#define XrmDestroyDatabase         _glfw_dyn_XrmDestroyDatabase
#define XrmGetResource             _glfw_dyn_XrmGetResource
#define XrmGetStringDatabase       _glfw_dyn_XrmGetStringDatabase
#define XrmInitialize              _glfw_dyn_XrmInitialize
#define XrmUniqueQuark             _glfw_dyn_XrmUniqueQuark
#define Xutf8LookupString          _glfw_dyn_Xutf8LookupString
#define Xutf8SetWMProperties       _glfw_dyn_Xutf8SetWMProperties
#define XwcLookupString            _glfw_dyn_XwcLookupString
