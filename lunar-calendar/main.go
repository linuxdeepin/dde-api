/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
	if !lib.UniqueOnSession(DBusDest) {
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
	dbus.SetAutoDestroyHandler(time.Second*100, nil)
	if err := dbus.Wait(); err != nil {
		logger.Warning("Lost Session DBus")
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
