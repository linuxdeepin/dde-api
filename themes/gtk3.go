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
	"io/ioutil"
	"os"
	"path"
	"sync"

	"pkg.deepin.io/gir/glib-2.0"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	gtk3GroupSettings = "Settings"
	gtk3KeyTheme      = "gtk-theme-name"
	gtk3KeyIcon       = "gtk-icon-theme-name"
	gtk3KeyCursor     = "gtk-cursor-theme-name"
)

var (
	gtk3Locker   sync.Mutex
	gtk3ConfFile = path.Join(os.Getenv("HOME"),
		".config", "gtk-3.0", "settings.ini")
)

func setGtk3Theme(name string) error {
	return setGtk3Prop(gtk3KeyTheme, name, gtk3ConfFile)
}

func setGtk3Icon(name string) error {
	return setGtk3Prop(gtk3KeyIcon, name, gtk3ConfFile)
}

func setGtk3Cursor(name string) error {
	return setGtk3Prop(gtk3KeyCursor, name, gtk3ConfFile)
}

func setGtk3Prop(key, value, file string) error {
	gtk3Locker.Lock()
	defer gtk3Locker.Unlock()

	if !dutils.IsFileExist(file) {
		err := os.MkdirAll(path.Dir(file), 0755)
		if err != nil {
			return err
		}

		err = dutils.CreateFile(file)
		if err != nil {
			return err
		}
	}

	kfile, err := dutils.NewKeyFileFromFile(file)
	if kfile == nil {
		return err
	}
	defer kfile.Free()

	if isGtk3PropEqual(key, value, kfile) {
		return nil
	}

	return doSetGtk3Prop(key, value, file, kfile)
}

func isGtk3PropEqual(key, value string, kfile *glib.KeyFile) bool {
	old, _ := kfile.GetString(gtk3GroupSettings, key)
	return old == value
}

func doSetGtk3Prop(key, value, file string, kfile *glib.KeyFile) error {
	kfile.SetString(gtk3GroupSettings, key, value)
	_, content, err := kfile.ToData()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(content), 0644)
}
