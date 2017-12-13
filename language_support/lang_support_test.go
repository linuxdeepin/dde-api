package language_support

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsePkgDepends(t *testing.T) {
	Convey("parsePkgDepends", t, func() {
		pkgDepnds, err := parsePkgDepends("testdata/pkg_depends")
		jsonData, _ := json.Marshal(pkgDepnds)
		t.Logf("%s", jsonData)

		So(err, ShouldBeNil)
		So(pkgDepnds, ShouldNotBeNil)
	})
}

func TestLangCodeFromLocale(t *testing.T) {
	Convey("langCodeFromLocale", t, func() {
		locale := langCodeFromLocale("zh_CN")
		So(locale, ShouldEqual, "zh-hans")

		locale = langCodeFromLocale("zh_SG")
		So(locale, ShouldEqual, "zh-hans")

		locale = langCodeFromLocale("zh_TW")
		So(locale, ShouldEqual, "zh-hant")

		locale = langCodeFromLocale("en_US")
		So(locale, ShouldEqual, "en")

		locale = langCodeFromLocale("en")
		So(locale, ShouldEqual, "en")

		locale = langCodeFromLocale("")
		So(locale, ShouldEqual, "")
	})
}

func TestExpendPkgPattern(t *testing.T) {
	Convey("expendPkgPattern", t, func() {
		pkgs := expendPkgPattern("[p]", "en_US")
		So(pkgs, ShouldResemble, []string{"[p]", "[p]en", "[p]enus", "[p]en-us"})

		pkgs = expendPkgPattern("[p]", "en")
		So(pkgs, ShouldResemble, []string{"[p]", "[p]en"})

		pkgs = expendPkgPattern("[p]", "")
		So(pkgs, ShouldResemble, []string{"[p]", "[p]"})

		pkgs = expendPkgPattern("[p]", "zh_CN")
		So(pkgs, ShouldResemble, []string{"[p]", "[p]zh", "[p]zhcn", "[p]zh-cn", "[p]zh-hans"})

		pkgs = expendPkgPattern("[p]", "zh_SG")
		So(pkgs, ShouldResemble, []string{"[p]", "[p]zh", "[p]zhsg", "[p]zh-sg", "[p]zh-hans"})

		pkgs = expendPkgPattern("[p]", "zh_TW")
		So(pkgs, ShouldResemble, []string{"[p]", "[p]zh", "[p]zhtw", "[p]zh-tw", "[p]zh-hant"})

		pkgs = expendPkgPattern("[p]", "zh_HK")
		So(pkgs, ShouldResemble, []string{"[p]", "[p]zh", "[p]zhhk", "[p]zh-hk", "[p]zh-hant"})

		pkgs = expendPkgPattern("[p]", "wa_BE@euro")
		So(pkgs, ShouldResemble, []string{"[p]", "[p]wa", "[p]wabe", "[p]wa-be", "[p]wa-euro", "[p]wa-be-euro"})
	})
}
