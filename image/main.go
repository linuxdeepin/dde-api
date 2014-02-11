package main

import (
	"dlib/dbus"
)

type DImage struct{}

func (dimg *DImage) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.api.Image",
		"/com/deepin/api/Image",
		"com.deepin.api.Image",
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			// TODO logFatal("deepin image api failed: %v", err)
		}
	}()

	dimg := &DImage{}
	err := dbus.InstallOnSession(dimg)
	if err != nil {
		panic(err)
	}
	dbus.DealWithUnhandledMessage()

	select {}
}
