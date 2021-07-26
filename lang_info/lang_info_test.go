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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupportedLocale(t *testing.T) {
	testName := "Test locale whether supported"
	tests := []struct {
		inputArgs string
		want      []string
	}{
		{
			inputArgs: "testdata/SUPPORTED",
			want:      nil,
		},
	}
	t.Run(testName, func(t *testing.T) {
		for _, tt := range tests {
			list, err := getSupportedLocaleList(tt.inputArgs)
			assert.NoError(t, err)
			assert.Equal(t, 475, len(list))
			assert.True(t, isItemInList("zh_CN.UTF-8", list))
			assert.True(t, !isItemInList("zh_CNN.UTF-8", list))
		}
	})
}

func TestLangInfo(t *testing.T) {
	testName1 := "Test getLangInfosFromFile"
	tests1 := []struct {
		inputArgs string
	}{
		{
			inputArgs: "testdata/language_info.json",
		},
	}
	t.Run(testName1, func(t *testing.T) {
		for _, tt := range tests1 {
			infos, err := getLangInfosFromFile(tt.inputArgs)
			assert.NoError(t, err)
			assert.Equal(t, 143, len(infos))
			_, err = infos.Get("zh_CNN")
			assert.Error(t, err)
		}
	})

	testName2 := "Test getLangInfoByLocale"
	tests2 := []struct {
		inputArgs1 string
		inputArgs2 string
	}{
		{
			inputArgs1: "zh_CN.UTF-8",
			inputArgs2: "testdata/language_info.json",
		},
	}
	t.Run(testName2, func(t *testing.T) {
		for _, tt := range tests2 {
			info, err := getLangInfoByLocale(tt.inputArgs1, tt.inputArgs2)
			assert.NoError(t, err)
			assert.Equal(t, "zh-hans", info.LangCode)
			assert.Equal(t, "CN", info.ToLangCode().CountryCode)
		}
	})
}
