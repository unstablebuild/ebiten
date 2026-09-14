// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build freebsd || linux || netbsd || openbsd

// Every libX11 entry point GLFW calls, as an X-macro so that one list drives
// the pointer declarations and the dlsym loop in x11_dynamic_linbsd.h/.c as
// well as the link-time check in x11_dynamic_linbsd_test.go. Kept free of any
// other definition so a translation unit can include it without picking up
// the renames.

#define _GLFW_X11_SYMBOLS(X) \
    X(XAllocClassHint) \
    X(XAllocSizeHints) \
    X(XAllocWMHints) \
    X(XChangeProperty) \
    X(XChangeWindowAttributes) \
    X(XCheckIfEvent) \
    X(XCheckTypedWindowEvent) \
    X(XCloseDisplay) \
    X(XCloseIM) \
    X(XConvertSelection) \
    X(XCreateColormap) \
    X(XCreateFontCursor) \
    X(XCreateIC) \
    X(XCreateRegion) \
    X(XCreateWindow) \
    X(XDefineCursor) \
    X(XDeleteContext) \
    X(XDeleteProperty) \
    X(XDestroyIC) \
    X(XDestroyRegion) \
    X(XDestroyWindow) \
    X(XDisplayKeycodes) \
    X(XEventsQueued) \
    X(XFilterEvent) \
    X(XFindContext) \
    X(XFlush) \
    X(XFree) \
    X(XFreeColormap) \
    X(XFreeCursor) \
    X(XFreeEventData) \
    X(XGetErrorText) \
    X(XGetEventData) \
    X(XGetICValues) \
    X(XGetIMValues) \
    X(XGetInputFocus) \
    X(XGetKeyboardMapping) \
    X(XGetScreenSaver) \
    X(XGetSelectionOwner) \
    X(XGetVisualInfo) \
    X(XGetWMNormalHints) \
    X(XGetWindowAttributes) \
    X(XGetWindowProperty) \
    X(XGrabPointer) \
    X(XIconifyWindow) \
    X(XInitThreads) \
    X(XInternAtom) \
    X(XLookupString) \
    X(XMapRaised) \
    X(XMapWindow) \
    X(XMoveResizeWindow) \
    X(XMoveWindow) \
    X(XNextEvent) \
    X(XOpenDisplay) \
    X(XOpenIM) \
    X(XPeekEvent) \
    X(XPending) \
    X(XQLength) \
    X(XQueryExtension) \
    X(XQueryPointer) \
    X(XRaiseWindow) \
    X(XResizeWindow) \
    X(XResourceManagerString) \
    X(XSaveContext) \
    X(XSelectInput) \
    X(XSendEvent) \
    X(XSetClassHint) \
    X(XSetErrorHandler) \
    X(XSetICFocus) \
    X(XSetInputFocus) \
    X(XSetLocaleModifiers) \
    X(XSetScreenSaver) \
    X(XSetSelectionOwner) \
    X(XSetWMHints) \
    X(XSetWMNormalHints) \
    X(XSetWMProtocols) \
    X(XSupportsLocale) \
    X(XSync) \
    X(XTranslateCoordinates) \
    X(XUndefineCursor) \
    X(XUngrabPointer) \
    X(XUnmapWindow) \
    X(XUnsetICFocus) \
    X(XWarpPointer) \
    X(XkbFreeKeyboard) \
    X(XkbFreeNames) \
    X(XkbGetMap) \
    X(XkbGetNames) \
    X(XkbGetState) \
    X(XkbKeycodeToKeysym) \
    X(XkbQueryExtension) \
    X(XkbSelectEventDetails) \
    X(XkbSetDetectableAutoRepeat) \
    X(XmbSetWMProperties) \
    X(XrmDestroyDatabase) \
    X(XrmGetResource) \
    X(XrmGetStringDatabase) \
    X(XrmInitialize) \
    X(XrmUniqueQuark) \
    X(Xutf8LookupString) \
    X(Xutf8SetWMProperties) \
    X(XwcLookupString)
