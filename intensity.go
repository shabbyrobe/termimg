package termimg

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/shabbyrobe/imgx/rgba"
)

type Intensity struct {
	Brightness uint8 // Minimum brightness for this intensity
	Rune       rune
}

func (p Intensity) String() string {
	bin := fmt.Sprintf("0x%02x", p.Brightness)
	return fmt.Sprintf("%s:%s", bin, string(p.Rune))
}

func (p Intensity) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *Intensity) UnmarshalText(text []byte) (err error) {
	parts := strings.SplitN(string(text), ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("termimg: invalid intensity %q, expected format <intesity>:<rune>", string(text))
	}

	v, err := strconv.ParseUint(parts[0], 0, 8)
	if err != nil {
		return err
	}

	p.Brightness = uint8(v)

	p.Rune, err = parseRune(parts[1])
	if err != nil {
		return err
	}

	return err
}

type IntensityRenderer struct {
	fg, bg      color.RGBA
	intensities []Intensity
	runes       [256]rune
}

// IntensityRendererFromChars constructs an IntensityRenderer using each char
// in chars as a step from 0 to 255. Intensity is calculated the same way as
// the 'V' or 'B' in 'HSV/B'.
//
// Some strings that have also worked (leading space is intentional):
//   " .·-:;+=xtm$X&@#M"
//   " ·-+xwXM"
//
// This is font-dependent though; it's worth exploring your own charset using
// a black and white gradient like this one: http://www.lagom.nl/lcd-test/gradient.php
//
func IntensityRendererFromChars(fg, bg color.RGBA, chars string) (*IntensityRenderer, error) {
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
	return NewIntensityRenderer(fg, bg, is)
}

func NewIntensityRenderer(fg, bg color.RGBA, intensities []Intensity) (*IntensityRenderer, error) {
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
		fg:          fg,
		bg:          bg,
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
	var sumV int32

	yN, xN, yOff := y0+8, x0+4, y0*img.Stride

	for y := y0; y < yN; y++ {
		for x := x0; x < xN; x++ {
			c := img.Vals[yOff+x]
			max := c.R
			if c.G > max {
				max = c.G
			} else if c.B > max {
				max = c.B
			}
			sumV += int32(max)
		}
	}

	result.FgColor = intr.fg
	result.BgColor = intr.bg
	result.Code = intr.runes[sumV>>5]
	return result
}
