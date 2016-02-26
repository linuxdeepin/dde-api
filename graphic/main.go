/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"os"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

var logger = log.NewLogger(dbusGraphicDest)

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSession(dbusGraphicDest) {
		logger.Warning("There already has an Graphic daemon running.")
		return
	}

	graphic := &Graphic{}
	err := dbus.InstallOnSession(graphic)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
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
