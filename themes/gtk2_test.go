/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package themes

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestGtk2Infos(t *testing.T) {
	Convey("Test gtk2 infos", t, func() {
		infos := gtk2FileReader("testdata/gtkrc-2.0")
		So(len(infos), ShouldEqual, 16)

		info := infos.Get("gtk-theme-name")
		So(info.value, ShouldEqual, "\"Paper\"")

		info.value = "\"Deepin\""
		So(info.value, ShouldEqual, "\"Deepin\"")

		infos = infos.Add("gtk2-test", "test")
		So(len(infos), ShouldEqual, 17)
	})

	Convey("Test nil infos", t, func() {
		var infos = gtk2FileReader("testdata/xxx")
		infos = infos.Add("gtk2-test", "test")
		So(len(infos), ShouldEqual, 1)
		info := infos.Get("gtk2-test")
		So(info.value, ShouldEqual, "test")

		err := gtk2FileWriter(infos, "testdata/tmp-gtk2rc")
		defer os.Remove("testdata/tmp-gtk2rc")
		So(err, ShouldBeNil)
	})
}
