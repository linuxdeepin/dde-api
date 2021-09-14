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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LocaleFile(t *testing.T) {
	finfo, err := NewLocaleFileInfo("testdata/locale.gen")
	assert.Nil(t, err)
	assert.Equal(t, len(finfo.Infos), 471)
	assert.Equal(t, len(finfo.GetEnabledLocales()), 5)

	// test locale valid
	assert.Equal(t, finfo.IsLocaleValid("zh_CN.UTF-8"), true)
	assert.Equal(t, finfo.IsLocaleValid("zh_CNN"), false)

	// enable
	finfo.EnableLocale("zh_CN.UTF-8")
	assert.Equal(t, len(finfo.GetEnabledLocales()), 5)
	finfo.EnableLocale("zh_TW.UTF-8")
	assert.Equal(t, len(finfo.GetEnabledLocales()), 6)
	var tmp = "/tmp/test_locale"
	err = finfo.Save(tmp)
	assert.Nil(t, err)

	finfo, err = NewLocaleFileInfo(tmp)
	assert.Nil(t, err)
	assert.Equal(t, len(finfo.Infos), 471)
	assert.Equal(t, len(finfo.GetEnabledLocales()), 6)

	// disable
	finfo.DisableLocale("zh_HK.UTF-8")
	assert.Equal(t, len(finfo.GetEnabledLocales()), 6)
	finfo.DisableLocale("zh_CN.UTF-8")
	assert.Equal(t, len(finfo.GetEnabledLocales()), 5)
	var tmp2 = "/tmp/test_locale2"
	err = finfo.Save(tmp2)
	assert.Nil(t, err)

	finfo, err = NewLocaleFileInfo(tmp2)
	assert.Nil(t, err)
	assert.Equal(t, len(finfo.Infos), 471)
	assert.Equal(t, len(finfo.GetEnabledLocales()), 5)
}
