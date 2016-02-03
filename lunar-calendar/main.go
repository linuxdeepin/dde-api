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
	"os"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

var (
	logger = log.NewLogger("api/LunarCalendar")
)

func main() {
	if !lib.UniqueOnSession(LUNAR_DEST) {
		logger.Warning("There already has an lunar-calendar running.")
		return
	}

	logger.SetRestartCommand("/usr/lib/deepin-api/lunar-calendar")

	m := NewManager()
	if err := dbus.InstallOnSession(m); err != nil {
		logger.Warning("LunarCalendar Install DBus Session Failed: ", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*1, nil)
	if err := dbus.Wait(); err != nil {
		logger.Warning("Lost Session DBus")
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
