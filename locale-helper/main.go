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
	"os/exec"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

const (
	dbusSender = "com.deepin.api.LocaleHelper"
	dbusPath   = "/com/deepin/api/LocaleHelper"
	dbusIFC    = "com.deepin.api.LocaleHelper"
)

type Helper struct {
	/**
	 * if failed, Success(false, reason), else Success(true, "")
	 **/
	Success func(bool, string)

	running bool
}

var (
	logger = log.NewLogger(dbusSender)
)

func main() {
	if !lib.UniqueOnSystem(dbusSender) {
		logger.Warning("There already has an LocaleHelper running...")
		return
	}

	logger.BeginTracing()
	defer logger.EndTracing()

	var h = &Helper{running: false}
	err := dbus.InstallOnSystem(h)
	if err != nil {
		logger.Error("Install system dbus failed:", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*30, func() bool {
		if h.running {
			return false
		}

		return true
	})

	err = dbus.Wait()
	if err != nil {
		logger.Error("Lost system dbus:", err)
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}

func (h *Helper) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusSender,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (h *Helper) doGenLocale() error {
	return exec.Command("/bin/sh", "-c", "locale-gen").Run()
}

// locales version <= 2.13
func (h *Helper) doGenLocaleWithParam(locale string) error {
	cmd := "locale-gen " + locale
	return exec.Command("/bin/sh", "-c", cmd).Run()
}
