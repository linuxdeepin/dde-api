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
	"pkg.deepin.io/lib/dbus"
)

type Validator struct{}

// GetDBusInfo implements dbus.DBusObject interface
func (validator *Validator) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DBusName,
		DBusPath,
		DBusInterface,
	}
}
