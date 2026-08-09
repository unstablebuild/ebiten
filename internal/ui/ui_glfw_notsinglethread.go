// Copyright 2020 The Ebiten Authors
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

//go:build !android && !ios && !js && !nintendosdk && !playstation5 && !ebitenginesinglethread && !ebitensinglethread

package ui

import (
	"fmt"
	"image"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2/internal/glfw"
	"github.com/hajimehoshi/ebiten/v2/internal/microsoftgdk"
)

func (u *UserInterface) updateInputState() error {
	var err error
	u.mainThread.Call(func() {
		err = u.updateInputStateImpl()
	})
	return err
}

func (u *UserInterface) KeyName(key Key) string {
	if !u.isRunning() {
		return ""
	}

	gk, ok := uiKeyToGLFWKey[key]
	if !ok {
		return ""
	}

	var name string
	u.mainThread.Call(func() {
		if u.isTerminated() {
			return
		}
		n, err := glfw.GetKeyName(gk, 0)
		if err != nil {
			u.setError(err)
			return
		}
		name = n
	})
	return name
}

// Monitor returns the window's current monitor. Returns nil if there is no current monitor yet.
func (u *UserInterface) Monitor() (*Monitor, bool) {
	if !u.isRunning() {
		m := u.getInitMonitor()
		return m, m != nil
	}
	var monitor *Monitor
	var ok bool
	u.mainThread.Call(func() {
		if u.isTerminated() {
			return
		}
		m, isOk, err := u.currentMonitor()
		if err != nil {
			u.setError(err)
			return
		}
		monitor = m
		ok = isOk
	})
	return monitor, ok
}

func (u *UserInterface) IsFullscreen() bool {
	if microsoftgdk.IsXbox() {
		return false
	}

	if u.isTerminated() {
		return false
	}
	if !u.isRunning() {
		return u.isInitFullscreen()
	}
	var fullscreen bool
	u.mainThread.Call(func() {
		if u.isTerminated() {
			return
		}
		b, err := u.isFullscreen()
		if err != nil {
			u.setError(err)
			return
		}
		fullscreen = b
	})
	return fullscreen
}

func (u *UserInterface) SetWindowBackgroundBlur(radius int) {
	if u.isTerminated() {
		return
	}
	if !u.isRunning() {
		u.setInitBackgroundBlur(radius)
		return
	}
	u.mainThread.Call(func() {
		err := u.setNativeBackgroundBlur(radius)
		if err != nil {
			u.setError(err)
			return
		}
	})
}

func (u *UserInterface) SetFullscreen(fullscreen bool) {
	if microsoftgdk.IsXbox() {
		return
	}

	if u.isTerminated() {
		return
	}
	if !u.isRunning() {
		u.setInitFullscreen(fullscreen)
		return
	}

	u.mainThread.Call(func() {
		if u.isTerminated() {
			return
		}
		f, err := u.isFullscreen()
		if err != nil {
			u.setError(err)
			return
		}
		if f == fullscreen {
			return
		}
		if err := u.setFullscreen(fullscreen); err != nil {
			u.setError(err)
			return
		}
	})
}

func (u *UserInterface) IsFocused() bool {
	if !u.isRunning() {
		return false
	}

	var focused bool
	u.mainThread.Call(func() {
		if u.isTerminated() {
			return
		}
		a, err := u.window.GetAttrib(glfw.Focused)
		if err != nil {
			u.setError(err)
			return
		}
		focused = a == glfw.True
	})
	return focused
}

func (u *UserInterface) SetFPSMode(mode FPSModeType) {
	if u.isTerminated() {
		return
	}
	if !u.isRunning() {
		u.m.Lock()
		defer u.m.Unlock()
		u.fpsMode = mode
		return
	}

	u.mainThread.Call(func() {
		if u.isTerminated() {
			return
		}
		if !u.fpsModeInited {
			u.fpsMode = mode
			return
		}
		if err := u.setFPSMode(mode); err != nil {
			u.setError(err)
			return
		}
	})
}

func (u *UserInterface) CursorMode() CursorMode {
	if u.isTerminated() {
		return 0
	}
	if !u.isRunning() {
		return u.getInitCursorMode()
	}

	var mode int
	u.mainThread.Call(func() {
		if u.isTerminated() {
			return
		}
		m, err := u.window.GetInputMode(glfw.CursorMode)
		if err != nil {
			u.setError(err)
			return
		}
		mode = m
	})

	var v CursorMode
	switch mode {
	case glfw.CursorNormal:
		v = CursorModeVisible
	case glfw.CursorHidden:
		v = CursorModeHidden
	case glfw.CursorDisabled:
		v = CursorModeCaptured
	default:
		panic(fmt.Sprintf("ui: invalid GLFW cursor mode: %d", mode))
	}
	return v
}

