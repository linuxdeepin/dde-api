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

func TestMergeThemeList(t *testing.T) {
	src := []string{"Deepin", "Adwaita", "Zukitwo"}
	target := []string{"Deepin", "Evolve"}
	ret := []string{"Deepin", "Adwaita", "Zukitwo", "Evolve"}

	assert.ElementsMatch(t, mergeThemeList(src, target), ret)
}

func TestSetQt4Theme(t *testing.T) {
	config := "/tmp/Trolltech.conf"
	assert.Equal(t, setQt4Theme(config), true)
	os.Remove(config)
}
