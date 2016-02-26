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
	"time"

	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var (
	logger = log.NewLogger(dbusDest)
)

func main() {
	if !lib.UniqueOnSystem(dbusDest) {
		logger.Warning("There already has an greeter-helper running.")
		return
	}

	logger.BeginTracing()
	defer logger.EndTracing()

	var m = new(Manager)
	if err := dbus.InstallOnSystem(m); err != nil {
		logger.Fatal("Install DBus Error:", err)
		return
	}
	dbus.DealWithUnhandledMessage()
	dbus.SetAutoDestroyHandler(time.Second*5, nil)

	err := dbus.Wait()
	if err != nil {
		logger.Warning("Lost DBus...")
		os.Exit(-1)
	}
	os.Exit(0)
}
