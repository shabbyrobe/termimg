package termimg

import (
	"github.com/shabbyrobe/imgx/rgba"
)

var (
	index256 rgba.Index
	index16  rgba.Index
)

func init() {
	idxr := rgba.NewRGBPrecacheIndexer(nil).(rgba.IndexUnmarshaler)

	var err error
	index256, err = idxr.UnmarshalIndex(index256Data)
	if err != nil {
		panic(err)
	}

	index16, err = idxr.UnmarshalIndex(index16Data)
	if err != nil {
		panic(err)
	}
}
