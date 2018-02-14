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
	"pkg.deepin.io/lib/calendar"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	DBusServiceName = "com.deepin.api.LunarCalendar"
	DBusPath        = "/com/deepin/api/LunarCalendar"
	DBusInterface   = "com.deepin.api.LunarCalendar"
)

type Manager struct {
	service *dbusutil.Service

	methods *struct {
		GetLunarInfoBySolar   func() `in:"year,month,day" out:"lunarDay,ok"`
		GetLunarMonthCalendar func() `in:"year,month,fill" out:"lunarMonth,ok"`
	}
}

func (m *Manager) GetDBusExportInfo() dbusutil.ExportInfo {
	return dbusutil.ExportInfo{
		Path:      DBusPath,
		Interface: DBusInterface,
	}
}

func NewManager(service *dbusutil.Service) *Manager {
	return &Manager{
		service: service,
	}
}

// GetLunarInfoBySolar 获取指定公历日期的农历信息
// year 公历年
// month 公历月
// day 公历日
func (m *Manager) GetLunarInfoBySolar(year, month, day int32) (calendar.LunarDayInfo, bool, *dbus.Error) {
	m.service.DelayAutoQuit()
	if info, ok := calendar.SolarToLunar(int(year), int(month), int(day)); !ok {
		return calendar.LunarDayInfo{}, false, nil
	} else {
		return info, true, nil
	}
}

// GetLunarMonthCalendar 获取指定指定公历月份的农历信息
// 第一项数据从周日开始
// year 公历年
// month 公历月
// fill 是否用上下月数据补齐首尾空缺
func (m *Manager) GetLunarMonthCalendar(year, month int32, fill bool) (LunarMonthInfo, bool, *dbus.Error) {
	m.service.DelayAutoQuit()
	logger.Infof("LUNAR DATE: %v %v %v", year, month, fill)
	if info, ok := getLunarMonthCalendar(int(year), int(month), fill); !ok {
		return LunarMonthInfo{}, false, nil
	} else {
		logger.Infof("Lunar Month Data: %v", info)
		return info, true, nil
	}
}
