package termimg

import (
	"image/color"
)

func term256AsRGB(tc uint8) (r, g, b uint8) {
	// FIXME: can this use termpalette?
	if tc < 16 {
		c := term16col[tc]
		return c.R, c.G, c.B

	} else if tc >= 232 {
		lv := ((tc - 232) * 10) + 8
		return lv, lv, lv

	} else {
		tc -= 16
		r = (tc / 36)
		g = (tc % 36) / 6
		b = (tc % 6)

		return r*40 + 55, g*40 + 55, b*40 + 55
	}
}

var term16col = [256]color.RGBA{
	{0, 0, 0, 255},
	{128, 0, 0, 255},
	{0, 128, 0, 255},
	{128, 128, 0, 255},
	{0, 0, 128, 255},
	{128, 0, 128, 255},
	{0, 128, 128, 255},
	{192, 192, 192, 255},
	{128, 128, 128, 255},
	{255, 0, 0, 255},
	{0, 255, 0, 255},
	{255, 255, 0, 255},
	{0, 0, 255, 255},
	{255, 0, 255, 255},
	{0, 255, 255, 255},
	{255, 255, 255, 255},
}
