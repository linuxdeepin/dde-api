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
	"fmt"
	"time"
)

func isYearValid(year int32) bool {
	if year > MaxYear || year < MinYear {
		logObj.Infof("Invalid Year: %d. Year Range(%d - %d)\n",
			year, MinYear, MaxYear)
		return false
	}

	return true
}

/**
 * 判断农历年闰月数
 * @param {Number} year 农历年
 * return 闰月数 （月份从1开始）
 */
func getLunarLeapYear(year int32) (int32, bool) {
	if !isYearValid(year) {
		return -1, false
	}

	info := lunarInfos[year-MinYear]
	//logObj.Info("LeapMonth: ", info.leapMonth)
	//logObj.Info("MonthNum: ", info.lunarMonthNum)
	//logObj.Info("ZhengMonth: ", info.zhengMonth)
	//logObj.Info("ZhengDay: ", info.zhengDay)
	return info.leapMonth, true
}

/**
 * 获取农历年每月的天数及一年的总天数
 */
func getLunarYearDays(year int32) ([]caDayInfo, int32, bool) {
	if !isYearValid(year) {
		return nil, -1, false
	}

	info := lunarInfos[year-MinYear]
	leapMonth := info.leapMonth
	monthData := fmt.Sprintf("%b", info.lunarMonthNum)
	//logObj.Info("Month Bianry before insert: ", monthData)
	tmp := ""
	l := len(monthData)
	//还原数据至16位,少于16位的在前面插入0（二进制存储时前面的0被忽略)
	for i := 0; i < 16-l; i++ {
		tmp += "0"
	}
	monthData = tmp + monthData
	//logObj.Info("Month Bianry after insert: ", monthData)

	monthNum := 0
	if leapMonth > 0 {
		monthNum = 13
	} else {
		monthNum = 12
	}

	yearDays := int32(0)
	monthDayInfos := []caDayInfo{}
	for i := 0; i < monthNum; i++ {
		tmp := caDayInfo{}
		if monthData[i] == '0' {
			yearDays += 29
			tmp.days = 29
		} else {
			yearDays += 30
			tmp.days = 30
		}
		// 让月份从1开始，不从0开始
		//t := i + 1
		// 处理闰月
		//if i >= leapMonth {
		//t -= 1
		//}
		//tmp.index = t
		tmp.index = int32(i)
		//tmp.index = int32(i) + 1
		monthDayInfos = append(monthDayInfos, tmp)
	}

	return monthDayInfos, yearDays, true
}

/**
 * 通过间隔天数查找农历日期
 */
func getLunarDateByBetween(year, between int32) (CaYearInfo, bool) {
	month := int32(-1)
	day := int32(-1)
	monthDayInfos, yearDays, ok := getLunarYearDays(year)
	if !ok {
		logObj.Info("Get Year Days Failed For Year: ", year)
		return CaYearInfo{year, month, day}, false
	}

	//leapMonth, _ := getLunarLeapYear(year)

	end := int32(0)
	if between > 0 {
		end = between
	} else {
		end = yearDays - between
	}
	//logObj.Info("Between: ", end)
	tmpDays := int32(0)
	for _, info := range monthDayInfos {
		tmpDays += info.days
		//logObj.Info("\tTmp: ", tmpDays)
		if tmpDays > end {
			month = info.index
			tmpDays = tmpDays - info.days
			break
		}
	}
	day = end - tmpDays + 1

	return CaYearInfo{year, month, day}, true
}

/**
 * 通过公历日期获取农历日期
 */
func getLunarDateBySolar(year, month, day int32) (CaYearInfo, bool) {
	if !isYearValid(year) {
		return CaYearInfo{-1, -1, -1}, false
	}

	info := lunarInfos[year-MinYear]
	zengMonth := info.zhengMonth
	zengDay := info.zhengDay
	between, _ := getDaysBetweenSolar(year, zengMonth, zengDay,
		year, month, day)
	if between == 0 { //正月初一
		return CaYearInfo{year, 1, 1}, true
	} else if between < 0 {
		year -= 1
	}
	return getLunarDateByBetween(year, int32(between))
}

/**
 * 计算两个公历日期之间的天数
 */
func getDaysBetweenSolar(year, month, day, year1, month1, day1 int32) (int64, bool) {
	date := time.Date(int(year), time.Month(month), int(day),
		0, 0, 0, 0, time.UTC).Unix()
	date1 := time.Date(int(year1), time.Month(month1), int(day1),
		0, 0, 0, 0, time.UTC).Unix()

	return (date1 - date) / 86400, true
}

