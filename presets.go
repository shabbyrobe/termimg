package termimg

import (
	"image/color"
)

func Default() RendererConfig {
	return PresetBitmapBlock()
}

func PresetBitmap() BitmapConfig {
	return PresetBitmapBlock()
}

func PresetHalfBlock() HalfBlockConfig {
	return HalfBlockConfig{}
}

func PresetIntensity() IntensityConfig {
	return PresetIntensityChar()
}

func PresetSimpleChar() SimpleConfig {
	return SimpleConfig{'X'}
}

func PresetSimpleBlock() SimpleConfig {
	return SimpleConfig{'█'}
}

func PresetBitmapBlock() BitmapConfig {
	return BitmapConfig{
		Default: Bitmap{lowerHalfBitmap, '▄'},

		Bitmaps: []Bitmap{
			// Each column of four bits maps to one of 8 4-pixel rows inferred from each
			// unicode character. In this example, the bits map to the 'U+259A QUADRANT UPPER
			// LEFT AND LOWER RIGHT' character, '▚'.
			//
			//       ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━▶  ■ ■ · ·
			//       ┃    ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━▶  ■ ■ · ·
			//       ┃    ┃    ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━▶  ■ ■ · ·
			//       ┃    ┃    ┃    ┏━━━━━━━━━━━━━━━━━━━━━━▶  ■ ■ · ·
			//       ┃    ┃    ┃    ┃    ┏━━━━━━━━━━━━━━━━━▶  · · ■ ■
			//       ┃    ┃    ┃    ┃    ┃    ┏━━━━━━━━━━━━▶  · · ■ ■
			//       ┃    ┃    ┃    ┃    ┃    ┃    ┏━━━━━━━▶  · · ■ ■
			//       ┃    ┃    ┃    ┃    ┃    ┃    ┃    ┏━━▶  · · ■ ■
			//     ┏━┻┓ ┏━┻┓ ┏━┻┓ ┏━┻┓ ┏━┻┓ ┏━┻┓ ┏━┻┓ ┏━┻┓
			//
			// {0b_1100_1100_1100_1100_0011_0011_0011_0011,   '▚'},
			//

			{0b_0000_0000_0000_0000_0000_0000_0000_0000, 0x00a0}, // NO_BREAK_SPACE

			{0b_0000_0000_0000_0000_0000_0000_0000_1111, '▁'}, // lower 1/8
			{0b_0000_0000_0000_0000_0000_0000_1111_1111, '▂'}, // lower 1/4
			{0b_0000_0000_0000_0000_0000_1111_1111_1111, '▃'},
			{0b_0000_0000_0000_0000_1111_1111_1111_1111, '▄'}, // lower 1/2
			{0b_0000_0000_0000_1111_1111_1111_1111_1111, '▅'},
			{0b_0000_0000_1111_1111_1111_1111_1111_1111, '▆'}, // lower 3/4
			{0b_0000_1111_1111_1111_1111_1111_1111_1111, '▇'},

			{0b_1110_1110_1110_1110_1110_1110_1110_1110, '▊'}, // left 3/4
			{0b_1100_1100_1100_1100_1100_1100_1100_1100, '▌'}, // left 1/2
			{0b_1000_1000_1000_1000_1000_1000_1000_1000, '▎'}, // left 1/4

			{0b_0000_0000_0000_0000_1100_1100_1100_1100, '▖'}, // quadrant lower left
			{0b_0000_0000_0000_0000_0011_0011_0011_0011, '▗'}, // quadrant lower right
			{0b_1100_1100_1100_1100_0000_0000_0000_0000, '▘'}, // quadrant upper left
			{0b_1100_1100_1100_1100_0011_0011_0011_0011, '▚'}, // diagonal 1/2
			{0b_0011_0011_0011_0011_0000_0000_0000_0000, '▝'}, // quadrant upper right

			// Line drawing subset: no double lines, no complex light lines
			// Simple light lines duplicated because there is no center pixel int the 4x8 matrix
			{0b_0000_0000_0000_1111_1111_0000_0000_0000, '━'}, // Heavy horizontal
			{0b_0110_0110_0110_0110_0110_0110_0110_0110, '┃'}, // Heavy vertical

			{0b_0000_0000_0000_0111_0111_0110_0110_0110, '┏'}, // Heavy down and right
			{0b_0000_0000_0000_1110_1110_0110_0110_0110, '┓'}, // Heavy down and left
			{0b_0110_0110_0110_0111_0111_0000_0000_0000, '┗'}, // Heavy up and right
			{0b_0110_0110_0110_1110_1110_0000_0000_0000, '┛'}, // Heavy up and left

			{0b_0110_0110_0110_0111_0111_0110_0110_0110, '┣'}, // Heavy vertical and right
			{0b_0110_0110_0110_1110_1110_0110_0110_0110, '┫'}, // Heavy vertical and left
			{0b_0000_0000_0000_1111_1111_0110_0110_0110, '┳'}, // Heavy down and horizontal
			{0b_0110_0110_0110_1111_1111_0000_0000_0000, '┻'}, // Heavy up and horizontal
			{0b_0110_0110_0110_1111_1111_0110_0110_0110, '╋'}, // Heavy cross

			{0b_0000_0000_0000_1100_1100_0000_0000_0000, '╸'}, // Bold horizontal left
			{0b_0000_0000_0000_0110_0110_0000_0000_0000, '╹'}, // Bold horizontal up
			{0b_0000_0000_0000_0011_0011_0000_0000_0000, '╺'}, // Bold horizontal right
			{0b_0000_0000_0000_0110_0110_0000_0000_0000, '╻'}, // Bold horizontal down

			{0b_0000_0110_0110_0000_0000_0110_0110_0000, '╏'}, // Heavy double dash vertical

			{0b_0000_0000_0000_1111_0000_0000_0000_0000, '─'}, // Light horizontal
			{0b_0000_0000_0000_0000_1111_0000_0000_0000, '─'}, //
			{0b_0100_0100_0100_0100_0100_0100_0100_0100, '│'}, // Light vertical
			{0b_0010_0010_0010_0010_0010_0010_0010_0010, '│'},

			{0b_0000_0000_0000_1110_0000_0000_0000_0000, '╴'}, // light left
			{0b_0000_0000_0000_0000_1110_0000_0000_0000, '╴'}, // light left
			{0b_0100_0100_0100_0100_0000_0000_0000_0000, '╵'}, // light up
			{0b_0010_0010_0010_0010_0000_0000_0000_0000, '╵'}, // light up
			{0b_0000_0000_0000_0011_0000_0000_0000_0000, '╶'}, // light right
			{0b_0000_0000_0000_0000_0011_0000_0000_0000, '╶'}, // light right
			{0b_0000_0000_0000_0000_0100_0100_0100_0100, '╷'}, // light down
			{0b_0000_0000_0000_0000_0010_0010_0010_0010, '╷'}, // light down

			{0b_0100_0100_0100_0100_0100_0100_0100_0100, 0x23a2}, // [ extension
			{0b_0010_0010_0010_0010_0010_0010_0010_0010, 0x23a5}, // ] extension

			{0b_0000_1111_0000_0000_0000_0000_0000_0000, 0x23ba}, // Horizontal scanline 1
			{0b_0000_0000_1111_0000_0000_0000_0000_0000, 0x23bb}, // Horizontal scanline 3
			{0b_0000_0000_0000_0000_0000_1111_0000_0000, 0x23bc}, // Horizontal scanline 7
			{0b_0000_0000_0000_0000_0000_0000_1111_0000, 0x23bd}, // Horizontal scanline 9

			{0b_0000_0000_0000_0110_0110_0000_0000_0000, 0x25aa}, // Black small square

			// ## Unused:

			// {0xffff0000, '▀'},  // upper 1/2; redundant with inverse lower 1/2
			// {0xffffffff, '█'},  // full; redundant with inverse space

			// {0xccccffff, '▙'},  // 3/4 redundant with inverse 1/4
			// {0xffffcccc, '▛'},  // 3/4 redundant
			// {0xffff3333, '▜'},  // 3/4 redundant
			// {0x3333cccc, '▞'},  // 3/4 redundant
			// {0x3333ffff, '▟'},  // 3/4 redundant

			// ## Geometrical shapes. Tricky because some of them are too wide.

			// {0x00ffff00, 0x25fe},  // Black medium small square
			// {0x11224488, 0x2571},  // diagonals
			// {0x88442211, 0x2572},
			// {0x99666699, 0x2573},
			// {0x000137f0, 0x25e2},  // Triangles
			// {0x0008cef0, 0x25e3},
			// {0x000fec80, 0x25e4},
			// {0x000f7310, 0x25e5},
		},
	}
}

