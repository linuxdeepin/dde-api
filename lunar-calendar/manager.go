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
