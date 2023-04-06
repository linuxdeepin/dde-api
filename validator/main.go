// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"os"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/linuxdeepin/go-lib"
	"github.com/linuxdeepin/go-lib/dbusutil"
	"github.com/linuxdeepin/go-lib/log"
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
