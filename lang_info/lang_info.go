// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package lang_info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type LangInfo struct {
	Locale      string `json:"Locale"`
	Description string `json:"Description"`
	LangCode    string `json:"LangCode"`
	CountryCode string `json:"CountryCode"`
}
type LangInfos []LangInfo

type langInfoGroup struct {
	Infos LangInfos `json:"LanguageList"`
}

type LangCodeInfo struct {
	LangCode    string
	CountryCode string
	Variant     string
}

const (
	langInfoFile      = "/usr/share/i18n/language_info.json"
	langSupportedFile = "/usr/share/i18n/SUPPORTED"
)

func IsSupportedLocale(locale string) bool {
	infos, err := GetSupportedLangInfos()
	if err != nil {
		return false
	}

	info, _ := infos.Get(locale)
	return (info != nil)
}

func GetSupportedLangInfos() (LangInfos, error) {
	allInfos, err := getLangInfosFromFile(langInfoFile)
	if err != nil {
		return nil, err
	}

	list, err := getSupportedLocaleList(langSupportedFile)
	if err != nil {
		return allInfos, nil
	}

	var infos LangInfos
	for _, info := range allInfos {
		if !isItemInList(info.Locale, list) {
			continue
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func GetLangCodeInfo(locale string) (*LangCodeInfo, error) {
	info, err := getLangInfoByLocale(locale, langInfoFile)
	if err != nil {
		return nil, err
	}
	return info.ToLangCode(), nil
}

func getSupportedLocaleList(config string) ([]string, error) {
	content, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")

	var list []string
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		array := strings.Split(line, " ")
		list = append(list, array[0])
	}
	return list, nil
}

func getLangInfosFromFile(config string) (LangInfos, error) {
	content, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var group langInfoGroup
	err = json.Unmarshal(content, &group)
	if err != nil {
		return nil, err
	}

	return group.Infos, nil
}

func getLangInfoByLocale(locale, config string) (*LangInfo, error) {
	infos, err := getLangInfosFromFile(config)
	if err != nil {
		return nil, err
	}

	info, err := infos.Get(locale)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (infos LangInfos) Get(locale string) (*LangInfo, error) {
	for _, info := range infos {
		if info.Locale == locale {
			return &info, nil
		}
	}
	return nil, fmt.Errorf("Invalid locale: %s", locale)
}

func (info *LangInfo) ToLangCode() *LangCodeInfo {
	var code = new(LangCodeInfo)
	code.LangCode = info.LangCode
	code.CountryCode = info.CountryCode

	array := strings.Split(strings.Split(info.Locale, ".")[0], "@")
	if len(array) > 1 {
		code.Variant = array[1]
	}
	return code
}

func isItemInList(item string, list []string) bool {
	for _, v := range list {
		if item == v {
			return true
		}
	}
	return false
}
