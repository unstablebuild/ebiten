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

//go:build !ebitenginesinglethread && !ebitensinglethread

package graphicscommand

import (
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2/internal/debug"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
)

func SetVsyncEnabled(enabled bool, graphicsDriver graphicsdriver.Graphics) {
	if enabled {
		atomic.StoreInt32(&vsyncEnabled, 1)
	} else {
		atomic.StoreInt32(&vsyncEnabled, 0)
	}

	runOnRenderThread(func() {
		graphicsDriver.SetVsyncEnabled(enabled)
	}, true)
}

// Flush flushes the command queue.
func (q *commandQueue) Flush(graphicsDriver graphicsdriver.Graphics, endFrame bool) error {
	if err := q.err.Load(); err != nil {
		return err.(error)
	}

	var sync bool
	// Disable asynchronous rendering when vsync is on, as this causes a rendering delay (#2822).
	if endFrame && atomic.LoadInt32(&vsyncEnabled) != 0 {
		sync = true
	}
	if !sync {
		for _, c := range q.commands {
			if c.NeedsSync() {
				sync = true
				break
			}
		}
	}

	logger := debug.SwitchLogger()

	var flushErr error
	runOnRenderThread(func() {
		defer logger.Flush()

		if err := q.flush(graphicsDriver, endFrame, logger); err != nil {
			if sync {
				flushErr = err
				return
			}
			q.err.Store(err)
			return
		}

		theCommandQueueManager.putCommandQueue(q)
	}, sync)

	if sync && flushErr != nil {
		return flushErr
	}

	return nil
}
