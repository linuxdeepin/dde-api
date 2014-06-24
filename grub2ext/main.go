/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
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
	"flag"
	"os"
	"pkg.linuxdeepin.com/lib"
	"pkg.linuxdeepin.com/lib/dbus"
	liblogger "pkg.linuxdeepin.com/lib/logger"
)

var logger = liblogger.NewLogger(grubDest)
var argDebug bool

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSystem(grubDest) {
		logger.Warning("There already has an Grub2 daemon running.")
		return
	}

	flag.BoolVar(&argDebug, "d", false, "debug mode")
	flag.BoolVar(&argDebug, "debug", false, "debug mode")
	flag.Parse()

	// configure logger
	logger.SetRestartCommand("/usr/lib/deepin-api/grub2ext", "--debug") // TODO: is still need?
	if argDebug {
		logger.Info("debug mode")
		logger.SetLogLevel(liblogger.LEVEL_DEBUG)
	}

	grub := NewGrub2Ext()
	err := dbus.InstallOnSystem(grub)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		os.Exit(1)
	}

	dbus.DealWithUnhandledMessage()

	// TODO
	// dbus.SetAutoDestroyHandler(300*time.Second, nil)

	if err := dbus.Wait(); err != nil {
		logger.Errorf("lost dbus session: %v", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
