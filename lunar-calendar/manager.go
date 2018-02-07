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
)

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

// GetLunarInfoBySolar 获取指定公历日期的农历信息
// year 公历年
// month 公历月
// day 公历日
func (m *Manager) GetLunarInfoBySolar(year, month, day int32) (calendar.LunarDayInfo, bool) {
	if info, ok := calendar.SolarToLunar(int(year), int(month), int(day)); !ok {
		return calendar.LunarDayInfo{}, false
	} else {
		return info, true
	}
}

// GetLunarMonthCalendar 获取指定指定公历月份的农历信息
// 第一项数据从周日开始
// year 公历年
// month 公历月
// fill 是否用上下月数据补齐首尾空缺
func (m *Manager) GetLunarMonthCalendar(year, month int32, fill bool) (LunarMonthInfo, bool) {
	logger.Infof("LUNAR DATE: %v %v %v", year, month, fill)
	if info, ok := getLunarMonthCalendar(int(year), int(month), fill); !ok {
		return LunarMonthInfo{}, false
	} else {
		logger.Infof("Lunar Month Data: %v", info)
		return info, true
	}
}
