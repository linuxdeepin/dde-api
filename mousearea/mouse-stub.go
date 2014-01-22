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
	"dlib/dbus"
)

type Manager struct {
	MouseArea          func(int32, int32)
	IdleTimeOut        func(int32, string)
	MotionCoordinate   func(string, int32, int32, int32)
	ButtonCoordinate   func(string, int32, int32, int32)
	KeyboardCoordinate func(string, int32, int32, int32)
	CancleAllArea      func(int32, int32, int32) //resolution changed
}

type IdleTick struct {
	IdleTimeOut func(string, int32)
}

const (
	MOUSE_AREA_DEST = "com.deepin.dde.api.MouseArea"
	MOUSE_AREA_PATH = "/com/deepin/dde/api/MouseArea"
	MOUSE_AREA_IFC  = "com.deepin.dde.api.MouseArea"

	IDLE_DEST = "com.deepin.dde.api.Idle"
	IDLE_PATH = "/com/deepin/dde/api/Idle"
	IDLE_IFC  = "com.deepin.dde.api.Idle"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		MOUSE_AREA_DEST,
		MOUSE_AREA_PATH,
		MOUSE_AREA_IFC,
	}
}

func (op *IdleTick) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		MOUSE_AREA_DEST,
		IDLE_PATH,
		IDLE_IFC,
	}
}
