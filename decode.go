package termimg

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"regexp"
	"unicode/utf8"

	"github.com/shabbyrobe/imgx/rgba"
	"github.com/shabbyrobe/imgx/termpalette"
)

func DecodeConfig(r io.Reader) (config image.Config, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return config, err
	}
	return DecodeConfigBytes(data)
}

func DecodeConfigBytes(data []byte) (config image.Config, err error) {
	config.ColorModel = color.RGBAModel

	cols, rows := decodeSize(data)
	config.Width = cols * cellW
	config.Height = rows * cellH

	return config, nil
}

// DecodeImage decodes a raw terminal image made of color escapes, runes and newlines into
// an rgba.Image.
//
// If patternSet is nil, DefaultPatternSet is used.
//
// If size is nil, it is inferred from the data. This will be slower.
//
func DecodeImage(rdr io.Reader, bit *BitmapRenderer, size *image.Point) (img *rgba.Image, err error) {
	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	return DecodeImageBytes(data, bit, size)
}

// DecodeImageBytes decodes a raw terminal image made of color escapes,
// runes and newlines into an rgba.Image.
//
// See DecodeImage()
//
func DecodeImageBytes(data []byte, bit *BitmapRenderer, size *image.Point) (img *rgba.Image, err error) {
	if bit == nil {
		bit = blockRenderer
	}

	var cols, rows int

	// If size is not passed, we need to scan the data once in order to determine it.
	if size != nil {
		cols = size.X / cellW
		rows = size.Y / cellH
	} else {
		cols, rows = decodeSize(data)
		size = &image.Point{
			X: cols * cellW,
			Y: rows * cellH,
		}
	}

	target := &decodeImageTarget{
		img: rgba.New(*size),
	}

	dec := &decoder{
		data:   data,
		bit:    bit,
		target: target,
		rows:   rows,
		cols:   cols,
	}

	if err := dec.decode(); err != nil {
		return nil, err
	}
	return target.img, nil
}

// DecodeCells decodes a raw terminal image made of color escapes, runes and newlines into
// a terimg.CellData.
//
// If bit is nil, BlockRenderer is used.
//
func DecodeCells(rdr io.Reader, bit *BitmapRenderer) (cells CellData, err error) {
	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		return cells, err
	}
	return DecodeCellsBytes(data, bit)
}

// DecodeCellsBytes decodes a raw terminal image made of color escapes,
// runes and newlines into an rgba.Image.
//
// See DecodeCells()
//
func DecodeCellsBytes(data []byte, bit *BitmapRenderer) (cells CellData, err error) {
	if bit == nil {
		bit = blockRenderer
	}
	cols, rows := decodeSize(data)

	target := &decodeCellsTarget{
		cells: CellDataFromTerm(cols, rows),
	}

	dec := &decoder{
		data:   data,
		bit:    bit,
		target: target,
		rows:   rows,
		cols:   cols,
	}
	if err := dec.decode(); err != nil {
		return cells, err
	}
	return target.cells, nil
}

var (
	ptnTrueColor = regexp.MustCompile(`^\x1b\[(?P<kind>[34])8;2;(?P<r>[0-9]+);(?P<g>[0-9]+);(?P<b>[0-9]+)m`)
	ptn256Color  = regexp.MustCompile(`^\x1b\[(?P<kind>[34])8;5;(?P<color>[0-9]+)m`)
	ptn16Color   = regexp.MustCompile(`^\x1b\[(?P<color>40|41|42|43|44|45|46|47|100|101|102|103|104|105|106|107|30|31|32|33|34|35|36|37|90|91|92|93|94|95|96|97)m`)

	sub256ColorKind  int
	sub256Color      int
	sub16Color       int
	subTrueColorKind int
	subTrueColorR    int
	subTrueColorG    int
	subTrueColorB    int
)

