package main

import (
	"dlib/dbus"
)

type Image struct{}

func (image *Image) GetDBusInfo() dbus.DBusInfo {
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

	image := &Image{}
	err := dbus.InstallOnSession(image)
	if err != nil {
		panic(err)
	}
	dbus.DealWithUnhandledMessage()

	select {}
}
