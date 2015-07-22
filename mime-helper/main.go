package main

import (
	"os"
	"time"

	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
)

var (
	logger = log.NewLogger("dde-api/mime")
)

func main() {
	if !lib.UniqueOnSession(dbusDest) {
		logger.Warning("There already has a dde-api/mime running...")
		return
	}

	gettext.InitI18n()
	gettext.Textdomain("dde-daemon")

	m := NewManager()
	err := dbus.InstallOnSession(m)
	if err != nil {
		logger.Error("Install mime dbus failed:", err)
		return
	}

	m.media, err = NewMedia()
	if err != nil {
		logger.Error("New Media failed:", err)
	} else {
		err := dbus.InstallOnSession(m.media)
		if err != nil {
			logger.Error("Install Media dbus failed:", err)
		}
	}

	dbus.DealWithUnhandledMessage()
	dbus.SetAutoDestroyHandler(time.Second*5, nil)

	err = dbus.Wait()
	if err != nil {
		logger.Error("dde-api/mime lost dbus:", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
