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

func TestMergeThemeList(t *testing.T) {
	Convey("Merge theme list", t, func() {
		src := []string{"Deepin", "Adwaita", "Zukitwo"}
		target := []string{"Deepin", "Evolve"}
		ret := []string{"Deepin", "Adwaita", "Zukitwo", "Evolve"}

		So(mergeThemeList(src, target), ShouldResemble, ret)
	})
}

func TestSetQt4Theme(t *testing.T) {
	Convey("Set qt4 theme", t, func() {
		config := "/tmp/Trolltech.conf"
		So(setQt4Theme(config), ShouldEqual, true)
		os.Remove(config)
	})
}
