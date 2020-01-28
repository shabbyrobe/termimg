package termimg

import (
	"fmt"
	"strconv"
	"strings"
)

func firstRune(s string) (r rune, n int, err error) {
	for n, r = range s {
		return r, n + 1, nil
	}
	return 0, 0, fmt.Errorf("termimg: rune not found")
}

func parseBitmap(s string) (p Bitmap, err error) {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, " ", "", -1)

	fields := strings.SplitN(s, ":", 2)
	if len(fields) != 2 {
		return p, fmt.Errorf("termimg: expected format '<rune>:<bitmap>'")
	}

	switch len(fields[0]) {
	case 0:
		return p, fmt.Errorf("termimg: expected format '<rune>:<bitmap>'")
	case 1:
		r, n, err := firstRune(fields[0])
		if err != nil {
			return p, err
		} else if n != len(fields[0]) {
			return p, fmt.Errorf("termimg: invalid bitmap %q: unexpected text after rune: format should be '<rune>:<bitmap>'", s)
		}
		p.Rune = r
	default:
		r, err := strconv.ParseInt(fields[0], 0, 0)
		if err != nil {
			return p, fmt.Errorf("termimg: invalid bitmap %q: could not parse rune as number; should be decimal or start with 0x, 0o or 0b: %w", s, err)
		}
		p.Rune = rune(r)
	}

	bmpstr := fields[1]
	bmpstr = strings.Replace(bmpstr, "_", "", -1)
	bmpstr = strings.Replace(bmpstr, ".", "0", -1)
	bmpstr = strings.TrimPrefix(bmpstr, "0b")
	bmpstr = strings.TrimPrefix(bmpstr, "0B")

	bmp, err := strconv.ParseInt(bmpstr, 2, 64)
	if err != nil {
		return p, fmt.Errorf("termimg: could not parse bitmap %q as binary number; should be in the format 0b0000", bmpstr)
	}

	if bmp < 0 || bmp >= 1<<32 {
		return p, fmt.Errorf("termimg: bitmap must be between 0 and (1<<31)-1; found %q", fields[1])
	}

	p.Bits = uint32(bmp)

	return p, nil
}
