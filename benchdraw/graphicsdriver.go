package benchdraw

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
	"github.com/hajimehoshi/ebiten/v2/internal/shaderir"
)

var _ graphicsdriver.Graphics = NopGraphicsDriver{}

// NopGraphicsDriver is a graphicsdriver.Graphics that does nothing.
type NopGraphicsDriver struct {
}

func (n NopGraphicsDriver) Initialize() error {
	return nil
}

func (n NopGraphicsDriver) Begin() error {
	return nil
}

func (n NopGraphicsDriver) End(present bool) error {
	return nil
}

func (n NopGraphicsDriver) SetTransparent(transparent bool) {
}

func (n NopGraphicsDriver) SetVertices(vertices []float32, indices []uint32) error {
	return nil
}

func (n NopGraphicsDriver) NewImage(width, height int) (
	graphicsdriver.Image, error,
) {
	return nil, errors.New("nop graphicsdriver used for flushing")
}

func (n NopGraphicsDriver) NewScreenFramebufferImage(width, height int) (
	graphicsdriver.Image, error,
) {
	return nil, errors.New("nop graphicsdriver used for flushing")
}

func (n NopGraphicsDriver) SetVsyncEnabled(enabled bool) {
}

func (n NopGraphicsDriver) NeedsClearingScreen() bool {
	return false
}

func (n NopGraphicsDriver) MaxImageSize() int {
	return 16384
}

func (n NopGraphicsDriver) NewShader(program *shaderir.Program) (
	graphicsdriver.Shader, error,
) {
	return nil, errors.New("nop graphicsdriver used for flushing")
}

func (n NopGraphicsDriver) DrawTriangles(
	dst graphicsdriver.ImageID, srcs [graphics.ShaderImageCount]graphicsdriver.ImageID,
	shader graphicsdriver.ShaderID, dstRegions []graphicsdriver.DstRegion,
	indexOffset int, blend graphicsdriver.Blend,
	uniforms []uint32, fillRule graphicsdriver.FillRule,
) error {
	return nil
}
