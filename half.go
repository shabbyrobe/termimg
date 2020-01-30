package termimg

import (
	"github.com/shabbyrobe/imgx/rgba"
)

type HalfBlockRenderer struct {
	BitmapRenderer
}

func (half *HalfBlockRenderer) cell(rend *imageRenderer, img *rgba.Image, x0, y0 int) (result Cell) {
	code := '▄'
	pattern := Bits(lowerHalfBitmap)
	return half.BitmapRenderer.cellForCode(rend, img, x0, y0, code, pattern)
}
