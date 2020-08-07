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

package main

import (
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
)

const (
	headerMaxLine = 6
)

var llocker sync.Mutex

type LocaleLangInfo struct {
	Enabled bool
	Line    string
	Locale  string
}

type LocaleLangInfos []LocaleLangInfo

type localeFileInfo struct {
	Header []string // file comments, the first 7 lines.
	Infos  LocaleLangInfos
}

func IsLocaleValid(locale string) bool {
	finfo, err := NewLocaleFileInfo(defaultLocaleGenFile)
	if err != nil {
		return false
	}

	return finfo.IsLocaleValid(locale)
}

func NewLocaleFileInfo(file string) (*localeFileInfo, error) {
	datas, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var finfo = &localeFileInfo{}
	finfo.marshal(string(datas))

	return finfo, nil
}

func (finfo *localeFileInfo) EnableLocale(locale string) {
	if finfo.IsLocaleEnabled(locale) {
		return
	}

	finfo.toggleLocale(locale, true)
}

func (finfo *localeFileInfo) DisableLocale(locale string) {
	if !finfo.IsLocaleEnabled(locale) {
		return
	}

	finfo.toggleLocale(locale, false)
}

func (finfo *localeFileInfo) IsLocaleValid(locale string) bool {
	return finfo.Infos.IsLocaleExist(locale)
}

func (finfo *localeFileInfo) IsLocaleEnabled(locale string) bool {
	return finfo.GetEnabledLocales().IsLocaleExist(locale)
}

func (finfo *localeFileInfo) GetEnabledLocales() LocaleLangInfos {
	var infos LocaleLangInfos
	for _, info := range finfo.Infos {
		if !info.Enabled {
			continue
		}

		infos = append(infos, info)
	}

	return infos
}

func (finfo *localeFileInfo) Save(file string) error {
	return writeContentToFile(file, finfo.unmarshal())
}

func (finfo *localeFileInfo) toggleLocale(locale string, enabled bool) {
	var infos LocaleLangInfos
	for _, info := range finfo.Infos {
		if info.Locale == locale {
			info.Enabled = enabled
		}
		infos = append(infos, info)
	}
	finfo.Infos = infos
}

func (finfo *localeFileInfo) marshal(content string) {
	var (
		header []string
		infos  LocaleLangInfos
	)
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if i < headerMaxLine {
			header = append(header, line)
			continue
		}

		// Marshal locale info
		if len(line) == 0 {
			continue
		}

		var info = LocaleLangInfo{
			Enabled: isUncommented(line),
			Line:    getLineContent(line),
		}
		info.Locale = strings.Split(info.Line, " ")[0]
		infos = append(infos, info)
	}

	finfo.Header = header
	finfo.Infos = infos
}

func (finfo *localeFileInfo) unmarshal() string {
	var content string

	for _, v := range finfo.Header {
		content += v + "\n"
	}

	var infoLen = len(finfo.Infos)
	for i, info := range finfo.Infos {
		if !info.Enabled {
			content += "# "
		}
		content += info.Line
		if i != infoLen-1 {
			content += "\n"
		}
	}

	return content
}

func (infos LocaleLangInfos) IsLocaleExist(locale string) bool {
	for _, info := range infos {
		if info.Locale == locale {
			return true
		}
	}

	return false
}

func isUncommented(line string) bool {
	var match = regexp.MustCompile(`^#`)
	return !match.MatchString(line)
}

func getLineContent(line string) string {
	strs := strings.Split(line, "#")
	var v string
	if len(strs) == 1 {
		v = strs[0]
	} else {
		v = strs[1]
	}

	return strings.TrimSpace(v)
}

func writeContentToFile(file, content string) error {
	llocker.Lock()
	defer llocker.Unlock()
	return ioutil.WriteFile(file, []byte(content), 0644)
}
