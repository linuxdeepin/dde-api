package language_support

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePkgDepends(t *testing.T) {
	pkgDepnds, err := parsePkgDepends("testdata/pkg_depends")
	jsonData, _ := json.Marshal(pkgDepnds)

	assert.Nil(t, err)
	assert.NotNil(t, pkgDepnds)
	assert.NotNil(t, jsonData)
}

func TestLangCodeFromLocale(t *testing.T) {
	locale := langCodeFromLocale("zh_CN")
	assert.Equal(t, locale, "zh-hans")

	locale = langCodeFromLocale("zh_SG")
	assert.Equal(t, locale, "zh-hans")

	locale = langCodeFromLocale("zh_TW")
	assert.Equal(t, locale, "zh-hant")

	locale = langCodeFromLocale("en_US")
	assert.Equal(t, locale, "en")

	locale = langCodeFromLocale("en")
	assert.Equal(t, locale, "en")

	locale = langCodeFromLocale("")
	assert.Equal(t, locale, "")
}

func TestExpendPkgPattern(t *testing.T) {
	pkgs := expendPkgPattern("[p]", "en_US")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]en", "[p]enus", "[p]en-us"})

	pkgs = expendPkgPattern("[p]", "en")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]en"})

	pkgs = expendPkgPattern("[p]", "")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]"})

	pkgs = expendPkgPattern("[p]", "zh_CN")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]zh", "[p]zhcn", "[p]zh-cn", "[p]zh-hans"})

	pkgs = expendPkgPattern("[p]", "zh_SG")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]zh", "[p]zhsg", "[p]zh-sg", "[p]zh-hans"})

	pkgs = expendPkgPattern("[p]", "zh_TW")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]zh", "[p]zhtw", "[p]zh-tw", "[p]zh-hant"})

	pkgs = expendPkgPattern("[p]", "zh_HK")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]zh", "[p]zhhk", "[p]zh-hk", "[p]zh-hant"})

	pkgs = expendPkgPattern("[p]", "wa_BE@euro")
	assert.ElementsMatch(t, pkgs, []string{"[p]", "[p]wa", "[p]wabe", "[p]wa-be", "[p]wa-euro", "[p]wa-be-euro"})
}
