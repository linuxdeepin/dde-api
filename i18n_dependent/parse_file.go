package i18n_dependent

import (
	"encoding/json"
	"io/ioutil"
	"pkg.deepin.io/dde/api/lang_info"
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
	dependents = append(dependents, categories.GetDependentInfos(
		"tr", locale)...)
	dependents = append(dependents, categories.GetDependentInfos(
		"wa", locale)...)
	// font
	dependents = append(dependents, categories.GetDependentInfos(
		"fn", locale)...)
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
		if len(codeInfo.LangCode) == 0 {
			break
		}
		pkgList = append(pkgList, info.PkgPull+codeInfo.LangCode)
	case formatTypeLCCC:
		if len(codeInfo.LangCode) == 0 {
			break
		}
		pkgList = append(pkgList, info.PkgPull+codeInfo.LangCode)

		if len(codeInfo.CountryCode) == 0 {
			break
		}
		pkgList = append(pkgList, info.PkgPull+codeInfo.LangCode+
			"-"+strings.ToLower(codeInfo.CountryCode))
	case formatTypeLCVA:
		if len(codeInfo.LangCode) == 0 {
			break
		}
		pkgList = append(pkgList, info.PkgPull+codeInfo.LangCode)

		if len(codeInfo.CountryCode) != 0 {
			pkgList = append(pkgList, info.PkgPull+
				codeInfo.LangCode+
				strings.ToLower(codeInfo.CountryCode))
		}

		if len(codeInfo.Variant) != 0 {
			pkgList = append(pkgList, info.PkgPull+
				codeInfo.LangCode+"-"+codeInfo.Variant)
		}
	}
	return pkgList
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
