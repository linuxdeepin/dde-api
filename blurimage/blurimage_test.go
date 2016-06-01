package blurimage

import (
	"github.com/disintegration/imaging"
	"testing"
)

func BenchmarkIsTooBright(b *testing.B) {
	b.ReportAllocs()

	img, err := imaging.Open("testdata/test1.jpg")
	if err != nil {
		return
	}

	for i := 0; i < b.N; i++ {
		isTooBright(img)
	}
}

func TestIsTooBright(t *testing.T) {
	img, err := imaging.Open("testdata/test1.jpg")
	if err != nil {
		t.Error(err)
	}

	if isTooBright(img) {
		t.Error("Judge for test1.jpg is not correct.")
	}

	img, err = imaging.Open("testdata/test2.jpg")
	if err != nil {
		t.Error(err)
	}

	if !isTooBright(img) {
		t.Error("Judge for test2.jpg is not correct.")
	}
}
