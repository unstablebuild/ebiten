// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2022 The Ebitengine Authors

package glfw

// InitHint is a no-op on Windows. The Unix implementation forwards to
// glfwInitHint (used for Cocoa menubar hints). ui_glfw.go calls it
// unconditionally; providing this stub lets GOOS=windows compile.
func InitHint(hint Hint, value int) {}
