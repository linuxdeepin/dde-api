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

const (
	DBusDest = "com.deepin.api.LunarCalendar"
	DBusPath = "/com/deepin/api/LunarCalendar"
	DBusIFC  = "com.deepin.api.LunarCalendar"
)

func (op *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest: DBusDest,
		ObjectPath: DBusPath,
		Interface: DBusIFC,
	}
}
