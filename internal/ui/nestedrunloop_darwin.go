// Copyright 2026 The Ebitengine Authors
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

//go:build !ios

package ui

import (
	"fmt"

	"github.com/ebitengine/purego/objc"

	"github.com/hajimehoshi/ebiten/v2/internal/cocoa"
	"github.com/hajimehoshi/ebiten/v2/internal/thread"
)

var (
	class_NSTimer   = objc.GetClass("NSTimer")
	class_NSRunLoop = objc.GetClass("NSRunLoop")
)

var (
	sel_nestedRunLoopTick     = objc.RegisterName("ebitengineNestedRunLoopTick:")
	sel_timerWithTimeInterval = objc.RegisterName("timerWithTimeInterval:target:selector:userInfo:repeats:")
	sel_mainRunLoop           = objc.RegisterName("mainRunLoop")
	sel_addTimerForMode       = objc.RegisterName("addTimer:forMode:")
)

// nestedRunLoopModes are the run-loop modes AppKit uses for nested
// event-tracking loops: menu tracking, window dragging and live
// resizing use the event-tracking mode, modal panels the modal mode.
// Regular event polling pumps only the default mode, so a timer
// scheduled on these modes never fires during normal operation.
var nestedRunLoopModes = []string{
	"NSEventTrackingRunLoopMode",
	"NSModalPanelRunLoopMode",
}

// installNestedRunLoopTicker schedules a timer that keeps producing
// frames while AppKit runs a nested event-tracking run loop.
// glfw.PollEvents does not return until such a loop ends, which stalls
// loopGame; the timer fires inside the nested loop and runs the game
// tick from there, so animations and background work continue while a
// menu is open or the window is being dragged or resized.
func (u *UserInterface) installNestedRunLoopTicker() error {
	class, err := objc.RegisterClass(
		"EbitengineNestedRunLoopTicker",
		objc.GetClass("NSObject"),
		nil,
		nil,
		[]objc.MethodDef{
			{
				Cmd: sel_nestedRunLoopTick,
				Fn: func(id objc.ID, cmd objc.SEL, timer objc.ID) {
					u.tickNestedRunLoop()
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("ui: registering EbitengineNestedRunLoopTicker failed: %w", err)
	}

	target := objc.ID(class).Send(sel_alloc).Send(sel_init)
	timer := objc.ID(class_NSTimer).Send(sel_timerWithTimeInterval,
		1.0/60, target, sel_nestedRunLoopTick, objc.ID(0), true)
	runLoop := objc.ID(class_NSRunLoop).Send(sel_mainRunLoop)
	for _, mode := range nestedRunLoopModes {
		nsMode := cocoa.NSString_alloc().InitWithUTF8String(mode)
		runLoop.Send(sel_addTimerForMode, timer, nsMode.ID)
	}
	return nil
}

// tickNestedRunLoop produces a single frame from inside a nested run
// loop. pollingEvents is true exactly while u.update is parked in
// glfw.PollEvents or glfw.WaitEvents, which guarantees no frame is in
// flight on this thread, so calling updateFrame here cannot re-enter
// the frame machinery. Restricted to the single-threaded scheme: in the
// multi-threaded one this timer fires on the main thread while the
// game loop owns frames on another thread.
func (u *UserInterface) tickNestedRunLoop() {
	if u.nestedTickInProgress || !u.pollingEvents {
		return
	}
	if _, ok := u.mainThread.(*thread.NoopThread); !ok {
		return
	}
	if !u.isRunning() || u.isTerminated() {
		return
	}
	u.nestedTickInProgress = true
	defer func() { u.nestedTickInProgress = false }()

	outsideWidth, outsideHeight, err := u.outsideSize()
	if err != nil {
		u.setError(err)
		return
	}
	deviceScaleFactor := 1.0
	m, ok, err := u.currentMonitor()
	if err != nil {
		u.setError(err)
		return
	}
	if ok {
		deviceScaleFactor = m.DeviceScaleFactor()
	}
	if err := u.context.updateFrame(u.graphicsDriver, outsideWidth, outsideHeight, deviceScaleFactor, u); err != nil {
		u.setError(err)
	}
}