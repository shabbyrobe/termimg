package termimg

import (
	"image"
	"math"
)

// StretchToCellSize returns the size an image needs to be scaled to in order to
// compensate for the distortion introduced by the terminal.
//
// StretchToCellSize always expands the image. Once you have the correct distortion, you
// can then resize based on your expected column width/height.
//
// The renderer breaks individual terminal cells up into a 4-by-8 grid, but depending on
// your terminal emulator or font size, they may be more like 8px-by-17px. If you render
// without compensating, the image will appear vertically stretched.
//
func StretchToCellSize(cellWidth, cellHeight float64, imgSize image.Point) image.Point {
	xPixW := cellWidth / 4
	if xPixW == 0 {
		xPixW = 1
	}
	yPixH := cellHeight / 8
	vRatio := xPixW / yPixH

	newW, newH := float64(imgSize.X), float64(imgSize.Y)
	if vRatio < 1 {
		newW = newW / vRatio
	} else {
		newH = newH * vRatio
	}

	return image.Point{int(math.Round(newW)), int(math.Round(newH))}
}
