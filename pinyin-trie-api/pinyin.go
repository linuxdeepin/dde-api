/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"dlib/dbus"
	"strconv"
	"strings"
	"unicode"
)

type Pinyin struct{}

const (
	PINYIN_DEST         = "com.deepin.dde.api.Pinyin"
	HANS_TO_PINYIN_PATH = "/com/deepin/dde/api/HansToPinyin"
	HANS_TO_PINYIN_IFC  = "com.deepin.dde.api.HansToPinyin"
)

func (m *Pinyin) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		PINYIN_DEST,
		HANS_TO_PINYIN_PATH,
		HANS_TO_PINYIN_IFC,
	}
}

func (m *Pinyin) PinyinFromKey(key string) []string {
	return getPinyinFromKey(key)
}

func getPinyinFromKey(key string) []string {
	rets := []string{}
	for _, c := range key {
		println("pinyin char:", string(c))
		if unicode.Is(unicode.Scripts["Han"], c) {
			array := getPinyinByHan(int64(c))
			if len(rets) == 0 {
				rets = array
				continue
			}
			rets = rangeArray(rets, array)
		} else {
			if (c >= 'a' && c <= 'z') ||
				(c >= 'A' && c <= 'Z') {
				array := []string{string(c)}
				if len(rets) == 0 {
					rets = array
				} else {
					rets = rangeArray(rets, array)
				}
			}
		}
	}

	return rets
}

func getPinyinByHan(han int64) []string {
	code := strconv.FormatInt(han, 16)
	value := PinyinDataMap[strings.ToUpper(code)]
	array := strings.Split(value, ";")
	return array
}

func rangeArray(a1, a2 []string) []string {
	rets := []string{}
	for _, v := range a1 {
		for _, r := range a2 {
			rets = append(rets, v+r)
		}
	}

	return rets
}
