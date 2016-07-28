/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package battery

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_parseStatus(t *testing.T) {
	Convey("parseStatus", t, func() {
		So(parseStatus("Unknown"), ShouldEqual, StatusUnknown)
		So(parseStatus("Charging"), ShouldEqual, StatusCharging)
		So(parseStatus("Discharging"), ShouldEqual, StatusDischarging)
		So(parseStatus("Not charging"), ShouldEqual, StatusNotCharging)
		So(parseStatus("Full"), ShouldEqual, StatusFull)
		So(parseStatus("Other"), ShouldEqual, StatusUnknown)
	})
}

func Test_GetDisplayStatus(t *testing.T) {
	Convey("GetDisplayStatus", t, func() {
		// one
		one := []Status{StatusDischarging}
		So(GetDisplayStatus(one), ShouldEqual, StatusDischarging)
		one[0] = StatusNotCharging
		So(GetDisplayStatus(one), ShouldEqual, StatusNotCharging)

		// two
		two := []Status{StatusFull, StatusFull}
		So(GetDisplayStatus(two), ShouldEqual, StatusFull)
		two[0] = StatusDischarging
		two[1] = StatusFull
		So(GetDisplayStatus(two), ShouldEqual, StatusDischarging)

		two[0] = StatusCharging
		two[1] = StatusFull
		So(GetDisplayStatus(two), ShouldEqual, StatusCharging)

		two[0] = StatusCharging
		two[1] = StatusDischarging
		So(GetDisplayStatus(two), ShouldEqual, StatusDischarging)
	})
}
