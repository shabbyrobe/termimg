package termimg

import (
	"image"
)

type Flag int

const (
	FlagMode256 Flag = 1 << iota
	FlagMode16
	FlagNoAlloc
	FlagHalfBlockOnly
	FlagNoReduce
)

// Encode img into an EscapeData as a series of escape codes and UTF-8 runes suitable for
// writing directly to the stdout of a VT100-compatible terminal. If 'into' is 'nil', an
// EscapeData is allocated.
//
//	var into termimg.EscapeData
//	if err := Encode(&into, img, flags, nil); err != nil {
//		// ...
//	}
//
// EscapeData can be reused to help control allocations. It will be grown if necessary,
// but never shrunk. To raise an error if an allocation would occur, pass FlagNoAlloc.
//
func Encode(into *EscapeData, img image.Image, flags Flag, patternSet *PatternSet) error {
	if patternSet == nil {
		patternSet = DefaultPatternSet
	}
	return newRenderer(patternSet).renderEscapes(into, img, flags)
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
func EncodeCells(into *CellData, img image.Image, flags Flag, patternSet *PatternSet) error {
	if patternSet == nil {
		patternSet = DefaultPatternSet
	}
	return newRenderer(patternSet).renderCells(into, img, flags)
}