func PresetIntensityChar() IntensityConfig {
	return IntensityConfig{
		Fg: color.RGBA{0xff, 0xff, 0xff, 0xff},
		Bg: color.RGBA{0x00, 0x00, 0x00, 0x00},

		Intensities: []Intensity{
			{Brightness: 0x00, Rune: ' '},
			{Brightness: 0x18, Rune: '.'},
			{Brightness: 0x28, Rune: ','},
			{Brightness: 0x30, Rune: ':'},
			{Brightness: 0x3c, Rune: '"'},
			{Brightness: 0x40, Rune: ';'},
			{Brightness: 0x4d, Rune: '!'},
			{Brightness: 0x64, Rune: 'r'},
			{Brightness: 0x6c, Rune: '<'},
			{Brightness: 0x74, Rune: 'l'},
			{Brightness: 0x78, Rune: 'c'},
			{Brightness: 0x7c, Rune: 'i'},
			{Brightness: 0x80, Rune: 'v'},
			{Brightness: 0x83, Rune: 'j'},
			{Brightness: 0x87, Rune: 'f'},
			{Brightness: 0x8f, Rune: 'J'},
			{Brightness: 0x97, Rune: 'C'},
			{Brightness: 0xa3, Rune: 'I'},
			{Brightness: 0xab, Rune: 'V'},
			{Brightness: 0xaf, Rune: 'k'},
			{Brightness: 0xb7, Rune: 'X'},
			{Brightness: 0xc3, Rune: 'P'},
			{Brightness: 0xc7, Rune: 'G'},
			{Brightness: 0xcb, Rune: 'U'},
			{Brightness: 0xd3, Rune: 'K'},
			{Brightness: 0xd7, Rune: 'O'},
			{Brightness: 0xdb, Rune: 'H'},
			{Brightness: 0xdf, Rune: 'D'},
			{Brightness: 0xe7, Rune: 'R'},
			{Brightness: 0xf7, Rune: 'W'},
			{Brightness: 0xfb, Rune: 'N'},
			{Brightness: 0xff, Rune: 'M'},
		},
	}
}

