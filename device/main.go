/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

var logger = log.NewLogger(deviceDest)

func main() {
	if !lib.UniqueOnSystem(deviceDest) {
		logger.Warning("dbus interface already exists", deviceDest)
		return
	}

	d := &Device{}
	err := dbus.InstallOnSystem(d)
	if err != nil {
		logger.Error("register dbus interface failed", err)
		return
	}

	dbus.SetAutoDestroyHandler(5*time.Second, func() bool {
		return true
	})

	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		logger.Error("lost dbus session", err)
	}
}
