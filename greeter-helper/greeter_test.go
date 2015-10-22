/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFormatLayout(t *testing.T) {
	Convey("Format layout", t, func() {
		So(formatLayout("us"), ShouldEqual, "")
		So(formatLayout("us;"), ShouldEqual, "us|")
		So(formatLayout("af;ps"), ShouldEqual, "af|ps")
	})
}

func TestFormatLayoutList(t *testing.T) {
	Convey("Format layout list", t, func() {
		var infos = []struct {
			list []string
			ret  string
		}{
			{
				list: []string{"us", "af"},
				ret:  "",
			},
			{
				list: []string{"us;", "af;ps"},
				ret:  "us| af|ps",
			},
		}

		for _, info := range infos {
			So(formatLayoutList(info.list), ShouldEqual, info.ret)
		}
	})
}

func TestDoSet(t *testing.T) {
	Convey("DoSet", t, func() {
		tmpFile := "tmp_test.ini"
		defer os.Remove(tmpFile)

		err := doSet(tmpFile, "deepin", kfKeyLayout, "us|")
		if err != nil {
			return
		}
		content, _ := ioutil.ReadFile(tmpFile)
		So(string(content), ShouldEqual, `[deepin]
KeyboardLayout=us|
`)

		So(doSet(tmpFile, "deepin", kfKeyLayoutList, "us| af|ps"),
			ShouldBeNil)
		content, _ = ioutil.ReadFile(tmpFile)
		So(string(content), ShouldEqual, `[deepin]
KeyboardLayout=us|
KeyboardLayoutList=us| af|ps
`)

		So(doSet(tmpFile, "deepin", kfKeyTheme, "sky"),
			ShouldBeNil)
		content, _ = ioutil.ReadFile(tmpFile)
		So(string(content), ShouldEqual, `[deepin]
KeyboardLayout=us|
KeyboardLayoutList=us| af|ps
GreeterTheme=sky
`)
	})
}
