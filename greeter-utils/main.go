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
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

var (
	logger = log.NewLogger(DEST)
)

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSystem(DEST) {
		logger.Warning("There already has an GreeterUtils running.")
		return
	}

	m := &Manager{}
	if err := dbus.InstallOnSystem(m); err != nil {
		logger.Fatal("Install DBus Error:", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*5, nil)
	if err := dbus.Wait(); err != nil {
		logger.Warning("Lost DBus...")
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
