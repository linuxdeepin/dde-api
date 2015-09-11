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
