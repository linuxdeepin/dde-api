/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
