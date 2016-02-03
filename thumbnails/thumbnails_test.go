/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package thumbnails

import (
	. "github.com/smartystreets/goconvey/convey"
	"pkg.deepin.io/dde/api/thumbnails/loader"
	"testing"
)

func TestCorrectSize(t *testing.T) {
	Convey("Test size correct", t, func() {
		So(correctSize(64), ShouldEqual, loader.SizeFlagSmall)
		So(correctSize(128), ShouldEqual, loader.SizeFlagNormal)
		So(correctSize(176), ShouldEqual, loader.SizeFlagNormal)
		So(correctSize(256), ShouldEqual, loader.SizeFlagLarge)
	})
}
