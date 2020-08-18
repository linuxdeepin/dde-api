package language_support

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsePkgDepends(t *testing.T) {
	Convey("parsePkgDepends", t, func(c C) {
		pkgDepnds, err := parsePkgDepends("testdata/pkg_depends")
		jsonData, _ := json.Marshal(pkgDepnds)
		t.Logf("%s", jsonData)

		c.So(err, ShouldBeNil)
		c.So(pkgDepnds, ShouldNotBeNil)
	})
}

func TestLangCodeFromLocale(t *testing.T) {
	Convey("langCodeFromLocale", t, func(c C) {
		locale := langCodeFromLocale("zh_CN")
		c.So(locale, ShouldEqual, "zh-hans")

		locale = langCodeFromLocale("zh_SG")
		c.So(locale, ShouldEqual, "zh-hans")

		locale = langCodeFromLocale("zh_TW")
		c.So(locale, ShouldEqual, "zh-hant")

		locale = langCodeFromLocale("en_US")
		c.So(locale, ShouldEqual, "en")

		locale = langCodeFromLocale("en")
		c.So(locale, ShouldEqual, "en")

		locale = langCodeFromLocale("")
		c.So(locale, ShouldEqual, "")
	})
}

func TestExpendPkgPattern(t *testing.T) {
	Convey("expendPkgPattern", t, func(c C) {
		pkgs := expendPkgPattern("[p]", "en_US")
		c.So(pkgs, ShouldResemble, []string{"[p]", "[p]en", "[p]enus", "[p]en-us"})

		pkgs = expendPkgPattern("[p]", "en")
		c.So(pkgs, ShouldResemble, []string{"[p]", "[p]en"})

		pkgs = expendPkgPattern("[p]", "")
		c.So(pkgs, ShouldResemble, []string{"[p]", "[p]"})

		pkgs = expendPkgPattern("[p]", "zh_CN")
		c.So(pkgs, ShouldResemble, []string{"[p]", "[p]zh", "[p]zhcn", "[p]zh-cn", "[p]zh-hans"})

		pkgs = expendPkgPattern("[p]", "zh_SG")
		c.So(pkgs, ShouldResemble, []string{"[p]", "[p]zh", "[p]zhsg", "[p]zh-sg", "[p]zh-hans"})

		pkgs = expendPkgPattern("[p]", "zh_TW")
		c.So(pkgs, ShouldResemble, []string{"[p]", "[p]zh", "[p]zhtw", "[p]zh-tw", "[p]zh-hant"})

		pkgs = expendPkgPattern("[p]", "zh_HK")
		c.So(pkgs, ShouldResemble, []string{"[p]", "[p]zh", "[p]zhhk", "[p]zh-hk", "[p]zh-hant"})

		pkgs = expendPkgPattern("[p]", "wa_BE@euro")
		c.So(pkgs, ShouldResemble, []string{"[p]", "[p]wa", "[p]wabe", "[p]wa-be", "[p]wa-euro", "[p]wa-be-euro"})
	})
}
