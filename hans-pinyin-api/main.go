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

type Manager struct{}

const (
	HANS_TO_PINYIN_DEST = "com.deepin.dde.api.HansToPinyin"
	HANS_TO_PINYIN_PATH = "/com/deepin/dde/api/HansToPinyin"
	HANS_TO_PINYIN_IFC  = "com.deepin.dde.api.HansToPinyin"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		HANS_TO_PINYIN_DEST,
		HANS_TO_PINYIN_PATH,
		HANS_TO_PINYIN_IFC,
	}
}

func (m *Manager) GetPinyinByHans(hans string) []string {
	print("key: ", hans, "\n")
	print("len: ", len(hans), "\n")
	if len(hans) != 3 {
		return nil
	}

	for _, c := range hans {
		print("c: ", c, "\n")
		if unicode.Is(unicode.Scripts["Han"], c) {
			code := strconv.FormatInt(int64(c), 16)
			print("str: ", code, "\n")
			value := PinyinDataMap[strings.ToUpper(code)]
			print("value: ", value, "\n")
                        array := strings.Split(value, ";")
			return array
		}
	}

	return nil
}

func main() {
	m := &Manager{}
	err := dbus.InstallOnSession(m)
	if err != nil {
		panic(err)
	}
	dbus.DealWithUnhandledMessage()

	select {}
}
