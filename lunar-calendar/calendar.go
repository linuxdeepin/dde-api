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
	"pkg.deepin.io/lib/calendar/util"
)

type DayInfo struct {
	Year  int32
	Month int32
	Day   int32
}

type LunarMonthInfo struct {
	FirstDayWeek int32
	Days         int32
	Datas        []calendar.LunarDayInfo
}

type SolarMonthInfo struct {
	FirstDayWeek int32
	Days         int32
	Datas        []DayInfo
}

/**
 * 获取指定公历月份的农历数据
 * year,month 公历年，月
 * fill 是否用上下月数据补齐首尾空缺，首例数据从周日开始
 */
func getLunarMonthCalendar(year, month int, fill bool) (LunarMonthInfo, bool) {
	solarMonth, _ := getSolarMonthCalendar(year, month, fill)
	datas := []calendar.LunarDayInfo{}
	for _, data := range solarMonth.Datas {
		lunarDay, ok := calendar.SolarToLunar(int(data.Year), int(data.Month), int(data.Day))
		if !ok {
			return LunarMonthInfo{}, false
		}
		datas = append(datas, lunarDay)
	}
	return LunarMonthInfo{solarMonth.FirstDayWeek, solarMonth.Days, datas}, true
}

/**
 * 公历某月日历
 * year,month 公历年，月
 * fill 是否用上下月数据补齐首尾空缺，首例数据从周日开始(7*6阵列)
 */

func getSolarMonthCalendar(year, month int, fill bool) (SolarMonthInfo, bool) {
	weekday := util.GetWeekday(year, month, 1)
	days := util.GetSolarMonthDays(year, month)
	// 本月的数据
	daysData := getMonthDays(year, month, 1, days)
	if fill {
		if weekday > 0 {
			preYear, preMonth := getPreMonth(year, month)
			// 前一个月的天数
			preDays := util.GetSolarMonthDays(preYear, preMonth)
			// 要补充上去的前一个月的数据
			preDaysData := getMonthDays(preYear, preMonth, preDays-weekday+1, preDays)
			daysData = append(preDaysData, daysData...)
		}
		nextYear, nextMonth := getNextMonth(year, month)
		count := 6*7 - (weekday + days)
		// 要补充上去的下一个月的数据
		nextDaysData := getMonthDays(nextYear, nextMonth, 1, count)
		daysData = append(daysData, nextDaysData...)
	}
	return SolarMonthInfo{int32(weekday), int32(days), daysData}, true
}

func getMonthDays(year, month, start, end int) []DayInfo {
	var list []DayInfo
	for day := start; day <= end; day++ {
		day := DayInfo{int32(year), int32(month), int32(day)}
		list = append(list, day)
	}
	return list
}

func getPreMonth(year, month int) (preYear, preMonth int) {
	if month == 1 {
		preYear = year - 1
		preMonth = 12
		return
	}
	preYear = year
	preMonth = month - 1
	return
}

func getNextMonth(year, month int) (nextYear, nextMonth int) {
	if month == 12 {
		nextYear = year + 1
		nextMonth = 1
		return
	}
	nextYear = year
	nextMonth = month + 1
	return
}
