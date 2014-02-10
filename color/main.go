package main

import (
	"dlib/dbus"
)

type Color struct{}

func (color *Color) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.api.Color",
		"/com/deepin/api/Color",
		"com.deepin.api.Color",
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			// TODO logFatal("deepin color api failed: %v", err)
		}
	}()

	color := &Color{}
	err := dbus.InstallOnSession(color)
	if err != nil {
		panic(err)
	}
	dbus.DealWithUnhandledMessage()

	select {}
}
