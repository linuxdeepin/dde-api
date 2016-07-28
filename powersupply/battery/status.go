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
	"strings"
)

type Status uint32

// /include/linux/power_supply.h
const (
	StatusUnknown Status = iota
	StatusCharging
	StatusDischarging
	StatusNotCharging
	StatusFull
)

var StatusMap = map[string]Status{
	"Unknown":      StatusUnknown,
	"Charging":     StatusCharging,
	"Discharging":  StatusDischarging,
	"Not charging": StatusNotCharging,
	"Full":         StatusFull,
}

func parseStatus(val string) Status {
	for k, v := range StatusMap {
		if strings.EqualFold(val, k) {
			return v
		}
	}
	return StatusUnknown
}

func (state Status) String() string {
	switch state {
	case StatusCharging:
		return "Charging"
	case StatusDischarging:
		return "Discharging"
	case StatusNotCharging:
		return "Not charging"
	case StatusFull:
		return "Full"
	default:
		return "Unknown"
	}
}

type batteryStatusSlice []Status

func (slice batteryStatusSlice) AllSame() bool {
	if len(slice) < 2 {
		return true
	}
	first := slice[0]
	for _, v := range slice[1:] {
		if v != first {
			return false
		}
	}
	return true
}

func (slice batteryStatusSlice) AnyEqual(val Status) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func GetDisplayStatus(slice []Status) Status {
	// 单个电池时, 唯一的电池的状态就是Display的状态
	if len(slice) == 1 {
		return slice[0]
	}
	statusSlice := batteryStatusSlice(slice)

	if statusSlice.AllSame() {
		return slice[0]
	}

	if statusSlice.AnyEqual(StatusDischarging) {
		return StatusDischarging
	}
	if statusSlice.AnyEqual(StatusCharging) {
		return StatusCharging
	}
	return StatusUnknown
}
