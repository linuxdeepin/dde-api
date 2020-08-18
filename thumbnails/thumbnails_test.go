/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package thumbnails

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"pkg.deepin.io/dde/api/thumbnails/loader"
)

func TestCorrectSize(t *testing.T) {
	Convey("Test size correct", t, func(c C) {
		c.So(correctSize(64), ShouldEqual, loader.SizeFlagSmall)
		c.So(correctSize(128), ShouldEqual, loader.SizeFlagNormal)
		c.So(correctSize(176), ShouldEqual, loader.SizeFlagNormal)
		c.So(correctSize(256), ShouldEqual, loader.SizeFlagLarge)
	})
}
