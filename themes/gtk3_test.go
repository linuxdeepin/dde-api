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