func init() {
	mustScanSubexps(ptnTrueColor,
		"kind", &subTrueColorKind,
		"r", &subTrueColorR,
		"g", &subTrueColorG,
		"b", &subTrueColorB)

	mustScanSubexps(ptn256Color,
		"kind", &sub256ColorKind,
		"color", &sub256Color)

	mustScanSubexps(ptn16Color, "color", &sub16Color)
}

type decoder struct {
	data       []byte
	bit        *BitmapRenderer
	target     decoderTarget
	rows, cols int
	i          int

	fg, bg       color.RGBA
	fgSet, bgSet bool
}

func (dec *decoder) readTrueColor(match [][]byte) error {
	if len(match[subTrueColorKind]) != 1 {
		return fmt.Errorf("termimg: decode expected color escape at byte %d", dec.i)
	}

	var target *color.RGBA
	switch match[subTrueColorKind][0] {
	case '3':
		target = &dec.fg
		dec.fgSet = true
	case '4':
		target = &dec.bg
		dec.bgSet = true
	default:
		return fmt.Errorf("termimg: decode expected color escape at byte %d", dec.i)
	}

	target.R = colorStringLookup[string(match[subTrueColorR])]
	target.G = colorStringLookup[string(match[subTrueColorG])]
	target.B = colorStringLookup[string(match[subTrueColorB])]
	target.A = 0xff

	return nil
}

func (dec *decoder) read16Color(match [][]byte) error {
	if len(match[sub16Color]) == 0 {
		return fmt.Errorf("termimg: decode expected color escape at byte %d", dec.i)
	}

	var num = colorStringLookup[string(match[sub16Color])]
	switch match[sub16Color][0] {
	case '3':
		dec.fgSet = true
		dec.fg = termpalette.Escape16FgColor[num]
	case '4':
		dec.bg = termpalette.Escape16BgColor[num]
		dec.bgSet = true
	default:
		return fmt.Errorf("termimg: decode expected color escape at byte %d", dec.i)
	}
	return nil
}

func (dec *decoder) read256Color(match [][]byte) error {
	if len(match[sub256ColorKind]) != 1 {
		return fmt.Errorf("termimg: decode expected color escape at byte %d", dec.i)
	}

	var target *color.RGBA
	switch match[sub256ColorKind][0] {
	case '3':
		target = &dec.fg
		dec.fgSet = true
	case '4':
		target = &dec.bg
		dec.bgSet = true
	default:
		return fmt.Errorf("termimg: decode expected color escape at byte %d", dec.i)
	}

	target.R, target.G, target.B = term256AsRGB(colorStringLookup[string(match[sub256Color])])
	target.A = 0xff
	return nil
}

func (dec *decoder) readColor() error {
	for dec.i < len(dec.data) {
		if m := ptnTrueColor.FindSubmatch(dec.data[dec.i:]); m != nil {
			if err := dec.readTrueColor(m); err != nil {
				return err
			}
			dec.i += len(m[0])

		} else if m := ptn256Color.FindSubmatch(dec.data[dec.i:]); m != nil {
			if err := dec.read256Color(m); err != nil {
				return err
			}
			dec.i += len(m[0])

		} else if m := ptn16Color.FindSubmatch(dec.data[dec.i:]); m != nil {
			if err := dec.read16Color(m); err != nil {
				return err
			}
			dec.i += len(m[0])

		} else if bytes.HasPrefix(dec.data[dec.i:], reset) {
			dec.fgSet, dec.bgSet = false, false
			dec.fg, dec.bg = color.RGBA{}, color.RGBA{}
			dec.i += len(reset)

		} else {
			return nil
		}
	}

	return io.EOF
}

