// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 The Ebitengine Authors

//go:build freebsd || linux || netbsd || openbsd

#include "internal_unix.h"

#include <dlfcn.h>
#include <stddef.h>

#define _GLFW_X11_DEFINE_SYMBOL(name) \
    _glfw_pfn_t_##name* _glfw_dyn_##name = NULL;
_GLFW_X11_SYMBOLS(_GLFW_X11_DEFINE_SYMBOL)
#undef _GLFW_X11_DEFINE_SYMBOL

static void* _glfwX11LibraryHandle = NULL;

GLFWbool _glfwLoadX11Library(void)
{
    if (_glfwX11LibraryHandle)
        return GLFW_TRUE;

#if defined(__CYGWIN__)
    void* handle = _glfw_dlopen("libX11-6.so");
#elif defined(__OpenBSD__) || defined(__NetBSD__)
    void* handle = _glfw_dlopen("libX11.so");
#else
    void* handle = _glfw_dlopen("libX11.so.6");
    if (!handle)
        handle = _glfw_dlopen("libX11.so");
#endif
    if (!handle)
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "X11: Failed to load libX11: %s", dlerror());
        return GLFW_FALSE;
    }

    const char* missing = NULL;

#define _GLFW_X11_LOAD_SYMBOL(name)                                      \
    _glfw_dyn_##name = (_glfw_pfn_t_##name*) _glfw_dlsym(handle, #name); \
    if (!_glfw_dyn_##name && !missing)                                   \
        missing = #name;
    _GLFW_X11_SYMBOLS(_GLFW_X11_LOAD_SYMBOL)
#undef _GLFW_X11_LOAD_SYMBOL

    if (missing)
    {
        _glfw_dlclose(handle);
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "X11: libX11 is missing the symbol %s", missing);
        return GLFW_FALSE;
    }

    _glfwX11LibraryHandle = handle;
    return GLFW_TRUE;
}
