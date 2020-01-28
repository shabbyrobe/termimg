package termimg

import (
	"fmt"
	"image/color"
	"math"

	"github.com/shabbyrobe/imgx/rgba"
)

type Intensity struct {
	Brightness uint8 // Minimum brightness for this intensity
	Rune       rune
}

type IntensityRenderer struct {
	intensities []Intensity
	runes       [256]rune
}

func IntensityRendererFromChars(chars string) (*IntensityRenderer, error) {
	if len(chars) == 0 {
		return nil, fmt.Errorf("termimg: no intensities")
	}

	// NOTE: this potentially slightly over-allocates; chars may contain UTF-8 codepoints.
	is := make([]Intensity, 0, len(chars))

	per := float64(256) / float64(len(chars))
	cur := float64(0)
	for _, c := range chars {
		is = append(is, Intensity{Brightness: uint8(math.Round(cur)), Rune: c})
		cur += per
	}
	return NewIntensityRenderer(is)
}

func NewIntensityRenderer(intensities []Intensity) (*IntensityRenderer, error) {
	if len(intensities) == 0 {
		return nil, fmt.Errorf("termimg: no intensities")
	}
	if intensities[0].Brightness != 0 {
		return nil, fmt.Errorf("termimg: intensity 1 must have brightness of 0")
	}

	first := intensities[0]
	last := first.Brightness
	for _, i := range intensities[1:] {
		if last >= i.Brightness {
			return nil, fmt.Errorf("termimg: intensity must be sorted by brightness and contain no duplicates")
		}
		last = i.Brightness
	}

	is := &IntensityRenderer{
		intensities: intensities,
	}

	r, last := first.Rune, first.Brightness
	for _, in := range intensities[1:] {
		for i := last; i < in.Brightness; i++ {
			is.runes[i] = r
		}
		r, last = in.Rune, in.Brightness
	}
	for i := int(last); i < 256; i++ {
		is.runes[i] = r
	}

	return is, nil
}

func (intr *IntensityRenderer) cell(rend *imageRenderer, img *rgba.Image, x0, y0 int) (result Cell) {
	var sumR, sumG, sumB, sumV uint

	yN, xN, yOff := y0+8, x0+4, y0*img.Stride

	for y := y0; y < yN; y++ {
		for x := x0; x < xN; x++ {
			c := img.Vals[yOff+x]
			sumR += uint(c.R)
			sumG += uint(c.G)
			sumB += uint(c.B)

			max := c.R
			if c.G > max {
				max = c.G
			} else if c.B > max {
				max = c.B
			}
			sumV += uint(max)
		}
	}

	result.FgColor = color.RGBA{
		R: uint8(sumR >> 5),
		G: uint8(sumG >> 5),
		B: uint8(sumB >> 5),
		A: 0xff,
	}
	result.Code = intr.runes[uint8(sumV>>5)]
	return result
}