func (dec *decoder) decode() error {
	var col, row int

	for dec.i < len(dec.data) {
		if err := dec.readColor(); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if dec.data[dec.i] == ' ' || dec.data[dec.i] == '\r' || dec.data[dec.i] == '\t' {
			dec.i++

		} else if dec.data[dec.i] == '\n' {
			if col > dec.cols {
				return fmt.Errorf("termimg: col exceeded width %d at byte %d", dec.cols, dec.i)
			}
			row++
			if row > dec.rows {
				return fmt.Errorf("termimg: row exceeded height %d at byte %d", dec.rows, dec.i)
			}
			col = 0
			dec.i++

		} else if dec.data[dec.i] == '\x1b' {
			return fmt.Errorf("termimg: decode found unexpected escape sequence at byte %d", dec.i)

		} else {
			rn, sz := utf8.DecodeRune(dec.data[dec.i:])
			if rn == utf8.RuneError {
				return fmt.Errorf("termimg: decode expected rune at byte %d", dec.i)
			}

			var found *Bitmap
			for _, b := range dec.bit.Bitmaps {
				if rn == b.Rune {
					found = &b
					break
				}
			}
			if found == nil && rn == dec.bit.Default.Rune {
				found = &dec.bit.Default
			}
			if found == nil {
				return fmt.Errorf("termimg: decode found rune %q byte %d, but this rune does not exist in the pattern set", string(rn), dec.i)
			}
			if !dec.fgSet {
				return fmt.Errorf("termimg: decode found a rune with no foreground color at byte %d", dec.i)
			}
			if !dec.bgSet {
				return fmt.Errorf("termimg: decode found a rune with no background color at byte %d", dec.i)
			}

			dec.target.set(col, row, dec.fg, dec.bg, found)
			dec.i += sz
			col++
		}
	}

	return nil
}

type decoderTarget interface {
	set(col, row int, fg, bg color.RGBA, bits *Bitmap)
}

type decodeCellsTarget struct {
	cells CellData
}

func (tgt *decodeCellsTarget) set(col, row int, fg, bg color.RGBA, bits *Bitmap) {
	tgt.cells.Cells[tgt.cells.Cols*row+col] = Cell{
		FgColor: fg,
		BgColor: bg,
		Code:    bits.Rune,
	}
}

type decodeImageTarget struct {
	img *rgba.Image
}

func (tgt *decodeImageTarget) set(col, row int, fg, bg color.RGBA, bits *Bitmap) {
	x, y := col*cellW, row*cellH

	n := Bits(1 << 31)
	for cellY := 0; cellY < 8; cellY++ {
		yoff := (y + cellY) * tgt.img.Stride
		for cellX := 0; cellX < 4; cellX++ {
			idx := yoff + x + cellX
			if bits.Bits&n == 0 {
				tgt.img.Vals[idx] = bg
			} else {
				tgt.img.Vals[idx] = fg
			}
			n >>= 1
		}
	}
}

const cellW, cellH = 4, 8

func decodeSize(data []byte) (cols, rows int) {
	var col int

	for i := 0; i < len(data); {
		for {
			if m := ptnTrueColor.Find(data[i:]); m != nil {
				i += len(m)
			} else if m := ptn256Color.Find(data[i:]); m != nil {
				i += len(m)
			} else if m := ptn16Color.Find(data[i:]); m != nil {
				i += len(m)
			} else if bytes.HasPrefix(data[i:], reset) {
				i += len(reset)
			} else {
				break
			}
		}

		if i >= len(data) {
			break
		}

		if data[i] == ' ' || data[i] == '\r' || data[i] == '\t' {
			i++
			continue

		} else if data[i] == '\n' {
			if col > cols {
				cols = col
			}
			rows++
			col = 0
			i++
			continue

		} else {
			rn, sz := utf8.DecodeRune(data[i:])
			_ = rn
			i += sz

			// If it's nothing else, we have to presume it's a rune:
			col++
		}
	}

	// Don't forget to account for the last row (now that we've possibly removed the last
	// newline!)
	if col > 0 {
		if col > cols {
			cols = col
		}
		rows++
	}

	return cols, rows
}

var colorStringLookup = make(map[string]uint8)

func init() {
	for i := 0; i < 256; i++ {
		colorStringLookup[fmt.Sprintf("%d", i)] = uint8(i)
	}
}
