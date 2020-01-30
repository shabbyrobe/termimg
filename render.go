package termimg

import (
	"fmt"
	"image"

	"github.com/shabbyrobe/imgx/rgba"
)

type imageRenderer struct {
	cellRenderer CellRenderer

	// Length 32 == 4x8 pixels per character cell.
	//
	// Bit layout for values:
	// 00..31 == count
	// 32..63 == ^color
	//
	// Add (1<<32) to add 1 to the count. This layout is used to allow sorting.
	// Color is inverted so that it sorts in the opposite order.
	colorsCount [32]uint64
}

func newImageRenderer(cr CellRenderer) *imageRenderer {
	rend := &imageRenderer{cellRenderer: cr}
	return rend
}

func (rend *imageRenderer) renderCells(into *CellData, rimg image.Image, flags Flag) error {
	if into == nil {
		*into = CellData{}
	}

	var (
		img, _ = rgba.Convert(rimg)
		size   = img.Bounds().Size()
		w, h   = size.X, size.Y
	)

	into.Cols, into.Rows = w/4, h/8
	max := into.Cols * into.Rows

	if cap(into.Cells) < max {
		if flags&NoAlloc != 0 {
			panic(fmt.Errorf("termimg: buffer size %d, expected %d", len(into.Cells), max))
		}
		into.Cells = make([]Cell, max)

	} else {
		into.Cells = into.Cells[:max]
	}

	cellRenderer := rend.cellRenderer

	n, xEnd, yEnd := 0, w-4, h-8
	for y := 0; y <= yEnd; y += 8 {
		for x := 0; x <= xEnd; x += 4 {
			into.Cells[n] = cellRenderer.cell(rend, img, x, y)
			n++
		}
	}

	return nil
}

func (rend *imageRenderer) renderEscapes(into *EscapeData, rimg image.Image, flags Flag) error {
	// Find the color channel (R, G or B) that has the biggest range of values for the current cell
	// Split this range in the middle and create a corresponding bitmap for the cell
	// Compare the bitmap to the assumed bitmaps for various unicode block graphics characters
	// Re-calculate the foreground and background colors for the chosen character.
	if into == nil {
		*into = EscapeData{}
	}

	img, _ := rgba.Convert(rimg)

	size := img.Bounds().Size()
	w, h := size.X, size.Y
	max := into.MaxSize(flags, w, h)

	into.Reset()

	if len(into.bits) < max {
		if flags&NoAlloc != 0 {
			panic(fmt.Errorf("termimg: buffer size %d, expected %d", len(into.bits), max))
		}
		into.bits = make([]byte, max)
	}

	cellRenderer := rend.cellRenderer
	xEnd, yEnd := w-4, h-8
	for y := 0; y <= yEnd; y += 8 {
		for x := 0; x <= xEnd; x += 4 {
			into.put(flags, cellRenderer.cell(rend, img, x, y))
		}

		// Don't print the last newline, so we can avoid scrolling when rendering video:
		if y < yEnd {
			into.nextRow()
		}
	}

	return nil
}
