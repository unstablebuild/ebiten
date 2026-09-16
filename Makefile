# Tests for the parts of the fork that diverge from upstream Ebitengine.
# Upstream's own suite is not run here: it needs a display and a GPU.

.PHONY: test test-e2e build-windows

# The keyboard-layout tests link Carbon to load layouts the user has not
# enabled, so they sit behind a build tag rather than in every build.
test:
	go test -tags glfwlayout ./internal/ui/ ./internal/glfw/ ./

# Needs a headless X server and setxkbmap: apt install xvfb x11-xkb-utils.
test-e2e:
	go test -tags e2e ./internal/glfw/

build-windows:
	GOOS=windows go build ./internal/glfw/