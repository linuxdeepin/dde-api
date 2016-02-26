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
	C "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) {
	C.TestingT(t)
}

func init() {
	C.Suite(NewManager())
}

func (m *Manager) TestYearValid(c *C.C) {
	c.Check(isYearValid(1340), C.Equals, false)
	c.Check(isYearValid(1990), C.Equals, true)
}

func (m *Manager) TestLeapYearMonth(c *C.C) {
	if month, ok := getLunarLeapYear(2014); !ok || month != 9 {
		c.Error("getLunarLeapYear failed")
		return
	}
}

func (m *Manager) TestLunarYearDays(c *C.C) {
	if info, days, ok := getLunarYearDays(2014); !ok {
		c.Error("getLunarYearDays failed")
		return
	} else {
		if days != 384 {
			c.Error("getLunarYearDays failed")
			return
		}

		tmp := []caDayInfo{
			caDayInfo{0, 29},
			caDayInfo{1, 30},
			caDayInfo{2, 29},
			caDayInfo{3, 30},
			caDayInfo{4, 29},
			caDayInfo{5, 30},
			caDayInfo{6, 29},
			caDayInfo{7, 30},
			caDayInfo{8, 30},
			caDayInfo{9, 29},
			caDayInfo{10, 30},
			caDayInfo{11, 29},
			caDayInfo{12, 30},
		}

		for i, v := range info {
			if v.index != tmp[i].index ||
				v.days != tmp[i].days {
				c.Error("getLunarYearDays failed")
				return
			}
		}
	}
}

func (m *Manager) TestLunarDateByBetween(c *C.C) {
	info, ok := getLunarDateByBetween(2014, 100)
	if !ok || info.Year != 2014 ||
		info.Month != 3 || info.Day != 13 {
		c.Error("getLunarDateByBetween failed")
		return
	}
}

func (m *Manager) TestLunarDateBySolar(c *C.C) {
	info, ok := getLunarDateBySolar(2014, 8, 8)
	if !ok || info.Year != 2014 ||
		info.Month != 6 || info.Day != 13 {
		c.Error("getLunarDateBySolar failed")
		return
	}
}

func (m *Manager) TestDaysBetweenSolar(c *C.C) {
	if days, ok := getDaysBetweenSolar(2014, 8, 1, 2015, 8, 1); !ok || days != 365 {
		c.Error("getDaysBetweenSolar failed")
		return
	}
}

func (m *Manager) TestDaysBetweenZheng(c *C.C) {
	if days, ok := getDaysBetweenZheng(2014, 8, 8); !ok || days != 243 {
		c.Error("getDaysBetweenZheng failed")
	}
}

func (m *Manager) TestFormatDay4(c *C.C) {
	c.Check(formatDayD4(1, 12), C.Equals, "d0112")
}

func (m *Manager) TestTermDate(c *C.C) {
	info, ok := getTermDate(2014, 8)
	if !ok || info.Year != 2014 ||
		info.Month != 5 || info.Day != 5 {
		c.Error("getTermDate failed")
		return
	}
}

func (m *Manager) TestYearZodiac(c *C.C) {
	if ret, ok := getYearZodiac(2014); !ok || ret != "马" {
		c.Error("getYearZodiac failed")
		return
	}
}

func (m *Manager) TestYearName(c *C.C) {
	if ret, ok := getLunarYearName(2014, 0); !ok || ret != "甲午" {
		c.Error("getLunarYearName failed")
		return
	}
}

func (m *Manager) TestMonthName(c *C.C) {
	if ret, ok := getLunarMonthName(2014, 5, 0); !ok || ret != "己巳" {
		c.Error("getLunarMonthName failed")
		return
	}
}

func (m *Manager) TestDayName(c *C.C) {
	if ret, ok := getLunarDayName(2014, 6, 8); !ok || ret != "庚戌" {
		c.Error("getLunarDayName failed")
		return
	}
}

func (m *Manager) TestSolarMonthDays(c *C.C) {
	if days, ok := getSolarMonthDays(2014, 8); !ok || days != 31 {
		c.Error("getSolarMonthDays failed")
		return
	}
}

func (m *Manager) TestLeapYear(c *C.C) {
	c.Check(isLeapYear(2000), C.Equals, true)
	c.Check(isLeapYear(2013), C.Equals, false)
}

func (m *Manager) TestLunarToSolar(c *C.C) {
	info, ok := lunarToSolar(2014, 8, 15)
	if !ok || info.Year != 2014 ||
		info.Month != 10 || info.Day != 8 {
		c.Error("lunarToSolar failed")
		return
	}
}

func (m *Manager) TestSolarToLunar(c *C.C) {
	info, ok := solarToLunar(2014, 9, 8)
	if !ok || info.GanZhiYear != "甲午" ||
		info.GanZhiMonth != "壬申" || info.GanZhiDay != "壬午" ||
		info.LunarMonthName != "八月" ||
		info.LunarDayName != "十五" ||
		info.LunarLeapMonth != 9 || info.Zodiac != "马" ||
		info.Term != "白露" ||
		info.SolarFestival != "国际扫盲日 国际新闻工作者日" ||
		info.LunarFestival != "中秋节" || info.Worktime != 2 {
		c.Error("solarToLunar failed")
		return
	}
}
