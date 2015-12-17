package themes

import (
	"io/ioutil"
	"os"
	"path"
	"sync"

	"gir/glib-2.0"
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
	if old == value {
		return true
	}
	return false
}

func doSetGtk3Prop(key, value, file string, kfile *glib.KeyFile) error {
	kfile.SetString(gtk3GroupSettings, key, value)
	_, content, err := kfile.ToData()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(content), 0644)
}
