package benchdraw

import (
	"github.com/hajimehoshi/ebiten/v2/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
	"github.com/hajimehoshi/ebiten/v2/internal/shaderir"
)

var _ graphicsdriver.Graphics = nopGraphicsDriver{}

type nopGraphicsDriver struct {
}

func (n nopGraphicsDriver) Initialize() error {
	return nil
}

func (n nopGraphicsDriver) Begin() error {
	return nil
}

func (n nopGraphicsDriver) End(present bool) error {
	return nil
}

func (n nopGraphicsDriver) SetTransparent(transparent bool) {
}

func (n nopGraphicsDriver) SetVertices(vertices []float32, indices []uint32) error {
	return nil
}

func (n nopGraphicsDriver) NewImage(width, height int) (
	graphicsdriver.Image, error,
) {
	return nopImage{}, nil
}

func (n nopGraphicsDriver) NewScreenFramebufferImage(width, height int) (
	graphicsdriver.Image, error,
) {
	return nopImage{}, nil
}

func (n nopGraphicsDriver) SetVsyncEnabled(enabled bool) {
}

func (n nopGraphicsDriver) NeedsClearingScreen() bool {
	return false
}

func (n nopGraphicsDriver) MaxImageSize() int {
	return 16384
}

func (n nopGraphicsDriver) NewShader(program *shaderir.Program) (
	graphicsdriver.Shader, error,
) {
	return nopShader{}, nil
}

func (n nopGraphicsDriver) DrawTriangles(
	dst graphicsdriver.ImageID, srcs [graphics.ShaderImageCount]graphicsdriver.ImageID,
	shader graphicsdriver.ShaderID, dstRegions []graphicsdriver.DstRegion,
	indexOffset int, blend graphicsdriver.Blend,
	uniforms []uint32, fillRule graphicsdriver.FillRule,
) error {
	return nil
}

type nopImage struct {
}

func (n nopImage) ID() graphicsdriver.ImageID {
	return 1
}

func (n nopImage) Dispose() {
}

func (n nopImage) ReadPixels(args []graphicsdriver.PixelsArgs) error {
	return nil
}

func (n nopImage) WritePixels(args []graphicsdriver.PixelsArgs) error {
	return nil
}

type nopShader struct {
}

func (nopShader) ID() graphicsdriver.ShaderID {
	return 1
}

func (nopShader) Dispose() {
}