func (u *UserInterface) SetCursorMode(mode CursorMode) {
	if u.isTerminated() {
		return
	}
	if !u.isRunning() {
		u.setInitCursorMode(mode)
		return
	}
	u.mainThread.Call(func() {
		if u.isTerminated() {
			return
		}
		if err := u.window.SetInputMode(glfw.CursorMode, driverCursorModeToGLFWCursorMode(mode)); err != nil {
			u.setError(err)
			return
		}
		if mode == CursorModeVisible {
			if err := u.window.SetCursor(glfwSystemCursors[u.getCursorShape()]); err != nil {
				u.setError(err)
				return
			}
		}
	})
}

func (u *UserInterface) CursorShape() CursorShape {
	return u.getCursorShape()
}

func (u *UserInterface) SetCursorShape(shape CursorShape) {
	if u.isTerminated() {
		return
	}

	old := u.setCursorShape(shape)
	if old == shape {
		return
	}
	if !u.isRunning() {
		return
	}
	u.mainThread.Call(func() {
		if u.isTerminated() {
			return
		}
		if err := u.window.SetCursor(glfwSystemCursors[shape]); err != nil {
			u.setError(err)
			return
		}
	})
}

func (u *UserInterface) updateGame() error {
	var unfocused bool

	// On Windows, the focusing state might be always false (#987).
	// On Windows, even if a window is in another workspace, vsync seems to work.
	// Then let's assume the window is always 'focused' as a workaround.
	if runtime.GOOS != "windows" {
		a, err := u.window.GetAttrib(glfw.Focused)
		if err != nil {
			return err
		}
		unfocused = a == glfw.False
	}

	var t1, t2 time.Time

	if unfocused {
		t1 = time.Now()
	}

	var outsideWidth, outsideHeight float64
	deviceScaleFactor := 1.0
	var err error
	if u.mainThread.Call(func() {
		outsideWidth, outsideHeight, err = u.update()
		if err != nil {
			return
		}
		m, ok, err := u.currentMonitor()
		if err != nil {
			return
		}
		if ok {
			deviceScaleFactor = m.DeviceScaleFactor()
		}
	}); err != nil {
		return err
	}

	if err := u.context.updateFrame(u.graphicsDriver, outsideWidth, outsideHeight, deviceScaleFactor, u); err != nil {
		return err
	}

	u.bufferOnceSwappedOnce.Do(func() {
		u.mainThread.Call(func() {
			u.bufferOnceSwapped = true
		})
	})

	if unfocused {
		t2 = time.Now()
	}

	// When a window is not focused or in another space, SwapBuffers might return immediately and CPU might be busy.
	// Mitigate this by sleeping (#982, #2521).
	if unfocused {
		d := t2.Sub(t1)
		const wait = time.Second / 60
		if d < wait {
			time.Sleep(wait - d)
		}
	}

	return nil
}

func (u *UserInterface) updateIconIfNeeded() error {
	// In the fullscreen mode, SetIcon fails (#1578).
	f, err := u.isFullscreen()
	if err != nil {
		return err
	}
	if f {
		return nil
	}

	imgs := u.getAndResetIconImages()
	// A 0-size slice and nil are distinguished here.
	// A 0-size slice means a user indicates to reset the icon.
	// On the other hand, nil means a user didn't update the icon state.
	if imgs == nil {
		return nil
	}

	var newImgs []image.Image
	if len(imgs) > 0 {
		newImgs = make([]image.Image, len(imgs))
	}
	for i, img := range imgs {
		// TODO: If img is not *ebiten.Image, this converting is not necessary.
		// However, this package cannot refer *ebiten.Image due to the package
		// dependencies.

		b := img.Bounds()
		rgba := image.NewRGBA(b)
		for j := b.Min.Y; j < b.Max.Y; j++ {
			for i := b.Min.X; i < b.Max.X; i++ {
				rgba.Set(i, j, img.At(i, j))
			}
		}
		newImgs[i] = rgba
	}

	// Catch a possible error at 'At' (#2647).
	if err := u.error(); err != nil {
		return err
	}

	u.mainThread.Call(func() {
		err = u.window.SetIcon(newImgs)
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *UserInterface) RunOnMainThread(f func()) {
	u.mainThread.Call(f)
}
