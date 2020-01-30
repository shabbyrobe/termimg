package termimg

import (
	"image/color"

	"github.com/shabbyrobe/imgx/rgba"
)

type SimpleRenderer struct {
	Code rune
}

func NewSimpleRenderer(code rune) *SimpleRenderer {
	return &SimpleRenderer{Code: code}
}

func (simp *SimpleRenderer) cell(rend *imageRenderer, img *rgba.Image, x0, y0 int) (result Cell) {
	var sumR, sumG, sumB uint

	yN, xN, yOff := y0+8, x0+4, y0*img.Stride

	for y := y0; y < yN; y++ {
		for x := x0; x < xN; x++ {
			c := img.Vals[yOff+x]
			sumR += uint(c.R)
			sumG += uint(c.G)
			sumB += uint(c.B)
		}
	}

	result.FgColor = color.RGBA{
		R: uint8(sumR >> 5),
		G: uint8(sumG >> 5),
		B: uint8(sumB >> 5),
		A: 0xff,
	}
	result.Code = simp.Code
	return result
}
