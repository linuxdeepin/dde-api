package scanner

import (
	. "github.com/smartystreets/goconvey/convey"
	"sort"
	"testing"
)

func TestListGtkTheme(t *testing.T) {
	Convey("List gtk theme", t, func() {
		list, err := ListGtkTheme("testdata/Themes")
		sort.Strings(list)
		So(list, ShouldResemble, []string{
			"testdata/Themes/Gtk1",
			"testdata/Themes/Gtk2"})
		So(err, ShouldBeNil)
	})
}

func TestListIconTheme(t *testing.T) {
	Convey("List icon theme", t, func() {
		list, err := ListIconTheme("testdata/Icons")
		sort.Strings(list)
		So(list, ShouldResemble, []string{
			"testdata/Icons/Icon1",
			"testdata/Icons/Icon2"})
		So(err, ShouldBeNil)
	})
}

func TestListCursorTheme(t *testing.T) {
	Convey("List cursor theme", t, func() {
		list, err := ListCursorTheme("testdata/Icons")
		sort.Strings(list)
		So(list, ShouldResemble, []string{
			"testdata/Icons/Icon1",
			"testdata/Icons/Icon2"})
		So(err, ShouldBeNil)
	})
}
