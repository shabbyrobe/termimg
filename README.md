Terminal Image Renderer
=======================

Go implementation of Stefan Haustein's
[TerminalImageViewer](https://github.com/stefanhaustein/TerminalImageViewer/)
algorithm.

NOTE: the API is a work in progress. The basic case (default bitmaps, no flags)
should not change much, but everything else is subject to change. If you
require stability, vendor this library into your repo.


## Quickstart

First, you need an `image.Image`. See `png.Decode` or `jpeg.Decode` to get something
quickly.

Secondly, depending on your use-case, you need to decide if you want to use
`termimg.CellData` or `termimg.EscapeData`.

`EscapeData` can be written directly to your terminal to show a rendered image.

`CellData` is used if you want access to the foreground and background color as `color.RGBA`
values, and the output character as a rune. This may be desirable if you are using something
like https://github.com/gdamore/tcell

To render an `image.Image` into an `EscapeData`, then print to stdout:

    var img image.Image // presuming you already have one of these...

	var data EscapeData
    if err := Encode(&data, img, 0, nil); err != nil {
        // ...
    }
    os.Stdout.Write(data.Value())

To render into a `CellData` into a `tcell.Screen`:

    var img image.Image // presuming you already have one of these...

    screen, err := tcell.NewScreen()
    err := screen.Init()
    screen.Clear()
    defer screen.Fini()

	var cells CellData
    if err := EncodeCells(&cells, img, 0, nil); err != nil {
        // ...
    }

	for x := 0; x < cells.Cols; x++ {
		for y := 0; y < cells.Rows; y++ {
			cell := cells.CellAt(x, y)
			style := tcell.StyleDefault.
				Foreground(tcell.NewRGBColor(cell.FgRGB32())).
				Background(tcell.NewRGBColor(cell.BgRGB32()))

			scr.SetCell(x, y, style, cell.Code)
		}
	}

    // Show the image for 2 seconds then quit:
    time.Sleep(2 * time.Second)

