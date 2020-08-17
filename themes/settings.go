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

// Theme settings.
package themes

import (
	"fmt"
	"os"
	"path"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
	"pkg.deepin.io/dde/api/themes/scanner"
	"pkg.deepin.io/gir/glib-2.0"
	"github.com/godbus/dbus"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	wmSchema        = "com.deepin.wrap.gnome.desktop.wm.preferences"
	metacitySchema  = "com.deepin.wrap.gnome.metacity"
	interfaceSchema = "com.deepin.wrap.gnome.desktop.interface"
	xsettingsSchema = "com.deepin.xsettings"

	xsKeyTheme      = "theme-name"
	xsKeyIconTheme  = "icon-theme-name"
	xsKeyCursorName = "gtk-cursor-theme-name"
)

func SetGtkTheme(name string) error {
	if !scanner.IsGtkTheme(getThemePath(name, scanner.ThemeTypeGtk, "themes")) {
		return fmt.Errorf("Invalid theme '%s'", name)
	}

	err := setGtk2Theme(name)
	if err != nil {
		return fmt.Errorf("set GTK2 Theme to '%s' failed", name)
	}

	err = setGtk3Theme(name)
	if err != nil {
		return fmt.Errorf("set GTK2 Theme to '%s' failed", name)
	}

	old := getXSettingsValue(xsKeyTheme)
	if old == name {
		return nil
	}

	if !setXSettingsKey(xsKeyTheme, name) {
		return fmt.Errorf("set theme to '%s' by xsettings failed",
			name)
	}

	if !setWMTheme(name) {
		setXSettingsKey(xsKeyTheme, old)
		return fmt.Errorf("set wm theme to '%s' failed", name)
	}

	if !setQTTheme(name) {
		setXSettingsKey(xsKeyTheme, old)
		setWMTheme(old)
		return fmt.Errorf("set qt theme to '%s' failed", name)
	}
	return nil
}

func SetIconTheme(name string) error {
	if !scanner.IsIconTheme(getThemePath(name, scanner.ThemeTypeIcon, "icons")) {
		return fmt.Errorf("Invalid theme '%s'", name)
	}

	err := setGtk2Icon(name)
	if err != nil {
		return fmt.Errorf("set GTK2 Icon to '%s' failed", name)
	}

	err = setGtk3Icon(name)
	if err != nil {
		return fmt.Errorf("set GTK3 Icon to '%s' failed", name)
	}

	old := getXSettingsValue(xsKeyIconTheme)
	if old == name {
		return nil
	}

	if !setXSettingsKey(xsKeyIconTheme, name) {
		return fmt.Errorf("set theme to '%s' by xsettings failed",
			name)
	}
	return nil
}

func SetCursorTheme(name string) error {
	if !scanner.IsCursorTheme(getThemePath(name, scanner.ThemeTypeCursor, "icons")) {
		return fmt.Errorf("invalid theme '%s'", name)
	}

	err := setGtk2Cursor(name)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	err = setGtk3Cursor(name)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	setDefaultCursor(name)

	old := getXSettingsValue(xsKeyCursorName)
	if old == name {
		return nil
	}

	if !setXSettingsKey(xsKeyCursorName, name) {
		return fmt.Errorf("set theme to '%s' by xsettings failed",
			name)
	}

	setQtCursor(name)
	setGtkCursor(name)
	setWMCursor(name)
	return nil
}

// set cursor theme for deepin-wm
func setWMCursor(name string) {
	ifc, _ := dutils.CheckAndNewGSettings(interfaceSchema)
	if ifc != nil {
		defer ifc.Unref()
		ifc.SetString("cursor-theme", name)
	}

	sessionBus, err := dbus.SessionBus()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to get session bus:", err)
		return
	}

	wmObj := wm.NewWm(sessionBus)
	err = wmObj.CursorTheme().Set(0, name)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to set wm cursorTheme:", err)
	}
}

func GetCursorTheme() string {
	return getXSettingsValue(xsKeyCursorName)
}

func getXSettingsValue(key string) string {
	xs, err := dutils.CheckAndNewGSettings(xsettingsSchema)
	if err != nil {
		return ""
	}
	defer xs.Unref()
	return xs.GetString(key)
}

func setXSettingsKey(key, value string) bool {
	xs, err := dutils.CheckAndNewGSettings(xsettingsSchema)
	if err != nil {
		return false
	}
	defer xs.Unref()
	return xs.SetString(key, value)
}

func setWMTheme(name string) bool {
	meta, _ := dutils.CheckAndNewGSettings(metacitySchema)
	if meta != nil {
		defer meta.Unref()
		meta.SetString("theme", name)
	}

	wm, err := dutils.CheckAndNewGSettings(wmSchema)
	if err != nil {
		return false
	}
	defer wm.Unref()
	return wm.SetString("theme", name)
}

func setQTTheme(name string) bool {
	config := path.Join(glib.GetUserConfigDir(), "Trolltech.conf")
	return setQt4Theme(config)
}

func setQt4Theme(config string) bool {
	value, _ := dutils.ReadKeyFromKeyFile(config, "Qt", "style", "")
	if value == "GTK+" {
		return true
	}
	return dutils.WriteKeyToKeyFile(config, "Qt", "style", "GTK+")
}

func setDefaultCursor(name string) bool {
	file := path.Join(os.Getenv("HOME"), ".icons", "default", "index.theme")
	err := os.MkdirAll(path.Dir(file), 0755)
	if err != nil {
		return false
	}

	value, _ := dutils.ReadKeyFromKeyFile(file, "Icon Theme", "Inherits", "")
	if value == name {
		return true
	}
	return dutils.WriteKeyToKeyFile(file, "Icon Theme", "Inherits", name)
}

func getThemePath(name, ty, key string) string {
	var dirs = []string{
		path.Join(os.Getenv("HOME"), ".local/share/", key),
		path.Join(os.Getenv("HOME"), "."+key),
		path.Join("/usr/local/share", key),
		path.Join("/usr/share", key),
	}

	for _, dir := range dirs {
		tmp := path.Join(dir, name)
		if !dutils.IsFileExist(tmp) {
			continue
		}

		switch ty {
		case scanner.ThemeTypeGtk, scanner.ThemeTypeIcon:
			return dutils.EncodeURI(path.Join(tmp, "index.theme"),
				dutils.SCHEME_FILE)
		case scanner.ThemeTypeCursor:
			return dutils.EncodeURI(path.Join(tmp, "cursor.theme"),
				dutils.SCHEME_FILE)
		}
	}
	return ""
}
