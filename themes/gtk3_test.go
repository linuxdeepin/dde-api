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
	dutils "pkg.deepin.io/lib/utils"
	"testing"
)

func TestGtk3Prop(t *testing.T) {
	Convey("Test gtk3 prop setting", t, func() {
		kfile, err := dutils.NewKeyFileFromFile("testdata/settings.ini")
		So(err, ShouldBeNil)
		defer kfile.Free()

		So(isGtk3PropEqual(gtk3KeyTheme, "Paper",
			kfile), ShouldEqual, true)
		So(isGtk3PropEqual("gtk-menu-images", "1",
			kfile), ShouldEqual, true)
		So(isGtk3PropEqual("gtk-modules", "gail:atk-bridge",
			kfile), ShouldEqual, true)
		So(isGtk3PropEqual("test-list", "1;2;3;",
			kfile), ShouldEqual, true)

		err = setGtk3Prop("test-gtk3", "test", "testdata/tmp-gtk3")
		defer os.Remove("testdata/tmp-gtk3")
		So(err, ShouldBeNil)
	})
}
