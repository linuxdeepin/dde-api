// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
