package termimg

import (
	"image/color"
	"strconv"

	"github.com/shabbyrobe/imgx/termpalette"
)

// The buffer passed to Cell.PutColor and Cell.PutCode must be at least this size
const CellMinBufSize = 64

type CellData struct {
	Cells []Cell
	Cols  int
	Rows  int
}

func CellDataFromPixels(w, h int) CellData {
	return CellDataFromTerm(w/4, h/8)
}

func CellDataFromTerm(cols, rows int) CellData {
	return CellData{
		Cols: cols, Rows: rows, Cells: make([]Cell, cols*rows),
	}
}

type Cell struct {
	FgColor color.RGBA
	BgColor color.RGBA
	Code    rune
}

func (c Cell) Fg256() uint8 {
	return uint8(termpalette.CodeInt[index256.NearestRGBAIndex(c.FgColor)])
}

func (c Cell) Bg256() uint8 {
	return uint8(termpalette.CodeInt[index256.NearestRGBAIndex(c.BgColor)])
}

func (c Cell) Fg16() uint8 {
	return uint8(termpalette.Escape16FgInt[index16.NearestRGBAIndex(c.FgColor)])
}

func (c Cell) Bg16() uint8 {
	return uint8(termpalette.Escape16BgInt[index16.NearestRGBAIndex(c.BgColor)])
}

func (c *Cell) PutCode(buf []byte) (n int) {
	// Avoid including utf8 package just for this:
	if c.Code < 128 {
		buf[0] = byte(c.Code)
		return 1

	} else if c.Code < 0x7ff {
		_ = buf[1]
		buf[0] = 0xc0 | byte(c.Code>>6)
		buf[1] = 0x80 | byte(c.Code&0x3f)
		return 2

	} else if c.Code < 0xffff {
		_ = buf[2]
		buf[0] = 0xe0 | (byte(c.Code >> 12))
		buf[1] = 0x80 | (byte(c.Code>>6) & 0x3f)
		buf[2] = 0x80 | (byte(c.Code & 0x3f))
		return 3

	} else if c.Code < 0x10ffff {
		_ = buf[3]
		buf[0] = 0xf0 | (byte(c.Code >> 18))
		buf[1] = 0x80 | (byte((c.Code >> 12) & 0x3f))
		buf[2] = 0x80 | (byte((c.Code >> 6) & 0x3f))
		buf[3] = 0x80 | (byte(c.Code & 0x3f))
		return 4

	} else {
		panic("invalid codepoint")
	}
}

func (c *Cell) PutFg(buf []byte) (n int) {
	// copy is miles faster than append
	r, g, b := c.FgColor.R, c.FgColor.G, c.FgColor.B
	n += copy(buf, fgPrefix)
	n += copy(buf[n:], colStr[r])
	buf[n] = ';'
	n++
	n += copy(buf[n:], colStr[g])
	buf[n] = ';'
	n++
	n += copy(buf[n:], colStr[b])
	buf[n] = 'm'
	n++
	return n
}

func (c *Cell) PutBg(buf []byte) (n int) {
	r, g, b := c.BgColor.R, c.BgColor.G, c.BgColor.B
	n += copy(buf, bgPrefix)
	n += copy(buf[n:], colStr[r])
	buf[n] = ';'
	n++
	n += copy(buf[n:], colStr[g])
	buf[n] = ';'
	n++
	n += copy(buf[n:], colStr[b])
	buf[n] = 'm'
	n++
	return n
}

func (c *Cell) PutFg256(buf []byte) (n int) {
	n += copy(buf, fg256Prefix)
	colorIndex := index256.NearestRGBAIndex(c.FgColor)
	n += copy(buf[n:], termpalette.Code[colorIndex])
	buf[n] = 'm'
	n++
	return n
}

func (c *Cell) PutBg256(buf []byte) (n int) {
	n += copy(buf, bg256Prefix)
	colorIndex := index256.NearestRGBAIndex(c.BgColor)
	n += copy(buf[n:], termpalette.Code[colorIndex])
	buf[n] = 'm'
	n++
	return n
}

func (c *Cell) PutFg16(buf []byte) (n int) {
	n += copy(buf, col16Prefix)
	colorIndex := index16.NearestRGBAIndex(c.FgColor)
	n += copy(buf[n:], termpalette.Escape16Fg[colorIndex])
	buf[n] = 'm'
	n++
	return n
}

func (c *Cell) PutBg16(buf []byte) (n int) {
	n += copy(buf, col16Prefix)
	colorIndex := index16.NearestRGBAIndex(c.BgColor)
	n += copy(buf[n:], termpalette.Escape16Bg[colorIndex])
	buf[n] = 'm'
	n++
	return n
}

var (
	reset    = []byte("\x1b[0m")
	nextRow  = []byte("\x1b[0m\n")
	bgPrefix = []byte("\x1b[48;2;")
	fgPrefix = []byte("\x1b[38;2;")

	bg256Prefix = []byte("\x1b[48;5;")
	fg256Prefix = []byte("\x1b[38;5;")

	col16Prefix = []byte("\x1b[")

	colStr = [256][]byte{}
)

func init() {
	for i := int64(0); i < 256; i++ {
		colStr[i] = []byte(strconv.FormatInt(i, 10))
	}
}
