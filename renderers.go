package termimg

var (
	DefaultRenderer          CellRenderer = BlockRenderer
	DefaultBitmapRenderer    CellRenderer = BlockRenderer
	DefaultIntensityRenderer CellRenderer // init()

	BlockRenderer     CellRenderer = blockRenderer
	HalfBlockRenderer CellRenderer = halfBlockRenderer
)

func init() {
	var err error
	DefaultIntensityRenderer, err = IntensityRendererFromChars(" .:;+=xX$&#")
	if err != nil {
		panic(err)
	}
}

var halfBlockRenderer = &BitmapRenderer{
	Default: Bitmap{lowerHalfBitmap, '▄'},
	Bitmaps: []Bitmap{{lowerHalfBitmap, '▄'}},
}

// blockRenderer is exported indirectly so it doesn't pollute the godoc
var blockRenderer = &BitmapRenderer{
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
