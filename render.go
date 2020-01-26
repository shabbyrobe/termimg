package termimg

import (
	"fmt"
	"image"
	"image/color"
	"math/bits"

	"github.com/shabbyrobe/imgx/rgba"
)

type renderer struct {
	patternSet *PatternSet

	// Length 32 == 4x8 pixels per character cell.
	//
	// Bit layout for values:
	// 00..31 == count
	// 32..63 == ^color
	//
	// Add (1<<32) to add 1 to the count. This layout is used to allow sorting.
	// Color is inverted so that it sorts in the opposite order.
	colorsCount [32]uint64
	colorsSize  int
}

func newRenderer(patternSet *PatternSet) *renderer {
	rend := &renderer{patternSet: patternSet}
	return rend
}

func (rend *renderer) renderCells(into *CellData, rimg image.Image, flags Flag) error {
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

	n, xEnd, yEnd := 0, w-4, h-8

	if flags&HalfBlockOnly != 0 {
		for y := 0; y <= yEnd; y += 8 {
			for x := 0; x <= xEnd; x += 4 {
				into.Cells[n] = rend.pixelCharDataForCode(img, x, y, '▄', lowerHalfBitmap)
				n++
			}
		}

	} else {
		for y := 0; y <= yEnd; y += 8 {
			for x := 0; x <= xEnd; x += 4 {
				into.Cells[n] = rend.pixelCharData(img, x, y)
				n++
			}
		}
	}

	return nil
}

func (rend *renderer) renderEscapes(into *EscapeData, rimg image.Image, flags Flag) error {
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

	xEnd, yEnd := w-4, h-8

	if flags&HalfBlockOnly == 0 {
		for y := 0; y <= yEnd; y += 8 {
			for x := 0; x <= xEnd; x += 4 {
				into.put(flags, rend.pixelCharData(img, x, y))
			}

			// Don't print the last newline, so we can avoid scrolling when rendering video:
			if y < yEnd {
				into.nextRow()
			}
		}

	} else {
		for y := 0; y <= yEnd; y += 8 {
			for x := 0; x <= xEnd; x += 4 {
				into.put(flags, rend.pixelCharDataForCode(img, x, y, '▄', lowerHalfBitmap))
			}

			// Don't print the last newline, so we can avoid scrolling when rendering video:
			if y < yEnd {
				into.nextRow()
			}
		}
	}

	return nil
}

// Return a CharData struct with the given code point and corresponding average fg and bg colors.
func (r *renderer) pixelCharDataForCode(img *rgba.Image, x0, y0 int, code rune, pattern uint32) (result Cell) {
	result.Code = code

	var (
		fgCount = uint16(0)
		bgCount = uint16(0)
		mask    = uint32(0b_1000_0000_0000_0000_0000_0000_0000_0000)

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

// Find the best character and colors for a 4x8 part of the image at the given position
func (rend *renderer) pixelCharData(img *rgba.Image, x0, y0 int) (result Cell) {
	var minr, ming, minb uint32 = 0xFF, 0xFF, 0xFF
	var maxr, maxg, maxb uint32 = 0, 0, 0

	rend.colorsSize = 0

	// Determine the minimum and maximum value for each color channel
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

			for i := 0; i < rend.colorsSize; i++ {
				if uint32(rend.colorsCount[i]) == colorInv {
					// Increment the count, which is stored in the high 32-bits:
					rend.colorsCount[i] += 0x1_0000_0000
					goto next
				}
			}

			rend.colorsCount[rend.colorsSize] = 0x1_0000_0000 | uint64(colorInv)
			rend.colorsSize++
		next:
		}

		yOff += img.Stride
	}

	var count2 uint32
	var maxCountColor1 uint32
	var maxCountColor2 uint32

	if rend.colorsSize == 1 {
		count2 = uint32(rend.colorsCount[0] >> 32)
		maxCountColor1 = ^uint32(rend.colorsCount[0])
		maxCountColor2 = maxCountColor1

	} else {
		var max1, max2 uint64
		for i := 0; i < rend.colorsSize; i++ {
			rc := rend.colorsCount[i]
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

	var setBits uint32 = 0 // Important - keep as uint32
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

		// Compute a bitmap using the given split and sum the color values for both buckets.
		yN, xN := y0+8, x0+4
		for y := y0; y < yN; y++ {
			for x := x0; x < xN; x++ {
				setBits <<= 1

				c := img.Vals[y*img.Stride+x]

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
		}
	}

	// Find the best bitmap match by counting the bits that don't match,
	// including the inverted bitmaps.
	var bestDiff = 8 // FIXME: why 8 and not 16? not sure, need to research.
	var best = rend.patternSet.Default
	var inverted bool

	for _, pattern := range rend.patternSet.Patterns {
		pbits := pattern.Bits
		diff := bits.OnesCount32(pbits ^ setBits)
		if diff < bestDiff {
			best, bestDiff, inverted = pattern, diff, false
		}

		// Invert the pattern and try again:
		pbits = ^pbits
		diff = bits.OnesCount32(pbits ^ setBits)
		if diff < bestDiff {
			best, bestDiff, inverted = pattern, diff, true
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

	return rend.pixelCharDataForCode(img, x0, y0, best.Rune, best.Bits)
}
