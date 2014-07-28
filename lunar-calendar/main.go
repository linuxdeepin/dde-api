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
	"pkg.linuxdeepin.com/lib/log"
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
