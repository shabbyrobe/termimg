package termimg

import (
	"fmt"
	"image"
	"image/color"
	"math/bits"

	"github.com/shabbyrobe/imgx/rgba"
)

const (
	lowerHalfBitmap = 0b_0000_0000_0000_0000_1111_1111_1111_1111
)

type Bits uint32

func (b Bits) Ones() int {
	return bits.OnesCount32(uint32(b))
}

type Bitmap struct {
	Bits Bits
	Rune rune
}

func (p Bitmap) String() string {
	bin := fmt.Sprintf("%032b", p.Bits)
	out := bin[:4]
	for i := 4; i < 32; i += 4 {
		out += "_" + bin[i:i+4]
	}
	return fmt.Sprintf("%s:%s", out, string(p.Rune))
}

func (p Bitmap) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *Bitmap) UnmarshalText(text []byte) (err error) {
	s := string(text)
	*p, err = parseBitmap(s)
	return err
}

type BitmapConfig struct {
	Bitmaps []Bitmap
	Default Bitmap
}

func (config BitmapConfig) Renderer() (Renderer, error) {
	return NewBitmapRenderer(config)
}

// BitmapRenderer renders each cell in the terminal using runes assocated with the closest
// matching grid of 4x8 pixels (Bitmap).
type BitmapRenderer struct {
	bitmaps       []Bitmap
	defaultBitmap Bitmap

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

func NewBitmapRenderer(config BitmapConfig) (*BitmapRenderer, error) {
	return &BitmapRenderer{
		bitmaps:       config.Bitmaps,
		defaultBitmap: config.Default,
	}, nil
}

func (bit *BitmapRenderer) Escapes(into *EscapeData, img image.Image, flags Flag) error {
	// XXX: intentional copy-pasta; see renderer.go for details

	into, rimg, w, h := prepareEscapes(into, img, flags)
	xEnd, yEnd := w-4, h-8
	for y := 0; y <= yEnd; y += 8 {
		for x := 0; x <= xEnd; x += 4 {
			into.put(flags, bit.cell(rimg, x, y))
		}

		// Don't print the last newline, so we can avoid scrolling when rendering video:
		if y < yEnd {
			into.nextRow()
		}
	}

	return nil
}

func (bit *BitmapRenderer) Cells(into *CellData, img image.Image, flags Flag) error {
	// XXX: intentional copy-pasta; see renderer.go for details

	into, rimg, w, h := prepareCells(into, img, flags)
	n, xEnd, yEnd := 0, w-4, h-8
	for y := 0; y <= yEnd; y += 8 {
		for x := 0; x <= xEnd; x += 4 {
			into.Cells[n] = bit.cell(rimg, x, y)
			n++
		}
	}

	return nil
}

// Find the best character and colors for a 4x8 part of the image at the given position
func (bit *BitmapRenderer) cell(img *rgba.Image, x0, y0 int) (result Cell) {
	// Find the color channel (R, G or B) that has the biggest range of values for the current cell
	// Split this range in the middle and create a corresponding bitmap for the cell
	// Compare the bitmap to the assumed bitmaps for various unicode block graphics characters
	// Re-calculate the foreground and background colors for the chosen character.

	// Determine the minimum and maximum value for each color channel:
	var minr, ming, minb uint32 = 0xFF, 0xFF, 0xFF
	var maxr, maxg, maxb uint32 = 0, 0, 0

	// Number of distinct colours we have found in this cell, used to
	// determine the end of the rend.colorsCount cache. Tracking the size
	// internal to this function obviates the need to zero the memory on
	// every call to cell().
	var colorsSize int

	yN, xN, yOff := y0+8, x0+4, y0*img.Stride

	for y := y0; y < yN; y++ {
		for x := x0; x < xN; x++ {
			var color uint32 = 0

			c := img.Vals[yOff+x]
			r, g, b := uint32(c.R), uint32(c.G), uint32(c.B)

			if r < minr {
				minr = r
			}
			if r > maxr {
				maxr = r
			}
			color = r

			if g < ming {
				ming = g
			}
			if g > maxg {
				maxg = g
			}
			color = (color << 8) | g

			if b < minb {
				minb = b
			}
			if b > maxb {
				maxb = b
			}
			color = (color << 8) | b

			// The original version of the algorithm appeared to sort by reverse count,
			// then by reverse colour. This may not be necessary, but should be tested
			// before being changed. The bit inversion is the "optimised" (i.e. "too
			// clever by half") version of the existing algorithm, not a deliberate choice
			// based on the output:
			colorInv := ^color

			for i := 0; i < colorsSize; i++ {
				if uint32(bit.colorsCount[i]) == colorInv {
					// Increment the count, which is stored in the high 32-bits:
					bit.colorsCount[i] += 0x1_0000_0000
					goto next
				}
			}

			bit.colorsCount[colorsSize] = 0x1_0000_0000 | uint64(colorInv)
			colorsSize++
		next:
		}

		yOff += img.Stride
	}

	var count2 uint32 // sum of the number of times the most common two colours appear in the 4x8 segment
	var maxCountColor1 uint32
	var maxCountColor2 uint32

	if colorsSize == 1 {
		count2 = uint32(bit.colorsCount[0] >> 32)
		maxCountColor1 = ^uint32(bit.colorsCount[0])
		maxCountColor2 = maxCountColor1

	} else {
		var max1, max2 uint64
		for i := 0; i < colorsSize; i++ {
			rc := bit.colorsCount[i]
			if rc > max1 {
				max1, max2 = rc, max1
			} else if rc > max2 {
				max2 = rc
			}
		}

		count2 = uint32(max1>>32 + max2>>32)
		maxCountColor1 = ^uint32(max1)
		maxCountColor2 = ^uint32(max2)
	}

	var setBits Bits

	// If the sum of the number of pixels containing max1 and max2 is more than half
	// the number of pixels, use 'direct' mode:
	var direct = count2 > (8*4)/2

	if direct {
		var maxR1, maxG1, maxB1 = (maxCountColor1 >> 16) & 0xff, (maxCountColor1 >> 8) & 0xff, (maxCountColor1 & 0xff)
		var maxR2, maxG2, maxB2 = (maxCountColor2 >> 16) & 0xff, (maxCountColor2 >> 8) & 0xff, (maxCountColor2 & 0xff)

		yOff := y0 * img.Stride

		for y := y0; y < yN; y++ {
			for x := x0; x < xN; x++ {
				setBits <<= 1
				var d1, d2 uint32

				c := img.Vals[yOff+x]
				r, g, b := uint32(c.R), uint32(c.G), uint32(c.B)

				cr1 := maxR1 - r
				cr2 := maxR2 - r
				d1 += cr1 * cr1
				d2 += cr2 * cr2

				cg1 := maxG1 - g
				cg2 := maxG2 - g
				d1 += cg1 * cg1
				d2 += cg2 * cg2

				cb1 := maxB1 - b
				cb2 := maxB2 - b
				d1 += cb1 * cb1
				d2 += cb2 * cb2

				if d1 > d2 {
					setBits |= 1
				}
			}

			yOff += img.Stride
		}

	} else {
		// Determine the color channel with the greatest range.
		// We just split at the middle of the interval instead of computing the median.
		var splitChannel byte
		var threshhold uint32

		rdiff, gdiff, bdiff := maxr-minr, maxg-ming, maxb-minb
		if rdiff >= gdiff && rdiff >= bdiff {
			splitChannel, threshhold = 'r', minr+(rdiff/2)
		} else if gdiff >= bdiff {
			splitChannel, threshhold = 'g', ming+(gdiff/2)
		} else {
			splitChannel, threshhold = 'b', minb+(bdiff/2)
		}

		yOff := y0 * img.Stride

		// Compute a bitmap using the given split and sum the color values for both buckets.
		yN, xN := y0+8, x0+4
		for y := y0; y < yN; y++ {
			for x := x0; x < xN; x++ {
				setBits <<= 1

				c := img.Vals[yOff+x]

				switch splitChannel {
				case 'r':
					if uint32(c.R) > threshhold {
						setBits |= 1
					}
				case 'g':
					if uint32(c.G) > threshhold {
						setBits |= 1
					}
				case 'b':
					if uint32(c.B) > threshhold {
						setBits |= 1
					}
				}
			}

			yOff += img.Stride
		}
	}

	// Find the best bitmap match by counting the bits that don't match,
	// including the inverted bitmaps.
	var bestDiff = 8 // FIXME: why 8 and not 16? not sure, need to research.
	var best = bit.defaultBitmap
	var inverted bool

	for _, bitmap := range bit.bitmaps {
		pbits := bitmap.Bits
		diff := (pbits ^ setBits).Ones()
		if diff < bestDiff {
			best, bestDiff, inverted = bitmap, diff, false
		}

		// Invert the pattern and try again:
		pbits = ^pbits
		diff = (pbits ^ setBits).Ones()
		if diff < bestDiff {
			best, bestDiff, inverted = bitmap, diff, true
		}
	}

	if direct {
		var result Cell
		if inverted {
			maxCountColor1, maxCountColor2 = maxCountColor2, maxCountColor1
		}
		result.Code = best.Rune
		result.FgColor = color.RGBA{
			R: uint8(maxCountColor2 >> 16),
			G: uint8(maxCountColor2 >> 8),
			B: uint8(maxCountColor2),
			A: 0xFF,
		}
		result.BgColor = color.RGBA{
			R: uint8(maxCountColor1 >> 16),
			G: uint8(maxCountColor1 >> 8),
			B: uint8(maxCountColor1),
			A: 0xFF,
		}
		return result
	}

	return bit.cellForCode(img, x0, y0, best.Rune, best.Bits)
}

// Return a Cell with the given code point and corresponding average fg and bg colors.
//
// NOTE: This is duplicated with the half-block renderer... I tried to share the code
// by making it a global function but got a 30% slowdown. WAT?
func (bit *BitmapRenderer) cellForCode(img *rgba.Image, x0, y0 int, code rune, pattern Bits) (result Cell) {
	result.Code = code

	var (
		fgCount = uint16(0)
		bgCount = uint16(0)
		mask    = Bits(0b_1000_0000_0000_0000_0000_0000_0000_0000)

		avgBgr, avgBgg, avgBgb uint16
		avgFgr, avgFgg, avgFgb uint16
	)

	yN, xN, yOff := y0+8, x0+4, y0*img.Stride

	for y := y0; y < yN; y++ {
		for x := x0; x < xN; x++ {
			c := img.Vals[yOff+x]

			if pattern&mask != 0 {
				avgFgr += uint16(c.R)
				avgFgg += uint16(c.G)
				avgFgb += uint16(c.B)
				fgCount++
			} else {
				avgBgr += uint16(c.R)
				avgBgg += uint16(c.G)
				avgBgb += uint16(c.B)
				bgCount++
			}
			mask = mask >> 1
		}
		yOff += img.Stride
	}

	// Calculate the average color value for each bucket
	if bgCount != 0 {
		result.BgColor = color.RGBA{
			R: uint8(avgBgr / bgCount),
			G: uint8(avgBgg / bgCount),
			B: uint8(avgBgb / bgCount),
			A: 0xFF,
		}
	}
	if fgCount != 0 {
		result.FgColor = color.RGBA{
			R: uint8(avgFgr / fgCount),
			G: uint8(avgFgg / fgCount),
			B: uint8(avgFgb / fgCount),
			A: 0xFF,
		}
	}
	return result
}
