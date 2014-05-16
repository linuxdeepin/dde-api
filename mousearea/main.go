/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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
	"dlib"
	"dlib/dbus"
	"dlib/logger"
	"os"
)

var (
	Logger = logger.NewLogger("dde-api/mousearea")
)

const (
	MouseAreaDest = "com.deepin.api.XMouseArea"
)

func main() {
	defer Logger.EndTracing()

	if !dlib.UniqueOnSession(MouseAreaDest) {
		Logger.Warning("There already has an XMouseArea daemon running.")
		return
	}

	// configure logger
	Logger.SetRestartCommand("/usr/lib/deepin-api/mousearea", "--debug")
	if stringInSlice("-d", os.Args) || stringInSlice("--debug", os.Args) {
		Logger.SetLogLevel(logger.LEVEL_DEBUG)
	}

	var err error
	if err != nil {
		Logger.Warning("New XGB Connection Failed")
		return
	}

	err = dbus.InstallOnSession(GetManager())
	if err != nil {
		Logger.Error("Install DBus Session Failed:", err)
		panic(err)
	}

	dbus.DealWithUnhandledMessage()
	//select {}
	if err = dbus.Wait(); err != nil {
		Logger.Error("lost dbus session:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
