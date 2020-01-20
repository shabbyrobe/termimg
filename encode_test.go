package termimg

import (
	"math/rand"
	"testing"

	"github.com/shabbyrobe/imgx/testimg"
)

func BenchmarkTextImageRender(b *testing.B) {
	r := rand.New(rand.NewSource(0))

	var data EscapeData
	data.SetBuffer(make([]byte, 10*1024*1024)) //  Should be big enough to avoid realloc

	var cells = CellDataFromPixels(512, 512)

	img := testimg.RandBlocks{W: 512, H: 512, BlockW: 1, BlockH: 1}.RGBA(r)
	b.Run("rgb-1x1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := Encode(&data, img, NoAlloc, nil); err != nil {
				panic(err)
			}
		}
	})

	b.Run("rgb-1x1-cells", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := EncodeCells(&cells, img, NoAlloc, nil); err != nil {
				panic(err)
			}
		}
	})

	img = testimg.RandBlocks{W: 512, H: 512, BlockW: 10, BlockH: 10}.RGBA(r)
	b.Run("rgb-10x10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := Encode(&data, img, NoAlloc, nil); err != nil {
				panic(err)
			}
		}
	})

	img = testimg.RandBlocks{W: 512, H: 512, BlockW: 1, BlockH: 1}.RGBA(r)
	b.Run("256-1x1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := Encode(&data, img, NoAlloc|Color256, nil); err != nil {
				panic(err)
			}
		}
	})

	img = testimg.RandBlocks{W: 512, H: 512, BlockW: 10, BlockH: 10}.RGBA(r)
	b.Run("256-10x10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := Encode(&data, img, NoAlloc|Color256, nil); err != nil {
				panic(err)
			}
		}
	})
}
