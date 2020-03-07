package termimg

import (
	"fmt"
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

type RendererConfig interface {
	Renderer() (Renderer, error)
}

// XXX: Implementation note: there are large slabs of copy-pasta boilerplate for the
// Escapes() and Cells() methods. This is needed to provide a proper, opt-in, zero-alloc
// strategy for downstream users. Previous attempts to avoid it that used interfaces and
// function pointers consistently showed small allocations in benchmarks.

type Renderer interface {
	EscapeRenderer
	CellRenderer
}

type EscapeRenderer interface {
	// Render img into an EscapeData as a series of escape codes and UTF-8 runes suitable for
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
	Escapes(into *EscapeData, img image.Image, flags Flag) error
}

type CellRenderer interface {
	// Render img into a CellData as a series of RGBA colors and UTF-8 runes, suitable for
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
	Cells(into *CellData, img image.Image, flags Flag) error
}

func prepareCells(into *CellData, rimg image.Image, flags Flag) (cells *CellData, img *rgba.Image, w, h int) {
	img, _ = rgba.Convert(rimg)
	size := img.Bounds().Size()
	w, h = size.X, size.Y

	if into == nil {
		if flags&NoAlloc != 0 {
			panic(fmt.Errorf("termimg: nil CellData, but NoAlloc set"))
		}
		*into = CellData{}
	}

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

	return into, img, w, h
}

func prepareEscapes(into *EscapeData, rimg image.Image, flags Flag) (cells *EscapeData, img *rgba.Image, w, h int) {
	img, _ = rgba.Convert(rimg)
	size := img.Bounds().Size()
	w, h = size.X, size.Y

	if into == nil {
		if flags&NoAlloc != 0 {
			panic(fmt.Errorf("termimg: nil EscapeData, but NoAlloc set"))
		}
		*into = EscapeData{}
	}
	max := into.MaxSize(flags, w, h)
	into.Reset()

	if len(into.bits) < max {
		if flags&NoAlloc != 0 {
			panic(fmt.Errorf("termimg: buffer size %d, expected %d", len(into.bits), max))
		}
		into.bits = make([]byte, max)
	}

	return into, img, w, h
}
