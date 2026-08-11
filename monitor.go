// Copyright 2023 The Ebitengine Authors
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

package ebiten

import (
	"github.com/hajimehoshi/ebiten/v2/internal/ui"
)

// MonitorType represents a monitor available to the system.
type MonitorType ui.Monitor

// Name returns the monitor's name. On Linux, this reports the monitors in xrandr format.
// On Windows, this reports "Generic PnP Monitor" for all monitors.
func (m *MonitorType) Name() string {
	return (*ui.Monitor)(m).Name()
}

// DeviceScaleFactor returns the device scale factor of the monitor.
//
// DeviceScaleFactor returns a meaningful value on high-DPI display environment,
// otherwise DeviceScaleFactor returns 1.
//
// DeviceScaleFactor might panic on init function on some devices like Android.
// Then, it is not recommended to call DeviceScaleFactor from init functions.
func (m *MonitorType) DeviceScaleFactor() float64 {
	return (*ui.Monitor)(m).DeviceScaleFactor()
}

// Size returns the size of the monitor in device-independent pixels.
// This is the same as the screen size in fullscreen mode.
// The returned value can be given to SetSize function if the perfectly fit fullscreen is needed.
//
// On mobiles, Size returns (0, 0) so far.
//
// Size's use cases are limited. If you are making a fullscreen application, you can use RunGame and
// the Game interface's Layout function instead. If you are making a not-fullscreen application but the application's
// behavior depends on the monitor size, Size is useful.
func (m *MonitorType) Size() (int, int) {
	return (*ui.Monitor)(m).Size()
}

// RefreshRate returns the refresh rate of the monitor's current video
// mode in Hz, or 0 when the platform does not expose it (browsers,
// mobiles, and consoles) or it is unknown.
//
// On displays with a variable refresh rate, RefreshRate reports the
// rate of the current video mode, which is typically the maximum rate.
func (m *MonitorType) RefreshRate() int {
	return (*ui.Monitor)(m).RefreshRate()
}

// SetMonitorChangedCallback sets a function invoked when the monitor
// the window is on changes, either because the window moved onto a
// different monitor or because the system's monitor configuration
// changed (a monitor was connected or disconnected, or changed video
// mode). The monitor argument may be nil when it cannot be determined.
//
// The callback is invoked on the main thread during event processing;
// it must not block and must not call functions that dispatch to the
// main thread. Passing nil removes the callback.
//
// SetMonitorChangedCallback is concurrent-safe. On browsers, mobiles,
// and consoles the callback is never invoked.
func SetMonitorChangedCallback(f func(monitor *MonitorType)) {
	if f == nil {
		ui.Get().SetMonitorChangedCallback(nil)
		return
	}
	ui.Get().SetMonitorChangedCallback(func(m *ui.Monitor) {
		f((*MonitorType)(m))
	})
}

// Monitor returns the current monitor.
func Monitor() *MonitorType {
	m, ok := ui.Get().Monitor()
	if !ok {
		return nil
	}
	return (*MonitorType)(m)
}

// SetMonitor sets the monitor that the window should be on. This can be called before or after Run.
func SetMonitor(monitor *MonitorType) {
	ui.Get().Window().SetMonitor((*ui.Monitor)(monitor))
}

// AppendMonitors returns the monitors reported by the system.
// On desktop platforms, there will always be at least one monitor appended and the first monitor in the slice will be the primary monitor.
// Any monitors added or removed will show up with subsequent calls to this function.
func AppendMonitors(monitors []*MonitorType) []*MonitorType {
	// TODO: This is not an efficient operation. It would be best if we could directly pass monitors directly into `ui.AppendMonitors`.
	for _, m := range ui.Get().AppendMonitors(nil) {
		monitors = append(monitors, (*MonitorType)(m))
	}
	return monitors
}
