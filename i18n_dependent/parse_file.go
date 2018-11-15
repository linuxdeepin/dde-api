/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package i18n_dependent

import (
	"encoding/json"
	"io/ioutil"
	"pkg.deepin.io/dde/api/lang_info"
	"regexp"
	"strings"
)

const (
	formatTypeNone int32 = 0
	//Format: %LCODE%
	formatTypeLC int32 = 1
	//Format: %LCODE% or %LCODE%-%CCODE%
	formatTypeLCCC int32 = 2
	//Format: %LCODE% or %LCODE%%CCODE% or %LCODE%-%VARIANT%
	formatTypeLCVA int32 = 3
)

type jsonDependentInfo struct {
	LangCode   string `json:"LangCode"`
	FormatType int32  `json:"FormatType"`
	Dependent  string `json:"DependentPkg"`
	PkgPull    string `json:"PkgPull"`
}
type jsonDependentInfos []jsonDependentInfo

type jsonDependentCategory struct {
	Category string             `json:"Category"`
	Infos    jsonDependentInfos `json:"PkgInfos"`
}
type jsonDependentCategories []jsonDependentCategory

type jsonDependentGroup struct {
	Categories jsonDependentCategories `json:"PkgDepends"`
}

func (categories jsonDependentCategories) GetAllDependentInfos(locale string) DependentInfos {
	var dependents DependentInfos
	// tr: translations
	dependents = append(dependents, categories.GetDependentInfos(
		"tr", locale)...)
	// wa: writing assistance
	dependents = append(dependents, categories.GetDependentInfos(
		"wa", locale)...)
	// fn: font
	dependents = append(dependents, categories.GetDependentInfos(
		"fn", locale)...)
	// im: input method. Ignore
	return dependents
}

func (categories jsonDependentCategories) GetDependentInfos(key, locale string) DependentInfos {
	infos := categories.GetInfos(key)
	if infos == nil {
		return nil
	}

	return infos.GetDependentInfos(locale)
}

func (categories jsonDependentCategories) GetInfos(key string) jsonDependentInfos {
	for _, category := range categories {
		if category.Category == key {
			return category.Infos
		}
	}
	return nil
}

func (infos jsonDependentInfos) GetDependentInfos(locale string) DependentInfos {
	var dependents DependentInfos
	for _, info := range infos {
		if len(info.LangCode) == 0 {
			dependents = append(dependents, DependentInfo{
				Dependent: info.Dependent,
				Packages:  info.GetPackages(locale),
			})
			continue
		}

		codeInfo, err := lang_info.GetLangCodeInfo(locale)
		if err != nil || codeInfo.LangCode != info.LangCode {
			continue
		}
		dependents = append(dependents, DependentInfo{
			Dependent: info.Dependent,
			Packages:  []string{info.PkgPull},
		})
	}
	return dependents
}

func (info *jsonDependentInfo) GetPackages(locale string) []string {
	codeInfo, err := lang_info.GetLangCodeInfo(locale)
	if err != nil {
		return nil
	}

	var pkgList []string
	switch info.FormatType {
	case formatTypeNone:
		pkgList = append(pkgList, info.PkgPull)
	case formatTypeLC:
		pkgList = append(pkgList, info.getPackagesByLangInfo(locale, codeInfo.LangCode,
			"", "")...)
	case formatTypeLCCC:
		pkgList = append(pkgList, info.getPackagesByLangInfo(locale, codeInfo.LangCode,
			codeInfo.CountryCode, "")...)
	case formatTypeLCVA:
		pkgList = append(pkgList, info.getPackagesByLangInfo(locale, codeInfo.LangCode,
			codeInfo.CountryCode, codeInfo.Variant)...)
	}
	return pkgList
}

var regUnderLine = regexp.MustCompile(`_`)

func (info *jsonDependentInfo) getPackagesByLangInfo(locale, langCode, countryCode, variant string) []string {
	var ret []string
	tmp := strings.Split(locale, ".")[0]
	tmp = strings.ToLower(tmp)
	// Fix for firefox-l10n, calligra-l10n
	ret = append(ret, info.PkgPull+regUnderLine.ReplaceAllString(tmp, "-"))
	ret = append(ret, info.PkgPull+regUnderLine.ReplaceAllString(tmp, ""))

	if len(langCode) == 0 {
		return ret
	}
	ret = append(ret, info.PkgPull+langCode)

	if len(countryCode) == 0 {
		return ret
	}
	countryCode = strings.ToLower(countryCode)
	ret = append(ret, info.PkgPull+langCode+"-"+countryCode)
	ret = append(ret, info.PkgPull+langCode+countryCode)

	if len(variant) == 0 {
		return ret
	}
	ret = append(ret, langCode+"-"+variant)
	return ret
}

func getDependentCategories(config string) (jsonDependentCategories, error) {
	content, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var group jsonDependentGroup
	err = json.Unmarshal(content, &group)
	if err != nil {
		return nil, err
	}
	return group.Categories, nil
}
