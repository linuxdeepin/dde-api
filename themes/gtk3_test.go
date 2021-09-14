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
	dutils "pkg.deepin.io/lib/utils"
)

func TestGtk3Prop(t *testing.T) {
	kfile, err := dutils.NewKeyFileFromFile("testdata/settings.ini")
	assert.Nil(t, err)
	defer kfile.Free()

	assert.Equal(t, isGtk3PropEqual(gtk3KeyTheme, "Paper",
		kfile), true)
	assert.Equal(t, isGtk3PropEqual("gtk-menu-images", "1",
		kfile), true)
	assert.Equal(t, isGtk3PropEqual("gtk-modules", "gail:atk-bridge",
		kfile), true)
	assert.Equal(t, isGtk3PropEqual("test-list", "1;2;3;",
		kfile), true)

	err = setGtk3Prop("test-gtk3", "test", "testdata/tmp-gtk3")
	defer os.Remove("testdata/tmp-gtk3")
	assert.Nil(t, err)
}
