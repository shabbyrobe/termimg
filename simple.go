package termimg

import (
	"image"
	"image/color"

	"github.com/shabbyrobe/imgx/rgba"
)

type SimpleConfig struct {
	Code rune
}

func (config SimpleConfig) Renderer() (Renderer, error) {
	return NewSimpleRenderer(config.Code), nil
}

type SimpleRenderer struct {
	Code rune
}

func NewSimpleRenderer(code rune) *SimpleRenderer {
	return &SimpleRenderer{Code: code}
}

func (simp *SimpleRenderer) Escapes(into *EscapeData, img image.Image, flags Flag) error {
	// XXX: intentional copy-pasta; see renderer.go for details

	into, rimg, w, h := prepareEscapes(into, img, flags)
	xEnd, yEnd := w-4, h-8
	for y := 0; y <= yEnd; y += 8 {
		for x := 0; x <= xEnd; x += 4 {
			into.put(flags, simp.cell(rimg, x, y))
		}

		// Don't print the last newline, so we can avoid scrolling when rendering video:
		if y < yEnd {
			into.nextRow()
		}
	}

	return nil
}

func (simp *SimpleRenderer) Cells(into *CellData, img image.Image, flags Flag) error {
	// XXX: intentional copy-pasta; see renderer.go for details

	into, rimg, w, h := prepareCells(into, img, flags)
	n, xEnd, yEnd := 0, w-4, h-8
	for y := 0; y <= yEnd; y += 8 {
		for x := 0; x <= xEnd; x += 4 {
			into.Cells[n] = simp.cell(rimg, x, y)
			n++
		}
	}

	return nil
}

func (simp *SimpleRenderer) cell(img *rgba.Image, x0, y0 int) (result Cell) {
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