/**
 * 计算农历日期离正月初一有多少天
 */
func getDaysBetweenZheng(year, month, day int32) (int32, bool) {
	monthDayInfos, _, ok := getLunarYearDays(year)
	if !ok {
		logObj.Info("Get Year Days Failed For Year: ", year)
		return -1, false
	}

	days := int32(0)
	for _, info := range monthDayInfos {
		if info.index < month {
			days += info.days
		} else {
			break
		}
	}

	return days + day - 1, true
}

func formatDayD4(month, day int32) string {
	monStr := ""
	dayStr := ""
	if month < 10 {
		monStr = fmt.Sprintf("0%d", month)
	} else {
		monStr = fmt.Sprintf("%d", month)
	}

	if day < 10 {
		dayStr = fmt.Sprintf("0%d", day)
	} else {
		dayStr = fmt.Sprintf("%d", day)
	}

	return fmt.Sprintf("d%s%s", monStr, dayStr)
}

/**
 * 某年的第n个节气为几日
 * 31556925974.7为地球公转周期，是毫秒
 * 1890年的正小寒点：01-05 16:02:31，1890年为基准点
 * year 公历年
 * n 第几个节气，从0小寒起算
 * 由于农历24节气交节时刻采用近似算法，可能存在少量误差(30分钟内)
 */
func getTermDate(year, n int32) (CaYearInfo, bool) {
	if !isYearValid(year) {
		return CaYearInfo{}, false
	}

	offset := 31556925974/1000*(int64(year)-1890) + int64(termInfo[n])*60 + time.Date(1890, 1, 5, 16, 2, 31, 0, time.UTC).Unix()
	y, m, d := time.Unix(offset, 0).Date()

	return CaYearInfo{int32(y), int32(m), int32(d)}, true
}

/**
 * 获取公历年一年的二十四节气
 * 返回key:日期，value:节气中文名
 */
func getYearTerm(year int32) map[string]string {
	logObj.Infof("YEAR: %v", year)
	res := make(map[string]string)
	month := int32(0)
	for i := int32(0); i < 24; i++ {
		if info, ok := getTermDate(year, i); !ok {
			continue
		} else {
			// 每个月中有两个节气
			month = i/2 + 1
			res[formatDayD4(month, info.Day)] = lunarData["solarTerm"][i]
		}
	}

	//logObj.Infof("TermList: %v", res)

	return res
}

/**
 * 获取生肖
 * 十二生肖，即：鼠、牛、虎、兔、龙、蛇、马、羊、猴、鸡、狗、猪
 * year: 干支所在年(默认以立春前的公历年作为基数)
 */
func getYearZodiac(year int32) (string, bool) {
	if !isYearValid(year) {
		return "", false
	}

	// 1890 属虎
	num := year - 1890 + 2 + 24 //参考干支纪年的计算，生肖对应地支
	//logObj.Info("zodiac num: ", num)
	return lunarData["zodiac"][num%12], true
}

/**
 * 计算天干地支
 * num 60进制中的位置(把60个天干地支，当成一个60进制的数)
 */
func cyclical(num int32) (string, bool) {
	return lunarData["heavenlyStems"][num%10] + lunarData["earthlyBranches"][num%12], true
}

/**
 * 获取干支纪年
 * year 干支所在年
 * offset 偏移量，默认为0，便于查询一个年跨两个干支纪年(以立春为分界线)
 */
func getLunarYearName(year, offset int32) (string, bool) {
	if !isYearValid(year) {
		return "", false
	}

	offset = offset | 0
	return cyclical(year - 1890 + 26 + offset)
}

/**
 * 获取干支纪月
 * year,month 公历年，干支所在月
 * offset 偏移量，默认为0，便于查询一个年跨两个干支纪年(以立春为分界线)
 */
func getLunarMonthName(year, month, offset int32) (string, bool) {
	if !isYearValid(year) {
		return "", false
	}

	offset = offset | 0
	return cyclical((year-1890)*12 + month + 12 + offset)
}

/**
 * 获取干支纪日
 * year,month,day 公历年，月，日
 */
