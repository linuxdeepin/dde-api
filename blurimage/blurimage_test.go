// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package blurimage

import (
	"os"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
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
	t.Skip("always skip")
	tests := []struct {
		Input  string
		Expect bool
	}{
		{
			"testdata/test1.jpg",
			false,
		},
		{
			"testdata/test2.jpg",
			true,
		},
	}
	for _, data := range tests {
		img, err := imaging.Open(data.Input)
		if err != nil {
			t.Error(err)
		}
		if !assert.Equal(t, data.Expect, isTooBright(img)) {
			t.Errorf("Judge for %s is not correct.", data.Input)
		}
	}
}

func TestBlurImage(t *testing.T) {
	tests := []struct {
		file  string
		sigma float64
		dest  string
	}{
		{
			"testdata/test1.jpg",
			20,
			"testdata/test1_blur.png",
		},
		{
			"testdata/test2.jpg",
			30,
			"testdata/test2_blur.png",
		},
	}
	for _, data := range tests {
		err := BlurImage(data.file, data.sigma, data.dest)
		assert.NoError(t, err)
		assert.FileExists(t, data.dest)
		_ = os.Remove(data.dest)
	}
}
