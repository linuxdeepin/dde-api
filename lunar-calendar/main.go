/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"os"
	"pkg.linuxdeepin.com/lib"
	"pkg.linuxdeepin.com/lib/dbus"
	Logger "pkg.linuxdeepin.com/lib/logger"
	"time"
)

var (
	logObj = Logger.NewLogger("api/LunarCalendar")
)

func main() {
	lib.UniqueOnSession(LUNAR_DEST)

	logObj.SetRestartCommand("/usr/lib/deepin-api/lunar-calendar")

	m := NewManager()
	if err := dbus.InstallOnSession(m); err != nil {
		logObj.Warning("LunarCalendar Install DBus Session Failed: ", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*1, func() bool {
		return true
	})
	if err := dbus.Wait(); err != nil {
		logObj.Warning("Lost Session DBus")
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
