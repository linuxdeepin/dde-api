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
)

const (
	MouseAreaDest = "com.deepin.api.XMouseArea"
)

var (
	logger = log.NewLogger(MouseAreaDest)
)

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSession(MouseAreaDest) {
		logger.Warning("There already has an XMouseArea daemon running.")
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
