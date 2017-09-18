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
