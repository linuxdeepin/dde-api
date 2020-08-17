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
	"time"

	"github.com/godbus/dbus"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger(DBusName)

// 此执行程序目前没有被使用和编译
func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSession(DBusName) {
		logger.Warning("Validator daemon is already running.")
		return
	}

	bus, err := dbus.SessionBus()
	if err != nil {
		logger.Error("failed to get session bus:", err)
		os.Exit(1)
	}

	service := dbusutil.NewService(bus)
	validator := &Validator{
		service: service,
	}

	err = service.Export(DBusPath, validator)
	if err != nil {
		logger.Error("failed to export dbus service:", err)
		os.Exit(1)
	}

	service.SetAutoQuitHandler(30*time.Second, nil)
	service.Wait()
}
