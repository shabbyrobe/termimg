package termimg

import (
	"image/color"
)

// EscapeData is used by Encode to store a rendered image which can be
// written directly to your terminal.
//
// The result of the last Encode is accessible using Value().
//
// EscapeData can be reused for multiple Encode() calls, but please note
// that the byte slice returned by Value() will change.
type EscapeData struct {
	bits []byte

	// If you add any more state to EscapeData, don't forget to add it to Reset():
	n          int
	firstOfRow bool
	lastBg     color.RGBA
	lastFg     color.RGBA
}

// Value returns the last image built into the EscapeData by Encode(), which can be
// written directly to your terminal.
//
// If the EscapeData is reused, the []byte slice returned will change; if this is
// undesirable, you will need to copy the bytes yourself.
func (t *EscapeData) Value() []byte {
	return t.bits[:t.n]
}

func (t *EscapeData) Preallocate(flags Flag, w, h int) {
	t.SetBuffer(make([]byte, t.MaxSize(flags, w, h)))
}

// MaxSize returns the largest possible buffer size
func (t *EscapeData) MaxSize(flags Flag, w, h int) int {
	return (h / 8) * t.maxRowSize(flags, w)
}

// SetBuffer gives EscapeData an existing scratch area to work with.
func (t *EscapeData) SetBuffer(buf []byte) {
	t.bits = buf
}

func (t *EscapeData) Reset() {
	t.n = 0
	t.firstOfRow = true
	t.lastFg = color.RGBA{}
	t.lastBg = color.RGBA{}
}

func (t *EscapeData) nextRow() {
	t.n += copy(t.bits[t.n:], nextRow)
	t.firstOfRow = true
}

func (t *EscapeData) put(flags Flag, cell Cell) {
	if flags&NoReduce != 0 || t.firstOfRow || t.lastBg != cell.BgColor {
		if flags&Color16 != 0 {
			t.n += cell.PutBg16(t.bits[t.n:])
		} else if flags&Color256 != 0 {
			t.n += cell.PutBg256(t.bits[t.n:])
		} else {
			t.n += cell.PutBg(t.bits[t.n:])
		}
		t.lastBg = cell.BgColor
	}

	if flags&NoReduce != 0 || t.firstOfRow || t.lastFg != cell.FgColor {
		if flags&Color16 != 0 {
			t.n += cell.PutFg16(t.bits[t.n:])
		} else if flags&Color256 != 0 {
			t.n += cell.PutFg256(t.bits[t.n:])
		} else {
			t.n += cell.PutFg(t.bits[t.n:])
		}
		t.lastFg = cell.FgColor
	}

	t.firstOfRow = false
	t.n += cell.PutCode(t.bits[t.n:])
}

func (t *EscapeData) maxPixelSize(flags Flag) int {
	var x EscapeData
	var b [128]byte // better not be bigger than this!
	var c Cell

	// XXX: hack here, when these were all 0xFF, 256 color mode seemed to produce shorter
	// output by falling into shorter escapes? Probably need to just work this out by
	// hand and hard code
	c.FgColor = color.RGBA{0xee, 0xdd, 0xcc, 0xff}
	c.BgColor = color.RGBA{0xee, 0xdd, 0xcc, 0xff}

	c.Code = 0x10fffe // 4-byte utf-8
	x.SetBuffer(b[:])
	x.put(flags|NoReduce, c)
	return x.n
}

func (t *EscapeData) maxRowSize(flags Flag, w int) int {
	return (w/4)*t.maxPixelSize(flags) + len(nextRow)
}
