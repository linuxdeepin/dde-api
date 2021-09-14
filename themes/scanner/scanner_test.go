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
