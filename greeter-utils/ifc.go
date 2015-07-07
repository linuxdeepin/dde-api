/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"os"
	"path"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
)

func (m *Manager) SetKbdLayout(group, layout string) {
	if !checkConfigDirValid() {
		return
	}

	layout = formatLayout(layout)
	if len(layout) < 1 {
		return
	}

	filename := path.Join(CONFIG_DIR, CONFIG_NAME)
	dutils.WriteKeyToKeyFile(filename, group, KEY_LAYOUT, layout)
}

func (m *Manager) SetKbdLayoutList(group string, list []string) {
	if !checkConfigDirValid() {
		return
	}

	tmp := []string{}
	for _, l := range list {
		l = formatLayout(l)
		if len(l) < 1 {
			continue
		}

		tmp = append(tmp, l)
	}
	list = tmp

	filename := path.Join(CONFIG_DIR, CONFIG_NAME)
	dutils.WriteKeyToKeyFile(filename, group, KEY_LAYOUT_LIST, list)
}

func (m *Manager) SetGreeterTheme(group, theme string) {
	if !checkConfigDirValid() {
		return
	}

	filename := path.Join(CONFIG_DIR, CONFIG_NAME)
	dutils.WriteKeyToKeyFile(filename, group, KEY_GREETER_THEME, theme)
}

func checkConfigDirValid() bool {
	if !dutils.IsFileExist(CONFIG_DIR) {
		err := os.MkdirAll(CONFIG_DIR, 0755)
		if err != nil {
			return false
		}
	}

	return true
}

func formatLayout(layout string) string {
	strs := strings.Split(layout, ";")
	if len(strs) != 2 {
		return ""
	}

	if len(strs[1]) < 1 {
		layout = strs[0]
	} else {
		layout = strs[0] + "\t" + strs[1]
	}

	return layout
}
