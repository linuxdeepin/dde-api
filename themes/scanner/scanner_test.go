// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package scanner

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListGtkTheme(t *testing.T) {
	list, err := ListGtkTheme("testdata/Themes")
	sort.Strings(list)
	assert.ElementsMatch(t, list, []string{
		"testdata/Themes/Gtk1",
		"testdata/Themes/Gtk2"})
	assert.Nil(t, err)
}

func TestListIconTheme(t *testing.T) {
	list, err := ListIconTheme("testdata/Icons")
	sort.Strings(list)
	assert.ElementsMatch(t, list, []string{
		"testdata/Icons/Icon1",
		"testdata/Icons/Icon2"})
	assert.Nil(t, err)
}

func TestListCursorTheme(t *testing.T) {
	list, err := ListCursorTheme("testdata/Icons")
	sort.Strings(list)
	assert.ElementsMatch(t, list, []string{
		"testdata/Icons/Icon1",
		"testdata/Icons/Icon2"})
	assert.Nil(t, err)
}

func TestThemeHidden(t *testing.T) {
	assert.Equal(t, isHidden("testdata/gtk_paper.theme", ThemeTypeGtk), false)
	assert.Equal(t, isHidden("testdata/gtk_paper_hidden.theme", ThemeTypeGtk), true)

	assert.Equal(t, isHidden("testdata/icon_deepin.theme", ThemeTypeIcon), false)
	assert.Equal(t, isHidden("testdata/icon_deepin_hidden.theme", ThemeTypeIcon), true)
}
