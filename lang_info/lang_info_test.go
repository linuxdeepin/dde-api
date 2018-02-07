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

package lang_info

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSupportedLocale(t *testing.T) {
	Convey("Test locale whether supported", t, func() {
		list, err := getSupportedLocaleList("testdata/SUPPORTED")
		So(err, ShouldEqual, nil)
		So(len(list), ShouldEqual, 475)

		So(isItemInList("zh_CN.UTF-8", list), ShouldEqual, true)
		So(isItemInList("zh_CNN.UTF-8", list), ShouldEqual, false)
	})
}

func TestLangInfo(t *testing.T) {
	Convey("Test language info", t, func() {
		infos, err := getLangInfosFromFile("testdata/language_info.json")
		So(err, ShouldEqual, nil)
		So(len(infos), ShouldEqual, 143)
		_, err = infos.Get("zh_CNN")
		So(err, ShouldNotEqual, nil)

		info, err := getLangInfoByLocale("zh_CN.UTF-8",
			"testdata/language_info.json")
		So(err, ShouldEqual, nil)
		So(info.LangCode, ShouldEqual, "zh-hans")
		So(info.ToLangCode().CountryCode, ShouldEqual, "CN")
	})
}
