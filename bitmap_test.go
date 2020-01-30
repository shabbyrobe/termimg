package termimg

import (
	"reflect"
	"testing"
)

func TestBitmapMarshal(t *testing.T) {
	b1 := Bitmap{
		Bits: 1<<32 - 1,
		Rune: 'q',
	}

	m, err := b1.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	var b2 Bitmap
	if err := b2.UnmarshalText(m); err != nil {
		t.Fatal(err)
	}

	reflect.DeepEqual(b1, b2)
}
