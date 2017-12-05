/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
)

const (
	dbusDest = "com.deepin.api.XEventMonitor"
)

var (
	logger = log.NewLogger(dbusDest)
)

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSession(dbusDest) {
		logger.Warning("There already has an XEventMonitor daemon running.")
		return
	}

	err := dbus.InstallOnSession(GetManager())
	if err != nil {
		logger.Error("Install DBus Session Failed:", err)
		panic(err)
	}

	dbus.DealWithUnhandledMessage()
	dbus.Emit(GetManager(), "CancelAllArea")

	if err = dbus.Wait(); err != nil {
		logger.Error("lost dbus session:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