// XXX: This will change at some point; it doesn't work well without
// the ability to disable the background color.
func PresetBrailleBitmap() *BitmapConfig {
	// Braille dots have a weird numbering scheme:
	//
	//     ╭───────▶ 1  ■ ·  4 ◄─────╮
	//     │ ╭─────▶ 2  · ■  5 ◄───╮ │
	//     │ │ ╭───▶ 3  ■ ■  6 ◄─╮ │ │
	//     │ │ │     7  ■ ·  8   │ │ │
	//     │ │ │        ▲ ▲      │ │ │
	//     │ │ │        │ │      6 5 4
	//     │ │ │  ╭─────│─╯      │ │ │
	//     │ │ │  │ ╭───╯        │ │ │
	//     │ │ │  │ │ ╭──────────╯ │ │
	//     1 2 3  │ │ │ ╭──────────╯ │
	//     │ │ │  │ │ │ │ ╭──────────╯
	//     │ │ │  │ │ │ │ │
	//     │ │ │  ▼ ▼ ▼ ▼ ▼
	//     │ │ │  8 7 6 5 4 3 2 1  ==  braille dot number
	//     │ │ │  0 1 1 1 0 1 0 1  ==  bit value
	//     │ │ │            ▲ ▲ ▲
	//     │ │ │            │ │ │
	//     │ │ ╰─────3──────╯ │ │
	//     │ ╰───────2────────╯ │
	//     ╰─────────1──────────╯
	//
	// Add the bit value to U+2800 to get the corresponding character.

	var bitmaps = make([]Bitmap, 0, 257)

	bitmaps = append(bitmaps, Bitmap{0, 0x00a0}) // NO_BREAK_SPACE

	for i := 0; i < 256; i++ {
		var bmp Bits

		if i&0b_0000_0001 != 0 { // 0, 0 == #1
			bmp |= 0b_1100_1100_0000_0000_0000_0000_0000_0000
		}
		if i&0b_0000_1000 != 0 { // 1, 0 == #4
			bmp |= 0b_0011_0011_0000_0000_0000_0000_0000_0000
		}
		if i&0b_0000_0010 != 0 { // 0, 1 == #2
			bmp |= 0b_0000_0000_1100_1100_0000_0000_0000_0000
		}
		if i&0b_0001_0000 != 0 { // 1, 1 == #5
			bmp |= 0b_0000_0000_0011_0011_0000_0000_0000_0000
		}
		if i&0b_0000_0100 != 0 { // 0, 2 == #3
			bmp |= 0b_0000_0000_0000_0000_1100_1100_0000_0000
		}
		if i&0b_0010_0000 != 0 { // 1, 2 == #6
			bmp |= 0b_0000_0000_0000_0000_0011_0011_0000_0000
		}
		if i&0b_0100_0000 != 0 { // 0, 3 == #7
			bmp |= 0b_0000_00000000_0000_0000_0000_1100_1100
		}
		if i&0b_1000_0000 != 0 { // 1, 3 == #8
			bmp |= 0b_0000_00000000_0000_0000_0000_0011_0011
		}

		bitmaps = append(bitmaps, Bitmap{
			Bits: bmp,
			Rune: rune(0x2800) + rune(i),
		})
	}

	// FIXME: this can be converted to its own renderer.
	return &BitmapConfig{
		Default: Bitmap{lowerHalfBitmap, 0x28E4},
		Bitmaps: bitmaps,
	}
}