func getLunarDayName(year, month, day int32) (string, bool) {
	if !isYearValid(year) {
		return "", false
	}

	//当日与1890/1/1 相差天数
	//1890/1/1与 1970/1/1 相差29219日, 1890/1/1 日柱为壬午日(60进制18)
	date := time.Date(int(year), time.Month(month), int(day),
		0, 0, 0, 0, time.UTC).Unix()
	dayCyclical := date/86400 + 29219 + 18
	return cyclical(int32(dayCyclical))
}

/**
 * 获取公历月份的天数
 */
func getSolarMonthDays(year, month int32) (int32, bool) {
	if !isYearValid(year) {
		return -1, false
	}

	monthDays := []int32{}
	if isLeapYear(year) {
		monthDays = []int32{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	} else {
		monthDays = []int32{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	}

	return monthDays[month-1], true
}

func isLeapYear(year int32) bool {
	return (year%4 == 0 && year%100 == 0) || year%400 == 0
}

/**
 * 将农历转换为公历
 * year,month,day 农历年，月(1-13，有闰月)，日
 */
func lunarToSolar(year, month, day int32) (CaYearInfo, bool) {
	if !isYearValid(year) {
		return CaYearInfo{-1, -1, -1}, false
	}

	between, _ := getDaysBetweenZheng(year, month, day)
	info := lunarInfos[year-MinYear]
	zengMonth := info.zhengMonth
	zengDay := info.zhengDay

	offDate := time.Date(int(year), time.Month(zengMonth), int(zengDay),
		0, 0, 0, 0, time.UTC).Unix() + int64(between)*86400
	newDate := time.Unix(offDate, 0)
	y, m, d := newDate.Date()

	return CaYearInfo{int32(y), int32(m), int32(d)}, true
}

/**
 * 将公历转换为农历
 */
func solarToLunar(year, month, day int32) (caLunarDayInfo, bool) {
	if !isYearValid(year) {
		return caLunarDayInfo{}, false
	}

	cacheObj.setCurrent(year)
	// 立春日期
	v, ok := cacheObj.getCache("term2")
	if !ok {
		info, _ := getTermDate(year, 2)
		v = cacheObj.setCache("term2", info)
	}
	term2 := v.(CaYearInfo)

	// 二十四节气
	v, ok = cacheObj.getCache("termList")
	if !ok {
		list := getYearTerm(year)
		v = cacheObj.setCache("termList", list)
	}
	termList := v.(map[string]string)

	//某月第一个节气开始日期
	firstTerm, _ := getTermDate(year, month*2)
	//干支所在年份
	ganZhiYear := int32(0)
	if month > 1 || (month == 1 && day >= term2.Day) {
		ganZhiYear = year
	} else {
		ganZhiYear = year - 1
	}
	//干支所在月份（以节气为界）
	ganZhiMonth := int32(0)
	if day >= firstTerm.Day {
		ganZhiMonth = month
	} else {
		ganZhiMonth = month - 1
	}

	lunarDate, _ := getLunarDateBySolar(year, month, day)
	lunarLeapMonth, _ := getLunarLeapYear(lunarDate.Year)
	lunarMonthName := ""
	if lunarLeapMonth > 0 && lunarLeapMonth == lunarDate.Month {
		lunarMonthName = "闰" + lunarData["monthCn"][lunarDate.Month-1] + "月"
	} else if lunarLeapMonth > 0 && lunarLeapMonth < lunarDate.Month {
		lunarMonthName = lunarData["monthCn"][lunarDate.Month-1] + "月"
	} else {
		lunarMonthName = lunarData["monthCn"][lunarDate.Month] + "月"
	}

	//农历节日判断
	lunarFtv := ""
	lunarTerm := ""
	lunarMonthInfos, _, _ := getLunarYearDays(lunarDate.Year)
	lunarMonthLen := int32(len(lunarMonthInfos))
	//除夕
	if int32(lunarDate.Month) == (lunarMonthLen-1) && lunarDate.Day == lunarMonthInfos[lunarMonthLen-1].days {
		lunarFtv = lunarFestival["d0100"]
		lunarTerm = termList["d0100"]
	} else if lunarLeapMonth > 0 && lunarDate.Month >= lunarLeapMonth {
		lunarFtv = lunarFestival[formatDayD4(lunarDate.Month, lunarDate.Day)]
		lunarTerm = termList[formatDayD4(lunarDate.Month, lunarDate.Day)]
	} else {
		lunarFtv = lunarFestival[formatDayD4(lunarDate.Month+1, lunarDate.Day)]
		lunarTerm = termList[formatDayD4(lunarDate.Month+1, lunarDate.Day)]
	}
	//logObj.Infof("Lunar Festival: %v, Term: %v", lunarFtv, resInfo.Term)

	// 返回结果
	resInfo := caLunarDayInfo{}

	//logObj.Info("GanZhiYear: ", ganZhiYear)
	zodiac, _ := getYearZodiac(ganZhiYear)
	resInfo.Zodiac = zodiac
	yearName, _ := getLunarYearName(ganZhiYear, 0)
	resInfo.GanZhiYear = yearName
	monthName, _ := getLunarMonthName(year, ganZhiMonth, 0)
	resInfo.GanZhiMonth = monthName
	dayName, _ := getLunarDayName(year, month, day)
	resInfo.GanZhiDay = dayName
	resInfo.Term = lunarTerm
	resInfo.LunarMonthName = lunarMonthName
	resInfo.LunarDayName = lunarData["dateCn"][lunarDate.Day-1]
	resInfo.LunarLeapMonth = lunarLeapMonth
	resInfo.SolarFestival = solarFestival[formatDayD4(month, day)]
	resInfo.LunarFestival = lunarFtv
	resInfo.Worktime = 0
	if m, ok := worktimeYearMap[fmt.Sprintf("y%d", year)]; ok {
		if v, ok := m[formatDayD4(month, day)]; ok {
			resInfo.Worktime = v
		}
	}

	return resInfo, true
}

/**
 * 获取指定公历月份的农历数据
 * year,month 公历年，月
 * fill 是否用上下月数据补齐首尾空缺，首例数据从周日开始
 */
func getLunarCalendar(year, month int32, fill bool) (caLunarMonthInfo, bool) {
	if !isYearValid(year) {
		return caLunarMonthInfo{}, false
	}

	solarData, _ := getSolarCalendar(year, month, fill)
	l := len(solarData.Datas)
	datas := []caLunarDayInfo{}
	for i := 0; i < l; i++ {
		data1 := solarData.Datas[i]
		tmp, _ := solarToLunar(data1.Year, data1.Month, data1.Day)
		datas = append(datas, tmp)
	}

	return caLunarMonthInfo{solarData.FirstDayWeek, solarData.Days, datas}, true
}

/**
 * 公历某月日历
 * year,month 公历年，月
 * fill 是否用上下月数据补齐首尾空缺，首例数据从周日开始(7*6阵列)
 */
func getSolarCalendar(year, month int32, fill bool) (caSolarMonthInfo, bool) {
	if !isYearValid(year) {
		return caSolarMonthInfo{}, false
	}

	date := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	week := int32(date.Weekday())
	days, _ := getSolarMonthDays(year, month)
	monthData := getMonthDatas(year, month, days, 1)

	if fill {
		if week > 0 { //前补
			// 获取前一个月的日期
			preYear := int32(0)
			preMonth := int32(0)
			if month-1 <= 0 {
				preYear = year - 1
				preMonth = 12
			} else {
				preMonth = month - 1
				preYear = year
			}

			preDays, _ := getSolarMonthDays(preYear, preMonth)
			preMonthData := getMonthDatas(preYear, preMonth,
				week, preDays-week+1)
			preMonthData = append(preMonthData, monthData...)
			monthData = preMonthData
		}

		if 7*6-len(monthData) != 0 { // 后补
			// 获取前一个月的日期
			nextYear := int32(0)
			nextMonth := int32(0)
			if month+1 > 12 {
				nextYear = year + 1
				nextMonth = 1
			} else {
				nextMonth = month + 1
				nextYear = year
			}

			fillLen := int32(7*6 - len(monthData))
			nextMonthData := getMonthDatas(nextYear, nextMonth,
				fillLen, 1)
			monthData = append(monthData, nextMonthData...)
		}
	}

	return caSolarMonthInfo{week, days, monthData}, true
}

func getMonthDatas(year, month, length, start int32) []CaYearInfo {
	monthDatas := []CaYearInfo{}

	if length < 1 {
		return monthDatas
	}

	k := start | 0
	for i := int32(0); i < length; i++ {
		tmp := CaYearInfo{year, month, k}
		monthDatas = append(monthDatas, tmp)
		k++
	}

	return monthDatas
}
