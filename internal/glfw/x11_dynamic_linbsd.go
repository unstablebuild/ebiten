// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build freebsd || linux || netbsd || openbsd

package glfw

// int _glfwLoadX11Library(void);
import "C"

import "errors"

// LoadX11LibraryForTesting runs the Xlib loader on its own, so that a test
// with no display can check that libX11 opens by SONAME and exports every
// entry point GLFW resolves.
func LoadX11LibraryForTesting() error {
	if C._glfwLoadX11Library() != 0 {
		return nil
	}
	select {
	case err := <-lastError:
		return err
	default:
		return errors.New("glfw: _glfwLoadX11Library failed without reporting an error")
	}
}
