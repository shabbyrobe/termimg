package termimg

import (
	"fmt"
	"image"
	"math/rand"
	"reflect"
	"testing"

	"github.com/shabbyrobe/imgx/rgba"
	"github.com/shabbyrobe/imgx/testimg"
)

func TestDecodeImage(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	for idx, tc := range []struct {
		name string
		img  image.Image
	}{
		{"rgb1x1", testimg.RandBlocks{W: 512, H: 512, BlockW: 1, BlockH: 1}.RGBA(r)},
		{"rgb4x4", testimg.RandBlocks{W: 512, H: 512, BlockW: 4, BlockH: 4}.RGBA(r)},
		{"rgb8x8", testimg.RandBlocks{W: 512, H: 512, BlockW: 8, BlockH: 8}.RGBA(r)},
		{"rgb16x16", testimg.RandBlocks{W: 512, H: 512, BlockW: 16, BlockH: 16}.RGBA(r)},
	} {
		t.Run(fmt.Sprintf("%s/%d", tc.name, idx), func(t *testing.T) {
			// Brute-force check of Decode: encodes to runes, then decodes back to an
			// image. Then takes that image, encodes to runes again, then back to an
			// image again, and makes sure the two images are the same.

			var data EscapeData
			if err := Encode(&data, tc.img, 0, nil); err != nil {
				panic(err)
			}

			first, err := DecodeImageBytes(data.Value(), nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			var back EscapeData
			if err := Encode(&back, first, 0, nil); err != nil {
				panic(err)
			}

			second, err := DecodeImageBytes(back.Value(), nil, nil)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(first.Vals, second.Vals) {
				t.Fatal()
			}
		})
	}
}

func TestDecodeCells(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	for idx, tc := range []struct {
		name string
		img  image.Image
	}{
		{"rgb8x8", testimg.RandBlocks{W: 16, H: 16, BlockW: 8, BlockH: 8}.RGBA(r)},
	} {
		t.Run(fmt.Sprintf("%s/%d", tc.name, idx), func(t *testing.T) {
			// Brute-force check of Decode: encodes to runes, then decodes back to an
			// image. Then takes that image, encodes to runes again, then back to an
			// image again, and makes sure the two images are the same.

			var data EscapeData
			if err := Encode(&data, tc.img, 0, nil); err != nil {
				panic(err)
			}

			var cells CellData
			if err := EncodeCells(&cells, tc.img, 0, nil); err != nil {
				panic(err)
			}

			back, err := DecodeCellsBytes(data.Value(), nil)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(cells, back) {
				t.Fatal()
			}
		})
	}
}

var (
	BenchRGBAImage *rgba.Image
)

func BenchmarkDecode(b *testing.B) {
	rng := rand.New(rand.NewSource(0))

	sz := image.Point{512, 512}
	img := testimg.RandBlocks{W: sz.X, H: sz.Y, BlockW: 1, BlockH: 1}.RGBA(rng)

	var data EscapeData
	data.SetBuffer(make([]byte, 10*1024*1024)) //  Should be big enough to avoid realloc
	if err := Encode(&data, img, NoAlloc, nil); err != nil {
		panic(err)
	}

	b.Run("rgb-1x1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var err error
			BenchRGBAImage, err = DecodeImageBytes(data.Value(), nil, &sz)
			if err != nil {
				panic(err)
			}
		}
	})
}
