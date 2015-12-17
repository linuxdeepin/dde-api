// Theme settings.
package themes

import (
	"fmt"
	"os"
	"path"

	"gir/glib-2.0"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	wmSchema        = "com.deepin.wrap.gnome.desktop.wm.preferences"
	metacitySchema  = "com.deepin.wrap.gnome.metacity"
	xsettingsSchema = "com.deepin.xsettings"

	xsKeyTheme      = "theme-name"
	xsKeyIconTheme  = "icon-theme-name"
	xsKeyCursorName = "gtk-cursor-theme-name"
)

func SetGtkTheme(name string) error {
	if !IsThemeInList(name, ListGtkTheme()) {
		return fmt.Errorf("Invalid theme '%s'", name)
	}

	old := getXSettingsValue(xsKeyTheme)
	if old == name {
		return nil
	}

	setGtk2Theme(name)
	setGtk3Theme(name)

	if !setXSettingsKey(xsKeyTheme, name) {
		return fmt.Errorf("Set theme to '%s' by xsettings failed",
			name)
	}

	if !setWMTheme(name) {
		setXSettingsKey(xsKeyTheme, old)
		return fmt.Errorf("Set wm theme to '%s' failed", name)
	}

	if !setQTTheme(name) {
		setXSettingsKey(xsKeyTheme, old)
		setWMTheme(old)
		return fmt.Errorf("Set qt theme to '%s' failed", name)
	}
	return nil
}

func SetIconTheme(name string) error {
	if !IsThemeInList(name, ListIconTheme()) {
		return fmt.Errorf("Invalid theme '%s'", name)
	}

	old := getXSettingsValue(xsKeyIconTheme)
	if old == name {
		return nil
	}

	setGtk2Icon(name)
	setGtk3Icon(name)

	if !setXSettingsKey(xsKeyIconTheme, name) {
		return fmt.Errorf("Set theme to '%s' by xsettings failed",
			name)
	}
	return nil
}

func SetCursorTheme(name string) error {
	if !IsThemeInList(name, ListCursorTheme()) {
		return fmt.Errorf("Invalid theme '%s'", name)
	}

	old := getXSettingsValue(xsKeyCursorName)
	if old == name {
		return nil
	}

	setGtk2Cursor(name)
	setGtk3Cursor(name)

	if !setXSettingsKey(xsKeyCursorName, name) {
		return fmt.Errorf("Set theme to '%s' by xsettings failed",
			name)
	}

	setDefaultCursor(name)

	return nil
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
