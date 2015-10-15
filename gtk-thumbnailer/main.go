package main

import (
	"os"
	"time"

	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

const (
	dbusDest = "com.deepin.api.GtkThumbnailer"
	dbusPath = "/com/deepin/api/GtkThumbnailer"
	dbusIFC  = "com.deepin.api.GtkThumbnailer"
)

var logger = log.NewLogger("api/GtkThumbnailer")

type Manager struct {
	running bool
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) Thumbnail(name, bg, dest string, width, height int32, force bool) error {
	m.running = true
	defer func() {
		m.running = false
	}()
	return doGenThumbnail(name, bg, dest, int(width), int(height), force)
}

func main() {
	if !lib.UniqueOnSession(dbusDest) {
		logger.Warning("There already has a gtk thumbnailer running...")
		return
	}

	err := initGtkEnv()
	if err != nil {
		logger.Error(err)
	}

	var m = new(Manager)
	m.running = false
	err = dbus.InstallOnSession(m)
	if err != nil {
		logger.Error("Install dbus session failed:", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*5, func() bool {
		if m.running {
			return false
		}
		return true
	})

	err = dbus.Wait()
	if err != nil {
		logger.Error("Lost dbus connect:", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
