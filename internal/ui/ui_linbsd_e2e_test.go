// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build (freebsd || (linux && !android) || netbsd || openbsd) && !nintendosdk && !playstation5 && e2e

package ui

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2/internal/glfw"
)

func TestWindowSizeIgnoresWrappedX11Size(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	startXvfb(t)

	// A standalone UserInterface leaves theUI for TestInitWithoutDisplay, and
	// terminating GLFW lets that test initialize it again without a display.
	u := &UserInterface{}
	if err := u.init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	t.Cleanup(func() { _ = glfw.Terminate() })
	u.setRunning(true)
	if err := u.createWindow(); err != nil {
		t.Fatalf("createWindow: %v", err)
	}
	if err := u.registerWindowFramebufferSizeCallback(); err != nil {
		t.Fatalf("registerWindowFramebufferSizeCallback: %v", err)
	}

	wantW, wantH, err := u.outsideSize()
	if err != nil {
		t.Fatalf("outsideSize: %v", err)
	}

	// A window manager that computes a negative client size, such as -8x-39
	// for an output that collapsed to nothing, sends it through
	// ConfigureWindow's unsigned 16-bit fields.
	resizeWindow(t, u, 65528, 65497)
	if w, h, err := u.window.GetSize(); err != nil || w != 65528 || h != 65497 {
		t.Fatalf("X server reports %dx%d (err %v), want the wrapped 65528x65497", w, h, err)
	}

	if w, h, err := u.outsideSize(); err != nil || w != wantW || h != wantH {
		t.Errorf("outsideSize() = %v, %v, %v; want the last valid size %v, %v, nil", w, h, err, wantW, wantH)
	}
	// The window size the API reports, and the fallback outsideSize uses.
	if w, h := u.origWindowWidthInDIP, u.origWindowHeightInDIP; float64(w) != wantW || float64(h) != wantH {
		t.Errorf("origWindowSizeInDIP = %d, %d; want the last valid size %v, %v", w, h, wantW, wantH)
	}

	resizeWindow(t, u, 800, 600)
	m, ok, err := u.currentMonitor()
	if err != nil || !ok {
		t.Fatalf("currentMonitor() = %v, %v, %v", m, ok, err)
	}
	s := m.DeviceScaleFactor()
	wantW, wantH = dipFromGLFWPixel(800, s), dipFromGLFWPixel(600, s)
	if w, h, err := u.outsideSize(); err != nil || w != wantW || h != wantH {
		t.Errorf("outsideSize() after a valid resize = %v, %v, %v; want %v, %v, nil", w, h, err, wantW, wantH)
	}
}

// resizeWindow resizes the X window and returns once GLFW has delivered the
// resulting framebuffer size event to the UserInterface.
func resizeWindow(t *testing.T, u *UserInterface, width, height int) {
	t.Helper()

	delivered := false
	if _, err := u.window.SetFramebufferSizeCallback(func(w *glfw.Window, fw, fh int) {
		u.defaultFramebufferSizeCallback(w, fw, fh)
		delivered = true
	}); err != nil {
		t.Fatalf("SetFramebufferSizeCallback: %v", err)
	}
	defer func() {
		if _, err := u.window.SetFramebufferSizeCallback(u.defaultFramebufferSizeCallback); err != nil {
			t.Errorf("SetFramebufferSizeCallback: %v", err)
		}
	}()

	if err := u.window.SetSize(width, height); err != nil {
		t.Fatalf("SetSize(%d, %d): %v", width, height, err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for !delivered {
		if time.Now().After(deadline) {
			t.Fatalf("no framebuffer size event after SetSize(%d, %d)", width, height)
		}
		if err := glfw.PollEvents(); err != nil {
			t.Fatalf("PollEvents: %v", err)
		}
		time.Sleep(time.Millisecond)
	}
}

// startXvfb runs a headless X server for the duration of the test and points
// DISPLAY at it. Xvfb picks a free display number, so this does not collide
// with the X servers of other packages' tests running in parallel.
func startXvfb(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("Xvfb"); err != nil {
		t.Skip("Xvfb is not installed")
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	cmd := exec.Command("Xvfb", "-displayfd", "3", "-screen", "0", "1024x768x24")
	cmd.ExtraFiles = []*os.File{w}
	err = cmd.Start()
	w.Close()
	if err != nil {
		t.Fatalf("start Xvfb: %v", err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})

	// Xvfb writes the display number once it accepts clients, and the read
	// fails instead if it exits first.
	display, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		t.Fatalf("read the Xvfb display number: %v", err)
	}
	t.Setenv("DISPLAY", ":"+strings.TrimSpace(display))
}
