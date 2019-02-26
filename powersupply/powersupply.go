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

package powersupply

import (
	"errors"
	"pkg.deepin.io/gir/gudev-1.0"
	"pkg.deepin.io/dde/api/powersupply/battery"
	"strings"
)

const (
	subsystemPowerSupply = "power_supply"
	propPsyOnline        = "POWER_SUPPLY_ONLINE"
)

var errClientIsNil = errors.New("gudev.Client is nil")

func isSystemPowerSupply(dev *gudev.Device) bool {
	scope := dev.GetSysfsAttr("scope")
	switch {
	case strings.EqualFold(scope, "device"):
		return false
	case strings.EqualFold(scope, "system"):
		return true
	default:
		return true
	}
}

func IsBattery(dev *gudev.Device) bool {
	attrType := dev.GetSysfsAttr("type")
	subsystem := dev.GetSubsystem()
	return subsystem == subsystemPowerSupply && strings.EqualFold(attrType, "battery")
}

func IsSystemBattery(dev *gudev.Device) bool {
	return IsBattery(dev) && isSystemPowerSupply(dev)
}

func IsMains(dev *gudev.Device) bool {
	subsystem := dev.GetSubsystem()
	attrType := dev.GetSysfsAttr("type")
	return subsystem == subsystemPowerSupply && strings.EqualFold(attrType, "mains")
}

func GetDevices(client *gudev.Client) []*gudev.Device {
	return client.QueryBySubsystem(subsystemPowerSupply)
}

func getClient() *gudev.Client {
	return gudev.NewClient([]string{subsystemPowerSupply})
}

// return exist, online, error
func ACOnline() (bool, bool, error) {
	client := getClient()
	if client == nil {
		return false, false, errClientIsNil
	}
	defer client.Unref()
	devices := GetDevices(client)
	defer func() {
		for _, dev := range devices {
			dev.Unref()
		}
	}()
	var ac *gudev.Device
	for _, dev := range devices {
		if IsMains(dev) {
			ac = dev
			break
		}
	}
	if ac == nil {
		return false, false, nil
	}
	if !ac.HasProperty(propPsyOnline) {
		return true, false, errors.New("no property " + propPsyOnline)
	}
	return true, ac.GetPropertyAsBoolean(propPsyOnline), nil
}

func GetSystemBatteryInfos() ([]*battery.BatteryInfo, error) {
	client := getClient()
	if client == nil {
		return nil, errClientIsNil
	}
	defer client.Unref()
	devices := GetDevices(client)
	defer func() {
		for _, dev := range devices {
			dev.Unref()
		}
	}()

	var ret []*battery.BatteryInfo
	for _, bat := range devices {
		if !IsSystemBattery(bat) {
			continue
		}
		batInfo := battery.GetBatteryInfo(bat)
		ret = append(ret, batInfo)
	}
	return ret, nil
}
