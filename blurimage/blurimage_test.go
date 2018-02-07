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
