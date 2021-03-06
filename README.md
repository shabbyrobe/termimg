Terminal Image Renderer
=======================

Go implementation of Stefan Haustein's
[TerminalImageViewer](https://github.com/stefanhaustein/TerminalImageViewer/)
algorithm.

NOTE: the API is a work in progress. The basic case (default bitmaps, no flags)
should not change much, but everything else is subject to change. If you
require stability, vendor this library into your repo.


## Examples

Default mode (using the default PatternSet):

![termimg-1](https://user-images.githubusercontent.com/288426/72950578-e78cff00-3ddf-11ea-80fc-b43194f1d65c.png)


Using the '▄' character only (via the `termimg.HalfBlockOnly` flag):

![termimg-2](https://user-images.githubusercontent.com/288426/72950585-ed82e000-3ddf-11ea-9972-f4a89941b989.png)


## Quickstart

First, you need an `image.Image`. See `png.Decode` or `jpeg.Decode` to get something
quickly.

Secondly, depending on your use-case, you need to decide if you want to use
`termimg.CellData` or `termimg.EscapeData`.

`EscapeData` can be written directly to your terminal to show a rendered image.

`CellData` is used if you want access to the foreground and background color as `color.RGBA`
values, and the output character as a rune. This may be desirable if you are using something
like https://github.com/gdamore/tcell

Thirdly, you need to decide which renderer you want to use. There are several:

- BitmapRenderer: the full color, "hi-res" TerminalImageViewer algorithm.
- IntensityRenderer: black and white renderer using different characters for intensity
- HalfBlockRenderer: color renderer using a unicode half-block. Requires background color.
- SimpleRenderer: color renderer using a single character. Doesn't require background color.

There are several presets available using the `Preset*()` functions. These examples will
use `PresetBitmapBlock()`, which uses the TerminalImageViewer algorithm and its pattern set.

To render an `image.Image` into an `EscapeData`, then print to stdout:

```go
var img image.Image // presuming you already have one of these...

var data termimg.EscapeData
var renderer, _ = termimg.PresetBitmapBlock().Renderer()
err := renderer.Escapes(&data, img, 0)
os.Stdout.Write(data.Value())

// Clean up afterwards:
os.Stdout.Write([]byte("\033[0m"))
os.Stdout.Write([]byte("\n"))
```

To render into a `CellData` into a `tcell.Screen`:

```go
var img image.Image // presuming you already have one of these...

screen, err := tcell.NewScreen()
err := screen.Init()
screen.Clear()
defer screen.Fini()

var cells termimg.CellData
var renderer, _ = termimg.PresetBitmapBlock().Renderer()
err := renderer.Cells(&cells, img, 0)

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
```
