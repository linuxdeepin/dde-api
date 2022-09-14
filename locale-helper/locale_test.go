// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
