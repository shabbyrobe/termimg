package termimg

import (
	"image"

	"github.com/shabbyrobe/imgx/rgba"
)

type Flag int

const (
	// Use only the 256-color terminal palette instead of true color.
	Color256 Flag = 1 << iota

	// Use only the 16-color terminal palette instead of true color. Takes precedence over
	// Color256.
	Color16

	// Raise an error if encoding into an EscapeData would cause a reallocation
	// of the buffer:
	NoAlloc

	// Do not compress runs of colors in the EscapeData output; every character
	// will have its color emitted.
	NoReduce
)

type CellRenderer interface {
	cell(r *imageRenderer, img *rgba.Image, x0, y0 int) (result Cell)
}

// Encode img into an EscapeData as a series of escape codes and UTF-8 runes suitable for
// writing directly to the stdout of a VT100-compatible terminal. If 'into' is 'nil', an
// EscapeData is allocated.
//
//	var into termimg.EscapeData
//	if err := Encode(&into, img, flags, nil); err != nil {
//		// ...
//	}
//
// EscapeData.Value() can be written directly to stdout, but take care to clean up
// afterwards:
//
//	os.Stdout.Write(into.Value())
//	os.Stdout.Write([]byte("\033[0m"))
// 	os.Stdout.Write([]byte("\n"))
//
// EscapeData can be reused to help control allocations. It will be grown if necessary,
// but never shrunk. To raise an error if an allocation would occur, pass FlagNoAlloc.
//
func Encode(into *EscapeData, img image.Image, flags Flag, cr CellRenderer) error {
	if cr == nil {
		cr = DefaultRenderer
	}
	var rend = imageRenderer{cellRenderer: cr} // FIXME: alloc
	return rend.renderEscapes(into, img, flags)
}

// Encode img into a CellData as a series of RGBA colors and UTF-8 runes, suitable for
// using with a library like github.com/gdamore/tcell. If 'into' is 'nil', a CellData
// is allocated.
//
//	var into termimg.CellData
//	if err := EncodeCells(&into, img, flags, nil); err != nil {
//		// ...
//	}
//
// CellData can be reused to help control allocations. It will be grown if necessary, but
// never shrunk. To raise an error if an allocation would occur, pass FlagNoAlloc.
//
func EncodeCells(into *CellData, img image.Image, flags Flag, cr CellRenderer) error {
	if cr == nil {
		cr = DefaultRenderer
	}
	var rend = imageRenderer{cellRenderer: cr} // FIXME: alloc
	return rend.renderCells(into, img, flags)
}
