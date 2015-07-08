/**
 * Copyright (c) 2015 Deepin, Inc.
 *               2015 Xu Shaohua
 *
 * Author:       Xu Shaohua<xushaohua@linuxdeepin.com>
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
	"time"

	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger(DBusName)

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSession(DBusName) {
		logger.Warning("Validator daemon is already running.")
		return
	}

	validator := &Validator{}
	err := dbus.InstallOnSession(validator)
	if err != nil {
		logger.Errorf("Failed to register dbus interface: %v", err)
		os.Exit(1)
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(30*time.Second, nil)
	if err := dbus.Wait(); err != nil {
		logger.Errorf("lost dbus session: %v", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
