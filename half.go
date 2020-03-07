package termimg

import (
	"image"

	"github.com/shabbyrobe/imgx/rgba"
)

type HalfBlockConfig struct{}

func (hc HalfBlockConfig) Renderer() (Renderer, error) {
	return &HalfBlockRenderer{}, nil
}

type HalfBlockRenderer struct {
	bit BitmapRenderer
}

func (half *HalfBlockRenderer) Escapes(into *EscapeData, img image.Image, flags Flag) error {
	// XXX: intentional copy-pasta; see renderer.go for details

	into, rimg, w, h := prepareEscapes(into, img, flags)
	xEnd, yEnd := w-4, h-8
	for y := 0; y <= yEnd; y += 8 {
		for x := 0; x <= xEnd; x += 4 {
			into.put(flags, half.cell(rimg, x, y))
		}

		// Don't print the last newline, so we can avoid scrolling when rendering video:
		if y < yEnd {
			into.nextRow()
		}
	}

	return nil
}

func (half *HalfBlockRenderer) Cells(into *CellData, img image.Image, flags Flag) error {
	// XXX: intentional copy-pasta; see renderer.go for details

	into, rimg, w, h := prepareCells(into, img, flags)
	n, xEnd, yEnd := 0, w-4, h-8
	for y := 0; y <= yEnd; y += 8 {
		for x := 0; x <= xEnd; x += 4 {
			into.Cells[n] = half.cell(rimg, x, y)
			n++
		}
	}

	return nil
}

func (half *HalfBlockRenderer) cell(img *rgba.Image, x0, y0 int) (result Cell) {
	code := '▄'
	pattern := Bits(lowerHalfBitmap)
	return half.bit.cellForCode(img, x0, y0, code, pattern)
}
