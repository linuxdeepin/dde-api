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
	"strconv"

	"pkg.deepin.io/lib/calendar"
	"pkg.deepin.io/lib/calendar/util"
)

type DayInfo struct {
	Year  int32
	Month int32
	Day   int32
}
type DayInfoList []DayInfo

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
func getLunarMonthCalendar(year, month int, fill bool) (LunarMonthInfo, SolarMonthInfo, bool) {
	solarMonth, ok := getSolarMonthCalendar(year, month, fill)
	if !ok {
		return LunarMonthInfo{}, SolarMonthInfo{}, false
	}
	var datas []calendar.LunarDayInfo
	for _, data := range solarMonth.Datas {
		lunarDay, ok := calendar.SolarToLunar(int(data.Year), int(data.Month), int(data.Day))
		if !ok {
			return LunarMonthInfo{}, SolarMonthInfo{}, false
		}
		datas = append(datas, lunarDay)
	}
	return LunarMonthInfo{solarMonth.FirstDayWeek, solarMonth.Days, datas}, solarMonth, true
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

func (days DayInfoList) GetIDList() (list []int64) {
	for _, day := range days {
		v, _ := strconv.ParseInt(fmt.Sprintf("%d%02d%02d",
			day.Year, day.Month, day.Day), 10, 64)
		list = append(list, v)
	}
	return
}
