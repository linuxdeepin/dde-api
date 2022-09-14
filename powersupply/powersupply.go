// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package powersupply

import (
	"errors"
	"strings"

	"github.com/linuxdeepin/dde-api/powersupply/battery"
	"github.com/linuxdeepin/go-gir/gudev-1.0"
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
