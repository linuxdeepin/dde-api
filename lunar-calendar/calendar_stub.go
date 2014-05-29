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

func (op *Manager) GetLunarDateBySolar(year, month, day int32) (CaYearInfo, bool, bool) {
	info, ok := getLunarDateBySolar(year, month, day)
	if !ok {
		return CaYearInfo{}, false, false
	}
	leapMonth, _ := getLunarLeapYear(year)
	isLeapMonth := false
	if leapMonth > 0 && leapMonth == info.Month {
		isLeapMonth = true
	} else if leapMonth > 0 && leapMonth > info.Month {
		info.Month += 1
	} else if leapMonth <= 0 {
		info.Month += 1
	}

	logObj.Infof("Date: %d - %d - %d\n\tIsLeapMonth: %v",
		info.Year, info.Month, info.Day, isLeapMonth)

	return info, isLeapMonth, true
}

func (op *Manager) GetSolarDateByLunar(year, month, day int32, isLeapMonth bool) (CaYearInfo, bool) {
	leapMonth, _ := getLunarLeapYear(year)
	if leapMonth <= 0 {
		isLeapMonth = false
	}
	if (leapMonth > 0 && month > leapMonth) || isLeapMonth {
		month = month
	} else {
		month -= 1
	}
	if info, ok := lunarToSolar(year, month, day); !ok {
		return CaYearInfo{}, false
	} else {
		return info, true
	}
}

func (op *Manager) GetLunarInfoBySolar(year, month, day int32) (caLunarDayInfo, bool) {
	if info, ok := solarToLunar(year, month, day); !ok {
		return caLunarDayInfo{}, false
	} else {
		return info, true
	}
}

func (op *Manager) GetSolarMonthCalendar(year, month int32, fill bool) (caSolarMonthInfo, bool) {
	logObj.Infof("SOLAR DATE: %v- %v- %v", year, month, fill)
	if info, ok := getSolarCalendar(year, month, fill); !ok {
		return caSolarMonthInfo{}, false
	} else {
		logObj.Infof("Solar Month Data: %v", info)
		return info, true
	}
}

func (op *Manager) GetLunarMonthCalendar(year, month int32, fill bool) (caLunarMonthInfo, bool) {
	logObj.Infof("LUNAR DATE: %v- %v- %v", year, month, fill)
	if info, ok := getLunarCalendar(year, month, fill); !ok {
		return caLunarMonthInfo{}, false
	} else {
		logObj.Infof("Lunar Month Data: %v", info)
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
