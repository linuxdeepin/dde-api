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

package themes

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGtk2Infos(t *testing.T) {
	infos := gtk2FileReader("testdata/gtkrc-2.0")
	assert.Equal(t, len(infos), 16)

	info := infos.Get("gtk-theme-name")
	assert.Equal(t, info.value, "\"Paper\"")

	info.value = "\"Deepin\""
	assert.Equal(t, info.value, "\"Deepin\"")

	infos = infos.Add("gtk2-test", "test")
	assert.Equal(t, len(infos), 17)

	infos = gtk2FileReader("testdata/xxx")
	infos = infos.Add("gtk2-test", "test")
	assert.Equal(t, len(infos), 1)
	info = infos.Get("gtk2-test")
	assert.Equal(t, info.value, "test")

	err := gtk2FileWriter(infos, "testdata/tmp-gtk2rc")
	defer os.Remove("testdata/tmp-gtk2rc")
	assert.Nil(t, err)
}
