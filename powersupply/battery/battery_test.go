/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package battery

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_parseStatus(t *testing.T) {
	Convey("parseStatus", t, func(c C) {
		c.So(parseStatus("Unknown"), ShouldEqual, StatusUnknown)
		c.So(parseStatus("Charging"), ShouldEqual, StatusCharging)
		c.So(parseStatus("Discharging"), ShouldEqual, StatusDischarging)
		c.So(parseStatus("Not charging"), ShouldEqual, StatusNotCharging)
		c.So(parseStatus("Full"), ShouldEqual, StatusFull)
		c.So(parseStatus("Other"), ShouldEqual, StatusUnknown)
	})
}

func Test_GetDisplayStatus(t *testing.T) {
	Convey("GetDisplayStatus", t, func(c C) {
		// one
		one := []Status{StatusDischarging}
		c.So(GetDisplayStatus(one), ShouldEqual, StatusDischarging)
		one[0] = StatusNotCharging
		c.So(GetDisplayStatus(one), ShouldEqual, StatusNotCharging)

		// two
		two := []Status{StatusFull, StatusFull}
		c.So(GetDisplayStatus(two), ShouldEqual, StatusFull)
		two[0] = StatusDischarging
		two[1] = StatusFull
		c.So(GetDisplayStatus(two), ShouldEqual, StatusDischarging)

		two[0] = StatusCharging
		two[1] = StatusFull
		c.So(GetDisplayStatus(two), ShouldEqual, StatusCharging)

		two[0] = StatusCharging
		two[1] = StatusDischarging
		c.So(GetDisplayStatus(two), ShouldEqual, StatusDischarging)
	})
}
