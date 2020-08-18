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

package themes

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMergeThemeList(t *testing.T) {
	Convey("Merge theme list", t, func(c C) {
		src := []string{"Deepin", "Adwaita", "Zukitwo"}
		target := []string{"Deepin", "Evolve"}
		ret := []string{"Deepin", "Adwaita", "Zukitwo", "Evolve"}

		c.So(mergeThemeList(src, target), ShouldResemble, ret)
	})
}

func TestSetQt4Theme(t *testing.T) {
	Convey("Set qt4 theme", t, func(c C) {
		config := "/tmp/Trolltech.conf"
		c.So(setQt4Theme(config), ShouldEqual, true)
		os.Remove(config)
	})
}
