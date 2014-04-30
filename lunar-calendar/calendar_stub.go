/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

type Manager struct{}

const (
	LUNAR_DEST = "com.deepin.api.LunarCalendar"
	LUNAR_PATH = "/com/deepin/api/LunarCalendar"
	LUNAR_IFC  = "com.deepin.api.LunarCalendar"
)

func (op *Manager) GetLunarDateBySolar(year, month, day int) (caYearInfo, bool) {
	if info, ok := getLunarDateBySolar(year, month, day); !ok {
		return caYearInfo{}, false
	} else {
		return info, true
	}
}

func (op *Manager) GetSolarDateByLunar(year, month, day int) (caYearInfo, bool) {
	if info, ok := lunarToSolar(year, month, day); !ok {
		return caYearInfo{}, false
	} else {
		return info, true
	}
}

func (op *Manager) GetLunarInfoBySolar(year, month, day int) (caLunarDayInfo, bool) {
	if info, ok := solarToLunar(year, month, day); !ok {
		return caLunarDayInfo{}, false
	} else {
		return info, true
	}
}

func (op *Manager) GetSolarMonthCalendar(year, month int, fill bool) (caSolarMonthInfo, bool) {
	if info, ok := getSolarCalendar(year, month, fill); !ok {
		return caSolarMonthInfo{}, false
	} else {
		return info, true
	}
}

func (op *Manager) GetLunarMonthCalendar(year, month int, fill bool) (caLunarMonthInfo, bool) {
	if info, ok := getLunarCalendar(year, month, fill); !ok {
		return caLunarMonthInfo{}, false
	} else {
		return info, true
	}
}

func (op *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		LUNAR_DEST,
		LUNAR_PATH,
		LUNAR_IFC,
	}
}

func NewManager() *Manager {
	m := &Manager{}

	return m
}
