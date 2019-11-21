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
	"fmt"

	"pkg.deepin.io/lib/calendar"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName = "com.deepin.api.LunarCalendar"
	dbusPath        = "/com/deepin/api/LunarCalendar"
	dbusInterface   = "com.deepin.api.LunarCalendar"
)

type Manager struct {
	service *dbusutil.Service

	methods *struct {
		GetLunarInfoBySolar   func() `in:"year,month,day" out:"lunarDay,ok"`
		GetLunarMonthCalendar func() `in:"year,month,fill" out:"lunarMonth,ok"`
		GetHuangLiDay         func() `in:"year,month,day" out:"json"`
		GetHuangLiMonth       func() `in:"year,month,fill" out:"json"`
		GetFestivalMonth      func() `in:"year,month" out:"json"`
	}
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
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
	logger.Debugf("LUNAR DATE: %v %v %v", year, month, fill)
	if info, _, ok := getLunarMonthCalendar(int(year), int(month), fill); !ok {
		return LunarMonthInfo{}, false, nil
	} else {
		logger.Debugf("Lunar Month Data: %v", info)
		return info, true, nil
	}
}

// GetHuangLiDay 获取指定公历日的黄历信息
func (m *Manager) GetHuangLiDay(year, month, day int32) (string, *dbus.Error) {
	m.service.DelayAutoQuit()
	info, ok := calendar.SolarToLunar(int(year), int(month), int(day))
	if !ok {
		return "", dbusutil.ToError(fmt.Errorf("invalid date: %d-%d-%d", year, month, day))
	}
	list := newHuangLiInfoList([]calendar.LunarDayInfo{info}, DayInfoList{DayInfo{
		Year:  year,
		Month: month,
		Day:   day,
	}})
	return list[0].String(), nil
}

// GetHuangLiMonth 获取指定公历月的黄历信息
func (m *Manager) GetHuangLiMonth(year, month int32, fill bool) (string, *dbus.Error) {
	m.service.DelayAutoQuit()
	lunarDays, solarDays, ok := getLunarMonthCalendar(int(year), int(month), fill)
	if !ok {
		return "", dbusutil.ToError(fmt.Errorf("invalid date: %d-%d", year, month))
	}
	list := newHuangLiInfoList(lunarDays.Datas, solarDays.Datas)
	var ret = HuangLiMonthInfo{
		FirstDayWeek: lunarDays.FirstDayWeek,
		Days:         lunarDays.Days,
		Datas:        list,
	}
	return ret.String(), nil
}

// GetFestivalMonth 获取指定公历月的假日信息
func (m *Manager) GetFestivalMonth(year, month int) (string, *dbus.Error) {
	list, err := newFestivalList(year, month)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	return list.String(), nil
}
