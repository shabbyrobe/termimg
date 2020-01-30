package termimg

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func parseRune(s string) (r rune, err error) {
	r, sz := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return r, fmt.Errorf("termimg: rune error in string %q", s)
	}
	if sz != len(s) {
		nv, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return r, fmt.Errorf("termimg: rune must be either a single UTF-8 character, or a Go-compatible number literal (0x, 0o, 0b supported), found %q", s)
		}
		r = rune(nv)
	}
	return r, nil
}

func parseBitmap(s string) (p Bitmap, err error) {
	const runeField, bitmapField = 1, 0

	fields := strings.SplitN(s, ":", 2)
	if len(fields) != 2 {
		return p, fmt.Errorf("termimg: expected format '<bitmap>:<rune>'")
	}

	if fields[runeField] == "" {
		return p, fmt.Errorf("termimg: expected format '<bitmap>:<rune>'")
	}

	p.Rune, err = parseRune(fields[runeField])
	if err != nil {
		return p, err
	}

	bmpstr := fields[bitmapField]
	bmpstr = strings.TrimSpace(bmpstr)
	bmpstr = strings.Replace(bmpstr, "_", "", -1)
	bmpstr = strings.Replace(bmpstr, ".", "0", -1)
	bmpstr = strings.TrimPrefix(bmpstr, "0b")
	bmpstr = strings.TrimPrefix(bmpstr, "0B")

	bmp, err := strconv.ParseInt(bmpstr, 2, 64)
	if err != nil {
		return p, fmt.Errorf("termimg: could not parse bitmap %q as binary number; should be in the format 0b0000", bmpstr)
	}

	if bmp < 0 || bmp >= 1<<32 {
		return p, fmt.Errorf("termimg: bitmap must be between 0 and (1<<31)-1; found %q", fields[bitmapField])
	}

	p.Bits = uint32(bmp)

	return p, nil
}
