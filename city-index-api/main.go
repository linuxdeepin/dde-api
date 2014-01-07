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
	"strings"
	"time"
)

type CityPinyin struct{}

const (
	_CITY_PINYIN_DEST = "com.deepin.dde.api.CityPinyin"
	_CITY_PINYIN_PATH = "/com/deepin/dde/api/CityPinyin"
	_CITY_PINYIN_IFC  = "com.deepin.dde.api.CityPinyin"
)

func (cp *CityPinyin) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_CITY_PINYIN_DEST,
		_CITY_PINYIN_PATH,
		_CITY_PINYIN_IFC,
	}
}

func (cp *CityPinyin) GetValues(key string) []string {
	if len(key) < 2 {
		return nil
	}

	values := []string{}
	tmp := strings.ToLower(key)
	for k, v := range CityPinyinMap {
		if strings.Contains(k, tmp) {
			values = append(values, v...)
		}
	}

	return values
}

func (cp *CityPinyin) GetValuesWithPinyin(key string) map[string][]string {
	if len(key) < 2 {
		return nil
	}

	values := make(map[string][]string)
	tmp := strings.ToLower(key)
	for k, v := range CityPinyinMap {
		if strings.Contains(k, tmp) {
			values[k] = append(values[k], v...)
		}
	}

	return values
}

func main() {
	cp := &CityPinyin{}
	dbus.InstallOnSession(cp)
	dbus.DealWithUnhandledMessage()

	timer := time.NewTimer(time.Second * 10)
	<-timer.C
}
